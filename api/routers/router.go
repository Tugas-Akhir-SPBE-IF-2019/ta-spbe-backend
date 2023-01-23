package routers

import (
	"net/http"
	"ta-spbe-backend/api/handlers"

	"github.com/go-chi/chi/v5"
)

func AssessmentRouter(handler handlers.AssessmentHandler) http.Handler {
	r := chi.NewRouter()

	r.Get("/", handler.GetAssessmentList)
	r.Get("/index", handler.GetAssessmentIndexList)
	r.Get("/{id}", handler.GetAssessmentResult)
	r.Post("/documents/upload", handler.UploadAssessmentDocument)
	r.Patch("/{id}/validate", handler.ValidateAssessmentResult)

	return r
}
