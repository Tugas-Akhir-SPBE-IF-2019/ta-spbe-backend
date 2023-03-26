package assessment

import (
	"log"
	"net/http"

	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/rest/response"
	"github.com/go-chi/chi/v5"

	apierror "github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/rest/error"
)

type AssessmentDocumentItem struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type GetAssessmentDocumentResponse struct {
	Items []AssessmentDocumentItem `json:"items"`
}

func (handler *assessmentHandler) GetSPBEAssessmentDocumentList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	indicatorAssessmentId := chi.URLParam(r, "id")

	assessmentDocumentList, err := handler.assessmentStore.FindAllDocumentsById(ctx, indicatorAssessmentId)
	if err != nil {
		log.Println(err)

		response.Error(w, apierror.InternalServerError())
		return
	}

	resp := GetAssessmentDocumentResponse{}
	var items []AssessmentDocumentItem
	documentDuplicates := make(map[string]bool) // Used because insert document is still buggy because repeating insert with indicator assessment
	for _, item := range assessmentDocumentList {
		if !documentDuplicates[item.Url] {
			documentDuplicates[item.Url] = true
			items = append(items, AssessmentDocumentItem{
				Name: item.Name,
				Url:  item.Url,
			},
			)
		}

	}

	resp.Items = items

	response.Respond(w, http.StatusOK, resp)
}
