package main

import (
	"database/sql"
	"log"
	"net/http"
	"ta-spbe-backend/api/handlers"
	"ta-spbe-backend/api/routers"
	"ta-spbe-backend/services"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/lib/pq"
)

func main() {
	assessmentService := services.NewAssessmentService()
	assessmentHandler := handlers.NewAssessmentHandler(assessmentService)

	connStr := "postgresql://ta-spbe:ta-spbe@db/ta-spbe?sslmode=disable"
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Route("/assessment", func(r chi.Router) {
		r.Mount("/", routers.AssessmentRouter(assessmentHandler))
	})
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Tugas Akhir Otomatisasi Penilaian Tingkat Kematangan Kebijakan SPBE IF 2019"))
	})

	port := 80
	log.Printf("Server started on port:%d!", port)

	http.ListenAndServe(":80", r)
}
