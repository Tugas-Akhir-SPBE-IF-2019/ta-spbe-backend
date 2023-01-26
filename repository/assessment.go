package repository

import (
	"context"
	"time"
)

type AssessmentStatus int

const (
	IN_PROGRESS = AssessmentStatus(1)
	COMPLETED   = AssessmentStatus(2)
	VALIDATED   = AssessmentStatus(3)
)

type AssessmentDetail struct {
	Id              string
	InstitutionName string
	Status          int
	SubmittedDate   time.Time
}

type IndicatorAssessmentInfo struct {
	Id              string
	IndicatorNumber int
	Status          int
}

type SupportDataDocumentInfo struct {
	Id           string
	DocumentName string
}

type AssessmentUploadDetail struct {
	AssessmentDetail        AssessmentDetail
	IndicatorAssessmentInfo IndicatorAssessmentInfo
	SupportDataDocumentInfo SupportDataDocumentInfo
	UserId                  string
}

type AssessmentRepository interface {
	FindAll(ctx context.Context) ([]*AssessmentDetail, error)
	FindAllPagination(ctx context.Context, offset int, limit int) ([]*AssessmentDetail, error)
	InsertUploadDocument(ctx context.Context, assessmentUploadDetail *AssessmentUploadDetail) error
}
