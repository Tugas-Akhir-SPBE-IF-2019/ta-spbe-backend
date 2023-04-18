package pgsql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Tugas-Akhir-SPBE-IF-2019/ta-spbe-backend/internal/store"
	"github.com/google/uuid"
)

type User struct {
	db *sql.DB
}

func NewUser(db *sql.DB) *User {
	return &User{db: db}
}

const userFindOneByEmailQuery = `SELECT id, email, name, role, COALESCE(contact_number, ''), COALESCE(linkedin_profile, ''), COALESCE(address, ''), COALESCE(profile_picture_link, '')
	FROM "users" WHERE email = $1
`

func (s *User) FindOneByEmail(ctx context.Context, email string) (*store.UserData, error) {
	user := &store.UserData{}

	row := s.db.QueryRowContext(ctx, userFindOneByEmailQuery, email)

	err := row.Scan(
		&user.ID, &user.Email, &user.Name, &user.Role, &user.ContactNumber, &user.LinkedinProfile, &user.Address, &user.ProfilePictureLink,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

const userFindOneByIDQuery = `SELECT id, email, name, role, COALESCE(contact_number, ''), COALESCE(linkedin_profile, ''), COALESCE(address, ''), COALESCE(profile_picture_link, '')
		FROM "users" WHERE id = $1
	`

func (s *User) FindOneByID(ctx context.Context, id string) (*store.UserData, error) {
	user := &store.UserData{}

	row := s.db.QueryRowContext(ctx, userFindOneByIDQuery, id)

	err := row.Scan(
		&user.ID, &user.Email, &user.Name, &user.Role, &user.ContactNumber, &user.LinkedinProfile, &user.Address, &user.ProfilePictureLink,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

const userFindEmailAndPasswordQuery = `SELECT id, email
		FROM "users" WHERE email = $1
	`

func (s *User) FindEmailAndPassword(ctx context.Context, email string) (*store.UserData, error) {
	user := &store.UserData{}

	row := s.db.QueryRowContext(ctx, userFindEmailAndPasswordQuery, email)

	err := row.Scan(
		&user.ID, &user.Email,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

const userFindAdminQuery = `SELECT id, email
		FROM "users" WHERE email = $1 AND role = 'admin'
	`

func (s *User) FindAdmin(ctx context.Context, email string) (*store.UserData, error) {
	user := &store.UserData{}

	row := s.db.QueryRowContext(ctx, userFindAdminQuery, email)

	err := row.Scan(
		&user.ID, &user.Email,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

const userInsert = `INSERT INTO
users(
	id, name, email, role, created_at
) values(
	$1, $2, $3, $4, $5	
)
`

func (s *User) Insert(ctx context.Context, user *store.UserData) error {
	insertStmt, err := s.db.PrepareContext(ctx, userInsert)
	if err != nil {
		return err
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}
	defer tx.Rollback()

	createdAt := time.Now().UTC()
	role := store.UserRole(store.RoleUser)
	_, err = tx.StmtContext(ctx, insertStmt).ExecContext(ctx,
		user.ID, user.Name, user.Email, role, createdAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}
	user.Role = string(role)
	user.CreatedAt = createdAt

	return nil

}

const userEvaluationFindByUserID = `SELECT ued.id, ued.user_id, ued.role, ued.institution_id, i.name, ued.evaluation_year, ued.created_at
	FROM user_evaluation_data ued
	LEFT JOIN institution i
	ON ued.institution_id = i.id 
	WHERE ued.user_id = $1
	ORDER by ued.created_at ASC
`

func (s *User) FindEvaluationDataByUserID(ctx context.Context, id string) ([]*store.UserEvaluationData, error) {
	userEvaluationList := []*store.UserEvaluationData{}

	rows, err := s.db.QueryContext(ctx, userEvaluationFindByUserID, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		userEvaluation := &store.UserEvaluationData{}
		err := rows.Scan(
			&userEvaluation.ID, &userEvaluation.UserID, &userEvaluation.Role,
			&userEvaluation.InstitutionID, &userEvaluation.InstitutionName,
			&userEvaluation.EvaluationYear, &userEvaluation.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		userEvaluationList = append(userEvaluationList, userEvaluation)
	}

	return userEvaluationList, nil
}

const userJobFindByUserID = `SELECT ujd.id, ujd.user_id, ujd.role, ujd.company, ujd.joined_date, ujd.created_at
	FROM user_job_data ujd
	WHERE ujd.user_id = $1
	ORDER by ujd.created_at ASC
`

func (s *User) FindJobDataByUserID(ctx context.Context, id string) ([]*store.UserJobData, error) {
	userEvaluationList := []*store.UserJobData{}

	rows, err := s.db.QueryContext(ctx, userJobFindByUserID, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		userJob := &store.UserJobData{}
		err := rows.Scan(
			&userJob.ID, &userJob.UserID, &userJob.Role,
			&userJob.Company, &userJob.JoinedDate,
			&userJob.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		userEvaluationList = append(userEvaluationList, userJob)
	}

	return userEvaluationList, nil
}

const userEvaluationInsert = `INSERT INTO
user_evaluation_data(
	id, user_id, role, institution_id, evaluation_year, created_at
) values(
	$1, $2, $3, $4, $5, $6
)
`

const userEvaluationUpdate = `UPDATE user_evaluation_data
	SET role = $2, institution_id = $3, evaluation_year = $4
	WHERE id = $1
`

func (s *User) InsertEvaluationData(ctx context.Context, userEvaluationData *store.UserEvaluationData) error {
	var insertStmt *sql.Stmt
	var err error

	if userEvaluationData.ID == "" {
		insertStmt, err = s.db.PrepareContext(ctx, userEvaluationInsert)
		if err != nil {
			return err
		}
	} else {
		insertStmt, err = s.db.PrepareContext(ctx, userEvaluationUpdate)
		if err != nil {
			return err
		}
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}
	defer tx.Rollback()

	if userEvaluationData.ID == "" {
		userEvaluationID := uuid.NewString()
		createdAt := time.Now().UTC()
		_, err = tx.StmtContext(ctx, insertStmt).ExecContext(ctx,
			userEvaluationID, userEvaluationData.UserID, userEvaluationData.Role,
			userEvaluationData.InstitutionID, userEvaluationData.EvaluationYear, createdAt,
		)
		if err != nil {
			return fmt.Errorf("failed to insert: %w", err)
		}
		userEvaluationData.ID = userEvaluationID
		userEvaluationData.CreatedAt = createdAt
	} else {
		_, err = tx.StmtContext(ctx, insertStmt).ExecContext(ctx,
			userEvaluationData.ID, userEvaluationData.Role,
			userEvaluationData.InstitutionID, userEvaluationData.EvaluationYear,
		)
		if err != nil {
			return fmt.Errorf("failed to insert: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}

	return nil
}

const userJobInsert = `INSERT INTO
user_job_data(
	id, user_id, role, company, joined_date, created_at
) values(
	$1, $2, $3, $4, $5, $6
)
`

const userJobUpdate = `UPDATE user_job_data
	SET	role = $2, company = $3, joined_date = $4
	WHERE id = $1
`

func (s *User) InsertJobData(ctx context.Context, userJobData *store.UserJobData) error {
	var insertStmt *sql.Stmt
	var err error

	if userJobData.ID == "" {
		insertStmt, err = s.db.PrepareContext(ctx, userJobInsert)
		if err != nil {
			return err
		}
	} else {
		insertStmt, err = s.db.PrepareContext(ctx, userJobUpdate)
		if err != nil {
			return err
		}
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}
	defer tx.Rollback()

	if userJobData.ID == "" {
		userJobID := uuid.NewString()
		createdAt := time.Now().UTC()
		_, err = tx.StmtContext(ctx, insertStmt).ExecContext(ctx,
			userJobID, userJobData.UserID, userJobData.Role,
			userJobData.Company, userJobData.JoinedDate, createdAt,
		)
		if err != nil {
			return fmt.Errorf("failed to insert: %w", err)
		}
		userJobData.ID = userJobID
		userJobData.CreatedAt = createdAt
	} else {
		_, err = tx.StmtContext(ctx, insertStmt).ExecContext(ctx,
			userJobData.ID, userJobData.Role,
			userJobData.Company, userJobData.JoinedDate,
		)
		if err != nil {
			return fmt.Errorf("failed to insert: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}

	return nil
}

const updateUser = `UPDATE users SET
	email = COALESCE($2, email),
	name = COALESCE($3, name),
	contact_number = COALESCE($4, contact_number),
	linkedin_profile = COALESCE($5, linkedin_profile),
	address = COALESCE($6, address),
	updated_at = $7
	WHERE id = $1
`

func (s *User) UpdateByID(ctx context.Context, user *store.UserData) error {
	updateStmt, err := s.db.PrepareContext(ctx, updateUser)
	if err != nil {
		return err
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}
	defer tx.Rollback()

	updatedAt := time.Now().UTC()
	_, err = tx.StmtContext(ctx, updateStmt).ExecContext(ctx,
		user.ID, user.Email, user.Name,
		user.ContactNumber, user.LinkedinProfile,
		user.Address, updatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to update: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}

	return nil
}

const updateUserWithPhoto = `UPDATE users SET
	email = COALESCE($2, email),
	name = COALESCE($3, name),
	contact_number = COALESCE($4, contact_number),
	linkedin_profile = COALESCE($5, linkedin_profile),
	address = COALESCE($6, address),
	profile_picture_link = COALESCE($7, profile_picture_link),
	updated_at = $8
	WHERE id = $1
`

func (s *User) UpdateWithPhotoByID(ctx context.Context, user *store.UserData) error {
	updateStmt, err := s.db.PrepareContext(ctx, updateUserWithPhoto)
	if err != nil {
		return err
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}
	defer tx.Rollback()

	updatedAt := time.Now().UTC()
	_, err = tx.StmtContext(ctx, updateStmt).ExecContext(ctx,
		user.ID, user.Email, user.Name,
		user.ContactNumber, user.LinkedinProfile,
		user.Address, user.ProfilePictureLink, updatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to update: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}

	return nil
}

const userCurrentInstitutionFindByUserID = `SELECT uci.id, uci.user_id, COALESCE(uci.institution_id, 0), COALESCE(uci.institution_name, ''), uci.role, uci.status, uci.created_at
	FROM user_current_institutions uci
	WHERE uci.user_id = $1
	ORDER by uci.created_at ASC
`

func (s *User) FindCurrentInstitutionDataByUserID(ctx context.Context, id string) ([]*store.UserCurrentInstitutionData, error) {
	userEvaluationList := []*store.UserCurrentInstitutionData{}

	rows, err := s.db.QueryContext(ctx, userCurrentInstitutionFindByUserID, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		userCurrentInstitution := &store.UserCurrentInstitutionData{}
		err := rows.Scan(
			&userCurrentInstitution.ID, &userCurrentInstitution.UserID, &userCurrentInstitution.InstitutionID,
			&userCurrentInstitution.InstitutionName, &userCurrentInstitution.Role,
			&userCurrentInstitution.Status, &userCurrentInstitution.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		userEvaluationList = append(userEvaluationList, userCurrentInstitution)
	}

	return userEvaluationList, nil
}

const userInstitutionInsert = `INSERT INTO
user_current_institutions(
	id, user_id, institution_id, institution_name, role, status, created_at
) values(
	$1, $2, $3, $4, $5, $6, $7
)
`

const userInstitutionWithoutIDInsert = `INSERT INTO
user_current_institutions(
	id, user_id, institution_name, role, status, created_at
) values(
	$1, $2, $3, $4, $5, $6
)
`

const userInstitutionUpdate = `UPDATE user_current_institutions
	SET institution_id = $2, institution_name = $3, role = $4
	WHERE id = $1
`

const userInstitutionWithoutIDUpdate = `UPDATE user_current_institutions
	SET institution_name = $2, role = $3
	WHERE id = $1
`

func (s *User) InsertCurrentInstitutionData(ctx context.Context, institutionData *store.UserCurrentInstitutionData) error {
	var insertStmt *sql.Stmt
	var err error
	var status string

	if institutionData.ID == "" {
		if institutionData.InstitutionID == 0 {
			insertStmt, err = s.db.PrepareContext(ctx, userInstitutionWithoutIDInsert)
			if err != nil {
				return err
			}
			status = "WAITING_FOR_VALIDATION"
		} else {
			insertStmt, err = s.db.PrepareContext(ctx, userInstitutionInsert)
			if err != nil {
				return err
			}
			status = "VALID"
		}
	} else {
		if institutionData.InstitutionID == 0 {
			insertStmt, err = s.db.PrepareContext(ctx, userInstitutionWithoutIDUpdate)
			if err != nil {
				return err
			}
		} else {
			insertStmt, err = s.db.PrepareContext(ctx, userInstitutionUpdate)
			if err != nil {
				return err
			}
		}
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}
	defer tx.Rollback()

	if institutionData.ID == "" {
		createdAt := time.Now().UTC()
		institutionDataID := uuid.NewString()
		if institutionData.InstitutionID == 0 {
			_, err = tx.StmtContext(ctx, insertStmt).ExecContext(ctx,
				institutionDataID, institutionData.UserID, institutionData.InstitutionName,
				institutionData.Role, status, createdAt,
			)
			if err != nil {
				return fmt.Errorf("failed to insert: %w", err)
			}
		} else {
			_, err = tx.StmtContext(ctx, insertStmt).ExecContext(ctx,
				institutionDataID, institutionData.UserID, institutionData.InstitutionID,
				institutionData.InstitutionName, institutionData.Role, status, createdAt,
			)
			if err != nil {
				return fmt.Errorf("failed to insert: %w", err)
			}
		}
		institutionData.ID = institutionDataID
		institutionData.CreatedAt = createdAt
	} else {
		if institutionData.InstitutionID == 0 {
			_, err = tx.StmtContext(ctx, insertStmt).ExecContext(ctx,
				institutionData.ID, institutionData.InstitutionName,
				institutionData.Role,
			)
			if err != nil {
				return fmt.Errorf("failed to insert: %w", err)
			}
		} else {
			_, err = tx.StmtContext(ctx, insertStmt).ExecContext(ctx,
				institutionData.ID, institutionData.InstitutionID,
				institutionData.InstitutionName, institutionData.Role,
			)
			if err != nil {
				return fmt.Errorf("failed to insert: %w", err)
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}

	return nil
}

const deleteCurrentInstitutionByIdQuery = `DELETE FROM user_current_institutions
	WHERE id = $1
`

func (s *User) DeleteCurrentInstitutionByID(ctx context.Context, id string) error {
	deleteStmt, err := s.db.PrepareContext(ctx, deleteCurrentInstitutionByIdQuery)
	if err != nil {
		return err
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.StmtContext(ctx, deleteStmt).ExecContext(ctx,
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to delete: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}

	return nil
}

const verifyCurrentInstitutionByIdQuery = `UPDATE user_current_institutions
	SET status = 'VALID'
	WHERE id = $1
`

func (s *User) VerifyInstitutionData(ctx context.Context, id string) error {
	deleteStmt, err := s.db.PrepareContext(ctx, verifyCurrentInstitutionByIdQuery)
	if err != nil {
		return err
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin tx: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.StmtContext(ctx, deleteStmt).ExecContext(ctx,
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to verify: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}

	return nil
}
