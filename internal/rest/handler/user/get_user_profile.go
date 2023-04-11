package user

import (
	"log"
	"net/http"

	apierror "github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/rest/error"
	userCtx "github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/rest/handler/context"
	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/rest/response"
)

type GetUserProfileResponse struct {
	Name               string `json:"name"`
	ContactNumber      string `json:"contact_number"`
	Email              string `json:"email"`
	LinkedinProfile    string `json:"linkedin_profile"`
	HouseAddress       string `json:"house_address"`
	ProfilePictureLink string `json:"profile_picture_link"`
}

func (handler *userHandler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userCred, ok := ctx.Value(userCtx.UserCtxKey).(userCtx.UserCtx)
	if !ok {
		response.Error(w, apierror.InternalServerError())
		return
	}

	user, err := handler.userStore.FindOneByID(ctx, userCred.ID)
	if err != nil {
		log.Println(err)
		response.Error(w, apierror.InternalServerError())
		return
	}

	resp := GetUserProfileResponse{
		Name:               user.Name,
		ContactNumber:      user.ContactNumber,
		Email:              user.Email,
		LinkedinProfile:    user.LinkedinProfile,
		HouseAddress:       user.Address,
		ProfilePictureLink: user.ProfilePictureLink,
	}

	response.Respond(w, http.StatusOK, resp)
}
