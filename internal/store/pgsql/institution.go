package pgsql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/store"
)

type Institution struct {
	db *sql.DB
}

func NewInstitution(db *sql.DB) *Institution {
	return &Institution{db: db}
}

const institutionFindAllQuery = `SELECT i.id, i.name, ic.category
		FROM institution i
		LEFT JOIN institution_category ic 
		ON i.category = ic.id `

func (s *Institution) FindAll(ctx context.Context, institutionName string) ([]*store.InstitutionDetail, error) {
	institutionList := []*store.InstitutionDetail{}
	var queryKeys []string
	var queryParams []interface{}

	query := institutionFindAllQuery
	if institutionName != "" {
		queryKeys = append(queryKeys, "Name")
		queryParams = append(queryParams, institutionName)
	}

	for index, key := range queryKeys {
		if index == 0 {
			query = query + "WHERE "
		} else {
			query = query + "AND "
		}

		switch key {
		case "Name":
			query = query + fmt.Sprintf(`i.name ILIKE '%%' || $%d || '%%' `, index+1)
		}
	}

	rows, err := s.db.QueryContext(ctx, query, queryParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		institution := &store.InstitutionDetail{}
		err := rows.Scan(
			&institution.ID,
			&institution.Name,
			&institution.Category,
		)
		if err != nil {
			return nil, err
		}
		institutionList = append(institutionList, institution)
	}

	return institutionList, nil
}

func (s *Institution) FindAllPagination(ctx context.Context, offset int, limit int, institutionName string) ([]*store.InstitutionDetail, error) {
	institutionList := []*store.InstitutionDetail{}
	var queryKeys []string
	var queryParams []interface{}

	queryParams = append(queryParams, offset, limit)

	query := institutionFindAllQuery
	if institutionName != "" {
		queryKeys = append(queryKeys, "Name")
		queryParams = append(queryParams, institutionName)
	}

	for index, key := range queryKeys {
		if index == 0 {
			query = query + "WHERE "
		} else {
			query = query + "AND "
		}

		// +3 is used because $1 and $2 are already used for pagination
		switch key {
		case "Name":
			query = query + fmt.Sprintf(`i.name ILIKE '%%' || $%d || '%%' `, index+3)
		}
	}

	query = query + `LIMIT $2 OFFSET $1 `

	rows, err := s.db.QueryContext(ctx, query, queryParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		institution := &store.InstitutionDetail{}
		err := rows.Scan(
			&institution.ID,
			&institution.Name,
			&institution.Category,
		)
		if err != nil {
			return nil, err
		}
		institutionList = append(institutionList, institution)
	}

	return institutionList, nil
}

const institutionFindByInstitutionNameQuery = `SELECT i.id, i.name, i.category
		FROM institution i
		WHERE i.name = $1`

func (s *Institution) FindByInstitutionName(ctx context.Context, institutionName string) (*store.InstitutionDetail, error) {
	institution := &store.InstitutionDetail{}

	row := s.db.QueryRowContext(ctx, institutionFindByInstitutionNameQuery, institutionName)
	err := row.Scan(
		&institution.ID, &institution.Name, &institution.Category,
	)
	if err != nil {
		return institution, err
	}

	return institution, nil
}

const insertInstitutionQuery = `INSERT INTO
institution(
	institution (category, name)
) values(
	$1, $2
)
`

func (s *Institution) Insert(ctx context.Context, institution *store.InstitutionDetail) error {
	insertStmt, err := s.db.PrepareContext(ctx, insertInstitutionQuery)
	if err != nil {
		return err
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.StmtContext(ctx, insertStmt).ExecContext(ctx,
		1, institution.Name, // category is still hardcoded; need change by making new category e.g USER_GENERATED
	)
	if err != nil {
		return fmt.Errorf("failed to insert: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}

	return nil
}
