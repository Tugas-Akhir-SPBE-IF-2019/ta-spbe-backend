package institution

import (
	"net/http"

	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/store"
)

type InstitutionHandler interface {
	GetInstitutionList(w http.ResponseWriter, r *http.Request)
}

type institutionHandler struct {
	institutionStore store.Institution
}

func NewInstitutionHandler(institutionStore store.Institution) InstitutionHandler {
	return &institutionHandler{
		institutionStore: institutionStore,
	}
}
