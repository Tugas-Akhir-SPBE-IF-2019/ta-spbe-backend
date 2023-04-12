package user

import (
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	apierror "github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/rest/error"
	userCtx "github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/rest/handler/context"
	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/rest/response"
	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/store"
	"github.com/google/uuid"
)

const MAX_UPLOAD_SIZE = 10240 * 10240 // 10MB

type UpdateUserProfileRequest struct {
	email                        string
	name                         string
	contactNumber                string
	linkedinProfile              string
	address                      string
	profilePictureLink           string
	profilePictureFile           multipart.File
	profilePictureFileHeader     *multipart.FileHeader
	profilePictureFileHeaderList []*multipart.FileHeader
}

type UpdateUserProfleResponse struct {
	Message string `json:"message"`
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

	fhs := r.MultipartForm.File["profile_picture"]
	for _, fh := range fhs {
		req.profilePictureFileHeaderList = append(req.profilePictureFileHeaderList, fh)
	}

	for _, profilePictureFileHeader := range req.profilePictureFileHeaderList {
		profilePicturetFile, err := profilePictureFileHeader.Open()
		if err != nil {
			log.Println(err.Error())
			response.Error(w, apierror.InternalServerError())
			return
		}
		defer profilePicturetFile.Close()

		uniqueId := uuid.New()
		filename := strings.Replace(uniqueId.String(), "-", "", -1)
		fileExt := filepath.Ext(profilePictureFileHeader.Filename)
		profilePicture := fmt.Sprintf("%s%s", filename, fileExt)
		profilePictureUrl := fmt.Sprintf("http://%s/static/%s", handler.apiCfg.Host, profilePicture)

		var dst *os.File
		dst, err = handler.filesystemClient.Create(fmt.Sprintf("./static/supporting-documents/%s", profilePicture))
		if err != nil {
			log.Println(err.Error())

			os.MkdirAll("./static/supporting-documents", os.ModePerm)
			dst, _ = handler.filesystemClient.Create(fmt.Sprintf("./static/supporting-documents/%s", profilePicture))
		}
		defer dst.Close()

		_, err = handler.filesystemClient.Copy(dst, profilePicturetFile)
		if err != nil {
			log.Println(err.Error())
			response.Error(w, apierror.InternalServerError())
			return
		}

		req.profilePictureLink = profilePictureUrl
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

	if user.ProfilePictureLink == "" {
		if err := handler.userStore.UpdateByID(ctx, user); err != nil {
			log.Println("error update user data: %w", err)
			response.Error(w, apierror.InternalServerError())
			return
		}
	} else {
		if err := handler.userStore.UpdateWithPhotoByID(ctx, user); err != nil {
			log.Println("error update user data: %w", err)
			response.Error(w, apierror.InternalServerError())
			return
		}
	}

	response.Respond(w, http.StatusOK, UpdateUserProfleResponse{
		Message: "success",
	})
}
