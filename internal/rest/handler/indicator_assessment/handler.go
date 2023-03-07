package indicatorassessment

import (
	"database/sql"
	"net/http"

	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/store"
	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/pkg/smtpmailer"
	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/pkg/whatsapp"
)

type IndicatorAssessmentHandler interface {
	GetIndicatorAssessmentIndexList(w http.ResponseWriter, r *http.Request)
	GetIndicatorAssessmentResultGetIndicatorAssessmentIndexList(w http.ResponseWriter, r *http.Request)
	ValidateIndicatorAssessmentResult(w http.ResponseWriter, r *http.Request)
	ResultCallback(w http.ResponseWriter, r *http.Request)
}

type indicatorAssessmentHandler struct {
	db                       *sql.DB
	indicatorAssessmentStore store.IndicatorAssessment
	userStore                store.User
	smtpMailer               smtpmailer.Client
	waClient                 whatsapp.Client
}

func NewIndicatorAssessmentHandler(db *sql.DB, indicatorAssessmentStore store.IndicatorAssessment, userStore store.User, smtpMailer smtpmailer.Client, waClient whatsapp.Client) IndicatorAssessmentHandler {
	return &indicatorAssessmentHandler{
		db:                       db,
		indicatorAssessmentStore: indicatorAssessmentStore,
		userStore:                userStore,
		smtpMailer:               smtpMailer,
		waClient:                 waClient,
	}
}
