package indicatorassessment

import (
	"net/http"
	"time"

	apierror "github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/rest/error"
	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/rest/response"

	"github.com/go-chi/chi/v5"
)

type IndicatorAssessmentResultItem struct {
	Domain          string `json:"domain"`
	Aspect          string `json:"aspect"`
	IndicatorNumber int    `json:"indicator_number"`
	Level           int    `json:"level"`
	Explanation     string `json:"explanation"`
	SupportDocument string `json:"support_document"`
	OldDocument     string `json:"old_document"`
	Proof           string `json:"proof"`
}

type IndicatorAssessmentResultResponse struct {
	InstitutionName  string                        `json:"institution_name"`
	SubmittedDate    time.Time                     `json:"submitted_date"`
	AssessmentStatus int                           `json:"assessment_status"`
	Result           IndicatorAssessmentResultItem `json:"result"`
}

func (handler *indicatorAssessmentHandler) GetIndicatorAssessmentResultGetIndicatorAssessmentIndexList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	indicatorAssessmentId := chi.URLParam(r, "id")
	indicatorAssessment, err := handler.indicatorAssessmentStore.FindIndicatorAssessmentResultById(ctx, indicatorAssessmentId)
	if err != nil {
		response.Error(w, apierror.NotFoundError("indicator assessment not found"))
		return
	}
	resp := IndicatorAssessmentResultResponse{
		InstitutionName:  indicatorAssessment.InstitutionName,
		SubmittedDate:    indicatorAssessment.SubmittedDate,
		AssessmentStatus: indicatorAssessment.AssessmentStatus,
		Result: IndicatorAssessmentResultItem{
			Domain:          indicatorAssessment.Result.Domain,
			Aspect:          indicatorAssessment.Result.Aspect,
			IndicatorNumber: indicatorAssessment.Result.IndicatorNumber,
			Level:           indicatorAssessment.Result.Level,
			Explanation:     indicatorAssessment.Result.Explanation,
			SupportDocument: indicatorAssessment.Result.SupportDocument,
			Proof:           indicatorAssessment.Result.Proof,
		},
	}

	response.Respond(w, http.StatusOK, resp)

}
