package user

import (
	"encoding/json"
	"log"
	"net/http"

	apierror "github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/rest/error"
	userCtx "github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/rest/handler/context"
	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/store"
	"github.com/google/uuid"

	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/rest/response"
)

type AddUserEvaluationDataRequest struct {
	Role           string `json:"role"`
	InstitutionID  int    `json:"institution_id"`
	EvaluationYear int    `json:"evaluation_year"`
}

type AddUserEvaluationDataResponse struct {
	Message string `json:"message"`
}

func (handler *userHandler) AddUserEvaluationData(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userCred, ok := ctx.Value(userCtx.UserCtxKey).(userCtx.UserCtx)
	if !ok {
		response.Error(w, apierror.InternalServerError())
		return
	}

	req := AddUserEvaluationDataRequest{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, apierror.BadRequestError(err.Error()))
		return
	}

	userEvaluationData := &store.UserEvaluationData{
		ID:             uuid.NewString(),
		UserID:         userCred.ID,
		Role:           req.Role,
		InstitutionID:  req.InstitutionID,
		EvaluationYear: req.EvaluationYear,
	}

	if err := handler.userStore.InsertEvaluationData(ctx, userEvaluationData); err != nil {
		log.Println("error insert new evaluation data: %w", err)
		response.Error(w, apierror.InternalServerError())
		return
	}

	resp := AddUserEvaluationDataResponse{
		Message: "success",
	}

	response.Respond(w, http.StatusCreated, resp)
}