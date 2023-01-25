package repository

import (
	"context"
	"time"
)

type AssessmentDetail struct {
	Id              string
	InstitutionName string
	Status          int
	SubmittedDate   time.Time
}

type AssessmentRepository interface {
	FindAll(ctx context.Context) ([]*AssessmentDetail, error)
	FindAllPagination(ctx context.Context, offset int, limit int) ([]*AssessmentDetail, error)
}
