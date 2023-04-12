package pgsql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/store"
	"github.com/google/uuid"
)

type IndicatorAssessment struct {
	db *sql.DB
}

func NewIndicatorAssessment(db *sql.DB) *IndicatorAssessment {
	return &IndicatorAssessment{db: db}
}

const indicatorAssessmentFindAllQuery = `SELECT a.institution_name, AVG(Cast(ia.level AS FLOAT)), a.created_at 
		FROM assessments a
		LEFT JOIN indicator_assessments ia 
		ON ia.assessment_id = a.id
		WHERE ia.level IS NOT NULL `

func (s *IndicatorAssessment) FindAll(ctx context.Context, queryInstitution string, startDate string, endDate string, indexMin float64, indexMax float64) ([]*store.IndicatorAssessmentDetail, error) {
	indicatorAssessmentList := []*store.IndicatorAssessmentDetail{}

	var queryKeys []string
	var queryParams []interface{}

	queryParams = append(queryParams, indexMin, indexMax)

	query := indicatorAssessmentFindAllQuery
	if queryInstitution != "" {
		queryKeys = append(queryKeys, "Institution")
		queryParams = append(queryParams, queryInstitution)
	}

	if startDate != "" {
		queryKeys = append(queryKeys, "StartDate")
		queryParams = append(queryParams, startDate)
	}

	if endDate != "" {
		queryKeys = append(queryKeys, "EndDate")
		queryParams = append(queryParams, endDate)
	}

	for index, key := range queryKeys {
		query = query + "AND "

		// +3 is used because $1 and $2 are already used for index
		switch key {
		case "Institution":
			query = query + fmt.Sprintf(`a.institution_name ILIKE '%%' || $%d || '%%' `, index+3)
		case "StartDate":
			query = query + fmt.Sprintf(`a.created_at >= $%d `, index+3)
		case "EndDate":
			query = query + fmt.Sprintf(`a.created_at <= $%d `, index+3)
		}
	}

	query = query + `GROUP BY a.institution_name, a.created_at HAVING AVG(ia.level) >= $1 and AVG(ia.level) <= $2 ORDER BY a.created_at desc `

	rows, err := s.db.QueryContext(ctx, query, queryParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		indicatorAssessment := &store.IndicatorAssessmentDetail{}
		err := rows.Scan(
			&indicatorAssessment.InstitutionName,
			&indicatorAssessment.SpbeIndex,
			&indicatorAssessment.SubmittedDate,
		)
		if err != nil {
			return nil, err
		}
		indicatorAssessmentList = append(indicatorAssessmentList, indicatorAssessment)
	}

	return indicatorAssessmentList, nil
}

func (s *IndicatorAssessment) FindAllPagination(ctx context.Context, offset int, limit int, queryInstitution string, startDate string, endDate string, indexMin float64, indexMax float64) ([]*store.IndicatorAssessmentDetail, error) {
	indicatorAssessmentList := []*store.IndicatorAssessmentDetail{}

	var queryKeys []string
	var queryParams []interface{}

	queryParams = append(queryParams, offset, limit, indexMin, indexMax)

	query := indicatorAssessmentFindAllQuery
	if queryInstitution != "" {
		queryKeys = append(queryKeys, "Institution")
		queryParams = append(queryParams, queryInstitution)
	}

	if startDate != "" {
		queryKeys = append(queryKeys, "StartDate")
		queryParams = append(queryParams, startDate)
	}

	if endDate != "" {
		queryKeys = append(queryKeys, "EndDate")
		queryParams = append(queryParams, endDate)
	}

	for index, key := range queryKeys {
		query = query + "AND "

		// +5 is used because $1 and $2 are already used for pagination and $3 and $4 are used for index
		switch key {
		case "Institution":
			query = query + fmt.Sprintf(`a.institution_name ILIKE '%%' || $%d || '%%' `, index+5)
		case "StartDate":
			query = query + fmt.Sprintf(`a.created_at >= $%d `, index+5)
		case "EndDate":
			query = query + fmt.Sprintf(`a.created_at <= $%d `, index+5)
		}
	}

	query = query + `GROUP BY a.institution_name, a.created_at HAVING AVG(ia.level) >= $3 and AVG(ia.level) <= $4 ORDER BY a.created_at DESC LIMIT $2 OFFSET $1 `

	rows, err := s.db.QueryContext(ctx, query, queryParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		indicatorAssessment := &store.IndicatorAssessmentDetail{}
		err := rows.Scan(
			&indicatorAssessment.InstitutionName,
			&indicatorAssessment.SpbeIndex,
			&indicatorAssessment.SubmittedDate,
		)
		if err != nil {
			return nil, err
		}
		indicatorAssessmentList = append(indicatorAssessmentList, indicatorAssessment)
	}

	return indicatorAssessmentList, nil
}

