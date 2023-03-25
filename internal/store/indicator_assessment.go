package store

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
	Domain              string
	Aspect              string
	IndicatorNumber     int
	Level               int
	Explanation         string
	SupportDocument     string
	SupportDocumentName string
	OldDocument         string
	Proof               string
}

type IndicatorAssessmentResultDetail struct {
	IndicatorAssessmentId string
	AssessmentId          string
	InstitutionName       string
	SubmittedDate         time.Time
	AssessmentStatus      int
	Result                IndicatorAssessmentResultInfo
	ResultFeedback        IndicatorAssessmentResultFeedback
}

type IndicatorAssessmentResultFeedback struct {
	Level    int
	Feedback string
}

type IndicatorAssessment interface {
	FindAll(ctx context.Context) ([]*IndicatorAssessmentDetail, error)
	FindAllPagination(ctx context.Context, offset int, limit int) ([]*IndicatorAssessmentDetail, error)
	FindIndicatorAssessmentResultById(ctx context.Context, id string) (IndicatorAssessmentResultDetail, error)
	FindIndicatorAssessmentResultByAssessmentId(ctx context.Context, id string) ([]*IndicatorAssessmentResultDetail, error)
	ValidateAssessmentResult(ctx context.Context, resultCorrect bool, indicatorAssessmentResult *IndicatorAssessmentResultDetail) error
	UpdateAssessmentResult(ctx context.Context, resultDetail *IndicatorAssessmentResultDetail) error
}
