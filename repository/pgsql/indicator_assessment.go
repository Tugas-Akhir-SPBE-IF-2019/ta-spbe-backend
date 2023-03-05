package pgsql

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"ta-spbe-backend/repository"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type indicatorAssessmentRepo struct {
	db *sql.DB
	ps map[string]*sql.Stmt
}

func NewIndicatorAssessmentRepo(db *sql.DB) (repository.IndicatorAssessmentRepository, error) {
	ps := make(map[string]*sql.Stmt, len(assessmentQueries))
	for key, query := range indicatorAssessmentQueries {
		stmt, err := prepareStmt(db, "indicatorAssessmentRepo", key, query)
		if err != nil {
			return nil, err
		}
		ps[key] = stmt
	}

	return &indicatorAssessmentRepo{db, ps}, nil
}

var indicatorAssessmentQueries = map[string]string{
	indicatorAssessmentFindAll:              indicatorAssessmentFindAllQuery,
	indicatorAssessmentFindAllPagination:    indicatorAssessmentFindAllPaginationQuery,
	indicatorAssessmentResultFindById:       indicatorAssessmentResultFindByIdQuery,
	insertAssessmentFeedback:                insertAssessmentFeedbackQuery,
	updateAssessmentFeedback:                updateAssessmentFeedbackQuery,
	updateValidatedAssessmentFeedbackStatus: updateValidatedAssessmentFeedbackStatusQuery,
	updateAssessmentStatus:                  updateAssessmentStatusQuery,
	updateIndicatorAssessmentResult:         updateIndicatorAssessmentResultQuery,
	insertSupportDataDocumentProof:          insertSupportDataDocumentProofQuery,
}

const indicatorAssessmentFindAll = "findAll"
const indicatorAssessmentFindAllQuery = `SELECT a.institution_name, ia.level, ia.created_at
		FROM indicator_assessments ia
		LEFT JOIN assessments a
		ON ia.assessment_id = a.id
		WHERE ia.level IS NOT NULL`

