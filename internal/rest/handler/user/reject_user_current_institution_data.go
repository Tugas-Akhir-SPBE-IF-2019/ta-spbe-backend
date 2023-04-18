package user

import (
	"fmt"
	"net/http"

	apierror "github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/rest/error"
	"github.com/go-chi/chi/v5"

	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/rest/response"
)

type RejectUserCurrentInstitutionDataResponse struct {
	Message string `json:"message"`
}

func (handler *userHandler) RejectUserCurrentInstitutionData(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userInstitutionId := chi.URLParam(r, "id")

	err := handler.userStore.DeleteCurrentInstitutionByID(ctx, userInstitutionId)
	if err != nil {
		response.Error(w, apierror.NotFoundError("user institution data not found"))
		return
	}

	resp := RejectUserCurrentInstitutionDataResponse{
		Message: "success",
	}

	subject := []byte("Your Request Was Rejected")
	message := []byte(fmt.Sprintf("Your Request to add a new institution was rejected. Please check your spelling or check if the institution that you want to add is already exist"))
	to := []string{
		"13519142@std.stei.itb.ac.id",
	}

	handler.smtpMailer.SendSimple(
		subject,
		message,
		to,
	)

	response.Respond(w, http.StatusOK, resp)
}
