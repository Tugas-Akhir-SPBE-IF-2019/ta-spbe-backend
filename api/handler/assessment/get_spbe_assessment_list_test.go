package assessment

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"net/http/httptest"
// 	"ta-spbe-backend/repository"
// 	"testing"
// 	"time"

// 	"github.com/go-chi/chi/v5"
// 	"github.com/stretchr/testify/assert"
// )

// type assessmentRepoMock struct{}

// var (
// 	findAllAssessmentRepoMock              func(ctx context.Context) ([]*repository.AssessmentDetail, error)
// 	findAllPaginationAssessmentRepoMock    func(ctx context.Context, offset int, limit int) ([]*repository.AssessmentDetail, error)
// 	insertUploadDocumentAssessmentRepoMock func(ctx context.Context, assessmentUploadDetail *repository.AssessmentUploadDetail) error
// )

// func (repo assessmentRepoMock) FindAll(ctx context.Context) ([]*repository.AssessmentDetail, error) {
// 	return findAllAssessmentRepoMock(ctx)
// }

// func (repo assessmentRepoMock) FindAllPagination(ctx context.Context, offset int, limit int) ([]*repository.AssessmentDetail, error) {
// 	return findAllPaginationAssessmentRepoMock(ctx, offset, limit)
// }

// func (repo assessmentRepoMock) InsertUploadDocument(ctx context.Context, assessmentUploadDetail *repository.AssessmentUploadDetail) error {
// 	return insertUploadDocumentAssessmentRepoMock(ctx, assessmentUploadDetail)
// }

// func TestGetSPBEAssessmentList_Success(t *testing.T) {
// 	timeNow := time.Now()

// 	rr := httptest.NewRecorder()
// 	r := chi.NewRouter()
// 	req, err := http.NewRequest(http.MethodGet, "/assessments", nil)
// 	if err != nil {
// 		t.Errorf("this is the error: %v\n", err)
// 	}

// 	findAllAssessmentRepoMock = func(ctx context.Context) ([]*repository.AssessmentDetail, error) {
// 		return []*repository.AssessmentDetail{
// 			{
// 				Id:              "test-id",
// 				InstitutionName: "test-institution-name",
// 				Status:          1,
// 				SubmittedDate:   timeNow,
// 			},
// 		}, nil
// 	}

// 	findAllPaginationAssessmentRepoMock = func(ctx context.Context, offset int, limit int) ([]*repository.AssessmentDetail, error) {
// 		return []*repository.AssessmentDetail{
// 			{
// 				Id:              "test-id",
// 				InstitutionName: "test-institution-name",
// 				Status:          1,
// 				SubmittedDate:   timeNow,
// 			},
// 		}, nil
// 	}

// 	mockAssessmentRepo := assessmentRepoMock{}
// 	r.Get("/assessments", GetSPBEAssessmentList(mockAssessmentRepo))
// 	r.ServeHTTP(rr, req)

// 	var response AssessmentListResponse
// 	err = json.Unmarshal(rr.Body.Bytes(), &response)

// 	assert.Nil(t, err)
// 	assert.NotNil(t, response)
// 	assert.EqualValues(t, rr.Code, http.StatusOK)
// 	assert.EqualValues(t, response.TotalItems, 1)
// 	assert.EqualValues(t, response.TotalPages, 1)
// }

// func TestGetSPBEAssessmentList_FailInvalidQueryParamPage(t *testing.T) {
// 	rr := httptest.NewRecorder()
// 	r := chi.NewRouter()
// 	req, err := http.NewRequest(http.MethodGet, "/assessments?page=invalidpage", nil)
// 	if err != nil {
// 		t.Errorf("this is the error: %v\n", err)
// 	}

// 	mockAssessmentRepo := assessmentRepoMock{}
// 	r.Get("/assessments", GetSPBEAssessmentList(mockAssessmentRepo))
// 	r.ServeHTTP(rr, req)

// 	assert.Nil(t, err)
// 	assert.EqualValues(t, rr.Code, http.StatusUnprocessableEntity)
// }

// func TestGetSPBEAssessmentList_FailInvalidQueryParamLimit(t *testing.T) {
// 	rr := httptest.NewRecorder()
// 	r := chi.NewRouter()
// 	req, err := http.NewRequest(http.MethodGet, "/assessments?limit=invalidlimit", nil)
// 	if err != nil {
// 		t.Errorf("this is the error: %v\n", err)
// 	}

// 	mockAssessmentRepo := assessmentRepoMock{}
// 	r.Get("/assessments", GetSPBEAssessmentList(mockAssessmentRepo))
// 	r.ServeHTTP(rr, req)

// 	assert.Nil(t, err)
// 	assert.EqualValues(t, rr.Code, http.StatusUnprocessableEntity)
// }

// func TestGetSPBEAssessmentList_FailAssessmentRepoFindAll(t *testing.T) {
// 	rr := httptest.NewRecorder()
// 	r := chi.NewRouter()
// 	req, err := http.NewRequest(http.MethodGet, "/assessments", nil)
// 	if err != nil {
// 		t.Errorf("this is the error: %v\n", err)
// 	}

// 	findAllAssessmentRepoMock = func(ctx context.Context) ([]*repository.AssessmentDetail, error) {
// 		return nil, fmt.Errorf("test-error")
// 	}
// 	mockAssessmentRepo := assessmentRepoMock{}
// 	r.Get("/assessments", GetSPBEAssessmentList(mockAssessmentRepo))
// 	r.ServeHTTP(rr, req)

// 	assert.Nil(t, err)
// 	assert.EqualValues(t, rr.Code, http.StatusInternalServerError)
// }

// func TestGetSPBEAssessmentList_FailAssessmentRepoFindAllPagination(t *testing.T) {
// 	rr := httptest.NewRecorder()
// 	r := chi.NewRouter()
// 	req, err := http.NewRequest(http.MethodGet, "/assessments", nil)
// 	if err != nil {
// 		t.Errorf("this is the error: %v\n", err)
// 	}

// 	findAllAssessmentRepoMock = func(ctx context.Context) ([]*repository.AssessmentDetail, error) {
// 		return nil, nil
// 	}
// 	findAllPaginationAssessmentRepoMock = func(ctx context.Context, offset int, limit int) ([]*repository.AssessmentDetail, error) {
// 		return nil, fmt.Errorf("test-error")
// 	}
// 	mockAssessmentRepo := assessmentRepoMock{}
// 	r.Get("/assessments", GetSPBEAssessmentList(mockAssessmentRepo))
// 	r.ServeHTTP(rr, req)

// 	assert.Nil(t, err)
// 	assert.EqualValues(t, rr.Code, http.StatusInternalServerError)
// }