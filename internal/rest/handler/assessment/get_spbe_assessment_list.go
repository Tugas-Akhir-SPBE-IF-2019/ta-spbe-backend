package assessment

import (
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/rest/response"

	apierror "github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/rest/error"
)

type GetAssessmentListRequest struct {
	institution  string
	startDateStr string
	startDate    time.Time
	endDateStr   string
	endDate      time.Time
	status       int
	statusStr    string
	pageStr      string
	limitStr     string
	page         int
	limit        int
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

	if r.statusStr != "" {
		r.status, err = strconv.Atoi(r.statusStr)
		if err != nil || (r.status < 0 && r.status > 3) {
			fieldErr = fieldErr.WithField("status", "status must be a positive integer between 0 and 3 inclusive")
		}
	} else {
		r.status = -1 // value for no filter
	}

	if r.startDateStr != "" {
		r.startDate, err = time.Parse(time.DateOnly, r.startDateStr)
		if err != nil {
			fieldErr = fieldErr.WithField("start_date", "start_date must be in the format of YYYY-MM-DD!")
		}
	}

	if r.endDateStr != "" {
		r.endDate, err = time.Parse(time.DateOnly, r.endDateStr)
		if err != nil {
			fieldErr = fieldErr.WithField("end_date", "end_date must be in the format of YYYY-MM-DD!")
		}
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

func (handler *assessmentHandler) GetSPBEAssessmentList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req := GetAssessmentListRequest{
		pageStr:      r.URL.Query().Get("page"),
		limitStr:     r.URL.Query().Get("limit"),
		institution:  r.URL.Query().Get("institution"),
		statusStr:    r.URL.Query().Get("status"),
		startDateStr: r.URL.Query().Get("start_date"),
		endDateStr:   r.URL.Query().Get("end_date"),
	}

	fieldErr := req.validate()
	if fieldErr != nil {
		response.FieldError(w, *fieldErr)
		return
	}

	assessmentListAll, err := handler.assessmentStore.FindAll(ctx, req.institution, req.status, req.startDateStr, req.endDateStr)
	if err != nil {
		log.Println(err)

		response.Error(w, apierror.InternalServerError())
		return
	}
	totalItems := len(assessmentListAll)

	offset := req.limit * (req.page - 1)
	assessmentList, err := handler.assessmentStore.FindAllPagination(ctx, offset, req.limit, req.institution, req.status, req.startDateStr, req.endDateStr)
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
