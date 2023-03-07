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
)

type AssessmentHandler interface {
	GetSPBEAssessmentList(w http.ResponseWriter, r *http.Request)
	UploadSPBEDocument(w http.ResponseWriter, r *http.Request)
}

type assessmentHandler struct {
	apiCfg           config.API
	db               *sql.DB
	assessmentStore  store.Assessment
	userStore        store.User
	smtpMailer       smtpmailer.Client
	filesystemClient filesystem.Client
	jsonClient       jsonmanipulator.Client
	messageQueue     messagequeue.Client
}

func NewAssessmentHandler(db *sql.DB, assessmentStore store.Assessment, apiCfg config.API, userStore store.User, smtpMailer smtpmailer.Client, filesystemClient filesystem.Client, jsonClient jsonmanipulator.Client, messageQueue messagequeue.Client) AssessmentHandler {
	return &assessmentHandler{
		apiCfg:           apiCfg,
		db:               db,
		assessmentStore:  assessmentStore,
		userStore:        userStore,
		smtpMailer:       smtpMailer,
		filesystemClient: filesystemClient,
		jsonClient:       jsonClient,
		messageQueue:     messageQueue,
	}
}