func (r *indicatorAssessmentRepo) FindAll(ctx context.Context) ([]*repository.IndicatorAssessmentDetail, error) {
	indicatorAssessmentList := []*repository.IndicatorAssessmentDetail{}

	tr := otel.Tracer("")
	_, span := tr.Start(ctx, "repository-indicator-asssessment-find-all")
	span.SetAttributes(attribute.Key("query").String(indicatorAssessmentFindAllQuery))
	defer span.End()

	rows, err := r.ps[indicatorAssessmentFindAll].QueryContext(ctx)
	err = fmt.Errorf("wow")
	if err != nil {
		log.Println("indicator assessment sql repo query context error: %w", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		indicatorAssessment := &repository.IndicatorAssessmentDetail{}
		err := rows.Scan(
			&indicatorAssessment.InstitutionName,
			&indicatorAssessment.SpbeIndex,
			&indicatorAssessment.SubmittedDate,
		)
		if err != nil {
			log.Println("indicator assessment sql repo scan error: %w", err)
			return nil, err
		}
		indicatorAssessmentList = append(indicatorAssessmentList, indicatorAssessment)
	}

	return indicatorAssessmentList, nil
}

const indicatorAssessmentFindAllPagination = "findAllPagination"
const indicatorAssessmentFindAllPaginationQuery = `SELECT a.institution_name, ia.level, ia.created_at
		FROM indicator_assessments ia
		LEFT JOIN assessments a
		ON ia.assessment_id = a.id
		WHERE ia.level IS NOT NULL
		LIMIT $2 OFFSET $1`

func (r *indicatorAssessmentRepo) FindAllPagination(ctx context.Context, offset int, limit int) ([]*repository.IndicatorAssessmentDetail, error) {
	indicatorAssessmentList := []*repository.IndicatorAssessmentDetail{}

	tr := otel.Tracer("")
	_, span := tr.Start(ctx, "repository-indicator-asssessment-find-all-pagination")
	span.SetAttributes(attribute.Key("query").String(indicatorAssessmentFindAllPaginationQuery))
	defer span.End()

	rows, err := r.ps[indicatorAssessmentFindAllPagination].QueryContext(ctx, offset, limit)
	if err != nil {
		log.Println("indicator assessment sql repo query context error: %w", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		indicatorAssessment := &repository.IndicatorAssessmentDetail{}
		err := rows.Scan(
			&indicatorAssessment.InstitutionName,
			&indicatorAssessment.SpbeIndex,
			&indicatorAssessment.SubmittedDate,
		)
		if err != nil {
			log.Println("indicator assessment sql repo scan error: %w", err)
			return nil, err
		}
		indicatorAssessmentList = append(indicatorAssessmentList, indicatorAssessment)
	}

	return indicatorAssessmentList, nil
}

const indicatorAssessmentResultFindById = "indicatorAssessmentResultFindById"
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

func (r *indicatorAssessmentRepo) FindIndicatorAssessmentResultById(ctx context.Context, id string) (repository.IndicatorAssessmentResultDetail, error) {
	assessmentResult := repository.IndicatorAssessmentResultDetail{}

	row := r.ps[indicatorAssessmentResultFindById].QueryRowContext(ctx, id)
	err := row.Scan(&assessmentResult.IndicatorAssessmentId, &assessmentResult.InstitutionName, &assessmentResult.SubmittedDate,
		&assessmentResult.AssessmentStatus, &assessmentResult.Result.Domain, &assessmentResult.Result.Aspect,
		&assessmentResult.Result.IndicatorNumber, &assessmentResult.Result.Level, &assessmentResult.Result.SupportDocument,
		&assessmentResult.Result.Explanation, &assessmentResult.Result.Proof,
	)
	if err != nil {
		log.Println("indicator assessment sql repo scan error: %w", err)
		return assessmentResult, err
	}

	return assessmentResult, nil
}

const insertAssessmentFeedback = "insertAssessmentFeedback"
const insertAssessmentFeedbackQuery = `INSERT INTO 
	indicator_assessment_feedbacks(
		id, indicator_assessment_id, level,
		feedback, created_at
	) values (
		$1, $2, $3, $4, $5
	)
`
const updateAssessmentFeedback = "updateAssessmentFeedback"
const updateAssessmentFeedbackQuery = `UPDATE indicator_assessment_feedbacks	
	SET level = $2, feedback = $3
	WHERE indicator_assessment_id = $1
`
const updateValidatedAssessmentFeedbackStatus = "updateValidatedAssessmentFeedbackStatus"
const updateValidatedAssessmentFeedbackStatusQuery = `UPDATE indicator_assessments
	SET status = $2
	WHERE id = $1
`

func (r *indicatorAssessmentRepo) ValidateAssessmentResult(ctx context.Context, resultCorrect bool, indicatorAssessmentResult *repository.IndicatorAssessmentResultDetail) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin validate assessment result tx: %w", err)
	}
	defer tx.Rollback()

	if repository.AssessmentStatus(indicatorAssessmentResult.AssessmentStatus) != repository.AssessmentStatus(repository.VALIDATED) {
		feedbackId := uuid.NewString()
		createdAt := time.Now().UTC()
		_, err = tx.StmtContext(ctx, r.ps[insertAssessmentFeedback]).ExecContext(ctx, feedbackId, indicatorAssessmentResult.IndicatorAssessmentId,
			indicatorAssessmentResult.ResultFeedback.Level, indicatorAssessmentResult.ResultFeedback.Feedback, createdAt)
		if err != nil {
			return fmt.Errorf("failed to insert indicator assessment feedback: %w", err)
		}

		if !resultCorrect {
			_, err = tx.StmtContext(ctx, r.ps[updateValidatedAssessmentFeedbackStatus]).ExecContext(ctx, indicatorAssessmentResult.IndicatorAssessmentId, repository.AssessmentStatus(repository.VALIDATED))
			if err != nil {
				return fmt.Errorf("failed to update indicator assessment validated status: %w", err)
			}
		}
	} else {
		if !resultCorrect {
			_, err = tx.StmtContext(ctx, r.ps[updateAssessmentFeedback]).ExecContext(ctx, indicatorAssessmentResult.IndicatorAssessmentId,
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

const updateAssessmentStatus = "updateAssessmentStatus"
const updateAssessmentStatusQuery = `UPDATE assessments
	SET status = $2, updated_at = $3
	WHERE id = $1
`
const updateIndicatorAssessmentResult = "updateIndicatorAssessmentResult"
const updateIndicatorAssessmentResultQuery = `UPDATE indicator_assessments
	SET status = $2, level = $3, explanation = $4, updated_at = $5
	WHERE id = $1
`
const insertSupportDataDocumentProof = "insertSupportDataDocumentProof"
const insertSupportDataDocumentProofQuery = `INSERT into
	support_data_document_proofs(
		id, indicator_assessment_id, support_data_document_id, proof, created_at
	) values(
		$1, $2, $3, $4, $5
	)
`

func (r *indicatorAssessmentRepo) UpdateAssessmentResult(ctx context.Context, resultDetail *repository.IndicatorAssessmentResultDetail) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to commit update assessment result tx: %w", err)
	}
	defer tx.Rollback()

	updatedAt := time.Now().UTC()
	_, err = tx.StmtContext(ctx, r.ps[updateAssessmentStatus]).ExecContext(ctx, resultDetail.AssessmentId,
		repository.AssessmentStatus(repository.COMPLETED), updatedAt)
	if err != nil {
		return fmt.Errorf("failed to update assessment: %w", err)
	}

	_, err = tx.StmtContext(ctx, r.ps[updateIndicatorAssessmentResult]).ExecContext(ctx,
		resultDetail.IndicatorAssessmentId, repository.AssessmentStatus(repository.COMPLETED),
		resultDetail.Result.Level, resultDetail.Result.Explanation, updatedAt)
	if err != nil {
		return fmt.Errorf("failed to update indicator assessment: %w", err)
	}

	supportDataDocumentProofId := uuid.NewString()
	_, err = tx.StmtContext(ctx, r.ps[insertSupportDataDocumentProof]).ExecContext(ctx,
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
