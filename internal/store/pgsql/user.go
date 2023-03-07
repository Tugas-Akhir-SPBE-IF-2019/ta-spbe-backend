package pgsql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/store"
)

type User struct {
	db *sql.DB
}

func NewUser(db *sql.DB) *User {
	return &User{db: db}
}

const userFindOneByEmailQuery = `SELECT id, email, name
	FROM "users" WHERE email = $1
`

func (s *User) FindOneByEmail(ctx context.Context, email string) (*store.UserData, error) {
	user := &store.UserData{}

	row := s.db.QueryRowContext(ctx, userFindOneByEmailQuery, email)

	err := row.Scan(
		&user.ID, &user.Email, &user.Name,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

const userFindOneByIDQuery = `SELECT id, email, name
		FROM "users" WHERE id = $1
	`

func (s *User) FindOneByID(ctx context.Context, id string) (*store.UserData, error) {
	user := &store.UserData{}

	row := s.db.QueryRowContext(ctx, userFindOneByIDQuery, id)

	err := row.Scan(
		&user.ID, &user.Email, &user.Name,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

const userFindEmailAndPasswordQuery = `SELECT id, email
		FROM "users" WHERE email = $1
	`

func (s *User) FindEmailAndPassword(ctx context.Context, email string) (*store.UserData, error) {
	user := &store.UserData{}

	row := s.db.QueryRowContext(ctx, userFindEmailAndPasswordQuery, email)

	err := row.Scan(
		&user.ID, &user.Email,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

const userFindAdminQuery = `SELECT id, email
		FROM "users" WHERE email = $1 AND role = 'admin'
	`

func (s *User) FindAdmin(ctx context.Context, email string) (*store.UserData, error) {
	user := &store.UserData{}

	row := s.db.QueryRowContext(ctx, userFindAdminQuery, email)

	err := row.Scan(
		&user.ID, &user.Email,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

const userInsert = `INSERT INTO
users(
	id, name, email, role, created_at
) values(
	$1, $2, $3, $4, $5	
)
`

func (s *User) Insert(ctx context.Context, user *store.UserData) error {
	insertStmt, err := s.db.PrepareContext(ctx, userInsert)
	if err != nil {
		return err
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}
	defer tx.Rollback()

	createdAt := time.Now().UTC()
	role := store.UserRole(store.RoleUser)
	_, err = tx.StmtContext(ctx, insertStmt).ExecContext(ctx,
		user.ID, user.Name, user.Email, role, createdAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}
	user.Role = string(role)
	user.CreatedAt = createdAt

	return nil

}
