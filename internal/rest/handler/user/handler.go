package user

import (
	"database/sql"
	"net/http"

	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/config"
	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/store"
	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/pkg/filesystem"
)

type UserHandler interface {
	GetUserProfile(w http.ResponseWriter, r *http.Request)
	GetUserEvaluationData(w http.ResponseWriter, r *http.Request)
	GetUserJobData(w http.ResponseWriter, r *http.Request)
	AddUserEvaluationData(w http.ResponseWriter, r *http.Request)
	AddUserJobData(w http.ResponseWriter, r *http.Request)
	UpdateUserProfile(w http.ResponseWriter, r *http.Request)
}

type userHandler struct {
	apiCfg    config.API
	db        *sql.DB
	userStore store.User
	filesystemClient filesystem.Client
}

func NewUserHandler(apiCfg config.API, db *sql.DB, userStore store.User, fileSystemClient filesystem.Client) UserHandler {
	return &userHandler{
		apiCfg:    apiCfg,
		db:        db,
		userStore: userStore,
		filesystemClient: fileSystemClient,
	}
}
