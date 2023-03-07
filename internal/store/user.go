package store

import (
	"context"
	"database/sql"
	"time"
)

type UserRole string

const (
	RoleAdmin = UserRole("ADMIN")
	RoleUser  = UserRole("USER")
)

type UserData struct {
	ID        string
	Name      string
	Email     string
	Role      string
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}

type User interface {
	FindOneByEmail(ctx context.Context, email string) (*UserData, error)
	FindOneByID(ctx context.Context, id string) (*UserData, error)
	FindEmailAndPassword(ctx context.Context, email string) (*UserData, error)
	FindAdmin(ctx context.Context, email string) (*UserData, error)
	Insert(ctx context.Context, user *UserData) error
}
