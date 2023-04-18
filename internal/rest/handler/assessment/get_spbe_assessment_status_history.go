package assessment

import (
	"log"
	"net/http"
	"time"

	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/rest/response"
	"github.com/go-chi/chi/v5"

	apierror "github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/rest/error"
)

type SupportDataDocumentItem struct {
	Name string `json:"name"`
	Type string `json:"type"`
	URL  string `json:"url"`
}

type AssessmentStatusHistoryItem struct {
	Status       int       `json:"status"`
	FinishedDate time.Time `json:"finished_date"`
}

type GetAssessmentStatusHistoryItem struct {
}

type GetAssessmentStatusHistoryResponse struct {
	Items               []AssessmentStatusHistoryItem `json:"status_histories"`
	SupportingDocuments []SupportDataDocumentItem     `json:"supporting_documents"`
}

func (handler *assessmentHandler) GetSPBEAssessmentStatusHistory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	assessmentId := chi.URLParam(r, "id")

	assessmentStatusHistoryList, err := handler.assessmentStore.FindAllStatusHistoryById(ctx, assessmentId)
	if err != nil {
		log.Println(err)

		response.Error(w, apierror.InternalServerError())
		return
	}

	resp := GetAssessmentStatusHistoryResponse{}
	items := make([]AssessmentStatusHistoryItem, len(assessmentStatusHistoryList))
	// statusBefore := -1 // Used because insert status history still buggy because of concurrency problem
	for idx, item := range assessmentStatusHistoryList {
		// if statusBefore != int(item.Status) {
		// 	statusBefore = int(item.Status)
		// 	items = append(items, AssessmentStatusHistoryItem{
		// 		Status:       int(item.Status),
		// 		FinishedDate: item.FinishedDate,
		// 	},
		// 	)
		// }
		items[idx] = AssessmentStatusHistoryItem{
			Status:       int(item.Status),
			FinishedDate: item.FinishedDate,
		}
	}

	supportDataDocumentList, err := handler.assessmentStore.FindAllDocumentsById(ctx, assessmentId)
	if err != nil {
		log.Println(err)

		response.Error(w, apierror.InternalServerError())
		return
	}

	supportDataDocumentItems := make([]SupportDataDocumentItem, len(supportDataDocumentList))
	for idx, item := range supportDataDocumentList {
		supportDataDocumentItems[idx] = SupportDataDocumentItem{
			Name: item.Name,
			Type: string(item.Type),
			URL:  item.Url,
		}
	}

	resp.Items = items
	resp.SupportingDocuments = supportDataDocumentItems

	response.Respond(w, http.StatusOK, resp)
}
