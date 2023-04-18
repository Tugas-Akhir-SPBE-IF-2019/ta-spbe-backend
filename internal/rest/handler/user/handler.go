package user

import (
	"database/sql"
	"net/http"

	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/config"
	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/store"
	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/pkg/filesystem"
	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/pkg/smtpmailer"
)

type UserHandler interface {
	GetUserProfile(w http.ResponseWriter, r *http.Request)
	GetUserEvaluationData(w http.ResponseWriter, r *http.Request)
	GetUserJobData(w http.ResponseWriter, r *http.Request)
	AddUserEvaluationData(w http.ResponseWriter, r *http.Request)
	AddUserJobData(w http.ResponseWriter, r *http.Request)
	AddUserInstitutionData(w http.ResponseWriter, r *http.Request)
	UpdateUserProfile(w http.ResponseWriter, r *http.Request)
	GetUserCurrentInstitutionData(w http.ResponseWriter, r *http.Request)
	DeleteUserCurrentInstitutionData(w http.ResponseWriter, r *http.Request)
	VerifyUserCurrentInstitutionData(w http.ResponseWriter, r *http.Request)
	RejectUserCurrentInstitutionData(w http.ResponseWriter, r *http.Request)
}

type userHandler struct {
	apiCfg           config.API
	db               *sql.DB
	userStore        store.User
	institutionStore store.Institution
	filesystemClient filesystem.Client
	smtpMailer       smtpmailer.Client
}

func NewUserHandler(apiCfg config.API, db *sql.DB, userStore store.User, institutionStore store.Institution, fileSystemClient filesystem.Client, smtpMailer smtpmailer.Client) UserHandler {
	return &userHandler{
		apiCfg:           apiCfg,
		db:               db,
		userStore:        userStore,
		institutionStore: institutionStore,
		filesystemClient: fileSystemClient,
		smtpMailer:       smtpMailer,
	}
}
