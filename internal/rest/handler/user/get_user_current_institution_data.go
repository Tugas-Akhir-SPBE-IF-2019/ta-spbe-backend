package user

import (
	"net/http"
	"time"

	apierror "github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/rest/error"
	userCtx "github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/rest/handler/context"

	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/rest/response"
)

type UserCurrentInstitutionDataItem struct {
	ID          string    `json:"id"`
	Institution string    `json:"institution_name"`
	Role        string    `json:"role"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

type GetUserCurrentInstitutionDataResponse struct {
	Items []UserCurrentInstitutionDataItem `json:"items"`
}

func (handler *userHandler) GetUserCurrentInstitutionData(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userCred, ok := ctx.Value(userCtx.UserCtxKey).(userCtx.UserCtx)
	if !ok {
		response.Error(w, apierror.InternalServerError())
		return
	}

	userCurrentInstitutionDataList, err := handler.userStore.FindCurrentInstitutionDataByUserID(ctx, userCred.ID)
	if err != nil {
		response.Error(w, apierror.NotFoundError("user evaluation data not found"))
		return
	}

	items := make([]UserCurrentInstitutionDataItem, len(userCurrentInstitutionDataList))
	for idx, userCurrentInstitution := range userCurrentInstitutionDataList {
		item := UserCurrentInstitutionDataItem{
			ID:          userCurrentInstitution.ID,
			Institution: userCurrentInstitution.InstitutionName,
			Role:        userCurrentInstitution.Role,
			Status:      userCurrentInstitution.Status,
			CreatedAt:   userCurrentInstitution.CreatedAt,
		}
		items[idx] = item
	}
	resp := GetUserCurrentInstitutionDataResponse{
		Items: items,
	}

	response.Respond(w, http.StatusOK, resp)
}
