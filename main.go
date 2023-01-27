package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	assessmenthandler "ta-spbe-backend/api/handler/assessment"
	indicatorassessmenthandler "ta-spbe-backend/api/handler/indicator_assessment"
	"ta-spbe-backend/config"
	"ta-spbe-backend/database"
	"ta-spbe-backend/repository/pgsql"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
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

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Tugas Akhir Otomatisasi Penilaian Tingkat Kematangan Kebijakan SPBE IF 2019"))
	})

	r.Route("/assessments", func(r chi.Router) {
		r.Get("/", assessmenthandler.GetSPBEAssessmentList(assessmentRepo))
		r.Get("/{id}", indicatorassessmenthandler.GetIndicatorAssessmentResult(indicatorAssessmentRepo))
		r.Get("/index", indicatorassessmenthandler.GetIndicatorAssessmentIndexList(indicatorAssessmentRepo))
		r.Post("/documents/upload", assessmenthandler.UploadSPBEDocument(assessmentRepo, cfg.API))
		r.Patch("/{id}/validate", indicatorassessmenthandler.ValidateIndicatorAssessmentResult(indicatorAssessmentRepo))
	})

	//static file serve (for testing purpose only)
	fs := http.FileServer(http.Dir("static/supporting-documents"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	log.Printf("Server is listening on port %d", cfg.API.Port)
	http.ListenAndServe(fmt.Sprintf(":%d", cfg.API.Port), r)
}
