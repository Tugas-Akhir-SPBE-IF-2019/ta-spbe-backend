package assessment

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	apierror "github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/rest/error"
	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/rest/response"
	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/store"

	waProto "go.mau.fi/whatsmeow/binary/proto"
	"google.golang.org/protobuf/proto"
)

type ProofItem struct {
	Text                string   `json:"text"`
	PictureUrlList      []string `json:"picture_url_list"`
	DocumentPageUrlList []string `json:"document_page_url_list"`
}

type IndicatorAssessmentDetail struct {
	ID     string `json:"id"`
	Number int    `json:"number"`
	Detail string `json:"detail"`
}

type DocumentProofDetail struct {
	Name                    string   `json:"name"`
	OriginalName            string   `json:"original_name"`
	Type                    string   `json:"type"`
	Text                    string   `json:"text"`
	Title                   string   `json:"title"`
	PictureFileList         []string `json:"picture_file_list"`
	SpecificPageDocumentURL []string `json:"specific_page_document_url"`
	DocumentPageList        []int    `json:"document_page_list"`
}

type ResultDetail struct {
	Level       int    `json:"level"`
	Explanation string `json:"explanation"`
}

type IndicatorAssessmentResultCallbackData struct {
	IndicatorAsssessment IndicatorAssessmentDetail `json:"indicator_assessment"`
	DocumentProof        []DocumentProofDetail     `json:"document_proof"`
	Result               ResultDetail              `json:"result"`
}

type IndicatorAssessmentResultCallbackRequest struct {
	UserId          string                                  `json:"user_id"`
	AssessmentId    string                                  `json:"assessment_id"`
	RecipientNumber string                                  `json:"recipient_number"`
	Data            []IndicatorAssessmentResultCallbackData `json:"data"`
}

type ValidateIndicatorAssessmentResultResponseS struct {
	Message string `json:"message"`
}

func (handler *assessmentHandler) ResultCallback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := IndicatorAssessmentResultCallbackRequest{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {

		response.Error(w, apierror.BadRequestError(err.Error()))
		return
	}

	var indicatorAssessmentUpdateResultList []store.IndicatorAssessmentUpdateResultDetail
	for _, data := range req.Data {
		var documentProofList []store.DocumentProofAssessmentUpdateResultDetail
		for _, documentProof := range data.DocumentProof {
			var specificPageDocumentURLList []string
			var pictureProofURLList []string
			for idx, specificPageDocument := range documentProof.SpecificPageDocumentURL {
				specificPageDocumentURLList = append(specificPageDocumentURLList, fmt.Sprintf("http://%s/static/%s", handler.apiCfg.Host, specificPageDocument))
				pictureProofURLList = append(pictureProofURLList, fmt.Sprintf("http://%s/static/%s", handler.apiCfg.Host, documentProof.PictureFileList[idx]))
			}
			proof := store.DocumentProofAssessmentUpdateResultDetail{
				Name:                    documentProof.Name,
				OriginalName:            documentProof.OriginalName,
				Type:                    documentProof.Type,
				Text:                    documentProof.Text,
				Title:                   documentProof.Title,
				PictureFileList:         pictureProofURLList,
				SpecificPageDocumentURL: specificPageDocumentURLList,
				DocumentPageList:        documentProof.DocumentPageList,
			}
			documentProofList = append(documentProofList, proof)
		}
		result := store.IndicatorAssessmentUpdateResultDetail{
			ID:            data.IndicatorAsssessment.ID,
			Number:        data.IndicatorAsssessment.Number,
			Detail:        data.IndicatorAsssessment.Detail,
			DocumentProof: documentProofList,
			Result: store.IndicatorResultAssessmentUpdateResultDetail{
				Level:       data.Result.Level,
				Explanation: data.Result.Explanation,
			},
		}
		indicatorAssessmentUpdateResultList = append(indicatorAssessmentUpdateResultList, result)
	}

	resultDetail := store.AssessmenUpdateResultDetail{
		IndicatorAssessmentList: indicatorAssessmentUpdateResultList,
	}
	err := handler.assessmentStore.UpdateAssessmentResult(ctx, &resultDetail)
	if err != nil {
		log.Println(err.Error())
		response.Error(w, apierror.InternalServerError())
		return
	}

	err = handler.assessmentStore.UpdateStatus(ctx, req.AssessmentId, store.AssessmentStatus(store.COMPLETED))
	if err != nil {
		log.Println(err.Error())
		response.Error(w, apierror.InternalServerError())
		return
	}

	user, err := handler.userStore.FindOneByID(ctx, req.UserId)
	if err != nil {
		log.Println(err)
		response.Error(w, apierror.InternalServerError())
		return
	}

	to := []string{user.Email}
	// TODO change email to not include level and explanation
	// TODO current level and explanation is hardcoded
	hardcoded_level := 4
	hardcoded_explanation := "still in development"
	subject, message := generateResultEmailContent(hardcoded_level, hardcoded_explanation)
	go func() {
		err := handler.smtpMailer.Send(subject, message, to, "result.html", map[string]string{"level": strconv.Itoa(hardcoded_level),
			"explanation": strings.ToUpper(hardcoded_explanation)})
		if err != nil {
			log.Println("error send email: %w", err)
		}
	}()

	protoMessage := generateResultWhatsAppMessage()

	err = handler.waClient.SendMessage(ctx, req.RecipientNumber, protoMessage)
	if err != nil {
		log.Println("error send whatsapp message: %w", err)
	}

	response.Respond(w, http.StatusNoContent, nil)
}

func generateResultEmailContent(level int, explanation string) (subject, message []byte) {
	subject = []byte("Hasil Otomatisasi Penilaian SPBE")
	message = []byte(fmt.Sprintf(`Terima kasih telah menggunakan Aplikasi Otomatisasi Penilaian SPBE. Berikut ini hasil penilaian anda
	Level: %d
	Penjelasan: %s`, level, explanation))

	return
}

func generateResultWhatsAppMessage() *waProto.Message {
	resultMessage := `*[OTOMATISASI PENILAIAN SPBE]*
	` + "```" + `Terima kasih telah menggunakan Aplikasi Otomatisasi Penilaian SPBE. Berikut ini hasil penilaian anda:` + "```" + `
	*- LEVEL*: ` + "```" + "4```" + `
	*- PENJELASAN*: ` + "```" + `Verifikasi dan validasi telah dilakukan terhadap penjelasan dan data dukung pada Indikator 10 Tingkat Kematangan Kebijakan Internal Tim Koordinasi SPBE pada Kementerian PANRB, dimana tercantum dalam PermenPANRB No xx tahun 2020, yaitu pada Pasal 11 halaman 9 tentang tugas dan fungsi Tim Koordinasi SPBE di lingkungan Kementerian PANRB sesuai data dukung 10.PermenPANRB-xx-2020.pdf` + "```" + `*(Fakta)*.
	` + "```" + `Berdasarkan penjelasan dan data dukung yang disampaikan, maka pengaturan tersebut telah memenuhi kekuatan hukum kebijakan mengikat secara internal, dan telah mencakup tugas dan fungsi Tim Koordinasi SPBE secara menyeluruh di lingkungan Kementerian PANRB, namun belum terdapat pengaturan arah koordinasi ataupun kolaborasi/kerja sama dengan Instansi lain di luar Kementerian PANRB` + "```" + `*(Analisis)*.
	` + "```" + `Hasil penilaian terhadap penjelasan dan data dukung menggambarkan tingkat kematangan 4 (empat).` + "```" + `*(Justifikasi Hasil)*.`
	return &waProto.Message{
		Conversation: proto.String(resultMessage),
	}
}
