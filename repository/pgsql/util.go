package pgsql

import (
	"database/sql"
	"fmt"
)

func prepareStmt(db *sql.DB, repoName, queryName, sql string) (*sql.Stmt, error) {
	stmt, err := db.Prepare(sql)
	if err != nil {
		return nil, fmt.Errorf("%s: error prepare statement %s: %w", repoName, queryName, err)
	}
	return stmt, nil
}
