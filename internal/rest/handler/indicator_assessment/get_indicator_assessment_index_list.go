package indicatorassessment

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

type GetIndicatorAssessmentIndexListRequest struct {
	institution  string
	startDateStr string
	startDate    time.Time
	endDateStr   string
	endDate      time.Time
	indexMinStr  string
	indexMin     float64
	indexMaxStr  string
	indexMax     float64
	pageStr      string
	limitStr     string
	page         int
	limit        int
}

func (r *GetIndicatorAssessmentIndexListRequest) validate() *apierror.FieldError {
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
		dayInt, _ := strconv.Atoi(r.endDateStr[8:10])
		dayInt++
		if dayInt < 10 {
			r.endDateStr = r.endDateStr[0:8] + "0" + strconv.Itoa(dayInt)
		} else {
			r.endDateStr = r.endDateStr[0:8] + strconv.Itoa(dayInt)
		}
	}

	if r.indexMinStr != "" {
		r.indexMin, err = strconv.ParseFloat(r.indexMinStr, 64)
		if err != nil {
			fieldErr = fieldErr.WithField("index_min", "index_min must be a number!")
		}
	} else {
		r.indexMin = -99999 // Hardcoded minimum value
	}

	if r.indexMaxStr != "" {
		r.indexMax, err = strconv.ParseFloat(r.indexMaxStr, 64)
		if err != nil {
			fieldErr = fieldErr.WithField("index_max", "index_max must be a number!")
		}
	} else {
		r.indexMax = 99999 // Hardcoded maximum value
	}

	if len(fieldErr.Fields) != 0 {
		return &fieldErr
	}

	return nil
}

type IndicatorAssessmentListResponse struct {
	TotalItems int                       `json:"total_items"`
	TotalPages int                       `json:"total_pages"`
	Items      []IndicatorAssessmentItem `json:"items"`
}

type IndicatorAssessmentItem struct {
	InstitutionName string    `json:"institution_name"`
	SpbeIndex       float64   `json:"spbe_index"`
	SubmittedDate   time.Time `json:"submitted_date"`
}

func (handler *indicatorAssessmentHandler) GetIndicatorAssessmentIndexList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req := GetIndicatorAssessmentIndexListRequest{
		pageStr:      r.URL.Query().Get("page"),
		limitStr:     r.URL.Query().Get("limit"),
		institution:  r.URL.Query().Get("institution"),
		startDateStr: r.URL.Query().Get("start_date"),
		endDateStr:   r.URL.Query().Get("end_date"),
		indexMinStr:  r.URL.Query().Get("index_min"),
		indexMaxStr:  r.URL.Query().Get("index_max"),
	}

	fieldErr := req.validate()
	if fieldErr != nil {
		response.FieldError(w, *fieldErr)
		return
	}

	indicatorAssessmentIndexListAll, err := handler.indicatorAssessmentStore.FindAll(ctx, req.institution, req.startDateStr, req.endDateStr, req.indexMin, req.indexMax)
	if err != nil {
		log.Println(err)
		response.Error(w, apierror.InternalServerError())
		return
	}
	totalItems := len(indicatorAssessmentIndexListAll)

	offset := req.limit * (req.page - 1)
	indicatorAssessmentIndexList, err := handler.indicatorAssessmentStore.FindAllPagination(ctx, offset, req.limit, req.institution, req.startDateStr, req.endDateStr, req.indexMin, req.indexMax)
	if err != nil {
		log.Println(err)
		response.Error(w, apierror.InternalServerError())
		return
	}
	itemsCount := len(indicatorAssessmentIndexList)

	resp := IndicatorAssessmentListResponse{}
	items := make([]IndicatorAssessmentItem, itemsCount)
	for idx, item := range indicatorAssessmentIndexList {
		items[idx] = IndicatorAssessmentItem{
			InstitutionName: item.InstitutionName,
			SpbeIndex:       item.SpbeIndex,
			SubmittedDate:   item.SubmittedDate,
		}
	}
	resp.TotalItems = totalItems
	resp.TotalPages = int(math.Ceil(float64(totalItems) / float64(req.limit)))
	resp.Items = items

	response.Respond(w, http.StatusOK, resp)
}
