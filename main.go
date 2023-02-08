package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	assessmenthandler "ta-spbe-backend/api/handler/assessment"
	authhandler "ta-spbe-backend/api/handler/auth"
	indicatorassessmenthandler "ta-spbe-backend/api/handler/indicator_assessment"
	apimiddleware "ta-spbe-backend/api/middleware"
	filesystem "ta-spbe-backend/service/file-system"
	jsonmanipulator "ta-spbe-backend/service/json-manipulator"
	mailerservice "ta-spbe-backend/service/mailer"
	messagequeue "ta-spbe-backend/service/message-queue"

	"ta-spbe-backend/config"
	"ta-spbe-backend/database"
	"ta-spbe-backend/repository/pgsql"
	"ta-spbe-backend/token"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
	"github.com/nsqio/go-nsq"
)

func main() {
	cfgPath := flag.String("c", "config.toml", "path to config file")
	cfg, err := config.LoadEnvFromFile(*cfgPath)
	if err != nil {
		log.Fatalln(err)
	}

	migrate := flag.Bool("migrate", cfg.DB.Migration, "do migration")

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		url.QueryEscape(cfg.DB.Username),
		url.QueryEscape(cfg.DB.Password),
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.Database,
	)

	db, err := sql.Open("postgres", dsn)

	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	if *migrate {
		err := database.Migrate(db, cfg.DB.Database)
		if err != nil {
			log.Fatalln(err)
		}
	}

	assessmentRepo, err := pgsql.NewAssessmentRepo(db)
	if err != nil {
		log.Fatalln(err)
	}

	indicatorAssessmentRepo, err := pgsql.NewIndicatorAssessmentRepo(db)
	if err != nil {
		log.Fatalln(err)
	}

	userRepo, err := pgsql.NewUserRepo(db)
	if err != nil {
		log.Fatalln(err)
	}

	nsqConfig := nsq.NewConfig()
	nsqdAddress := fmt.Sprintf("%s:%d", cfg.MessageBroker.Host, cfg.MessageBroker.Port)
	nsqProducer, err := nsq.NewProducer(nsqdAddress, nsqConfig)
	if err != nil {
		log.Fatalln(err)
	}

	messageQueue := &messagequeue.NSQ{
		Producer: nsqProducer,
	}
	mailerService, err := mailerservice.NewSimpleMailer(cfg.SMTPClient)
	if err != nil {
		log.Fatalln(err)
	}
	jsonEC := jsonmanipulator.EncoderDecoder{}

	fileSystemIO := filesystem.IO{}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Tugas Akhir Otomatisasi Penilaian Tingkat Kematangan Kebijakan SPBE IF 2019"))
	})

	jwt := token.NewJWT(cfg.JWT)
	r.Route("/auth", func(r chi.Router) {
		r.Get("/", token.HandleMain)
		r.Post("/google", authhandler.Google(cfg.OAuth))
		r.Get("/google/callback", authhandler.GoogleCallback(userRepo, cfg.OAuth, jwt))
	})

	authMW := apimiddleware.Auth(jwt)
	r.Get("/assessments/index", indicatorassessmenthandler.GetIndicatorAssessmentIndexList(indicatorAssessmentRepo))
	r.Route("/assessments", func(r chi.Router) {
		r.Use(authMW)
		r.Get("/", assessmenthandler.GetSPBEAssessmentList(assessmentRepo))
		r.Get("/{id}", indicatorassessmenthandler.GetIndicatorAssessmentResult(indicatorAssessmentRepo))
		r.Post("/documents/upload", assessmenthandler.UploadSPBEDocument(assessmentRepo, userRepo, messageQueue, mailerService, fileSystemIO, jsonEC, cfg.API, cfg.SMTPClient))
		r.Patch("/{id}/validate", indicatorassessmenthandler.ValidateIndicatorAssessmentResult(indicatorAssessmentRepo))
	})

	r.Post("/assessments/result/callback", indicatorassessmenthandler.ResultCallback(indicatorAssessmentRepo, userRepo, mailerService))

	//static file serve (for testing purpose only)
	fs := http.FileServer(http.Dir("static/supporting-documents"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	log.Printf("Server is listening on port %d", cfg.API.Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.API.Port), r); err != nil {
		log.Println(err)

	}
}
