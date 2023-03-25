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

const assessmentFindAllQuery = `SELECT a.id, a.institution_name, a.status, a.created_at
	FROM assessments a `

func (s *Assessment) FindAll(ctx context.Context, queryInstitution string, status int, startDate string, endDate string) ([]*store.AssessmentDetail, error) {
	assessmentList := []*store.AssessmentDetail{}
	var queryKeys []string
	var queryParams []interface{}

	query := assessmentFindAllQuery
	if queryInstitution != "" {
		queryKeys = append(queryKeys, "Institution")
		queryParams = append(queryParams, queryInstitution)
	}

	if status != -1 {
		queryKeys = append(queryKeys, "Status")
		queryParams = append(queryParams, status)
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
		if index == 0 {
			query = query + "WHERE "
		} else {
			query = query + "AND "
		}

		switch key {
		case "Institution":
			query = query + fmt.Sprintf(`a.institution_name ILIKE '%%' || $%d || '%%' `, index+1)
		case "Status":
			query = query + fmt.Sprintf(`ia.status = $%d `, index+1)
		case "StartDate":
			query = query + fmt.Sprintf(`ia.created_at >= $%d `, index+1)
		case "EndDate":
			query = query + fmt.Sprintf(`ia.created_at <= $%d `, index+1)
		}
	}

	rows, err := s.db.QueryContext(ctx, query, queryParams...)
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

func (r *Assessment) FindAllPagination(ctx context.Context, offset int, limit int, queryInstitution string, status int, startDate string, endDate string) ([]*store.AssessmentDetail, error) {
	assessmentList := []*store.AssessmentDetail{}
	var queryKeys []string
	var queryParams []interface{}

	queryParams = append(queryParams, offset, limit)

	query := assessmentFindAllQuery
	if queryInstitution != "" {
		queryKeys = append(queryKeys, "Institution")
		queryParams = append(queryParams, queryInstitution)
	}

	if status != -1 {
		queryKeys = append(queryKeys, "Status")
		queryParams = append(queryParams, status)
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
		if index == 0 {
			query = query + "WHERE "
		} else {
			query = query + "AND "
		}

		// +3 is used because $1 and $2 are already used for pagination
		switch key {
		case "Institution":
			query = query + fmt.Sprintf(`a.institution_name ILIKE '%%' || $%d || '%%' `, index+3)
		case "Status":
			query = query + fmt.Sprintf(`ia.status = $%d `, index+3)
		case "StartDate":
			query = query + fmt.Sprintf(`ia.created_at >= $%d `, index+3)
		case "EndDate":
			query = query + fmt.Sprintf(`ia.created_at <= $%d `, index+3)
		}
	}

	query = query + `LIMIT $2 OFFSET $1 `

	rows, err := r.db.QueryContext(ctx, query, queryParams...)
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

	assessmentCreatedAt := time.Now().UTC()
	if assessmentUploadDetail.AssessmentDetail.Id == "" {
		assessmentId := uuid.NewString()

		assessmentInsertStmt, err := s.db.PrepareContext(ctx, assessmentInsertQuery)
		_, err = tx.StmtContext(ctx, assessmentInsertStmt).ExecContext(ctx,
			assessmentId, assessmentUploadDetail.UserId, store.AssessmentStatus(store.IN_PROGRESS), assessmentUploadDetail.AssessmentDetail.InstitutionName, assessmentCreatedAt)
		if err != nil {
			return fmt.Errorf("failed to insert assessment: %w", err)
		}
		assessmentUploadDetail.AssessmentDetail.Id = assessmentId
		assessmentUploadDetail.AssessmentDetail.SubmittedDate = assessmentCreatedAt
	}

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
		indicatorAssessmentId, indicatorId, assessmentUploadDetail.AssessmentDetail.Id, store.AssessmentStatus(store.IN_PROGRESS), 0, assessmentCreatedAt)
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
