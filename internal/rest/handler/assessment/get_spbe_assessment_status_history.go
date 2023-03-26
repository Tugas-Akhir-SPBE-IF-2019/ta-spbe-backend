package assessment

import (
	"log"
	"net/http"
	"time"

	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/rest/response"
	"github.com/go-chi/chi/v5"

	apierror "github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/rest/error"
)

type AssessmentStatusHistoryItem struct {
	Status       int       `json:"status"`
	FinishedDate time.Time `json:"finished_date"`
}

type GetAssessmentStatusHistoryResponse struct {
	Items []AssessmentStatusHistoryItem `json:"items"`
}

func (handler *assessmentHandler) GetSPBEAssessmentStatusHistory(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	indicatorAssessmentId := chi.URLParam(r, "id")

	assessmentStatusHistoryList, err := handler.assessmentStore.FindAllStatusHistoryById(ctx, indicatorAssessmentId)
	if err != nil {
		log.Println(err)

		response.Error(w, apierror.InternalServerError())
		return
	}

	resp := GetAssessmentStatusHistoryResponse{}
	var items []AssessmentStatusHistoryItem
	statusBefore := -1 // Used because insert status history still buggy because of concurrency problem
	for _, item := range assessmentStatusHistoryList {
		if statusBefore != int(item.Status) {
			statusBefore = int(item.Status)
			items = append(items, AssessmentStatusHistoryItem{
				Status:       int(item.Status),
				FinishedDate: item.FinishedDate,
			},
			)
		}

	}

	resp.Items = items

	response.Respond(w, http.StatusOK, resp)
}
