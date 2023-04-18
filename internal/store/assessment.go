package store

import (
	"context"
	"time"
)

type AssessmentStatus int
type SupportDocumentType string

const (
	IN_PROGRESS     = AssessmentStatus(1)
	COMPLETED       = AssessmentStatus(2)
	VALIDATED       = AssessmentStatus(3)
	NEW_DOCUMENT    = SupportDocumentType("NEW_DOCUMENT")
	OLD_DOCUMENT    = SupportDocumentType("OLD_DOCUMENT")
	MEETING_MINUTES = SupportDocumentType("MEETING_MINUTES")
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
	Type                 SupportDocumentType
}

type AssessmentUploadDetail struct {
	AssessmentDetail            AssessmentDetail
	IndicatorAssessmentInfo     IndicatorAssessmentInfo
	IndicatorAssessmentInfoList []IndicatorAssessmentInfo
	SupportDataDocumentInfoList []SupportDataDocumentInfo
	UserId                      string
}

type AssessmentStatusHistoryDetail struct {
	Status                      AssessmentStatus
	SupportDataDocumentInfoList []SupportDataDocumentInfo
	FinishedDate                time.Time
}

type AssessmentDocumentDetail struct {
	Id   string
	Name string
	Url  string
	Type SupportDocumentType
}

type IndicatorAssessmentUpdateResultDetail struct {
	ID            string
	Number        int
	Detail        string
	DocumentProof []DocumentProofAssessmentUpdateResultDetail
	Result        IndicatorResultAssessmentUpdateResultDetail
}

type DocumentProofAssessmentUpdateResultDetail struct {
	Name                    string
	OriginalName            string
	Type                    string
	Text                    string
	Title                   string
	PictureFileList         []string
	SpecificPageDocumentURL []string
	DocumentPageList        []int
}

type IndicatorResultAssessmentUpdateResultDetail struct {
	Level       int
	Explanation string
}

type AssessmenUpdateResultDetail struct {
	IndicatorAssessmentList []IndicatorAssessmentUpdateResultDetail
}

type Assessment interface {
	FindAll(ctx context.Context, queryInstitution string, status int, startDate string, endDate string) ([]*AssessmentDetail, error)
	FindAllPagination(ctx context.Context, offset int, limit int, queryInstitution string, status int, startDate string, endDate string) ([]*AssessmentDetail, error)
	InsertUploadDocument(ctx context.Context, assessmentUploadDetail *AssessmentUploadDetail) error
	UpdateAssessmentResult(ctx context.Context, resultDetail *AssessmenUpdateResultDetail) error
	UpdateStatus(ctx context.Context, assessmentId string, status AssessmentStatus) error
	FindAllStatusHistoryById(ctx context.Context, assessmentId string) ([]*AssessmentStatusHistoryDetail, error)
	FindAllDocumentsById(ctx context.Context, assessmentId string) ([]*AssessmentDocumentDetail, error)
}
