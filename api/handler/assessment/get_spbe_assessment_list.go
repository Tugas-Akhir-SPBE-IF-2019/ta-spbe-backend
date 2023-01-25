package assessment

import (
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"ta-spbe-backend/api/response"
	"ta-spbe-backend/repository"
	"time"

	apierror "ta-spbe-backend/api/error"
)

type GetAssessmentListRequest struct {
	pageStr  string
	limitStr string
	page     int
	limit    int
}

func (r *GetAssessmentListRequest) validate() *apierror.FieldError {
	var err error
	fieldErr := apierror.NewFieldError()

	r.pageStr = strings.TrimSpace(r.pageStr)
	r.limitStr = strings.TrimSpace(r.limitStr)

	if r.pageStr == "" {
		r.pageStr = "1"
	}
	if r.limitStr == "" {
		r.limitStr = "50"
	}

	r.page, err = strconv.Atoi(r.pageStr)
	if err != nil || r.page < 0 {
		fieldErr = fieldErr.WithField("page", "page must be a positive integer")
	}

	r.limit, err = strconv.Atoi(r.limitStr)
	if err != nil || r.limit < 0 {
		fieldErr = fieldErr.WithField("limit", "limit must be a positive integer")
	}

	if len(fieldErr.Fields) != 0 {
		return &fieldErr
	}

	return nil
}

type AssessmentListResponse struct {
	TotalItems int              `json:"total_items"`
	TotalPages int              `json:"total_pages"`
	Items      []AssessmentItem `json:"items"`
}

type AssessmentItem struct {
	Id              string    `json:"id"`
	InstitutionName string    `json:"institution_name"`
	Status          int       `json:"status"`
	SubmittedDate   time.Time `json:"submitted_date"`
}

func GetSPBEAssessmentList(assessmentRepo repository.AssessmentRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		req := GetAssessmentListRequest{
			pageStr:  r.URL.Query().Get("page"),
			limitStr: r.URL.Query().Get("limit"),
		}

		fieldErr := req.validate()
		if fieldErr != nil {
			response.FieldError(w, *fieldErr)
			return
		}

		assessmentListAll, err := assessmentRepo.FindAll(ctx)
		if err != nil {
			log.Println(err)

			response.Error(w, apierror.InternalServerError())
			return
		}
		totalItems := len(assessmentListAll)

		offset := req.limit * (req.page - 1)
		assessmentList, err := assessmentRepo.FindAllPagination(ctx, offset, req.limit)
		if err != nil {
			log.Println(err)

			response.Error(w, apierror.InternalServerError())
			return
		}
		itemsCount := len(assessmentList)

		resp := AssessmentListResponse{}
		items := make([]AssessmentItem, itemsCount)
		for idx, item := range assessmentList {
			items[idx] = AssessmentItem{
				Id:              item.Id,
				InstitutionName: item.InstitutionName,
				Status:          item.Status,
				SubmittedDate:   item.SubmittedDate,
			}
		}
		resp.TotalItems = totalItems
		resp.TotalPages = int(math.Ceil(float64(totalItems) / float64(req.limit)))
		resp.Items = items

		response.Respond(w, http.StatusOK, resp)
	}
}
