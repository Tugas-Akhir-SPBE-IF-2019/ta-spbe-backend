package institution

import (
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/rest/response"

	apierror "github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/rest/error"
)

type GetInstitutionListRequest struct {
	name     string
	pageStr  string
	limitStr string
	page     int
	limit    int
}

func (r *GetInstitutionListRequest) validate() *apierror.FieldError {
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

type InstitutionItem struct {
	ID       int    `json:"id"`
	Name     string `json:"institution_name"`
	Category string `json:"institution_category"`
}

type InstitutionListResponse struct {
	TotalItems int               `json:"total_items"`
	TotalPages int               `json:"total_pages"`
	Items      []InstitutionItem `json:"items"`
}

func (handler *institutionHandler) GetInstitutionList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req := GetInstitutionListRequest{
		pageStr:  r.URL.Query().Get("page"),
		limitStr: r.URL.Query().Get("limit"),
		name:     r.URL.Query().Get("name"),
	}

	fieldErr := req.validate()
	if fieldErr != nil {
		response.FieldError(w, *fieldErr)
		return
	}

	institutionListAll, err := handler.institutionStore.FindAll(ctx, req.name)
	if err != nil {
		log.Println(err)

		response.Error(w, apierror.InternalServerError())
		return
	}
	totalItems := len(institutionListAll)

	offset := req.limit * (req.page - 1)
	institutionList, err := handler.institutionStore.FindAllPagination(ctx, offset, req.limit, req.name)
	if err != nil {
		log.Println(err)

		response.Error(w, apierror.InternalServerError())
		return
	}
	itemsCount := len(institutionList)

	resp := InstitutionListResponse{}
	items := make([]InstitutionItem, itemsCount)
	for idx, item := range institutionList {
		items[idx] = InstitutionItem{
			ID:       item.ID,
			Name:     item.Name,
			Category: item.Category,
		}
	}
	resp.TotalItems = totalItems
	resp.TotalPages = int(math.Ceil(float64(totalItems) / float64(req.limit)))
	resp.Items = items

	response.Respond(w, http.StatusOK, resp)
}
