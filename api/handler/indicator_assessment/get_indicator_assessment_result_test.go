package indicatorassessment

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"ta-spbe-backend/repository"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestGetIndicatorAssessmentResult_Success(t *testing.T) {
	testIndicatorAssessmentId := "70968d00-0959-4400-b82f-5b968a11abfd"

	rr := httptest.NewRecorder()
	r := chi.NewRouter()
	req, err := http.NewRequest(http.MethodGet, "/assessments/"+testIndicatorAssessmentId, nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}

	findIndicatorAssessmentResultByIdIndicatorAssessmentRepoMock = func(ctx context.Context, id string) (repository.IndicatorAssessmentResultDetail, error) {
		return repository.IndicatorAssessmentResultDetail{
			IndicatorAssessmentId: testIndicatorAssessmentId,
		}, nil
	}

	mockIndicatorAssessmentRepo := indicatorAssessmentRepoMock{}
	r.Get("/assessments/{id}", GetIndicatorAssessmentResult(mockIndicatorAssessmentRepo))
	r.ServeHTTP(rr, req)

	var response IndicatorAssessmentResultResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)

	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.EqualValues(t, http.StatusOK, rr.Code)
}

func TestGetIndicatorAssessmentResult_FailResultNotFound(t *testing.T) {
	testIndicatorAssessmentId := "70968d00-0959-4400-b82f-5b968a11abfd"

	rr := httptest.NewRecorder()
	r := chi.NewRouter()
	req, err := http.NewRequest(http.MethodGet, "/assessments/"+testIndicatorAssessmentId, nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}

	findIndicatorAssessmentResultByIdIndicatorAssessmentRepoMock = func(ctx context.Context, id string) (repository.IndicatorAssessmentResultDetail, error) {
		return repository.IndicatorAssessmentResultDetail{
		}, fmt.Errorf("test error")
	}

	mockIndicatorAssessmentRepo := indicatorAssessmentRepoMock{}
	r.Get("/assessments/{id}", GetIndicatorAssessmentResult(mockIndicatorAssessmentRepo))
	r.ServeHTTP(rr, req)

	var response IndicatorAssessmentResultResponse
	err = json.Unmarshal(rr.Body.Bytes(), &response)

	assert.Nil(t, err)
	assert.NotNil(t, response)
	assert.EqualValues(t, http.StatusNotFound, rr.Code)
}
