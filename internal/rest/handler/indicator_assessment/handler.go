package indicatorassessment

import (
	"database/sql"
	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/store"
	"net/http"
)

type IndicatorAssessmentHandler interface {
	GetIndicatorAssessmentIndexList(w http.ResponseWriter, r *http.Request)
}

type indicatorAssessmentHandler struct {
	db                       *sql.DB
	indicatorAssessmentStore store.IndicatorAssessment
}

func NewIndicatorAssessmentHandler(db *sql.DB, indicatorAssessmentStore store.IndicatorAssessment) IndicatorAssessmentHandler {
	return &indicatorAssessmentHandler{
		db:                       db,
		indicatorAssessmentStore: indicatorAssessmentStore,
	}
}
