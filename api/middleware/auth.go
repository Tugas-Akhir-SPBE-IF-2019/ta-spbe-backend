package middleware

import (
	"context"
	"net/http"
	"strings"

	apierror "ta-spbe-backend/api/error"
	userCtx "ta-spbe-backend/api/handler/context"
	"ta-spbe-backend/api/response"
	"ta-spbe-backend/repository"
	"ta-spbe-backend/token"
)

func Auth(jwt token.JWT) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

			ctx := context.WithValue(r.Context(), userCtx.UserCtxKey, userCtx.UserCtx{
				ID:   claim.UserID,
				Role: repository.UserRole(claim.Role),
			})
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
