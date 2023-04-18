package assessment

import (
	"database/sql"
	"net/http"

	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/config"
	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/store"
	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/pkg/filesystem"
	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/pkg/jsonmanipulator"
	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/pkg/messagequeue"
	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/pkg/smtpmailer"
	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/pkg/whatsapp"
)

type AssessmentHandler interface {
	GetSPBEAssessmentList(w http.ResponseWriter, r *http.Request)
	UploadSPBEDocument(w http.ResponseWriter, r *http.Request)
	GetSPBEAssessmentStatusHistory(w http.ResponseWriter, r *http.Request)
	GetSPBEAssessmentDocumentList(w http.ResponseWriter, r *http.Request)
	ResultCallback(w http.ResponseWriter, r *http.Request)
}

type assessmentHandler struct {
	apiCfg           config.API
	db               *sql.DB
	assessmentStore  store.Assessment
	indicatoreAssessmentStore store.IndicatorAssessment
	userStore        store.User
	smtpMailer       smtpmailer.Client
	filesystemClient filesystem.Client
	jsonClient       jsonmanipulator.Client
	messageQueue     messagequeue.Client
	waClient         whatsapp.Client
}

func NewAssessmentHandler(db *sql.DB, assessmentStore store.Assessment, indicatorAssessmentStore store.IndicatorAssessment,apiCfg config.API, userStore store.User, smtpMailer smtpmailer.Client, filesystemClient filesystem.Client, jsonClient jsonmanipulator.Client, messageQueue messagequeue.Client, waClient whatsapp.Client) AssessmentHandler {
	return &assessmentHandler{
		apiCfg:           apiCfg,
		db:               db,
		assessmentStore:  assessmentStore,
		indicatoreAssessmentStore: indicatorAssessmentStore,
		userStore:        userStore,
		smtpMailer:       smtpMailer,
		filesystemClient: filesystemClient,
		jsonClient:       jsonClient,
		messageQueue:     messageQueue,
		waClient:         waClient,
	}
}
