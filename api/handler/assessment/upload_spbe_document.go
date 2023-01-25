package assessment

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	apierror "ta-spbe-backend/api/error"
	"ta-spbe-backend/api/response"
	"ta-spbe-backend/repository"

	"github.com/google/uuid"
)

const MAX_UPLOAD_SIZE = 10240 * 10240 // 10MB

type UploadSpbeDocumentResponse struct {
	Message      string `json:"string"`
	AssessmentId string `json:"assessment_id"`
}

func UploadSPBEDocument(assessmentRepo repository.AssessmentRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// ctx := r.Context()

		r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)

		if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
			response.Error(w, apierror.BadRequestError("The uploaded file is too big! The maximum allowed size is 1MB"))
			return
		}

		file, fileHeader, err := r.FormFile("supporting_document")
		if err != nil {
			response.Error(w, apierror.BadRequestError("Supporting document is missing"))
			return
		}

		log.Println(r.FormValue("institution_name"))
		log.Println(r.FormValue("indicator_number"))

		defer file.Close()

		uniqueId := uuid.New()
		filename := strings.Replace(uniqueId.String(), "-", "", -1)
		fileExt := filepath.Ext(fileHeader.Filename)
		supportingDocument := fmt.Sprintf("%s%s", filename, fileExt)

		dst, err := os.Create(fmt.Sprintf("./static/supporting-documents/%s", supportingDocument))
		if err != nil {
			response.Error(w, apierror.InternalServerError())
			return
		}

		defer dst.Close()

		_, err = io.Copy(dst, file)
		if err != nil {
			response.Error(w, apierror.InternalServerError())
			return
		}

		resp := UploadSpbeDocumentResponse{
			Message:      "Document has been successfully uploaded",
			AssessmentId: uniqueId.String(),
		}

		response.Respond(w, http.StatusCreated, resp)
	}
}
