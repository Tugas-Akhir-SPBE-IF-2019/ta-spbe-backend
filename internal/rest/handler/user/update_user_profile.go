package user

import (
	"log"
	"net/http"

	apierror "github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/rest/error"
	userCtx "github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/rest/handler/context"
	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/rest/response"
	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/store"
)

const MAX_UPLOAD_SIZE = 10240 * 10240 // 10MB

type UpdateUserProfileRequest struct {
	email              string
	name               string
	contactNumber      string
	linkedinProfile    string
	address            string
	profilePictureLink string
}

func (handler *userHandler) UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)
	req := UpdateUserProfileRequest{
		email:           r.FormValue("email"),
		name:            r.FormValue("name"),
		contactNumber:   r.FormValue("contact_number"),
		linkedinProfile: r.FormValue("linkedin_profile"),
		address:         r.FormValue("address"),
	}

	userCred, ok := ctx.Value(userCtx.UserCtxKey).(userCtx.UserCtx)
	if !ok {
		response.Error(w, apierror.InternalServerError())
		return
	}

	user := &store.UserData{
		ID:                 userCred.ID,
		Name:               req.name,
		Email:              req.email,
		ContactNumber:      req.contactNumber,
		LinkedinProfile:    req.linkedinProfile,
		Address:            req.address,
		ProfilePictureLink: req.profilePictureLink,
	}

	if err := handler.userStore.UpdateByID(ctx, user); err != nil {
		log.Println("error update user data: %w", err)
		response.Error(w, apierror.InternalServerError())
		return
	}

	response.Respond(w, http.StatusCreated, nil)
}
