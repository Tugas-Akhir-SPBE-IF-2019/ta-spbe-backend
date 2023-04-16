package store

import (
	"context"
	"time"
)

type IndicatorAssessmentDetail struct {
	InstitutionName string
	SpbeIndex       float64
	SubmittedDate   time.Time
}

type IndicatorAssessmentSupportDocumentInfo struct {
	Name                    string
	URL                     string
	Type                    string
	Proof                   string
	ImageURL                string
	SpecificPageDocumentURL string
}

type IndicatorAssessmentResultInfo struct {
	Domain              string
	Aspect              string
	IndicatorNumber     int
	Level               int
	Explanation         string
	SupportDocument     string
	SupportDocumentName string
	SupportDocumentList []IndicatorAssessmentSupportDocumentInfo
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

type IndicatorAssessmentProofData struct {
	ID                    string
	IndicatorAssessmentID string
	ImageURL              string
	DocumentURL           string
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

type IndicatorData struct {
	ID              string
	IndicatorNumber int
	Aspect          string
	Domain          string
	Detail          string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type IndicatorAssessment interface {
	FindAll(ctx context.Context, queryInstitution string, startDate string, endDate string, indexMin float64, indexMax float64) ([]*IndicatorAssessmentDetail, error)
	FindAllPagination(ctx context.Context, offset int, limit int, queryInstitution string, startDate string, endDate string, indexMin float64, indexMax float64) ([]*IndicatorAssessmentDetail, error)
	FindIndicatorAssessmentResultById(ctx context.Context, id string) (IndicatorAssessmentResultDetail, error)
	FindIndicatorAssessmentResultByAssessmentId(ctx context.Context, id string) ([]*IndicatorAssessmentResultDetail, error)
	ValidateAssessmentResult(ctx context.Context, resultCorrect bool, indicatorAssessmentResult *IndicatorAssessmentResultDetail) error
	UpdateAssessmentResult(ctx context.Context, resultDetail *IndicatorAssessmentResultDetail) error
	InsertProofData(ctx context.Context, proofData *IndicatorAssessmentProofData) error
	FindProofDataByIndicatorAssessmentId(ctx context.Context, id string) ([]*IndicatorAssessmentProofData, error)
	FindIndicatorDetailByIndicatorNumber(ctx context.Context, indicatorNumber int) (IndicatorData, error)
}
