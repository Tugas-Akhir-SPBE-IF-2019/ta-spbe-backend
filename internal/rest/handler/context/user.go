package context

import (
	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/store"
)

type CtxKey string

const UserCtxKey = CtxKey("user_ctx")

type UserCtx struct {
	ID   string
	Role store.UserRole
}
