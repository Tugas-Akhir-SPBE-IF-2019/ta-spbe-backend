package pgsql

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"ta-spbe-backend/repository"
	"time"

	"github.com/google/uuid"
)

type assessmentRepo struct {
	db *sql.DB
	ps map[string]*sql.Stmt
}

func NewAssessmentRepo(db *sql.DB) (repository.AssessmentRepository, error) {
	ps := make(map[string]*sql.Stmt, len(assessmentQueries))
	for key, query := range assessmentQueries {
		stmt, err := prepareStmt(db, "assessmentRepo", key, query)
		if err != nil {
			return nil, err
		}
		ps[key] = stmt
	}

	return &assessmentRepo{db, ps}, nil
}

var assessmentQueries = map[string]string{
	assessmentFindAll:                       assessmentFindAllQuery,
	assessmentFindAllPagination:             assessmentFindAllPaginationQuery,
	assessmentInsert:                        assessmentInsertQuery,
	getIndicatorIdByIndicatorNumber:         getIndicatorIdByIndicatorNumberQuery,
	indicatorAssessmentUploadDocumentInsert: indicatorAssessmentUploadDocumentInsertQuery,
	supportDataDocumentUploadInsert:         supportDataDocumentUploadInsertQuery,
}

const assessmentFindAll = "findAll"
const assessmentFindAllQuery = `SELECT a.id, a.institution_name, a.status, a.created_at
		FROM assessments a`

func (r *assessmentRepo) FindAll(ctx context.Context) ([]*repository.AssessmentDetail, error) {
	assessmentList := []*repository.AssessmentDetail{}

	rows, err := r.ps[assessmentFindAll].QueryContext(ctx)
	if err != nil {
		log.Println("assessment sql repo query context error: %w", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		assessment := &repository.AssessmentDetail{}
		err := rows.Scan(
			&assessment.Id,
			&assessment.InstitutionName,
			&assessment.Status,
			&assessment.SubmittedDate,
		)
		if err != nil {
			log.Println("assessment sql repo scan error: %w", err)
			return nil, err
		}
		assessmentList = append(assessmentList, assessment)
	}

	return assessmentList, nil
}

const assessmentFindAllPagination = "findAllPagination"
const assessmentFindAllPaginationQuery = `SELECT a.id, a.institution_name, a.status, a.created_at
		FROM assessments a LIMIT $2 OFFSET $1`

func (r *assessmentRepo) FindAllPagination(ctx context.Context, offset int, limit int) ([]*repository.AssessmentDetail, error) {
	assessmentList := []*repository.AssessmentDetail{}

	rows, err := r.ps[assessmentFindAllPagination].QueryContext(ctx, offset, limit)
	if err != nil {
		log.Println("assessment sql repo query context error: %w", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		assessment := &repository.AssessmentDetail{}
		err := rows.Scan(
			&assessment.Id,
			&assessment.InstitutionName,
			&assessment.Status,
			&assessment.SubmittedDate,
		)
		if err != nil {
			log.Println("assessment sql repo scan error: %w", err)
			return nil, err
		}
		assessmentList = append(assessmentList, assessment)
	}

	return assessmentList, nil
}

const assessmentInsert = "assessmentInsert"
const assessmentInsertQuery = `INSERT into
	assessments(
		id, user_id, status, institution_name, created_at
	) values(
		$1, $2, $3, $4, $5
	)
`
const getIndicatorIdByIndicatorNumber = "getIndicatorIdByIndicatorNumber"
const getIndicatorIdByIndicatorNumberQuery = `SELECT i.id
	FROM indicators i
	WHERE i.indicator_number = $1
`
const indicatorAssessmentUploadDocumentInsert = "indicatorAssessmentUploadDocumentInsert"
const indicatorAssessmentUploadDocumentInsertQuery = `INSERT into
	indicator_assessments(
		id, indicator_id,assessment_id, status, level, created_at
	) values(
		$1, $2, $3, $4, $5, $6
	)
`
const supportDataDocumentUploadInsert = "supportDataDocumentUploadInsert"
const supportDataDocumentUploadInsertQuery = `INSERT into
	support_data_documents(
		id, indicator_assessment_id, document_name, document_url, created_at
	) values(
		$1, $2, $3, $4, $5
	)
`

func (r *assessmentRepo) InsertUploadDocument(ctx context.Context, assessmentUploadDetail *repository.AssessmentUploadDetail) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin insert upload document tx: %w", err)
	}
	defer tx.Rollback()
	assessmentId := uuid.NewString()
	assessmentCreatedAt := time.Now().UTC()

	_, err = tx.StmtContext(ctx, r.ps[assessmentInsert]).ExecContext(ctx,
		assessmentId, assessmentUploadDetail.UserId, repository.AssessmentStatus(repository.IN_PROGRESS), assessmentUploadDetail.AssessmentDetail.InstitutionName, assessmentCreatedAt)
	if err != nil {
		return fmt.Errorf("failed to insert assessment: %w", err)
	}
	assessmentUploadDetail.AssessmentDetail.Id = assessmentId
	assessmentUploadDetail.AssessmentDetail.SubmittedDate = assessmentCreatedAt

	var indicatorId string
	row := tx.StmtContext(ctx, r.ps[getIndicatorIdByIndicatorNumber]).QueryRowContext(ctx, assessmentUploadDetail.IndicatorAssessmentInfo.IndicatorNumber)
	err = row.Scan(&indicatorId)
	if err != nil {
		return fmt.Errorf("failed to get indicator id: %w", err)
	}

	indicatorAssessmentId := uuid.NewString()
	_, err = tx.StmtContext(ctx, r.ps[indicatorAssessmentUploadDocumentInsert]).ExecContext(ctx,
		indicatorAssessmentId, indicatorId, assessmentId, repository.AssessmentStatus(repository.IN_PROGRESS), 0, assessmentCreatedAt)
	if err != nil {
		return fmt.Errorf("failed to insert indicator assessment: %w", err)
	}
	assessmentUploadDetail.IndicatorAssessmentInfo.Id = indicatorAssessmentId

	supportDataDocumentId := uuid.NewString()
	_, err = tx.StmtContext(ctx, r.ps[supportDataDocumentUploadInsert]).ExecContext(ctx,
		supportDataDocumentId, indicatorAssessmentId, assessmentUploadDetail.SupportDataDocumentInfo.DocumentName, assessmentUploadDetail.SupportDataDocumentInfo.DocumentName, assessmentCreatedAt)
	if err != nil {
		return fmt.Errorf("failed to insert support data document: %w", err)
	}
	assessmentUploadDetail.SupportDataDocumentInfo.Id = supportDataDocumentId

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit insert upload document tx: %w", err)
	}

	return nil
}
