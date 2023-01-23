package services

type AssessmentService interface {}

type assessmentService struct{}

func NewAssessmentService() AssessmentService {
	return assessmentService{}
}