package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/config"
	userCtx "github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/rest/handler/context"

	apierror "github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/rest/error"
	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/rest/response"
	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/store"
	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/pkg/token"
)

func JWTAuth(jwt token.JWT, devCfg config.DevSettings) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var ctx context.Context
			if devCfg.Auth || r.Header.Get("Authorization") != "" {
				authHeader := r.Header.Get("Authorization")
				if authHeader == "" {
					response.Error(w, apierror.UnauthorizedError("Authorization is missing"))
					return
				}

				authSplit := strings.Split(authHeader, " ")
				if authSplit[0] != "Bearer" {
					response.Error(w, apierror.UnauthorizedError("Wrong authorization scheme"))
					return
				}

				tokenString := authSplit[1]
				claim, err := jwt.GetClaims(tokenString)
				if err != nil {
					response.Error(w, apierror.UnauthorizedError("Unauthorized"))
					return
				}

				ctx = context.WithValue(r.Context(), userCtx.UserCtxKey, userCtx.UserCtx{
					ID:   claim.UserID,
					Role: store.UserRole(claim.Role),
				})
			} else {
				ctx = context.WithValue(r.Context(), userCtx.UserCtxKey, userCtx.UserCtx{
					ID:   "ccd52961-fa4e-43ba-a6df-a4c97849d899",
					Role: store.RoleAdmin,
				})
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
