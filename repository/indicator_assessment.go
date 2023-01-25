package repository

import (
	"context"
	"time"
)

type IndicatorAssessmentDetail struct {
	InstitutionName string
	SpbeIndex       int
	SubmittedDate   time.Time
}

type IndicatorAssessmentRepository interface {
	FindAll(ctx context.Context) ([]*IndicatorAssessmentDetail, error)
	FindAllPagination(ctx context.Context, offset int, limit int) ([]*IndicatorAssessmentDetail, error)
}