const indicatorAssessmentResultFindByIdQuery = `SELECT ia.id, a.institution_name, ia.created_at, ia.status, i.domain, i.aspect, i.indicator_number, 
		ia.level, sdd.document_url, sdd.document_original_name, COALESCE(ia.explanation, ''),  COALESCE(sddp.proof, '')
		FROM indicator_assessments ia
		LEFT JOIN assessments a
		ON ia.assessment_id = a.id
		LEFT JOIN indicators i
		ON ia.indicator_id = i.id
		LEFT JOIN support_data_documents sdd
		ON ia.id = sdd.indicator_assessment_id
		LEFT JOIN support_data_document_proofs sddp
		ON ia.id = sddp.indicator_assessment_id
		WHERE ia.id = $1`

func (s *IndicatorAssessment) FindIndicatorAssessmentResultById(ctx context.Context, id string) (store.IndicatorAssessmentResultDetail, error) {
	assessmentResult := store.IndicatorAssessmentResultDetail{}

	row := s.db.QueryRowContext(ctx, indicatorAssessmentResultFindByIdQuery, id)
	err := row.Scan(&assessmentResult.IndicatorAssessmentId, &assessmentResult.InstitutionName, &assessmentResult.SubmittedDate,
		&assessmentResult.AssessmentStatus, &assessmentResult.Result.Domain, &assessmentResult.Result.Aspect,
		&assessmentResult.Result.IndicatorNumber, &assessmentResult.Result.Level, &assessmentResult.Result.SupportDocument,
		&assessmentResult.Result.SupportDocumentName, &assessmentResult.Result.Explanation, &assessmentResult.Result.Proof,
	)
	if err != nil {
		return assessmentResult, err
	}

	return assessmentResult, nil
}

const indicatorAssessmentResultFindByAssessmentIdQuery = `SELECT ia.id, a.institution_name, ia.created_at, ia.status, i.domain, i.aspect, i.indicator_number, 
		ia.level, sdd.document_url, sdd.document_original_name, COALESCE(ia.explanation, ''),  COALESCE(sddp.proof, '')
		FROM indicator_assessments ia
		LEFT JOIN assessments a
		ON ia.assessment_id = a.id
		LEFT JOIN indicators i
		ON ia.indicator_id = i.id
		LEFT JOIN support_data_documents sdd
		ON ia.id = sdd.indicator_assessment_id
		LEFT JOIN support_data_document_proofs sddp
		ON ia.id = sddp.indicator_assessment_id
		WHERE a.id = $1
		ORDER BY i.indicator_number ASC `

func (s *IndicatorAssessment) FindIndicatorAssessmentResultByAssessmentId(ctx context.Context, id string) ([]*store.IndicatorAssessmentResultDetail, error) {
	assessmentResultList := []*store.IndicatorAssessmentResultDetail{}

	rows, err := s.db.QueryContext(ctx, indicatorAssessmentResultFindByAssessmentIdQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		assessmentResult := &store.IndicatorAssessmentResultDetail{}
		err := rows.Scan(&assessmentResult.IndicatorAssessmentId, &assessmentResult.InstitutionName, &assessmentResult.SubmittedDate,
			&assessmentResult.AssessmentStatus, &assessmentResult.Result.Domain, &assessmentResult.Result.Aspect,
			&assessmentResult.Result.IndicatorNumber, &assessmentResult.Result.Level, &assessmentResult.Result.SupportDocument,
			&assessmentResult.Result.SupportDocumentName, &assessmentResult.Result.Explanation, &assessmentResult.Result.Proof,
		)
		if err != nil {
			return nil, err
		}
		assessmentResultList = append(assessmentResultList, assessmentResult)
	}

	return assessmentResultList, nil
}

const insertAssessmentFeedbackQuery = `INSERT INTO 
	indicator_assessment_feedbacks(
		id, indicator_assessment_id, level,
		feedback, created_at
	) values (
		$1, $2, $3, $4, $5
	)
`
const updateAssessmentFeedbackQuery = `UPDATE indicator_assessment_feedbacks	
	SET level = $2, feedback = $3
	WHERE indicator_assessment_id = $1
`
const updateValidatedAssessmentFeedbackStatusQuery = `UPDATE indicator_assessments
	SET status = $2
	WHERE id = $1
`

