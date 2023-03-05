package indicatorassessment

// import (
// 	"bytes"
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"net/http/httptest"
// 	"ta-spbe-backend/repository"
// 	"testing"

// 	"github.com/go-chi/chi/v5"
// 	"github.com/stretchr/testify/assert"
// )

// func TestValidateIndicatorAssessmentResult_Success(t *testing.T) {
// 	reqBody := ValidateIndicatorAssessmentResultRequest{
// 		ResultCorrect: false,
// 		CorrectLevel:  1,
// 		Explanation:   "test-explanation",
// 	}
// 	reqBodyBytes, _ := json.Marshal(reqBody)

// 	testIndicatorAssessmentId := "70968d00-0959-4400-b82f-5b968a11abfd"

// 	rr := httptest.NewRecorder()
// 	r := chi.NewRouter()
// 	req, err := http.NewRequest(http.MethodPost, "/assessments/"+testIndicatorAssessmentId+"/validate", bytes.NewBuffer(reqBodyBytes))
// 	if err != nil {
// 		t.Errorf("this is the error: %v\n", err)
// 	}

// 	findIndicatorAssessmentResultByIdIndicatorAssessmentRepoMock = func(ctx context.Context, id string) (repository.IndicatorAssessmentResultDetail, error) {
// 		return repository.IndicatorAssessmentResultDetail{
// 			IndicatorAssessmentId: testIndicatorAssessmentId,
// 			AssessmentStatus:      int(repository.AssessmentStatus(repository.COMPLETED)),
// 		}, nil
// 	}
// 	validateAssessmentResultIndicatorAssessmentRepoMock = func(ctx context.Context, resultCorrect bool, indicatorAssessmentResult *repository.IndicatorAssessmentResultDetail) error {
// 		return nil
// 	}

// 	mockIndicatorAssessmentRepo := indicatorAssessmentRepoMock{}
// 	r.Post("/assessments/{id}/validate", ValidateIndicatorAssessmentResult(mockIndicatorAssessmentRepo))
// 	r.ServeHTTP(rr, req)

// 	var response ValidateIndicatorAssessmentResultResponse
// 	err = json.Unmarshal(rr.Body.Bytes(), &response)

// 	assert.Nil(t, err)
// 	assert.NotNil(t, response)
// 	assert.EqualValues(t, http.StatusOK, rr.Code)
// }

// func TestValidateIndicatorAssessmentResult_SuccessResultCorrect(t *testing.T) {
// 	reqBody := ValidateIndicatorAssessmentResultRequest{
// 		ResultCorrect: true,
// 	}
// 	reqBodyBytes, _ := json.Marshal(reqBody)

// 	testIndicatorAssessmentId := "70968d00-0959-4400-b82f-5b968a11abfd"

// 	rr := httptest.NewRecorder()
// 	r := chi.NewRouter()
// 	req, err := http.NewRequest(http.MethodPost, "/assessments/"+testIndicatorAssessmentId+"/validate", bytes.NewBuffer(reqBodyBytes))
// 	if err != nil {
// 		t.Errorf("this is the error: %v\n", err)
// 	}

// 	findIndicatorAssessmentResultByIdIndicatorAssessmentRepoMock = func(ctx context.Context, id string) (repository.IndicatorAssessmentResultDetail, error) {
// 		return repository.IndicatorAssessmentResultDetail{
// 			IndicatorAssessmentId: testIndicatorAssessmentId,
// 			AssessmentStatus:      int(repository.AssessmentStatus(repository.COMPLETED)),
// 		}, nil
// 	}
// 	validateAssessmentResultIndicatorAssessmentRepoMock = func(ctx context.Context, resultCorrect bool, indicatorAssessmentResult *repository.IndicatorAssessmentResultDetail) error {
// 		return nil
// 	}

// 	mockIndicatorAssessmentRepo := indicatorAssessmentRepoMock{}
// 	r.Post("/assessments/{id}/validate", ValidateIndicatorAssessmentResult(mockIndicatorAssessmentRepo))
// 	r.ServeHTTP(rr, req)

// 	var response ValidateIndicatorAssessmentResultResponse
// 	err = json.Unmarshal(rr.Body.Bytes(), &response)

// 	assert.Nil(t, err)
// 	assert.NotNil(t, response)
// 	assert.EqualValues(t, http.StatusOK, rr.Code)
// }

// func TestValidateIndicatorAssessmentResult_FailInvalidRequestBody(t *testing.T) {
// 	reqBodyJsonInvalid := "{invalid:json}"
// 	testIndicatorAssessmentId := "70968d00-0959-4400-b82f-5b968a11abfd"

// 	rr := httptest.NewRecorder()
// 	r := chi.NewRouter()
// 	req, err := http.NewRequest(http.MethodPost, "/assessments/"+testIndicatorAssessmentId+"/validate", bytes.NewBufferString(reqBodyJsonInvalid))
// 	if err != nil {
// 		t.Errorf("this is the error: %v\n", err)
// 	}

// 	mockIndicatorAssessmentRepo := indicatorAssessmentRepoMock{}
// 	r.Post("/assessments/{id}/validate", ValidateIndicatorAssessmentResult(mockIndicatorAssessmentRepo))
// 	r.ServeHTTP(rr, req)

// 	var response ValidateIndicatorAssessmentResultResponse
// 	err = json.Unmarshal(rr.Body.Bytes(), &response)

// 	assert.Nil(t, err)
// 	assert.NotNil(t, response)
// 	assert.EqualValues(t, http.StatusBadRequest, rr.Code)
// }

