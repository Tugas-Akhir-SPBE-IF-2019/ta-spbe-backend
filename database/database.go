package database

import (
	"database/sql"
	"embed"
	"fmt"
	"net/url"
	"ta-spbe-backend/config"

	"github.com/golang-migrate/migrate/v4/source/iofs"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
)

func Open(cfg config.Postgres) (*sql.DB, error) {
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s",
		url.QueryEscape(cfg.Username),
		url.QueryEscape(cfg.Password),
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)

	conCfg, err := pgx.ParseConfig(connString)

	if err != nil {
		return nil, fmt.Errorf("Parse config failed: %w", err)
	}

	db := stdlib.OpenDB(*conCfg)

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("Ping failed: %w", err)
	}

	return db, nil
}

//go:embed migration/*.sql
var MigrationFiles embed.FS

const MigrationFilesPath = "migration"

func Migrate(db *sql.DB, databaseName string) error {
	d, err := iofs.New(MigrationFiles, MigrationFilesPath)
	if err != nil {
		return fmt.Errorf("failed to prepare migration files: %w", err)
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("Migrate failed: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", d, databaseName, driver)
	if err != nil {
		return fmt.Errorf("Migrate failed: %w", err)
	}

	if err = m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			return nil
		}
		return fmt.Errorf("Migrate failed: %w", err)
	}

	return nil
}
