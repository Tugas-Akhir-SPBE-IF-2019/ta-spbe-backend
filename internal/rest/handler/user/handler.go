package user

import (
	"database/sql"
	"net/http"

	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/store"
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
	db        *sql.DB
	userStore store.User
}

func NewUserHandler(db *sql.DB, userStore store.User) UserHandler {
	return &userHandler{
		db:        db,
		userStore: userStore,
	}
}