// func TestValidateIndicatorAssessmentResult_FailIndicatorAssessmentNotFound(t *testing.T) {
// 	reqBody := ValidateIndicatorAssessmentResultRequest{
// 		ResultCorrect: false,
// 		CorrectLevel:  1,
// 		Explanation:   "test-explanation",
// 	}
// 	reqBodyBytes, _ := json.Marshal(reqBody)

// 	testIndicatorAssessmentId := "70968d00-0959-4400-b82f-5b968a11abfd"

// 	rr := httptest.NewRecorder()
// 	r := chi.NewRouter()
// 	req, err := http.NewRequest(http.MethodPost, "/assessments/"+testIndicatorAssessmentId+"/validate", bytes.NewBuffer(reqBodyBytes))
// 	if err != nil {
// 		t.Errorf("this is the error: %v\n", err)
// 	}

// 	findIndicatorAssessmentResultByIdIndicatorAssessmentRepoMock = func(ctx context.Context, id string) (repository.IndicatorAssessmentResultDetail, error) {
// 		return repository.IndicatorAssessmentResultDetail{
// 		}, fmt.Errorf("test error")
// 	}

// 	mockIndicatorAssessmentRepo := indicatorAssessmentRepoMock{}
// 	r.Post("/assessments/{id}/validate", ValidateIndicatorAssessmentResult(mockIndicatorAssessmentRepo))
// 	r.ServeHTTP(rr, req)

// 	var response ValidateIndicatorAssessmentResultResponse
// 	err = json.Unmarshal(rr.Body.Bytes(), &response)

// 	assert.Nil(t, err)
// 	assert.NotNil(t, response)
// 	assert.EqualValues(t, http.StatusNotFound, rr.Code)
// }

// func TestValidateIndicatorAssessmentResult_FailStatusInProgress(t *testing.T) {
// 	reqBody := ValidateIndicatorAssessmentResultRequest{
// 		ResultCorrect: false,
// 		CorrectLevel:  1,
// 		Explanation:   "test-explanation",
// 	}
// 	reqBodyBytes, _ := json.Marshal(reqBody)

// 	testIndicatorAssessmentId := "70968d00-0959-4400-b82f-5b968a11abfd"

// 	rr := httptest.NewRecorder()
// 	r := chi.NewRouter()
// 	req, err := http.NewRequest(http.MethodPost, "/assessments/"+testIndicatorAssessmentId+"/validate", bytes.NewBuffer(reqBodyBytes))
// 	if err != nil {
// 		t.Errorf("this is the error: %v\n", err)
// 	}

// 	findIndicatorAssessmentResultByIdIndicatorAssessmentRepoMock = func(ctx context.Context, id string) (repository.IndicatorAssessmentResultDetail, error) {
// 		return repository.IndicatorAssessmentResultDetail{
// 			IndicatorAssessmentId: testIndicatorAssessmentId,
// 			AssessmentStatus:      int(repository.AssessmentStatus(repository.IN_PROGRESS)),
// 		}, nil
// 	}

// 	mockIndicatorAssessmentRepo := indicatorAssessmentRepoMock{}
// 	r.Post("/assessments/{id}/validate", ValidateIndicatorAssessmentResult(mockIndicatorAssessmentRepo))
// 	r.ServeHTTP(rr, req)

// 	var response ValidateIndicatorAssessmentResultResponse
// 	err = json.Unmarshal(rr.Body.Bytes(), &response)

// 	assert.Nil(t, err)
// 	assert.NotNil(t, response)
// 	assert.EqualValues(t, http.StatusBadRequest, rr.Code)
// }

// func TestValidateIndicatorAssessmentResult_FailValidate(t *testing.T) {
// 	reqBody := ValidateIndicatorAssessmentResultRequest{
// 		ResultCorrect: false,
// 		CorrectLevel:  1,
// 		Explanation:   "test-explanation",
// 	}
// 	reqBodyBytes, _ := json.Marshal(reqBody)

// 	testIndicatorAssessmentId := "70968d00-0959-4400-b82f-5b968a11abfd"

// 	rr := httptest.NewRecorder()
// 	r := chi.NewRouter()
// 	req, err := http.NewRequest(http.MethodPost, "/assessments/"+testIndicatorAssessmentId+"/validate", bytes.NewBuffer(reqBodyBytes))
// 	if err != nil {
// 		t.Errorf("this is the error: %v\n", err)
// 	}

// 	findIndicatorAssessmentResultByIdIndicatorAssessmentRepoMock = func(ctx context.Context, id string) (repository.IndicatorAssessmentResultDetail, error) {
// 		return repository.IndicatorAssessmentResultDetail{
// 			IndicatorAssessmentId: testIndicatorAssessmentId,
// 			AssessmentStatus:      int(repository.AssessmentStatus(repository.COMPLETED)),
// 		}, nil
// 	}
// 	validateAssessmentResultIndicatorAssessmentRepoMock = func(ctx context.Context, resultCorrect bool, indicatorAssessmentResult *repository.IndicatorAssessmentResultDetail) error {
// 		return fmt.Errorf("test error")
// 	}

// 	mockIndicatorAssessmentRepo := indicatorAssessmentRepoMock{}
// 	r.Post("/assessments/{id}/validate", ValidateIndicatorAssessmentResult(mockIndicatorAssessmentRepo))
// 	r.ServeHTTP(rr, req)

// 	var response ValidateIndicatorAssessmentResultResponse
// 	err = json.Unmarshal(rr.Body.Bytes(), &response)

// 	assert.Nil(t, err)
// 	assert.NotNil(t, response)
// 	assert.EqualValues(t, http.StatusInternalServerError, rr.Code)
// }