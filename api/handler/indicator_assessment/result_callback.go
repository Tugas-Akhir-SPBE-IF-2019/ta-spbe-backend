package indicatorassessment

import (
	"encoding/json"
	"log"
	"net/http"
	apierror "ta-spbe-backend/api/error"
	"ta-spbe-backend/api/response"
	"ta-spbe-backend/repository"
)

type IndicatorAssessmentResultCallbackRequest struct {
	UserId                string `json:"user_id"`
	AssessmentId          string `json:"assessment_id"`
	IndicatorAssessmentId string `json:"indicator_assessment_id"`
	Level                 int    `json:"level"`
	Explanation           string `json:"explanation"`
	SupportDataDocumentId string `json:"support_data_document_id"`
	Proof                 string `json:"proof"`
}

type ValidateIndicatorAssessmentResultResponseS struct {
	Message string `json:"message"`
}

func ResultCallback(indicatorAssessmentRepo repository.IndicatorAssessmentRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		req := IndicatorAssessmentResultCallbackRequest{}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Error(w, apierror.BadRequestError(err.Error()))
			return
		}

		result := repository.IndicatorAssessmentResultDetail{
			AssessmentId:          req.AssessmentId,
			IndicatorAssessmentId: req.IndicatorAssessmentId,
			Result: repository.IndicatorAssessmentResultInfo{
				Level:           req.Level,
				Explanation:     req.Explanation,
				SupportDocument: req.SupportDataDocumentId,
				Proof:           req.Proof,
			},
		}

		err := indicatorAssessmentRepo.UpdateAssessmentResult(ctx, &result)
		if err != nil {
			log.Println(err.Error())
			response.Error(w, apierror.InternalServerError())
			return
		}

		response.Respond(w, http.StatusNoContent, nil)
	}
}
