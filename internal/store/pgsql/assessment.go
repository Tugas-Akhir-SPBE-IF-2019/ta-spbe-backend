package pgsql

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/store"
	"github.com/google/uuid"
)

type Assessment struct {
	db *sql.DB
}

func NewAssessment(db *sql.DB) *Assessment {
	return &Assessment{db: db}
}

const assessmentFindAllQuery = `SELECT ia.id, a.institution_name, ia.status, ia.created_at
	FROM assessments a
	RIGHT JOIN indicator_assessments ia
	ON ia.assessment_id = a.id`

func (s *Assessment) FindAll(ctx context.Context) ([]*store.AssessmentDetail, error) {
	assessmentList := []*store.AssessmentDetail{}

	rows, err := s.db.QueryContext(ctx, assessmentFindAllQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		assessment := &store.AssessmentDetail{}
		err := rows.Scan(
			&assessment.Id,
			&assessment.InstitutionName,
			&assessment.Status,
			&assessment.SubmittedDate,
		)
		if err != nil {
			return nil, err
		}
		assessmentList = append(assessmentList, assessment)
	}

	return assessmentList, nil
}

const assessmentFindAllPaginationQuery = `SELECT ia.id, a.institution_name, ia.status, ia.created_at
	FROM assessments a 
	RIGHT JOIN indicator_assessments ia
	ON ia.assessment_id = a.id
	LIMIT $2 OFFSET $1`

func (r *Assessment) FindAllPagination(ctx context.Context, offset int, limit int) ([]*store.AssessmentDetail, error) {
	assessmentList := []*store.AssessmentDetail{}

	rows, err := r.db.QueryContext(ctx, assessmentFindAllPaginationQuery, offset, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		assessment := &store.AssessmentDetail{}
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

func (s *Assessment) InsertUploadDocument(ctx context.Context, assessmentUploadDetail *store.AssessmentUploadDetail) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin insert upload document tx: %w", err)
	}
	defer tx.Rollback()
	assessmentId := uuid.NewString()
	assessmentCreatedAt := time.Now().UTC()

	assessmentInsertStmt, err := s.db.PrepareContext(ctx, assessmentInsertQuery)
	_, err = tx.StmtContext(ctx, assessmentInsertStmt).ExecContext(ctx,
		assessmentId, assessmentUploadDetail.UserId, store.AssessmentStatus(store.IN_PROGRESS), assessmentUploadDetail.AssessmentDetail.InstitutionName, assessmentCreatedAt)
	if err != nil {
		return fmt.Errorf("failed to insert assessment: %w", err)
	}
	assessmentUploadDetail.AssessmentDetail.Id = assessmentId
	assessmentUploadDetail.AssessmentDetail.SubmittedDate = assessmentCreatedAt

	var indicatorId string
	getIndicatorIdByIndicatorNumberStmt, err := s.db.PrepareContext(ctx, getIndicatorIdByIndicatorNumberQuery)
	row := tx.StmtContext(ctx, getIndicatorIdByIndicatorNumberStmt).QueryRowContext(ctx, assessmentUploadDetail.IndicatorAssessmentInfo.IndicatorNumber)
	err = row.Scan(&indicatorId)
	if err != nil {
		return fmt.Errorf("failed to get indicator id: %w", err)
	}

	indicatorAssessmentId := uuid.NewString()
	indicatorAssessmentUploadDocumentInsertStmt, err := s.db.PrepareContext(ctx, indicatorAssessmentUploadDocumentInsertQuery)
	_, err = tx.StmtContext(ctx, indicatorAssessmentUploadDocumentInsertStmt).ExecContext(ctx,
		indicatorAssessmentId, indicatorId, assessmentId, store.AssessmentStatus(store.IN_PROGRESS), 0, assessmentCreatedAt)
	if err != nil {
		return fmt.Errorf("failed to insert indicator assessment: %w", err)
	}
	assessmentUploadDetail.IndicatorAssessmentInfo.Id = indicatorAssessmentId

	supportDataDocumentId := uuid.NewString()
	supportDataDocumentUploadInsertStmt, err := s.db.PrepareContext(ctx, supportDataDocumentUploadInsertQuery)
	_, err = tx.StmtContext(ctx, supportDataDocumentUploadInsertStmt).ExecContext(ctx,
		supportDataDocumentId, indicatorAssessmentId, assessmentUploadDetail.SupportDataDocumentInfo.DocumentName, assessmentUploadDetail.SupportDataDocumentInfo.DocumentUrl, assessmentCreatedAt)
	if err != nil {
		return fmt.Errorf("failed to insert support data document: %w", err)
	}
	assessmentUploadDetail.SupportDataDocumentInfo.Id = supportDataDocumentId

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit insert upload document tx: %w", err)
	}

	return nil
}
