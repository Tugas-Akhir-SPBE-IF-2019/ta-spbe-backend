package auth

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	apierror "github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/rest/error"
	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/rest/response"
	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/store"
	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/pkg/token"

	"github.com/google/uuid"
)

type RegisterGoogleCallbackResponse struct {
	Message string `json:"message"`
}

type UserInfo struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type LoginResponse struct {
	AccessToken  TokenItem `json:"access_token"`
	RefreshToken TokenItem `json:"refresh_token"`
}

func (handler *authHandler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	stateCookie, err := r.Cookie("state")
	if err != nil {
		log.Println("google callback, missing cookie: %w", err)
		response.Error(w, apierror.UnauthorizedError("Unauthorized"))
		return
	}

	if state := r.FormValue("state"); state != stateCookie.Value {
		log.Println("google callback, state from form doesn't match with state from cookie")
		response.Error(w, apierror.UnauthorizedError("Unauthorized"))
		return
	}

	oAuthConfig := newOAuthConfig(handler.cfgOAuth, SCOPE_EMAIL, SCOPE_PROFILE)

	code := r.FormValue("code")
	tokenOAuth, err := oAuthConfig.Exchange(r.Context(), code)
	if err != nil {
		log.Println("google callback, error exchange code for token: %w", err)
		response.Error(w, apierror.InternalServerError())
		return
	}

	res, err := http.Get(fmt.Sprintf("https://www.googleapis.com/oauth2/v2/userinfo?access_token=%s", tokenOAuth.AccessToken))
	if err != nil {
		log.Println("google callback, error get user info: %w", err)
		response.Error(w, apierror.InternalServerError())
		return
	}

	payload := UserInfo{}
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		log.Println("google callback, error decode payload: %w", err)
		response.Error(w, apierror.InternalServerError())
		return
	}

	user, err := handler.userStore.FindOneByEmail(r.Context(), payload.Email)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			response.Error(w, apierror.InternalServerError())
			return
		}

		user = &store.UserData{
			ID:    uuid.NewString(),
			Name:  payload.Name,
			Email: payload.Email,
		}

		if err := handler.userStore.Insert(r.Context(), user); err != nil {
			log.Println("google callback, error insert newuser: %w", err)
			response.Error(w, apierror.InternalServerError())
			return
		}

	}

	claim := token.JWTClaim{UserID: user.ID}
	accessToken, err := handler.jwt.CreateAccessToken(claim)
	if err != nil {
		response.Error(w, apierror.InternalServerError())
		return
	}

	refreshToken, err := handler.jwt.CreateRefreshToken(claim)
	if err != nil {
		response.Error(w, apierror.InternalServerError())
		return
	}

	resp := LoginResponse{
		AccessToken: TokenItem{
			Token:     accessToken.Token,
			Scheme:    accessToken.Scheme,
			ExpiresAt: accessToken.ExpiresAt,
		},
		RefreshToken: TokenItem{
			Token:     refreshToken.Token,
			Scheme:    refreshToken.Scheme,
			ExpiresAt: refreshToken.ExpiresAt,
		},
	}

	response.Respond(w, http.StatusCreated, resp)
}
