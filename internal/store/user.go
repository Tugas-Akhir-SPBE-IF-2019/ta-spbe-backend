package store

import (
	"context"
	"database/sql"
	"time"
)

type UserRole string

const (
	RoleAdmin = UserRole("ADMIN")
	RoleUser  = UserRole("USER")
)

type UserData struct {
	ID                 string
	Name               string
	Email              string
	Role               string
	ContactNumber      string
	LinkedinProfile    string
	Address            string
	ProfilePictureLink string
	CreatedAt          time.Time
	UpdatedAt          sql.NullTime
}

type UserEvaluationData struct {
	ID              string
	UserID          string
	Role            string
	InstitutionID   int
	InstitutionName string
	EvaluationYear  int
	CreatedAt       time.Time
	UpdatedAt       sql.NullTime
}

type UserJobData struct {
	ID         string
	UserID     string
	Role       string
	Company    string
	JoinedDate int
	CreatedAt  time.Time
}

type UserCurrentInstitutionData struct {
	ID              string
	UserID          string
	InstitutionID   int
	InstitutionName string
	Role            string
	Status          string
	CreatedAt       time.Time
}

type User interface {
	FindOneByEmail(ctx context.Context, email string) (*UserData, error)
	FindOneByID(ctx context.Context, id string) (*UserData, error)
	FindEmailAndPassword(ctx context.Context, email string) (*UserData, error)
	FindAdmin(ctx context.Context, email string) (*UserData, error)
	Insert(ctx context.Context, user *UserData) error
	FindEvaluationDataByUserID(ctx context.Context, id string) ([]*UserEvaluationData, error)
	FindJobDataByUserID(ctx context.Context, id string) ([]*UserJobData, error)
	FindCurrentInstitutionDataByUserID(ctx context.Context, id string) ([]*UserCurrentInstitutionData, error)
	InsertEvaluationData(ctx context.Context, evaluationData *UserEvaluationData) error
	InsertJobData(ctx context.Context, jobData *UserJobData) error
	InsertCurrentInstitutionData(ctx context.Context, institutionData *UserCurrentInstitutionData) error
	DeleteCurrentInstitutionByID(ctx context.Context, id string) error
	UpdateByID(ctx context.Context, user *UserData) error
	UpdateWithPhotoByID(ctx context.Context, user *UserData) error
	VerifyInstitutionData(ctx context.Context, id string) error
}
