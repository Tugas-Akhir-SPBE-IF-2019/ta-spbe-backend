package auth

import (
	"database/sql"
	"net/http"

	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/config"
	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/store"
	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/pkg/token"
)

type AuthHandler interface {
	Google(w http.ResponseWriter, r *http.Request)
	GoogleCallback(w http.ResponseWriter, r *http.Request)
	GoogleValidate(w http.ResponseWriter, r *http.Request)
}

type authHandler struct {
	db        *sql.DB
	userStore store.User
	cfgOAuth  config.OAuth
	jwt       token.JWT
}

func NewAuthHandler(db *sql.DB, userStore store.User, cfgOAuth config.OAuth, jwt token.JWT) AuthHandler {
	return &authHandler{
		db:        db,
		userStore: userStore,
		cfgOAuth:  cfgOAuth,
		jwt:       jwt,
	}
}
