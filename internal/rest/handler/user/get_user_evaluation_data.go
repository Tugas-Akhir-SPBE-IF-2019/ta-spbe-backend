package user

import (
	"net/http"
	"time"

	apierror "github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/rest/error"
	userCtx "github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/rest/handler/context"

	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/rest/response"
)

type UserEvaluationDataItem struct {
	Role            string    `json:"role"`
	InstitutionName string    `json:"institution_name"`
	EvaluationYear  int       `json:"evaluation_year"`
	CreatedAt       time.Time `json:"created_at"`
}

type GetUserEvaluationDataResponse struct {
	Items []UserEvaluationDataItem `json:"items"`
}

func (handler *userHandler) GetUserEvaluationData(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userCred, ok := ctx.Value(userCtx.UserCtxKey).(userCtx.UserCtx)
	if !ok {
		response.Error(w, apierror.InternalServerError())
		return
	}

	userEvaluationDataList, err := handler.userStore.FindEvaluationDataByUserID(ctx, userCred.ID)
	if err != nil {
		response.Error(w, apierror.NotFoundError("user evaluation data not found"))
		return
	}

	items := make([]UserEvaluationDataItem, len(userEvaluationDataList))
	for idx, userEvaluationData := range userEvaluationDataList {
		item := UserEvaluationDataItem{
			Role:            userEvaluationData.Role,
			InstitutionName: userEvaluationData.InstitutionName,
			EvaluationYear:  userEvaluationData.EvaluationYear,
			CreatedAt:       userEvaluationData.CreatedAt,
		}
		items[idx] = item
	}
	resp := GetUserEvaluationDataResponse{
		Items: items,
	}

	response.Respond(w, http.StatusOK, resp)
}
