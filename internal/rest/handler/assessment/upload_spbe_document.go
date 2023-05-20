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
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"google.golang.org/protobuf/proto"
)

const MAX_UPLOAD_SIZE = 1024 * 1024 * 10 // 10MB

type UploadSpbeDocumentRequest struct {
	institutionName                 string
	indicatorNumberStr              string
	indicatorNumbersStr             []string
	indicatorNumber                 int
	indicatorNumbers                []int
	phoneNumberStr                  string
	supportingDocumentFile          multipart.File
	supportingDocumentFileHeader    *multipart.FileHeader
	supportingDocumenFileHeaderList []*multipart.FileHeader
	oldDocumentFile                 multipart.File
	oldDocumentFileHeader           *multipart.FileHeader
	oldDocumentFileHeaderList       []*multipart.FileHeader
	meetingMinutesFile              multipart.File
	meetingMinutesFileHeader        *multipart.FileHeader
	meetingMinutestFileHeaderList   []*multipart.FileHeader
}

type IndicatorAssessmentProducerItem struct {
	Id     string `json:"id"`
	Number int    `json:"number"`
	Detail string `json:"detail"`
}

type SupportDocumentProducerItem struct {
	Name         string `json:"name"`
	OriginalName string `json:"original_name"`
	Type         string `json:"type"`
}

type UploadProducerMessage struct {
	Name                    string                            `json:"name"`
	Content                 string                            `json:"content"`
	UserId                  string                            `json:"user_id"`
	RecipientNumber         string                            `json:"recipient_number"`
	AssessmentId            string                            `json:"assessment_id"`
	IndicatorAssessmentList []IndicatorAssessmentProducerItem `json:"indicator_assessment_list"`
	SupportDocumentList     []SupportDocumentProducerItem     `json:"support_document_list"`
	IndicatorDetail         string                            `json:"indicator_detail"`
	InstitutionName         string                            `json:"institution_name"`
	Timestamp               string                            `json:"timestamp"`
}

func (req *UploadSpbeDocumentRequest) validate(r *http.Request) *apierror.FieldError {
	fieldErr := apierror.NewFieldError()

	if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
		fieldErr = fieldErr.WithField("supporting_document", "the uploaded file is too big! The maximum allowed size is 1MB")
	}

	fhs := r.MultipartForm.File["supporting_document[]"]
	for _, fh := range fhs {
		req.supportingDocumenFileHeaderList = append(req.supportingDocumenFileHeaderList, fh)
	}
	if len(req.supportingDocumenFileHeaderList) == 0 {
		fieldErr = fieldErr.WithField("supporting_document", "supporting document is missing")
	}

	fhs = r.MultipartForm.File["old_document[]"]
	for _, fh := range fhs {
		req.oldDocumentFileHeaderList = append(req.oldDocumentFileHeaderList, fh)
	}

	fhs = r.MultipartForm.File["meeting_minutes[]"]
	for _, fh := range fhs {
		req.meetingMinutestFileHeaderList = append(req.meetingMinutestFileHeaderList, fh)
	}

	req.institutionName = strings.TrimSpace(req.institutionName)
	if req.institutionName == "" {
		fieldErr = fieldErr.WithField("institution_name", "institution name is missing")
	}

	for _, indicatorNumberStr := range req.indicatorNumbersStr {
		indicatorNumber, err := strconv.Atoi(indicatorNumberStr)
		if err != nil {
			fieldErr = fieldErr.WithField("indicator_number", "indicator number must be a positive integer ranging between 1 and 10")
		}
		req.indicatorNumbers = append(req.indicatorNumbers, indicatorNumber)
	}

	if len(fieldErr.Fields) != 0 {
		return &fieldErr
	}

	return nil
}

type DocumentInfo struct {
	Name string `json:"name"`
	Url  string `json:"url"`
	Type string `json:"type"`
}

type UploadSpbeDocumentResponse struct {
	Message          string         `json:"string"`
	AssessmentId     string         `json:"assessment_id"`
	DocumentInfoList []DocumentInfo `json:"documents_info"`
}

