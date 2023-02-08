package context

import "ta-spbe-backend/repository"

type CtxKey string

const UserCtxKey = CtxKey("user_ctx")

type UserCtx struct {
	ID   string
	Role repository.UserRole
}
