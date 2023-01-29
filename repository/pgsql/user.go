package pgsql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"ta-spbe-backend/repository"
)

type userRepo struct {
	db *sql.DB
	ps *userPreparedStatement
}

type userPreparedStatement struct {
	findOneByID          *sql.Stmt
	findOneByEmail       *sql.Stmt
	findEmailAndPassword *sql.Stmt
	findAdmin            *sql.Stmt
	insert               *sql.Stmt
}

func NewUserRepo(db *sql.DB) (repository.UserRepository, error) {
	repo := &userRepo{
		db: db,
		ps: &userPreparedStatement{},
	}

	err := repo.prepareStatement()
	if err != nil {
		return nil, fmt.Errorf("User Repo: %w", err)
	}

	return repo, nil
}

func (r *userRepo) prepareStatement() error {
	var err error
	repoName := "userRepo"

	if r.ps.findOneByID, err = prepareStmt(r.db, repoName, "findOneByID", userFindOneByIDQuery); err != nil {
		return err
	}

	if r.ps.findOneByEmail, err = prepareStmt(r.db, repoName, "findOneByEmail", userFindOneByEmailQuery); err != nil {
		return err
	}

	if r.ps.findEmailAndPassword, err = prepareStmt(r.db, repoName, "findEmailAndPassword", userFindEmailAndPasswordQuery); err != nil {
		return err
	}

	if r.ps.findAdmin, err = prepareStmt(r.db, repoName, "findAdmin", userFindAdminQuery); err != nil {
		return err
	}

	if r.ps.insert, err = prepareStmt(r.db, repoName, "insert", userInsert); err != nil {
		return err
	}

	return nil
}

const userFindOneByEmailQuery = `SELECT id, email, name
	FROM "users" WHERE email = $1
`

func (r *userRepo) FindOneByEmail(ctx context.Context, email string) (*repository.User, error) {
	user := &repository.User{}

	row := r.ps.findOneByEmail.QueryRowContext(ctx, email)

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

func (r *userRepo) FindOneByID(ctx context.Context, id string) (*repository.User, error) {
	user := &repository.User{}

	row := r.ps.findOneByID.QueryRowContext(ctx, id)

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

func (r *userRepo) FindEmailAndPassword(ctx context.Context, email string) (*repository.User, error) {
	user := &repository.User{}

	row := r.ps.findEmailAndPassword.QueryRowContext(ctx, email)

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

func (r *userRepo) FindAdmin(ctx context.Context, email string) (*repository.User, error) {
	user := &repository.User{}

	row := r.ps.findAdmin.QueryRowContext(ctx, email)

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

func (r *userRepo) Insert(ctx context.Context, user *repository.User) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}
	defer tx.Rollback()

	createdAt := time.Now().UTC()
	role := repository.UserRole(repository.RoleUser)
	_, err = tx.StmtContext(ctx, r.ps.insert).ExecContext(ctx,
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
