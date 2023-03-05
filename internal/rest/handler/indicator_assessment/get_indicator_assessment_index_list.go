package indicatorassessment

import (
	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/rest/response"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	apierror "github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/rest/error"
)

type GetIndicatorAssessmentIndexListRequest struct {
	pageStr  string
	limitStr string
	page     int
	limit    int
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
	SpbeIndex       int       `json:"spbe_index"`
	SubmittedDate   time.Time `json:"submitted_date"`
}

func (handler *indicatorAssessmentHandler) GetIndicatorAssessmentIndexList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req := GetIndicatorAssessmentIndexListRequest{
		pageStr:  r.URL.Query().Get("page"),
		limitStr: r.URL.Query().Get("limit"),
	}

	fieldErr := req.validate()
	if fieldErr != nil {
		response.FieldError(w, *fieldErr)
		return
	}

	indicatorAssessmentIndexListAll, err := handler.indicatorAssessmentStore.FindAll(ctx)
	if err != nil {
		log.Println(err)
		response.Error(w, apierror.InternalServerError())
		return
	}
	totalItems := len(indicatorAssessmentIndexListAll)

	offset := req.limit * (req.page - 1)
	indicatorAssessmentIndexList, err := handler.indicatorAssessmentStore.FindAllPagination(ctx, offset, req.limit)
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
