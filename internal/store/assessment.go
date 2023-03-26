package store

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
	Id                   string
	DocumentName         string
	DocumentUrl          string
	OriginalDocumentName string
}

type AssessmentUploadDetail struct {
	AssessmentDetail            AssessmentDetail
	IndicatorAssessmentInfo     IndicatorAssessmentInfo
	SupportDataDocumentInfoList []SupportDataDocumentInfo
	UserId                      string
}

type AssessmentStatusHistoryDetail struct {
	Status       AssessmentStatus
	FinishedDate time.Time
}

type Assessment interface {
	FindAll(ctx context.Context, queryInstitution string, status int, startDate string, endDate string) ([]*AssessmentDetail, error)
	FindAllPagination(ctx context.Context, offset int, limit int, queryInstitution string, status int, startDate string, endDate string) ([]*AssessmentDetail, error)
	InsertUploadDocument(ctx context.Context, assessmentUploadDetail *AssessmentUploadDetail) error
	UpdateStatus(ctx context.Context, assessmentId string, status AssessmentStatus) error
	FindAllStatusHistoryById(ctx context.Context, assessmentId string) ([]*AssessmentStatusHistoryDetail, error)
}
