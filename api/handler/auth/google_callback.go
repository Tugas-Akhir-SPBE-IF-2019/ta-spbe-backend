package auth

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	apierror "ta-spbe-backend/api/error"
	"ta-spbe-backend/api/response"
	"ta-spbe-backend/config"
	"ta-spbe-backend/repository"
	"ta-spbe-backend/token"

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

func GoogleCallback(
	userRepo repository.UserRepository,
	cfgOAuth config.OAuth,
	jwt token.JWT,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		oAuthConfig := newOAuthConfig(cfgOAuth, SCOPE_EMAIL, SCOPE_PROFILE)

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

		user, err := userRepo.FindOneByEmail(r.Context(), payload.Email)
		if err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				response.Error(w, apierror.InternalServerError())
				return
			}

			user = &repository.User{
				ID:         uuid.NewString(),
				Name:       payload.Name,
				Email:      payload.Email,
			}

			if err := userRepo.Insert(r.Context(), user); err != nil {
				log.Println("google callback, error insert newuser: %w", err)
				response.Error(w, apierror.InternalServerError())
				return
			}

		}

		claim := token.JWTClaim{UserID: user.ID}
		accessToken, err := jwt.CreateAccessToken(claim)
		if err != nil {
			response.Error(w, apierror.InternalServerError())
			return
		}

		refreshToken, err := jwt.CreateRefreshToken(claim)
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
}
