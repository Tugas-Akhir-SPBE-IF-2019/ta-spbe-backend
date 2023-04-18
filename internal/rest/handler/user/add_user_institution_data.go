package user

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	apierror "github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/rest/error"
	userCtx "github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/rest/handler/context"
	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/store"

	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/rest/response"
)

type AddUserInstitutionRequest struct {
	ID              string `json:"id"`
	InstitutionName string `json:"institution_name"`
	Role            string `json:"role"`
}

type AddUserInstitutionResponse struct {
	Message string `json:"message"`
}

func (handler *userHandler) AddUserInstitutionData(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userCred, ok := ctx.Value(userCtx.UserCtxKey).(userCtx.UserCtx)
	if !ok {
		response.Error(w, apierror.InternalServerError())
		return
	}

	req := []AddUserInstitutionRequest{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, apierror.BadRequestError(err.Error()))
		return
	}

	for _, data := range req {
		userInstitutionData := &store.UserCurrentInstitutionData{
			ID:              data.ID,
			UserID:          userCred.ID,
			Role:            data.Role,
			InstitutionName: data.InstitutionName,
		}

		institution, err := handler.institutionStore.FindByInstitutionName(ctx, data.InstitutionName)
		if err != nil && err != sql.ErrNoRows {
			log.Println(err.Error())
			response.Error(w, apierror.InternalServerError())
			return
		}

		userInstitutionData.InstitutionID = institution.ID

		if err := handler.userStore.InsertCurrentInstitutionData(ctx, userInstitutionData); err != nil {
			log.Println("error insert new user institution data: %w", err)
			response.Error(w, apierror.InternalServerError())
			return
		}

		if userInstitutionData.InstitutionID == 0 {
			subject := []byte("User Wants to Add a New Institution")
			message := []byte(fmt.Sprintf("User with id %s  wants to add a new Institution: %s. To approve, click this link http://localhost/users/institution/%s/approve . To reject click this link: http://localhost/users/institution/%s/reject", userCred.ID, userInstitutionData.InstitutionName, userInstitutionData.ID, userInstitutionData.ID))
			to := []string{
				"13519142@std.stei.itb.ac.id",
			}
			handler.smtpMailer.SendSimple(
				subject,
				message,
				to,
			)
		}
	}

	resp := AddUserInstitutionResponse{
		Message: "success",
	}

	response.Respond(w, http.StatusCreated, resp)
}
