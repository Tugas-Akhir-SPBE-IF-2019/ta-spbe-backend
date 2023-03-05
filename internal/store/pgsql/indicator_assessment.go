package pgsql

import (
	"context"
	"database/sql"
	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/store"
	"log"
)

type IndicatorAssessment struct {
	db *sql.DB
}

func NewIndicatorAssessment(db *sql.DB) *IndicatorAssessment {
	return &IndicatorAssessment{db: db}
}

const indicatorAssessmentFindAllQuery = `SELECT a.institution_name, ia.level, ia.created_at
		FROM indicator_assessments ia
		LEFT JOIN assessments a
		ON ia.assessment_id = a.id
		WHERE ia.level IS NOT NULL`

func (s *IndicatorAssessment) FindAll(ctx context.Context) ([]*store.IndicatorAssessmentDetail, error) {
	indicatorAssessmentList := []*store.IndicatorAssessmentDetail{}

	rows, err := s.db.QueryContext(ctx, indicatorAssessmentFindAllQuery)
	if err != nil {
		log.Println("indicator assessment sql repo query context error: %w", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		indicatorAssessment := &store.IndicatorAssessmentDetail{}
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

const indicatorAssessmentFindAllPaginationQuery = `SELECT a.institution_name, ia.level, ia.created_at
		FROM indicator_assessments ia
		LEFT JOIN assessments a
		ON ia.assessment_id = a.id
		WHERE ia.level IS NOT NULL
		LIMIT $2 OFFSET $1`

func (s *IndicatorAssessment) FindAllPagination(ctx context.Context, offset int, limit int) ([]*store.IndicatorAssessmentDetail, error) {
	indicatorAssessmentList := []*store.IndicatorAssessmentDetail{}

	rows, err := s.db.QueryContext(ctx, indicatorAssessmentFindAllPaginationQuery, offset, limit)
	if err != nil {
		log.Println("indicator assessment sql repo query context error: %w", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		indicatorAssessment := &store.IndicatorAssessmentDetail{}
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
