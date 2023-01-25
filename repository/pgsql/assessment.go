package pgsql

import (
	"context"
	"database/sql"
	"log"
	"ta-spbe-backend/repository"
)

type assessmentRepo struct {
	db *sql.DB
	ps map[string]*sql.Stmt
}

func NewAssessmentRepo(db *sql.DB) (repository.AssessmentRepository, error) {
	ps := make(map[string]*sql.Stmt, len(assessmentQueries))
	for key, query := range assessmentQueries {
		stmt, err := prepareStmt(db, "assessmentRepo", key, query)
		if err != nil {
			return nil, err
		}
		ps[key] = stmt
	}

	return &assessmentRepo{db, ps}, nil
}

var assessmentQueries = map[string]string{
	assessmentFindAll: assessmentFindAllQuery,
	assessmentFindAllPagination: assessmentFindAllPaginationQuery,
}

const assessmentFindAll = "findAll"
const assessmentFindAllQuery = `SELECT a.id, a.institution_name, a.status, a.created_at
		FROM assessments a`

func (r *assessmentRepo) FindAll(ctx context.Context) ([]*repository.AssessmentDetail, error) {
	assessmentList := []*repository.AssessmentDetail{}

	rows, err := r.ps[assessmentFindAll].QueryContext(ctx)
	if err != nil {
		log.Println("assessment sql repo query context error: %w", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		assessment := &repository.AssessmentDetail{}
		err := rows.Scan(
			&assessment.Id,
			&assessment.InstitutionName,
			&assessment.Status,
			&assessment.SubmittedDate,
		)
		if err != nil {
			log.Println("assessment sql repo scan error: %w", err)
			return nil, err
		}
		assessmentList = append(assessmentList, assessment)
	}

	return assessmentList, nil
}


const assessmentFindAllPagination = "findAllPagination"
const assessmentFindAllPaginationQuery = `SELECT a.id, a.institution_name, a.status, a.created_at
		FROM assessments a LIMIT $2 OFFSET $1`

func (r *assessmentRepo) FindAllPagination(ctx context.Context, offset int, limit int) ([]*repository.AssessmentDetail, error) {
	assessmentList := []*repository.AssessmentDetail{}

	rows, err := r.ps[assessmentFindAllPagination].QueryContext(ctx, offset, limit)
	if err != nil {
		log.Println("assessment sql repo query context error: %w", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		assessment := &repository.AssessmentDetail{}
		err := rows.Scan(
			&assessment.Id,
			&assessment.InstitutionName,
			&assessment.Status,
			&assessment.SubmittedDate,
		)
		if err != nil {
			log.Println("assessment sql repo scan error: %w", err)
			return nil, err
		}
		assessmentList = append(assessmentList, assessment)
	}

	return assessmentList, nil
}
