package assessment

import (
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	apierror "ta-spbe-backend/api/error"
	userCtx "ta-spbe-backend/api/handler/context"
	"ta-spbe-backend/api/response"
	"ta-spbe-backend/config"
	"ta-spbe-backend/repository"
	"ta-spbe-backend/service"
	"time"

	"github.com/google/uuid"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"google.golang.org/protobuf/proto"
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

func UploadSPBEDocument(
	assessmentRepo repository.AssessmentRepository,
	userRepo repository.UserRepository,
	mq service.MessageQueue,
	mailer service.Mailer,
	fsIO service.FileSystem,
	jsonEC service.JsonManipulator,
	waClient service.WhatsApp,
	apiCfg config.API,
	smtpCfg config.SMTPClient,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		supportingDocumentUrl := fmt.Sprintf("http://%s/static/%s", apiCfg.Host, supportingDocument)

		dst, err := fsIO.Create(fmt.Sprintf("./static/supporting-documents/%s", supportingDocument))
		if err != nil {
			log.Println(err.Error())
			response.Error(w, apierror.InternalServerError())
			return
		}
		defer dst.Close()

		_, err = fsIO.Copy(dst, req.supportingDocumentFile)
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
			UserId: userCred.ID,
		}
		err = assessmentRepo.InsertUploadDocument(ctx, &assessmentUploadDetail)

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

		producerPayload, err := jsonEC.Marshal(msg)
		if err != nil {
			log.Println(err)
			response.Error(w, apierror.InternalServerError())
			return
		}

		err = mq.Produce(topic, producerPayload)
		if err != nil {
			log.Println(err)
			response.Error(w, apierror.InternalServerError())
			return
		}

		user, err := userRepo.FindOneByID(ctx, userCred.ID)
		if err != nil {
			log.Println(err)
			response.Error(w, apierror.InternalServerError())
			return
		}

		to := []string{user.Email}
		subject, message := generateEmailContent()
		go func() {
			err := mailer.Send(subject, message, to, "upload.html", map[string]string{"username": "Conor"})
			if err != nil {
				log.Println("error send email: %w", err)
			}
		}()

		resp := UploadSpbeDocumentResponse{
			Message:      "Document has been successfully uploaded",
			AssessmentId: assessmentUploadDetail.AssessmentDetail.Id,
			DocumentUrl:  supportingDocumentUrl,
		}

		// 		const initialMessage = `*[OTOMATISASI PENILAIAN SPBE]*

		// ` + "```" + `Terima kasih telah menggunakan Aplikasi Otomatisasi Penilaian SPBE. Hasil penilaian anda akan keluar dalam beberapa saat lagi.` + "```"
		// protoMessage := &waProto.Message{
		// 	Conversation: proto.String(initialMessage),
		// }
		protoMessage := &waProto.Message{
			TemplateMessage: &waProto.TemplateMessage{
				HydratedTemplate: &waProto.TemplateMessage_HydratedFourRowTemplate{
					Title: &waProto.TemplateMessage_HydratedFourRowTemplate_HydratedTitleText{
						HydratedTitleText: "[OTOMATISASI PENILAIAN SPBE]",
					},
					TemplateId:          proto.String("template-id"),
					HydratedContentText: proto.String("Terima kasih telah menggunakan Aplikasi Otomatisasi Penilaian SPBE. Hasil penilaian anda akan keluar dalam beberapa saat lagi."),
					HydratedFooterText:  proto.String("APLIKASI OTOMATISASI PENILAIAN TINGKAT KEMATANGAN SPBE"),
					HydratedButtons: []*waProto.HydratedTemplateButton{

						// This for URL button
						{
							Index: proto.Uint32(1),
							HydratedButton: &waProto.HydratedTemplateButton_UrlButton{
								UrlButton: &waProto.HydratedTemplateButton_HydratedURLButton{
									DisplayText: proto.String("Otomatisasi Penilaian SPBE"),
									Url:         proto.String("https://fb.me/this"),
								},
							},
						},

						// This for call button
						{
							Index: proto.Uint32(2),
							HydratedButton: &waProto.HydratedTemplateButton_CallButton{
								CallButton: &waProto.HydratedTemplateButton_HydratedCallButton{
									DisplayText: proto.String("Hubungi Kami"),
									PhoneNumber: proto.String("1234567890"),
								},
							},
						},

						// This is just a quick reply
						{
							Index: proto.Uint32(3),
							HydratedButton: &waProto.HydratedTemplateButton_QuickReplyButton{
								QuickReplyButton: &waProto.HydratedTemplateButton_HydratedQuickReplyButton{
									DisplayText: proto.String("Quick reply"),
									Id:          proto.String("quick-id"),
								},
							},
						},
					},
				},
			},
		}

		err = waClient.SendMessage(ctx, req.phoneNumberStr, protoMessage)
		if err != nil {
			response.Error(w, apierror.InternalServerError())
			return
		}

		response.Respond(w, http.StatusCreated, resp)
	}
}

func generateEmailContent() (subject, message []byte) {
	subject = []byte("Otomatisasi Penilaian SPBE")
	message = []byte(fmt.Sprintf("Terima kasih telah menggunakan Aplikasi Otomatisasi Penilaian SPBE. Hasil penilaian anda akan keluar dalam beberapa saat lagi."))

	return
}
