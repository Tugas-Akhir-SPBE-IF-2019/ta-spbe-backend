package indicatorassessment

import (
	"encoding/json"
	"log"
	"net/http"

	apierror "github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/rest/error"
	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/rest/response"
	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/store"

	"github.com/go-chi/chi/v5"
)

type ValidateIndicatorAssessmentResultRequest struct {
	IndicatorNumber int    `json:"indicator_number"`
	ResultCorrect   bool   `json:"result_correct"`
	CorrectLevel    int    `json:"correct_level"`
	Explanation     string `json:"explanation"`
}

type ValidateIndicatorAssessmentResultResponse struct {
	Message string `json:"message"`
}

func (handler *indicatorAssessmentHandler) ValidateIndicatorAssessmentResult(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	reqList := []ValidateIndicatorAssessmentResultRequest{}
	assessmentId := chi.URLParam(r, "id")

	if err := json.NewDecoder(r.Body).Decode(&reqList); err != nil {
		response.Error(w, apierror.BadRequestError(err.Error()))
		return
	}

	indicatorAssessmentResultList, err := handler.indicatorAssessmentStore.FindIndicatorAssessmentResultByAssessmentId(ctx, assessmentId)
	if err != nil {
		response.Error(w, apierror.NotFoundError("indicator assessment not found"))
		return
	}

	for _, req := range reqList {
		for _, indicatorAssessmentResult := range indicatorAssessmentResultList {
			if req.IndicatorNumber == indicatorAssessmentResult.Result.IndicatorNumber {
				if indicatorAssessmentResult.AssessmentStatus == int(store.AssessmentStatus(store.IN_PROGRESS)) {
					response.Error(w, apierror.BadRequestError("indicator assessment result is still in progress"))
					return
				}

				if !req.ResultCorrect {
					indicatorAssessmentResult.ResultFeedback.Level = req.CorrectLevel
					indicatorAssessmentResult.ResultFeedback.Feedback = req.Explanation
				}

				err = handler.indicatorAssessmentStore.ValidateAssessmentResult(ctx, req.ResultCorrect, indicatorAssessmentResult)
				if err != nil {
					log.Println(err.Error())
					response.Error(w, apierror.InternalServerError())
					return
				}
				break
			}
		}
	}

	resp := ValidateIndicatorAssessmentResultResponse{
		Message: "validation success",
	}
	response.Respond(w, http.StatusOK, resp)
}
