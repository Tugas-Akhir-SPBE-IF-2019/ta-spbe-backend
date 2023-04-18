package user

import (
	"encoding/json"
	"log"
	"net/http"

	apierror "github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/rest/error"
	userCtx "github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/rest/handler/context"
	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/store"

	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/rest/response"
)

type AddUserJobDataRequest struct {
	ID         string `json:"id"`
	Role       string `json:"role"`
	Company    string `json:"company"`
	JoinedYear int    `json:"joined_year"`
}

type AddUserJobDataResponse struct {
	Message string `json:"message"`
}

func (handler *userHandler) AddUserJobData(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userCred, ok := ctx.Value(userCtx.UserCtxKey).(userCtx.UserCtx)
	if !ok {
		response.Error(w, apierror.InternalServerError())
		return
	}

	req := []AddUserJobDataRequest{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, apierror.BadRequestError(err.Error()))
		return
	}

	for _, data := range req {
		userJobData := &store.UserJobData{
			ID:         data.ID,
			UserID:     userCred.ID,
			Role:       data.Role,
			Company:    data.Company,
			JoinedDate: data.JoinedYear,
		}

		if err := handler.userStore.InsertJobData(ctx, userJobData); err != nil {
			log.Println("error insert new job data: %w", err)
			response.Error(w, apierror.InternalServerError())
			return
		}
	}

	resp := AddUserJobDataResponse{
		Message: "success",
	}

	response.Respond(w, http.StatusCreated, resp)
}
