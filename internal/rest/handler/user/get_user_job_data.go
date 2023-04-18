package user

import (
	"net/http"
	"time"

	apierror "github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/rest/error"
	userCtx "github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/rest/handler/context"

	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/rest/response"
)

type UserJobDataItem struct {
	ID         string    `json:"id"`
	Role       string    `json:"role"`
	Company    string    `json:"company"`
	JoinedYear int       `json:"joined_year"`
	CreatedAt  time.Time `json:"created_at"`
}

type GetUserJobDataResponse struct {
	Items []UserJobDataItem `json:"items"`
}

func (handler *userHandler) GetUserJobData(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userCred, ok := ctx.Value(userCtx.UserCtxKey).(userCtx.UserCtx)
	if !ok {
		response.Error(w, apierror.InternalServerError())
		return
	}

	userJobDataList, err := handler.userStore.FindJobDataByUserID(ctx, userCred.ID)
	if err != nil {
		response.Error(w, apierror.NotFoundError("user evaluation data not found"))
		return
	}

	items := make([]UserJobDataItem, len(userJobDataList))
	for idx, userJobData := range userJobDataList {
		item := UserJobDataItem{
			ID:         userJobData.ID,
			Role:       userJobData.Role,
			Company:    userJobData.Company,
			JoinedYear: userJobData.JoinedDate,
			CreatedAt:  userJobData.CreatedAt,
		}
		items[idx] = item
	}
	resp := GetUserJobDataResponse{
		Items: items,
	}

	response.Respond(w, http.StatusOK, resp)
}
