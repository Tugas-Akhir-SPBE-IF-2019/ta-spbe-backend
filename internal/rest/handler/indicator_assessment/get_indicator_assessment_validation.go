package indicatorassessment

import (
	"net/http"

	apierror "github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/rest/error"
	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/rest/response"

	"github.com/go-chi/chi/v5"
)

type IndicatorAssessmentResultValidationResponse struct {
	IndicatorNumber int    `json:"indicator_number"`
	ResultCorrect   bool   `json:"result_correct"`
	CorrectLevel    int    `json:"correct_level"`
	Explanation     string `json:"explanation"`
}

func (handler *indicatorAssessmentHandler) GetIndicatorAssessmentResultValidation(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	indicatorAssessmentId := chi.URLParam(r, "id")
	indicatorAssessmentList, err := handler.indicatorAssessmentStore.FindIndicatorAssessmentResultByAssessmentId(ctx, indicatorAssessmentId)
	if err != nil {
		response.Error(w, apierror.NotFoundError("indicator assessment not found"))
		return
	}

	result := make([]IndicatorAssessmentResultItem, len(indicatorAssessmentList))
	for idx, indicatorAssessment := range indicatorAssessmentList {
		resultItem := IndicatorAssessmentResultItem{
			Domain:          indicatorAssessment.Result.Domain,
			Aspect:          indicatorAssessment.Result.Aspect,
			IndicatorNumber: indicatorAssessment.Result.IndicatorNumber,
			Level:           indicatorAssessment.Result.Level,
			Explanation:     indicatorAssessment.Result.Explanation,
		}

		supportDocumentProofList := make([]SupportDocumentProofInfo, len(indicatorAssessment.Result.SupportDocumentList))
		for index, supportDocumentItem := range indicatorAssessment.Result.SupportDocumentList {
			supportDocumentProof := SupportDocumentProofInfo{
				Name:                 supportDocumentItem.Name,
				URL:                  supportDocumentItem.URL,
				Type:                 supportDocumentItem.Type,
				Proof:                supportDocumentItem.Proof,
				ImageURL:             supportDocumentItem.ImageURL,
				ProofPageDocumentURL: supportDocumentItem.SpecificPageDocumentURL,
			}
			supportDocumentProofList[index] = supportDocumentProof
		}

		resultItem.SupportDocumentProofList = supportDocumentProofList
		result[idx] = resultItem
	}
	resp := []IndicatorAssessmentResultValidationResponse{}
	for _, item := range result {
		resp = append(resp, IndicatorAssessmentResultValidationResponse{
			IndicatorNumber: item.IndicatorNumber,
			ResultCorrect:   false,
			CorrectLevel:    item.Level,
			Explanation:     "hardcoded",
		})
	}

	response.Respond(w, http.StatusOK, resp)
}