func (handler *assessmentHandler) UploadSPBEDocument(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)
	req := UploadSpbeDocumentRequest{
		institutionName: r.FormValue("institution_name"),
		phoneNumberStr:  r.FormValue("phone_number"),
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
	var supportDocumentProducerItemList []SupportDocumentProducerItem
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
			Type:                 store.NEW_DOCUMENT,
		})
		supportDocumentProducerItemList = append(supportDocumentProducerItemList, SupportDocumentProducerItem{
			Name:         supportingDocument,
			OriginalName: originalDocumentName,
			Type:         string(store.NEW_DOCUMENT),
		})
		documentInfoList = append(documentInfoList, DocumentInfo{
			Name: originalDocumentName,
			Type: string(store.NEW_DOCUMENT),
			Url:  supportingDocumentUrl,
		})

	}

	for _, supportingDocumenFileHeader := range req.oldDocumentFileHeaderList {
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
			Type:                 store.OLD_DOCUMENT,
		})
		supportDocumentProducerItemList = append(supportDocumentProducerItemList, SupportDocumentProducerItem{
			Name:         supportingDocument,
			OriginalName: originalDocumentName,
			Type:         string(store.OLD_DOCUMENT),
		})
		documentInfoList = append(documentInfoList, DocumentInfo{
			Name: originalDocumentName,
			Type: string(store.OLD_DOCUMENT),
			Url:  supportingDocumentUrl,
		})

	}

	for _, supportingDocumenFileHeader := range req.meetingMinutestFileHeaderList {
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
			Type:                 store.MEETING_MINUTES,
		})
		supportDocumentProducerItemList = append(supportDocumentProducerItemList, SupportDocumentProducerItem{
			Name:         supportingDocument,
			OriginalName: originalDocumentName,
			Type:         string(store.MEETING_MINUTES),
		})
		documentInfoList = append(documentInfoList, DocumentInfo{
			Name: originalDocumentName,
			Type: string(store.MEETING_MINUTES),
			Url:  supportingDocumentUrl,
		})

	}

	indicatorAssessmentInfoList := make([]store.IndicatorAssessmentInfo, len(req.indicatorNumbers))
	for idx, indicatorNumber := range req.indicatorNumbers {
		indicatorAssessmentInfoList[idx] = store.IndicatorAssessmentInfo{
			IndicatorNumber: indicatorNumber,
		}
	}
	assessmentUploadDetail := store.AssessmentUploadDetail{
		AssessmentDetail: store.AssessmentDetail{
			InstitutionName: req.institutionName,
		},
		IndicatorAssessmentInfoList: indicatorAssessmentInfoList,
		SupportDataDocumentInfoList: supportDataDocumentInfoList,
		UserId:                      userCred.ID,
	}

	err := handler.assessmentStore.InsertUploadDocument(ctx, &assessmentUploadDetail)
	if err != nil {
		log.Println(err)
		response.Error(w, apierror.InternalServerError())
		return
	}

	indicatorAssessmentProducerItemList := make([]IndicatorAssessmentProducerItem, len(assessmentUploadDetail.IndicatorAssessmentInfoList))
	for idx, indicatorAssessmentInfo := range assessmentUploadDetail.IndicatorAssessmentInfoList {
		indicatorData, err := handler.indicatoreAssessmentStore.FindIndicatorDetailByIndicatorNumber(ctx, indicatorAssessmentInfo.IndicatorNumber)
		if err != nil {
			log.Println(err)
			response.Error(w, apierror.InternalServerError())
			return
		}
		indicatorAssessmentProducerItem := IndicatorAssessmentProducerItem{
			Id:     indicatorAssessmentInfo.Id,
			Number: indicatorAssessmentInfo.IndicatorNumber,
			Detail: indicatorData.Detail,
		}
		indicatorAssessmentProducerItemList[idx] = indicatorAssessmentProducerItem
	}

	topic := "SPBE_Assessment"
	msg := UploadProducerMessage{
		Name:                    assessmentUploadDetail.AssessmentDetail.InstitutionName,
		Content:                 assessmentUploadDetail.SupportDataDocumentInfoList[0].Id, // WIP need to be updated later
		UserId:                  userCred.ID,
		RecipientNumber:         req.phoneNumberStr,
		AssessmentId:            assessmentUploadDetail.AssessmentDetail.Id,
		InstitutionName:         req.institutionName,
		IndicatorAssessmentList: indicatorAssessmentProducerItemList,
		SupportDocumentList:     supportDocumentProducerItemList,
		Timestamp:               time.Now().UTC().String(),
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
	// }

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

func generateEmailContent() (subject, message []byte) {
	subject = []byte("Otomatisasi Penilaian SPBE")
	message = []byte(fmt.Sprintf("Terima kasih telah menggunakan Aplikasi Otomatisasi Penilaian SPBE. Hasil penilaian anda akan keluar dalam beberapa saat lagi."))

	return
}

func generateWhatsAppMessage() *waProto.Message {
	return &waProto.Message{
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
}
