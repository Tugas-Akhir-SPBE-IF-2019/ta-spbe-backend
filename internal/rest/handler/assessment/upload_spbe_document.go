package assessment

import (
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	apierror "github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/rest/error"
	userCtx "github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/rest/handler/context"
	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/rest/response"
	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/store"

	"github.com/google/uuid"
)

const MAX_UPLOAD_SIZE = 10240 * 10240 // 10MB

type UploadSpbeDocumentRequest struct {
	institutionName              string
	indicatorNumberStr           string
	indicatorNumber              int
	phoneNumberStr               string
	supportingDocumentFile       multipart.File
	supportingDocumentFileHeader *multipart.FileHeader
	oldDocumentFile              multipart.File
	oldDocumentFileHeader        *multipart.FileHeader
}

type UploadProducerMessage struct {
	Name                  string
	Content               string
	UserId                string
	RecipientNumber       string
	AssessmentId          string
	IndicatorAssessmentId string
	Timestamp             string
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

	// FOR TESTING PURPOSE
	if req.phoneNumberStr == "" {
		req.phoneNumberStr = "6285157017311"
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

func (handler *assessmentHandler) UploadSPBEDocument(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)
	req := UploadSpbeDocumentRequest{
		institutionName:    r.FormValue("institution_name"),
		indicatorNumberStr: r.FormValue("indicator_number"),
		phoneNumberStr:     r.FormValue("phone_number"),
	}

	fieldErr := req.validate(r)
	if fieldErr != nil {
		response.FieldError(w, *fieldErr)
		return
	}

	userCred, ok := r.Context().Value(userCtx.UserCtxKey).(userCtx.UserCtx)
	if !ok {
		response.Error(w, apierror.InternalServerError())
		return
	}

	defer req.supportingDocumentFile.Close()

	uniqueId := uuid.New()
	filename := strings.Replace(uniqueId.String(), "-", "", -1)
	fileExt := filepath.Ext(req.supportingDocumentFileHeader.Filename)
	supportingDocument := fmt.Sprintf("%s%s", filename, fileExt)
	supportingDocumentUrl := fmt.Sprintf("http://%s/static/%s", handler.apiCfg.Host, supportingDocument)

	var dst *os.File
	dst, err := handler.filesystemClient.Create(fmt.Sprintf("./static/supporting-documents/%s", supportingDocument))
	if err != nil {
		log.Println(err.Error())

		os.MkdirAll("./static/supporting-documents", os.ModePerm)
		dst, _ = handler.filesystemClient.Create(fmt.Sprintf("./static/supporting-documents/%s", supportingDocument))
	}
	defer dst.Close()

	_, err = handler.filesystemClient.Copy(dst, req.supportingDocumentFile)
	if err != nil {
		log.Println(err.Error())
		response.Error(w, apierror.InternalServerError())
		return
	}

	assessmentUploadDetail := store.AssessmentUploadDetail{
		AssessmentDetail: store.AssessmentDetail{
			InstitutionName: req.institutionName,
		},
		IndicatorAssessmentInfo: store.IndicatorAssessmentInfo{
			IndicatorNumber: req.indicatorNumber,
		},
		SupportDataDocumentInfo: store.SupportDataDocumentInfo{
			DocumentName: supportingDocument,
			DocumentUrl:  supportingDocumentUrl,
		},
		UserId: userCred.ID,
	}
	err = handler.assessmentStore.InsertUploadDocument(ctx, &assessmentUploadDetail)

	topic := "SPBE_Assessment"
	msg := UploadProducerMessage{
		Name:                  assessmentUploadDetail.AssessmentDetail.InstitutionName,
		Content:               assessmentUploadDetail.SupportDataDocumentInfo.Id,
		UserId:                userCred.ID,
		RecipientNumber:       req.phoneNumberStr,
		AssessmentId:          assessmentUploadDetail.AssessmentDetail.Id,
		IndicatorAssessmentId: assessmentUploadDetail.IndicatorAssessmentInfo.Id,
		Timestamp:             time.Now().UTC().String(),
	}

	producerPayload, err := handler.jsonClient.Marshal(msg)
	if err != nil {
		log.Println(err)
		response.Error(w, apierror.InternalServerError())
		return
	}

	err = handler.messageQueue.Produce(topic, producerPayload)
	if err != nil {
		log.Println(err)
		response.Error(w, apierror.InternalServerError())
		return
	}

	user, err := handler.userStore.FindOneByID(ctx, userCred.ID)
	if err != nil {
		log.Println(err)
		response.Error(w, apierror.InternalServerError())
		return
	}

	to := []string{user.Email}
	subject, message := generateEmailContent()
	go func() {
		err := handler.smtpMailer.Send(subject, message, to, "upload.html", map[string]string{"username": "Conor"})
		if err != nil {
			log.Println("error send email: %w", err)
		}
	}()

	resp := UploadSpbeDocumentResponse{
		Message:      "Document has been successfully uploaded",
		AssessmentId: assessmentUploadDetail.AssessmentDetail.Id,
		DocumentUrl:  supportingDocumentUrl,
	}

	response.Respond(w, http.StatusCreated, resp)
}

func generateEmailContent() (subject, message []byte) {
	subject = []byte("Otomatisasi Penilaian SPBE")
	message = []byte(fmt.Sprintf("Terima kasih telah menggunakan Aplikasi Otomatisasi Penilaian SPBE. Hasil penilaian anda akan keluar dalam beberapa saat lagi."))

	return
}
