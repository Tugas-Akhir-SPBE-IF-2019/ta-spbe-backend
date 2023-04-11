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

type GoogleValidateRequest struct {
	AccessToken string `json:"access_token"`
}

func (handler *authHandler) GoogleValidate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := GoogleValidateRequest{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, apierror.BadRequestError(err.Error()))
		return
	}

	res, err := http.Get(fmt.Sprintf("https://www.googleapis.com/oauth2/v2/userinfo?access_token=%s", req.AccessToken))
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

	user, err := handler.userStore.FindOneByEmail(ctx, payload.Email)
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

		if err := handler.userStore.Insert(ctx, user); err != nil {
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
