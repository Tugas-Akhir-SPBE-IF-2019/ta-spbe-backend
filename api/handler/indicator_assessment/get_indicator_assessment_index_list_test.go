package indicatorassessment

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"ta-spbe-backend/repository"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

type indicatorAssessmentRepoMock struct{}

var (
	findAllIndicatorAssessmentRepoMock                           func(ctx context.Context) ([]*repository.IndicatorAssessmentDetail, error)
	findAllPaginationIndicatorAssessmentRepoMock                 func(ctx context.Context, offset int, limit int) ([]*repository.IndicatorAssessmentDetail, error)
	findIndicatorAssessmentResultByIdIndicatorAssessmentRepoMock func(ctx context.Context, id string) (repository.IndicatorAssessmentResultDetail, error)
	validateAssessmentResultIndicatorAssessmentRepoMock          func(ctx context.Context, resultCorrect bool, indicatorAssessmentResult *repository.IndicatorAssessmentResultDetail) error
	updateAssessmentResultIndicatorAssessmentRepoMock            func(ctx context.Context, resultDetail *repository.IndicatorAssessmentResultDetail) error
)

func (repo indicatorAssessmentRepoMock) FindAll(ctx context.Context) ([]*repository.IndicatorAssessmentDetail, error) {
	return findAllIndicatorAssessmentRepoMock(ctx)
}

func (repo indicatorAssessmentRepoMock) FindAllPagination(ctx context.Context, offset int, limit int) ([]*repository.IndicatorAssessmentDetail, error) {
	return findAllPaginationIndicatorAssessmentRepoMock(ctx, offset, limit)
}

func (repo indicatorAssessmentRepoMock) FindIndicatorAssessmentResultById(ctx context.Context, id string) (repository.IndicatorAssessmentResultDetail, error) {
	return findIndicatorAssessmentResultByIdIndicatorAssessmentRepoMock(ctx, id)
}

func (repo indicatorAssessmentRepoMock) ValidateAssessmentResult(ctx context.Context, resultCorrect bool, indicatorAssessmentResult *repository.IndicatorAssessmentResultDetail) error {
	return validateAssessmentResultIndicatorAssessmentRepoMock(ctx, resultCorrect, indicatorAssessmentResult)
}

func (repo indicatorAssessmentRepoMock) UpdateAssessmentResult(ctx context.Context, resultDetail *repository.IndicatorAssessmentResultDetail) error {
	return updateAssessmentResultIndicatorAssessmentRepoMock(ctx, resultDetail)
}

func TestGetIndicatorAssessmentIndexList_Success(t *testing.T) {
	timeNow := time.Now()

	rr := httptest.NewRecorder()
	r := chi.NewRouter()
	req, err := http.NewRequest(http.MethodGet, "/assessments/index", nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}

	findAllIndicatorAssessmentRepoMock = func(ctx context.Context) ([]*repository.IndicatorAssessmentDetail, error) {
		return []*repository.IndicatorAssessmentDetail{
			{
				InstitutionName: "test-institution-name",
				SpbeIndex:       1,
				SubmittedDate:   timeNow,
			},
		}, nil
	}

	findAllPaginationIndicatorAssessmentRepoMock = func(ctx context.Context, offset int, limit int) ([]*repository.IndicatorAssessmentDetail, error) {
		return []*repository.IndicatorAssessmentDetail{
			{
				InstitutionName: "test-institution-name",
				SpbeIndex:       1,
				SubmittedDate:   timeNow,
			},
		}, nil
	}

	mockIndicatorAssessmentRepo := indicatorAssessmentRepoMock{}
	r.Get("/assessments/index", GetIndicatorAssessmentIndexList(mockIndicatorAssessmentRepo))
	r.ServeHTTP(rr, req)

	var response IndicatorAssessmentListResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)

	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.EqualValues(t, http.StatusOK, rr.Code)
	assert.EqualValues(t, 1, response.TotalItems)
	assert.EqualValues(t, 1, response.TotalPages)
}

func TestGetIndicatorAssessmentIndexList_FaildInvalidQueryParamPage(t *testing.T) {
	rr := httptest.NewRecorder()
	r := chi.NewRouter()
	req, err := http.NewRequest(http.MethodGet, "/assessments/index?page=invalidpage", nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}

	mockIndicatorAssessmentRepo := indicatorAssessmentRepoMock{}
	r.Get("/assessments/index", GetIndicatorAssessmentIndexList(mockIndicatorAssessmentRepo))
	r.ServeHTTP(rr, req)

	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestGetIndicatorAssessmentIndexList_FailInvalidQueryParamLimit(t *testing.T) {
	rr := httptest.NewRecorder()
	r := chi.NewRouter()
	req, err := http.NewRequest(http.MethodGet, "/assessments/index?limit=invalidlimit", nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}

	mockIndicatorAssessmentRepo := indicatorAssessmentRepoMock{}
	r.Get("/assessments/index", GetIndicatorAssessmentIndexList(mockIndicatorAssessmentRepo))
	r.ServeHTTP(rr, req)

	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusUnprocessableEntity, rr.Code)
}

func TestGetIndicatorAssessmentIndexList_FailAssessmentRepoFindAll(t *testing.T) {
	rr := httptest.NewRecorder()
	r := chi.NewRouter()
	req, err := http.NewRequest(http.MethodGet, "/assessments/index", nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}

	findAllIndicatorAssessmentRepoMock = func(ctx context.Context) ([]*repository.IndicatorAssessmentDetail, error) {
		return nil, fmt.Errorf("test-error")
	}
	mockIndicatorAssessmentRepo := indicatorAssessmentRepoMock{}
	r.Get("/assessments/index", GetIndicatorAssessmentIndexList(mockIndicatorAssessmentRepo))
	r.ServeHTTP(rr, req)

	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, rr.Code)
}

func TestGetIndicatorAssessmentIndexList_FailAssessmentRepoFindAllPagination(t *testing.T) {
	rr := httptest.NewRecorder()
	r := chi.NewRouter()
	req, err := http.NewRequest(http.MethodGet, "/assessments/index", nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}

	findAllIndicatorAssessmentRepoMock = func(ctx context.Context) ([]*repository.IndicatorAssessmentDetail, error) {
		return nil, nil
	}
	findAllPaginationIndicatorAssessmentRepoMock = func(ctx context.Context, offset int, limit int) ([]*repository.IndicatorAssessmentDetail, error) {
		return nil, fmt.Errorf("test-error")
	}
	mockIndicatorAssessmentRepo := indicatorAssessmentRepoMock{}
	r.Get("/assessments/index", GetIndicatorAssessmentIndexList(mockIndicatorAssessmentRepo))
	r.ServeHTTP(rr, req)

	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, rr.Code)
}
