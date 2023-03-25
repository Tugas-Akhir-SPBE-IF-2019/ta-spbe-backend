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

const indicatorAssessmentFindAllQuery = `SELECT a.institution_name, ia.level, ia.created_at
		FROM indicator_assessments ia
		LEFT JOIN assessments a
		ON ia.assessment_id = a.id
		WHERE ia.level IS NOT NULL`

func (s *IndicatorAssessment) FindAll(ctx context.Context) ([]*store.IndicatorAssessmentDetail, error) {
	indicatorAssessmentList := []*store.IndicatorAssessmentDetail{}

	rows, err := s.db.QueryContext(ctx, indicatorAssessmentFindAllQuery)
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

const indicatorAssessmentFindAllPaginationQuery = `SELECT a.institution_name, ia.level, ia.created_at
		FROM indicator_assessments ia
		LEFT JOIN assessments a
		ON ia.assessment_id = a.id
		WHERE ia.level IS NOT NULL
		LIMIT $2 OFFSET $1`

func (s *IndicatorAssessment) FindAllPagination(ctx context.Context, offset int, limit int) ([]*store.IndicatorAssessmentDetail, error) {
	indicatorAssessmentList := []*store.IndicatorAssessmentDetail{}

	rows, err := s.db.QueryContext(ctx, indicatorAssessmentFindAllPaginationQuery, offset, limit)
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
		ia.level, sdd.document_url, COALESCE(ia.explanation, ''),  COALESCE(sddp.proof, '')
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
		&assessmentResult.Result.Explanation, &assessmentResult.Result.Proof,
	)
	if err != nil {
		return assessmentResult, err
	}

	return assessmentResult, nil
}

const indicatorAssessmentResultFindByAssessmentIdQuery = `SELECT ia.id, a.institution_name, ia.created_at, ia.status, i.domain, i.aspect, i.indicator_number, 
		ia.level, sdd.document_url, COALESCE(ia.explanation, ''),  COALESCE(sddp.proof, '')
		FROM indicator_assessments ia
		LEFT JOIN assessments a
		ON ia.assessment_id = a.id
		LEFT JOIN indicators i
		ON ia.indicator_id = i.id
		LEFT JOIN support_data_documents sdd
		ON ia.id = sdd.indicator_assessment_id
		LEFT JOIN support_data_document_proofs sddp
		ON ia.id = sddp.indicator_assessment_id
		WHERE a.id = $1`

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
			&assessmentResult.Result.Explanation, &assessmentResult.Result.Proof,
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

const updateAssessmentStatusQuery = `UPDATE assessments
	SET status = $2, updated_at = $3
	WHERE id = $1
`
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
	updateAssessmentStatusStmt, err := s.db.PrepareContext(ctx, updateAssessmentStatusQuery)
	_, err = tx.StmtContext(ctx, updateAssessmentStatusStmt).ExecContext(ctx, resultDetail.AssessmentId,
		store.AssessmentStatus(store.COMPLETED), updatedAt)
	if err != nil {
		return fmt.Errorf("failed to update assessment: %w", err)
	}

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
