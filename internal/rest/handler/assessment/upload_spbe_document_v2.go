package assessment

import (
	"fmt"
	"log"
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

// const MAX_UPLOAD_SIZE = 10240 * 10240 // 10MB

// type UploadSpbeDocumentRequest struct {
// 	institutionName                 string
// 	indicatorNumberStr              string
// 	indicatorNumbersStr             []string
// 	indicatorNumber                 int
// 	indicatorNumbers                []int
// 	phoneNumberStr                  string
// 	supportingDocumentFile          multipart.File
// 	supportingDocumentFileHeader    *multipart.FileHeader
// 	supportingDocumenFileHeaderList []*multipart.FileHeader
// 	oldDocumentFile                 multipart.File
// 	oldDocumentFileHeader           *multipart.FileHeader
// }

// type UploadProducerMessage struct {
// 	Name                  string
// 	Content               string
// 	UserId                string
// 	RecipientNumber       string
// 	AssessmentId          string
// 	IndicatorAssessmentId string
// 	Filename              string
// 	OriginalFilename      string
// 	Timestamp             string
// 	IndicatorNumber       string
// 	IndicatorDetail       string
// 	InstitutionName       string
// }

// func (req *UploadSpbeDocumentRequest) validate(r *http.Request) *apierror.FieldError {
// 	fieldErr := apierror.NewFieldError()

// 	if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
// 		fieldErr = fieldErr.WithField("supporting_document", "the uploaded file is too big! The maximum allowed size is 1MB")
// 	}

// 	fhs := r.MultipartForm.File["supporting_document[]"]
// 	for _, fh := range fhs {
// 		req.supportingDocumenFileHeaderList = append(req.supportingDocumenFileHeaderList, fh)
// 	}
// 	if len(req.supportingDocumenFileHeaderList) == 0 {
// 		fieldErr = fieldErr.WithField("supporting_document", "supporting document is missing")
// 	}

// 	req.institutionName = strings.TrimSpace(req.institutionName)
// 	if req.institutionName == "" {
// 		fieldErr = fieldErr.WithField("institution_name", "institution name is missing")
// 	}

// 	for _, indicatorNumberStr := range req.indicatorNumbersStr {
// 		indicatorNumber, err := strconv.Atoi(indicatorNumberStr)
// 		if err != nil {
// 			fieldErr = fieldErr.WithField("indicator_number", "indicator number must be a positive integer ranging between 1 and 10")
// 		}
// 		req.indicatorNumbers = append(req.indicatorNumbers, indicatorNumber)
// 	}

// 	// FOR TESTING PURPOSE
// 	if req.phoneNumberStr == "" {
// 		req.phoneNumberStr = "6285157017311"
// 	}

// 	if len(fieldErr.Fields) != 0 {
// 		return &fieldErr
// 	}

// 	return nil
// }

// type DocumentInfo struct {
// 	Name string `json:"name"`
// 	Url  string `json:"url"`
// }

// type UploadSpbeDocumentResponse struct {
// 	Message          string         `json:"string"`
// 	AssessmentId     string         `json:"assessment_id"`
// 	DocumentInfoList []DocumentInfo `json:"documents_info"`
// }

func (handler *assessmentHandler) UploadSPBEDocumentV2(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)
	req := UploadSpbeDocumentRequest{
		institutionName: r.FormValue("institution_name"),
		phoneNumberStr: r.FormValue("phone_number"),
	}

	r.ParseForm()
	req.indicatorNumbersStr = r.Form["indicator_number[]"]

	fieldErr := req.validate(r)
	if fieldErr != nil {
		response.FieldError(w, *fieldErr)
		return
	}

	userCred, ok := ctx.Value(userCtx.UserCtxKey).(userCtx.UserCtx)
	if !ok {
		response.Error(w, apierror.InternalServerError())
		return
	}

	var documentInfoList []DocumentInfo
	var supportDataDocumentInfoList []store.SupportDataDocumentInfo
	for _, supportingDocumenFileHeader := range req.supportingDocumenFileHeaderList {
		supportingDocumentFile, err := supportingDocumenFileHeader.Open()
		if err != nil {
			log.Println(err.Error())
			response.Error(w, apierror.InternalServerError())
			return
		}
		defer supportingDocumentFile.Close()

		uniqueId := uuid.New()
		filename := strings.Replace(uniqueId.String(), "-", "", -1)
		fileExt := filepath.Ext(supportingDocumenFileHeader.Filename)
		originalDocumentName := supportingDocumenFileHeader.Filename
		supportingDocument := fmt.Sprintf("%s%s", filename, fileExt)
		supportingDocumentUrl := fmt.Sprintf("http://%s/static/%s", handler.apiCfg.Host, supportingDocument)

		var dst *os.File
		dst, err = handler.filesystemClient.Create(fmt.Sprintf("./static/supporting-documents/%s", supportingDocument))
		if err != nil {
			log.Println(err.Error())

			os.MkdirAll("./static/supporting-documents", os.ModePerm)
			dst, _ = handler.filesystemClient.Create(fmt.Sprintf("./static/supporting-documents/%s", supportingDocument))
		}
		defer dst.Close()

		_, err = handler.filesystemClient.Copy(dst, supportingDocumentFile)
		if err != nil {
			log.Println(err.Error())
			response.Error(w, apierror.InternalServerError())
			return
		}

		supportDataDocumentInfoList = append(supportDataDocumentInfoList, store.SupportDataDocumentInfo{
			DocumentName:         supportingDocument,
			DocumentUrl:          supportingDocumentUrl,
			OriginalDocumentName: originalDocumentName,
		})
		documentInfoList = append(documentInfoList, DocumentInfo{
			Name: originalDocumentName,
			Url:  supportingDocumentUrl,
		})

	}

	assessmentUploadDetail := store.AssessmentUploadDetail{
		AssessmentDetail: store.AssessmentDetail{
			InstitutionName: req.institutionName,
		},
		SupportDataDocumentInfoList: supportDataDocumentInfoList,
		UserId:                      userCred.ID,
	}
	for _, indicatorNumber := range req.indicatorNumbers {
		assessmentUploadDetail.IndicatorAssessmentInfo.IndicatorNumber = indicatorNumber
		err := handler.assessmentStore.InsertUploadDocument(ctx, &assessmentUploadDetail)
		if err != nil {
			log.Println(err)
			response.Error(w, apierror.InternalServerError())
			return
		}

		indicatorData, err := handler.indicatoreAssessmentStore.FindIndicatorDetailByIndicatorNumber(ctx, indicatorNumber)
		if err != nil {
			log.Println(err)
			response.Error(w, apierror.InternalServerError())
			return
		}

		topic := "SPBE_Assessment"
		msg := UploadProducerMessage{
			Name:                  assessmentUploadDetail.AssessmentDetail.InstitutionName,
			Content:               assessmentUploadDetail.SupportDataDocumentInfoList[0].Id, // WIP need to be updated later
			UserId:                userCred.ID,
			RecipientNumber:       req.phoneNumberStr,
			AssessmentId:          assessmentUploadDetail.AssessmentDetail.Id,
			IndicatorAssessmentId: assessmentUploadDetail.IndicatorAssessmentInfo.Id,
			Filename:              assessmentUploadDetail.SupportDataDocumentInfoList[0].DocumentName, // WIP need to be updated later
			OriginalFilename:      supportDataDocumentInfoList[0].OriginalDocumentName,
			Timestamp:             time.Now().UTC().String(),
			IndicatorNumber:       strconv.Itoa(indicatorNumber),
			IndicatorDetail:       indicatorData.Detail,
			InstitutionName:       req.institutionName,
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
		Message:          "Document has been successfully uploaded",
		AssessmentId:     assessmentUploadDetail.AssessmentDetail.Id,
		DocumentInfoList: documentInfoList,
	}

	protoMessage := generateWhatsAppMessage()
	err = handler.waClient.SendMessage(ctx, req.phoneNumberStr, protoMessage)
	if err != nil {
		response.Error(w, apierror.InternalServerError())
		return
	}

	response.Respond(w, http.StatusCreated, resp)
}

// func generateEmailContent() (subject, message []byte) {
// 	subject = []byte("Otomatisasi Penilaian SPBE")
// 	message = []byte(fmt.Sprintf("Terima kasih telah menggunakan Aplikasi Otomatisasi Penilaian SPBE. Hasil penilaian anda akan keluar dalam beberapa saat lagi."))

// 	return
// }

// func generateWhatsAppMessage() *waProto.Message {
// 	return &waProto.Message{
// 		TemplateMessage: &waProto.TemplateMessage{
// 			HydratedTemplate: &waProto.TemplateMessage_HydratedFourRowTemplate{
// 				Title: &waProto.TemplateMessage_HydratedFourRowTemplate_HydratedTitleText{
// 					HydratedTitleText: "[OTOMATISASI PENILAIAN SPBE]",
// 				},
// 				TemplateId:          proto.String("template-id"),
// 				HydratedContentText: proto.String("Terima kasih telah menggunakan Aplikasi Otomatisasi Penilaian SPBE. Hasil penilaian anda akan keluar dalam beberapa saat lagi."),
// 				HydratedFooterText:  proto.String("APLIKASI OTOMATISASI PENILAIAN TINGKAT KEMATANGAN SPBE"),
// 				HydratedButtons: []*waProto.HydratedTemplateButton{

// 					// This for URL button
// 					{
// 						Index: proto.Uint32(1),
// 						HydratedButton: &waProto.HydratedTemplateButton_UrlButton{
// 							UrlButton: &waProto.HydratedTemplateButton_HydratedURLButton{
// 								DisplayText: proto.String("Otomatisasi Penilaian SPBE"),
// 								Url:         proto.String("https://fb.me/this"),
// 							},
// 						},
// 					},

// 					// This for call button
// 					{
// 						Index: proto.Uint32(2),
// 						HydratedButton: &waProto.HydratedTemplateButton_CallButton{
// 							CallButton: &waProto.HydratedTemplateButton_HydratedCallButton{
// 								DisplayText: proto.String("Hubungi Kami"),
// 								PhoneNumber: proto.String("1234567890"),
// 							},
// 						},
// 					},

// 					// This is just a quick reply
// 					{
// 						Index: proto.Uint32(3),
// 						HydratedButton: &waProto.HydratedTemplateButton_QuickReplyButton{
// 							QuickReplyButton: &waProto.HydratedTemplateButton_HydratedQuickReplyButton{
// 								DisplayText: proto.String("Quick reply"),
// 								Id:          proto.String("quick-id"),
// 							},
// 						},
// 					},
// 				},
// 			},
// 		},
// 	}
// }
