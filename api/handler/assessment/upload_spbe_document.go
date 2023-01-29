package assessment

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	apierror "ta-spbe-backend/api/error"
	"ta-spbe-backend/api/response"
	"ta-spbe-backend/config"
	"ta-spbe-backend/repository"

	"github.com/google/uuid"
)

const MAX_UPLOAD_SIZE = 10240 * 10240 // 10MB

type UploadSpbeDocumentRequest struct {
	institutionName              string
	indicatorNumberStr           string
	indicatorNumber              int
	supportingDocumentFile       multipart.File
	supportingDocumentFileHeader *multipart.FileHeader
	oldDocumentFile              multipart.File
	oldDocumentFileHeader        *multipart.FileHeader
}

func (req *UploadSpbeDocumentRequest) validate(r *http.Request) *apierror.FieldError {
	var err error
	fieldErr := apierror.NewFieldError()

	if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
		fieldErr = fieldErr.WithField("supporting_document", "the uploaded file is too big! The maximum allowed size is 1MB")
	}

	req.supportingDocumentFile, req.supportingDocumentFileHeader, err = r.FormFile("supporting_document")
	if err != nil {
		fieldErr = fieldErr.WithField("supporting_document", "supporting document is missing")
	}

	req.institutionName = strings.TrimSpace(req.institutionName)
	if req.institutionName == "" {
		fieldErr = fieldErr.WithField("institution_name", "institution name is missing")
	}

	req.indicatorNumberStr = strings.TrimSpace(req.indicatorNumberStr)
	if req.indicatorNumberStr == "" {
		fieldErr = fieldErr.WithField("indicator_number", "indicator number is missing")
	}

	req.indicatorNumber, err = strconv.Atoi(req.indicatorNumberStr)
	if err != nil || req.indicatorNumber < 1 && req.indicatorNumber > 10 {
		fieldErr = fieldErr.WithField("indicator_number", "indicator number must be a positive integer ranging between 1 and 10")
	}

	if len(fieldErr.Fields) != 0 {
		return &fieldErr
	}

	return nil
}

type UploadSpbeDocumentResponse struct {
	Message      string `json:"string"`
	AssessmentId string `json:"assessment_id"`
	DocumentUrl  string `json:"document_url"`
}

func UploadSPBEDocument(assessmentRepo repository.AssessmentRepository, apiCfg config.API) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)
		req := UploadSpbeDocumentRequest{
			institutionName:    r.FormValue("institution_name"),
			indicatorNumberStr: r.FormValue("indicator_number"),
		}

		fieldErr := req.validate(r)
		if fieldErr != nil {
			response.FieldError(w, *fieldErr)
			return
		}

		defer req.supportingDocumentFile.Close()

		uniqueId := uuid.New()
		filename := strings.Replace(uniqueId.String(), "-", "", -1)
		fileExt := filepath.Ext(req.supportingDocumentFileHeader.Filename)
		supportingDocument := fmt.Sprintf("%s%s", filename, fileExt)
		supportingDocumentUrl := fmt.Sprintf("http://%s/static/%s", apiCfg.Host, supportingDocument)

		dst, err := os.Create(fmt.Sprintf("./static/supporting-documents/%s", supportingDocument))
		if err != nil {
			log.Println(err.Error())
			response.Error(w, apierror.InternalServerError())
			return
		}
		defer dst.Close()

		_, err = io.Copy(dst, req.supportingDocumentFile)
		if err != nil {
			log.Println(err.Error())
			response.Error(w, apierror.InternalServerError())
			return
		}

		assessmentUploadDetail := repository.AssessmentUploadDetail{
			AssessmentDetail: repository.AssessmentDetail{
				InstitutionName: req.institutionName,
			},
			IndicatorAssessmentInfo: repository.IndicatorAssessmentInfo{
				IndicatorNumber: req.indicatorNumber,
			},
			SupportDataDocumentInfo: repository.SupportDataDocumentInfo{
				DocumentName: supportingDocument,
				DocumentUrl:  supportingDocumentUrl,
			},
			UserId: "ccd52961-fa4e-43ba-a6df-a4c97849d899",
		}
		err = assessmentRepo.InsertUploadDocument(ctx, &assessmentUploadDetail)

		resp := UploadSpbeDocumentResponse{
			Message:      "Document has been successfully uploaded",
			AssessmentId: assessmentUploadDetail.AssessmentDetail.Id,
			DocumentUrl:  supportingDocumentUrl,
		}

		response.Respond(w, http.StatusCreated, resp)
	}
}
