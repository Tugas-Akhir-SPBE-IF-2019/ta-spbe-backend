package user

import (
	"net/http"

	apierror "github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/rest/error"
	"github.com/go-chi/chi/v5"

	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/rest/response"
)

type VerifyUserCurrentInstitutionDataResponse struct {
	Message string `json:"message"`
}

func (handler *userHandler) VerifyUserCurrentInstitutionData(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userInstitutionId := chi.URLParam(r, "id")

	err := handler.userStore.VerifyInstitutionData(ctx, userInstitutionId)
	if err != nil {
		response.Error(w, apierror.NotFoundError("user institution data not found"))
		return
	}

	resp := VerifyUserCurrentInstitutionDataResponse{
		Message: "success",
	}

	response.Respond(w, http.StatusOK, resp)
}