func (s *IndicatorAssessment) ValidateAssessmentResult(ctx context.Context, resultCorrect bool, indicatorAssessmentResult *store.IndicatorAssessmentResultDetail) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin validate assessment result tx: %w", err)
	}
	defer tx.Rollback()

	if store.AssessmentStatus(indicatorAssessmentResult.AssessmentStatus) != store.AssessmentStatus(store.VALIDATED) {
		feedbackId := uuid.NewString()
		createdAt := time.Now().UTC()

		insertAssessmentFeedbackStmt, err := s.db.PrepareContext(ctx, insertAssessmentFeedbackQuery)
		_, err = tx.StmtContext(ctx, insertAssessmentFeedbackStmt).ExecContext(ctx, feedbackId, indicatorAssessmentResult.IndicatorAssessmentId,
			indicatorAssessmentResult.ResultFeedback.Level, indicatorAssessmentResult.ResultFeedback.Feedback, createdAt)
		if err != nil {
			return fmt.Errorf("failed to insert indicator assessment feedback: %w", err)
		}

		if !resultCorrect {
			updateValidatedAssessmentFeedbackStmt, err := s.db.PrepareContext(ctx, updateValidatedAssessmentFeedbackStatusQuery)
			_, err = tx.StmtContext(ctx, updateValidatedAssessmentFeedbackStmt).ExecContext(ctx, indicatorAssessmentResult.IndicatorAssessmentId, store.AssessmentStatus(store.VALIDATED))
			if err != nil {
				return fmt.Errorf("failed to update indicator assessment validated status: %w", err)
			}
		}
	} else {
		if !resultCorrect {
			updateAssessmentFeedbackStmt, err := s.db.PrepareContext(ctx, updateAssessmentFeedbackQuery)
			_, err = tx.StmtContext(ctx, updateAssessmentFeedbackStmt).ExecContext(ctx, indicatorAssessmentResult.IndicatorAssessmentId,
				indicatorAssessmentResult.ResultFeedback.Level, indicatorAssessmentResult.ResultFeedback.Feedback)
			if err != nil {
				return fmt.Errorf("failed to update indicator assessment feedback: %w", err)
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit validate assessment result tx: %w", err)
	}

	return nil
}

const updateIndicatorAssessmentResultQuery = `UPDATE indicator_assessments
	SET status = $2, level = $3, explanation = $4, updated_at = $5
	WHERE id = $1
`
const insertSupportDataDocumentProofQuery = `INSERT into
	support_data_document_proofs(
		id, indicator_assessment_id, support_data_document_id, proof, created_at
	) values(
		$1, $2, $3, $4, $5
	)
`

func (s *IndicatorAssessment) UpdateAssessmentResult(ctx context.Context, resultDetail *store.IndicatorAssessmentResultDetail) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to commit update assessment result tx: %w", err)
	}
	defer tx.Rollback()

	updatedAt := time.Now().UTC()
	// updateAssessmentStatusStmt, err := s.db.PrepareContext(ctx, updateAssessmentStatusQuery)
	// _, err = tx.StmtContext(ctx, updateAssessmentStatusStmt).ExecContext(ctx, resultDetail.AssessmentId,
	// 	store.AssessmentStatus(store.COMPLETED), updatedAt)
	// if err != nil {
	// 	return fmt.Errorf("failed to update assessment: %w", err)
	// }

	updateIndicatorAssessmentResultStmt, err := s.db.PrepareContext(ctx, updateIndicatorAssessmentResultQuery)
	_, err = tx.StmtContext(ctx, updateIndicatorAssessmentResultStmt).ExecContext(ctx,
		resultDetail.IndicatorAssessmentId, store.AssessmentStatus(store.COMPLETED),
		resultDetail.Result.Level, resultDetail.Result.Explanation, updatedAt)
	if err != nil {
		return fmt.Errorf("failed to update indicator assessment: %w", err)
	}

	supportDataDocumentProofId := uuid.NewString()
	insertSupportDataDocumentProofStmt, err := s.db.PrepareContext(ctx, insertSupportDataDocumentProofQuery)
	_, err = tx.StmtContext(ctx, insertSupportDataDocumentProofStmt).ExecContext(ctx,
		supportDataDocumentProofId, resultDetail.IndicatorAssessmentId, resultDetail.Result.SupportDocument,
		resultDetail.Result.Proof, updatedAt)
	if err != nil {
		return fmt.Errorf("failed to insert support data document proof: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit update assessment result tx: %w", err)
	}

	return nil
}

const indicatorAssessmentProofInsert = `INSERT INTO
indicator_assessment_proofs(
	id, indicator_assessment_id, image_url, document_url, created_at
) values(
	$1, $2, $3, $4, $5	
)
`

func (s *IndicatorAssessment) InsertProofData(ctx context.Context, proofData *store.IndicatorAssessmentProofData) error {
	insertStmt, err := s.db.PrepareContext(ctx, indicatorAssessmentProofInsert)
	if err != nil {
		return err
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}
	defer tx.Rollback()

	createdAt := time.Now().UTC()
	_, err = tx.StmtContext(ctx, insertStmt).ExecContext(ctx,
		proofData.ID, proofData.IndicatorAssessmentID, proofData.ImageURL, proofData.DocumentURL, createdAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}
	proofData.CreatedAt = createdAt

	return nil

}

const proofDataFindAllByIndicatorAssessmentIdQuery = `SELECT iap.id, iap.indicator_assessment_id, iap.image_url, iap.document_url, iap.created_at
	FROM indicator_assessment_proofs iap
	WHERE iap.indicator_assessment_id = $1
`

func (s *IndicatorAssessment) FindProofDataByIndicatorAssessmentId(ctx context.Context, id string) ([]*store.IndicatorAssessmentProofData, error) {
	proofResultList := []*store.IndicatorAssessmentProofData{}

	rows, err := s.db.QueryContext(ctx, proofDataFindAllByIndicatorAssessmentIdQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		proofResult := &store.IndicatorAssessmentProofData{}
		err := rows.Scan(
			&proofResult.ID, &proofResult.IndicatorAssessmentID,
			&proofResult.ImageURL, &proofResult.DocumentURL, &proofResult.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		proofResultList = append(proofResultList, proofResult)
	}

	return proofResultList, nil
}
