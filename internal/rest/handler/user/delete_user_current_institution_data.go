package user

import (
	"net/http"

	apierror "github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/rest/error"
	"github.com/go-chi/chi/v5"

	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/rest/response"
)

type DeleteUserCurrentInstitutionDataResponse struct {
	Message string `json:"message"`
}

func (handler *userHandler) DeleteUserCurrentInstitutionData(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userInstitutionId := chi.URLParam(r, "id")

	err := handler.userStore.DeleteCurrentInstitutionByID(ctx, userInstitutionId)
	if err != nil {
		response.Error(w, apierror.NotFoundError("user institution data not found"))
		return
	}

	resp := DeleteUserCurrentInstitutionDataResponse{
		Message: "success",
	}

	response.Respond(w, http.StatusOK, resp)
}
