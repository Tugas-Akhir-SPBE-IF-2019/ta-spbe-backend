package responses

type AssessmentListResponse struct {
	TotalItems int          `json:"total_items"`
	TotalPages int          `json:"total_pages"`
	Items      []Assessment `json:"items"`
}

type AssessmentIndexListResponse struct {
	TotalItems int               `json:"total_items"`
	TotalPages int               `json:"total_pages"`
	Items      []AssessmentIndex `json:"items"`
}

type AssessmentResultResponse struct {
	InstitutionName   string           `json:"institution_name"`
	SubmittedDate     string           `json:"submitted_date"`
	AssesssmentStatus int              `json:"assessment_status"`
	Result            AssessmentResult `json:"result"`
	Validated         bool             `json:"validated"`
}

type AssessmentDocumentUploadResponse struct {
	Message string `json:"message"`
	AssessmentId string `json:"assessment_id"`
}

type Assessment struct {
	Id              string `json:"id"`
	InstitutionName string `json:"institution_name"`
	Status          int    `json:"status"`
	SubmittedDate   string `json:"submitted_date"`
}

type AssessmentIndex struct {
	InstitutionName string  `json:"institution_name"`
	SpbeIndex       float64 `json:"spbe_index"`
	SubmittedDate   string  `json:"submitted_date"`
}

type AssessmentResult struct {
	Domain             string `json:"domain"`
	Aspect             string `json:"aspect"`
	IndicatorNumber    int    `json:"indicator_number"`
	Level              int    `json:"level"`
	Explanation        string `json:"explanation"`
	SupportingDocument string `json:"supporting_document"`
	OldDocument        string `json:"old_document"`
	Proof              string `json:"proof"`
}
