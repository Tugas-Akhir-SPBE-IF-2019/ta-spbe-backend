package indicatorassessment

import (
	"net/http"
	"time"

	apierror "github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/rest/error"
	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/rest/response"

	"github.com/go-chi/chi/v5"
)

type SupportDocumentProofInfo struct {
	Name                 string `json:"name"`
	URL                  string `json:"document_url"`
	Type                 string `json:"type"`
	Proof                string `json:"proof"`
	ImageURL             string `json:"proof_image_url"`
	ProofPageDocumentURL string `json:"proof_page_document_url"`
}

type IndicatorAssessmentResultItem struct {
	Domain                   string                     `json:"domain"`
	Aspect                   string                     `json:"aspect"`
	IndicatorNumber          int                        `json:"indicator_number"`
	Level                    int                        `json:"level"`
	Explanation              string                     `json:"explanation"`
	SupportDocumentProofList []SupportDocumentProofInfo `json:"support_document_proof"`
}

type IndicatorAssessmentResultResponse struct {
	InstitutionName  string                          `json:"institution_name"`
	SubmittedDate    time.Time                       `json:"submitted_date"`
	AssessmentStatus int                             `json:"assessment_status"`
	Result           []IndicatorAssessmentResultItem `json:"result"`
}

func (handler *indicatorAssessmentHandler) GetIndicatorAssessmentResultGetIndicatorAssessmentIndexList(w http.ResponseWriter, r *http.Request) {
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
	resp := IndicatorAssessmentResultResponse{
		InstitutionName:  indicatorAssessmentList[0].InstitutionName,
		SubmittedDate:    indicatorAssessmentList[0].SubmittedDate,
		AssessmentStatus: indicatorAssessmentList[0].AssessmentStatus,
		Result:           result,
	}

	response.Respond(w, http.StatusOK, resp)
}
