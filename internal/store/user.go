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

type User struct {
	ID        string
	Name      string
	Email     string
	Role      string
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}

type UserRepository interface {
	FindOneByEmail(ctx context.Context, email string) (*User, error)
	FindOneByID(ctx context.Context, id string) (*User, error)
	FindEmailAndPassword(ctx context.Context, email string) (*User, error)
	FindAdmin(ctx context.Context, email string) (*User, error)
	Insert(ctx context.Context, user *User) error
}
