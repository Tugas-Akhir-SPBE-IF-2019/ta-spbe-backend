package handlers

import (
	"encoding/json"
	"net/http"
	"ta-spbe-backend/api/responses"
	"ta-spbe-backend/services"

	"github.com/go-chi/chi/v5"
)

type AssessmentHandler struct {
	assessmentService services.AssessmentService
}

func (hanlder *AssessmentHandler) GetAssessmentList(w http.ResponseWriter, r *http.Request) {
	assessment := []responses.Assessment{
		{
			Id:              "940c6ac1-3e0a-4316-8526-43aaf8120cbf",
			InstitutionName: "Kabupaten Lamongan",
			Status:          1,
			SubmittedDate:   "2019-10-12T07:20:50.52Z",
		},
		{
			Id:              "810c6ac1-3e0a-4316-8526-43aaf8120cab",
			InstitutionName: "Kabupaten Bantul",
			Status:          2,
			SubmittedDate:   "2022-10-12T07:20:50.52Z",
		},
	}

	response := responses.AssessmentListResponse{
		TotalItems: 2,
		TotalPages: 1,
		Items:      assessment,
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func NewAssessmentHandler(assessmentService services.AssessmentService) AssessmentHandler {
	return AssessmentHandler{
		assessmentService: assessmentService,
	}
}

func (handler *AssessmentHandler) GetAssessmentIndexList(w http.ResponseWriter, r *http.Request) {
	assessmentIndex := []responses.AssessmentIndex{
		{
			InstitutionName: "Kabupaten Lamongan",
			SpbeIndex:       2.3,
			SubmittedDate:   "2019-10-12T07:20:50.52Z",
		},
		{
			InstitutionName: "Kabupaten Bantul",
			SpbeIndex:       3.9,
			SubmittedDate:   "2022-10-12T07:20:50.52Z",
		},
	}

	response := responses.AssessmentIndexListResponse{
		TotalItems: 2,
		TotalPages: 1,
		Items:      assessmentIndex,
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (handler *AssessmentHandler) GetAssessmentResult(w http.ResponseWriter, r *http.Request) {
	assessmentId := chi.URLParam(r, "id")
	if assessmentId != "940c6ac1-3e0a-4316-8526-43aaf8120cbf" {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(nil)
		return
	}
	assessmentResult := responses.AssessmentResult{
		Domain:          "Kebijakan Internal",
		Aspect:          "Kebijakan Tata Kelola SPBE",
		IndicatorNumber: 1,
		Level:           3,
		Explanation: `Verifikasi dan validasi telah dilakukan terhadap penjelasan dan data dukung pada Indikator 1 Tingkat Kematangan Kebijakan Internal Arsitektur SPBE Instansi Pusat/Pemerintah Daerah pada Kabupaten Lamongan, dimana tercantum dalam Peraturan Bupati Lamongan Nomor 27 Tahun 2021, yaitu pada Pasal 5, 6, dan 7 halaman 9-10 tentang Sistem Pemerintahan Berbasis Elektronik sesuai data dukung Pasal 5,6,7 Pada PERBUP Nomor 27 Tahun 2021 .pdf (Fakta). 

		Berdasarkan penjelasan dan data dukung yang disampaikan, maka pengaturan tersebut telah memenuhi kekuatan hukum kebijakan mengikat secara internal, dan telah mencakup secara lengkap pengaturan mengenai referensi Arsitektur dan domain Arsitektur SPBE (Proses Bisnis, Data dan Informasi, Infrastruktur SPBE, Aplikasi SPBE, Keamanan SPBE, dan Layanan SPBE)  di lingkungan Kabupaten Lamongan, namun belum terdapat pengaturan integrasi SPBE antar Instansi Pusat, antar Pemerintah Daerah, dan/atau antar Instansi Pusat dan Pemerintah Daerah di luar Kabupaten Lamongan (Analisis). 
		
		Hasil penilaian terhadap penjelasan dan data dukung menggambarkan tingkat kematangan 3 (tiga). (Justifikasi Hasil).`,
		SupportingDocument: "PERBUP Nomor 27 Tahun 2021 Revisi .pdf",
		OldDocument:        "PERBUP Nomor 27 Tahun 2021 .pdf",
		Proof:              "<p><b>Ini buktinya</b></p>",
	}

	response := responses.AssessmentResultResponse{
		InstitutionName:   "Kabupaten Lamongan",
		SubmittedDate:     "2019-10-12T07:20:50.52Z",
		AssesssmentStatus: 2,
		Result:            assessmentResult,
		Validated:         false,
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (handler *AssessmentHandler) UploadAssessmentDocument(w http.ResponseWriter, r *http.Request) {
	response := responses.AssessmentDocumentUploadResponse{
		Message:      "success",
		AssessmentId: "940c6ac1-3e0a-4316-8526-43aaf8120cbf",
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (handler *AssessmentHandler) ValidateAssessmentResult(w http.ResponseWriter, r *http.Request) {
	response := responses.BaseCreateResponse{
		Message: "Validation success",
	}

	assessmentId := chi.URLParam(r, "id")
	if assessmentId != "940c6ac1-3e0a-4316-8526-43aaf8120cbf" {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(nil)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
