package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Tugas Akhir Otomatisasi Penilaian Tingkat Kematangan Kebijakan SPBE IF 2019"))
	})
	fmt.Println("Server started!")
	http.ListenAndServe(":80", r)
}
