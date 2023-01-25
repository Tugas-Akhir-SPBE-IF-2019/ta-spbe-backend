package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"ta-spbe-backend/api/handlers"
	"ta-spbe-backend/api/routers"
	"ta-spbe-backend/config"
	"ta-spbe-backend/database"
	"ta-spbe-backend/services"

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

	assessmentService := services.NewAssessmentService()
	assessmentHandler := handlers.NewAssessmentHandler(assessmentService)

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",
		url.QueryEscape(cfg.DB.Username),
		url.QueryEscape(cfg.DB.Password),
		cfg.DB.Host,
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

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Route("/assessment", func(r chi.Router) {
		r.Mount("/", routers.AssessmentRouter(assessmentHandler))
	})
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Tugas Akhir Otomatisasi Penilaian Tingkat Kematangan Kebijakan SPBE IF 2019"))
	})

	log.Printf("Server is listening on port %d", cfg.API.Port)
	http.ListenAndServe(fmt.Sprintf(":%d", cfg.API.Port), r)
}
