package pgsql

import (
	"context"
	"database/sql"
	"log"
	"ta-spbe-backend/repository"
)

type indicatorAssessmentRepo struct {
	db *sql.DB
	ps map[string]*sql.Stmt
}

func NewIndicatorAssessmentRepo(db *sql.DB) (repository.IndicatorAssessmentRepository, error) {
	ps := make(map[string]*sql.Stmt, len(assessmentQueries))
	for key, query := range indicatorAssessmentQueries {
		stmt, err := prepareStmt(db, "indicatorAssessmentRepo", key, query)
		if err != nil {
			return nil, err
		}
		ps[key] = stmt
	}

	return &indicatorAssessmentRepo{db, ps}, nil
}

var indicatorAssessmentQueries = map[string]string{
	indicatorAssessmentFindAll:           indicatorAssessmentFindAllQuery,
	indicatorAssessmentFindAllPagination: indicatorAssessmentFindAllPaginationQuery,
}

const indicatorAssessmentFindAll = "findAll"
const indicatorAssessmentFindAllQuery = `SELECT a.institution_name, ia.level, ia.created_at
		FROM indicator_assessments ia
		LEFT JOIN assessments a
		ON ia.assessment_id = a.id
		WHERE ia.level IS NOT NULL`

func (r *indicatorAssessmentRepo) FindAll(ctx context.Context) ([]*repository.IndicatorAssessmentDetail, error) {
	indicatorAssessmentList := []*repository.IndicatorAssessmentDetail{}

	rows, err := r.ps[indicatorAssessmentFindAll].QueryContext(ctx)
	if err != nil {
		log.Println("indicator assessment sql repo query context error: %w", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		indicatorAssessment := &repository.IndicatorAssessmentDetail{}
		err := rows.Scan(
			&indicatorAssessment.InstitutionName,
			&indicatorAssessment.SpbeIndex,
			&indicatorAssessment.SubmittedDate,
		)
		if err != nil {
			log.Println("indicator assessment sql repo scan error: %w", err)
			return nil, err
		}
		indicatorAssessmentList = append(indicatorAssessmentList, indicatorAssessment)
	}

	return indicatorAssessmentList, nil
}

const indicatorAssessmentFindAllPagination = "findAllPagination"
const indicatorAssessmentFindAllPaginationQuery = `SELECT a.institution_name, ia.level, ia.created_at
		FROM indicator_assessments ia
		LEFT JOIN assessments a
		ON ia.assessment_id = a.id
		WHERE ia.level IS NOT NULL
		LIMIT $2 OFFSET $1`

func (r *indicatorAssessmentRepo) FindAllPagination(ctx context.Context, offset int, limit int) ([]*repository.IndicatorAssessmentDetail, error) {
	indicatorAssessmentList := []*repository.IndicatorAssessmentDetail{}

	rows, err := r.ps[indicatorAssessmentFindAllPagination].QueryContext(ctx, offset, limit)
	if err != nil {
		log.Println("indicator assessment sql repo query context error: %w", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		indicatorAssessment := &repository.IndicatorAssessmentDetail{}
		err := rows.Scan(
			&indicatorAssessment.InstitutionName,
			&indicatorAssessment.SpbeIndex,
			&indicatorAssessment.SubmittedDate,
		)
		if err != nil {
			log.Println("indicator assessment sql repo scan error: %w", err)
			return nil, err
		}
		indicatorAssessmentList = append(indicatorAssessmentList, indicatorAssessment)
	}

	return indicatorAssessmentList, nil
}
