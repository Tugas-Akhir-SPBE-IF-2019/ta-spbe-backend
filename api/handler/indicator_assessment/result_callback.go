package indicatorassessment

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	apierror "ta-spbe-backend/api/error"
	"ta-spbe-backend/api/response"
	"ta-spbe-backend/repository"
	"ta-spbe-backend/service"
)

type IndicatorAssessmentResultCallbackRequest struct {
	UserId                string `json:"user_id"`
	AssessmentId          string `json:"assessment_id"`
	IndicatorAssessmentId string `json:"indicator_assessment_id"`
	Level                 int    `json:"level"`
	Explanation           string `json:"explanation"`
	SupportDataDocumentId string `json:"support_data_document_id"`
	Proof                 string `json:"proof"`
}

type ValidateIndicatorAssessmentResultResponseS struct {
	Message string `json:"message"`
}

func ResultCallback(indicatorAssessmentRepo repository.IndicatorAssessmentRepository, userRepo repository.UserRepository, mailer service.Mailer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		req := IndicatorAssessmentResultCallbackRequest{}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Error(w, apierror.BadRequestError(err.Error()))
			return
		}

		result := repository.IndicatorAssessmentResultDetail{
			AssessmentId:          req.AssessmentId,
			IndicatorAssessmentId: req.IndicatorAssessmentId,
			Result: repository.IndicatorAssessmentResultInfo{
				Level:           req.Level,
				Explanation:     req.Explanation,
				SupportDocument: req.SupportDataDocumentId,
				Proof:           req.Proof,
			},
		}

		err := indicatorAssessmentRepo.UpdateAssessmentResult(ctx, &result)
		if err != nil {
			log.Println(err.Error())
			response.Error(w, apierror.InternalServerError())
			return
		}

		user, err := userRepo.FindOneByID(ctx, req.UserId)
		if err != nil {
			log.Println(err)
			response.Error(w, apierror.InternalServerError())
			return
		}

		to := []string{user.Email}
		subject, message := generateEmailContent(req.Level, req.Explanation)
		go func() {
			err := mailer.Send(subject, message, to, "result.html", map[string]string{"level": strconv.Itoa(req.Level),
				"explanation": strings.ToUpper(req.Explanation)})
			if err != nil {
				log.Println("error send email: %w", err)
			}
		}()

		response.Respond(w, http.StatusNoContent, nil)
	}
}

func generateEmailContent(level int, explanation string) (subject, message []byte) {
	subject = []byte("Hasil Otomatisasi Penilaian SPBE")
	message = []byte(fmt.Sprintf(`Terima kasih telah menggunakan Aplikasi Otomatisasi Penilaian SPBE. Berikut ini hasil penilaian anda
	Level: %d
	Penjelasan: %s`, level, explanation))

	return
}
