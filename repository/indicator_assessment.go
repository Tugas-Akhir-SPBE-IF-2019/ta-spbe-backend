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

type IndicatorAssessmentResultInfo struct {
	Domain          string
	Aspect          string
	IndicatorNumber int
	Level           int
	Explanation     string
	SupportDocument string
	OldDocument     string
	Proof           string
}

type IndicatorAssessmentResultDetail struct {
	InstitutionName  string
	SubmittedDate    time.Time
	AssessmentStatus int
	Result           IndicatorAssessmentResultInfo
	Validated        bool
}

type IndicatorAssessmentRepository interface {
	FindAll(ctx context.Context) ([]*IndicatorAssessmentDetail, error)
	FindAllPagination(ctx context.Context, offset int, limit int) ([]*IndicatorAssessmentDetail, error)
	FindIndicatorAssessmentResultById(ctx context.Context, id string) (IndicatorAssessmentResultDetail, error)
}
