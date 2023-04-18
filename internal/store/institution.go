package store

import (
	"context"
	"time"
)

type InstitutionDetail struct {
	ID        int
	Name      string
	Category  string
	CreatedAt time.Time
}

type Institution interface {
	FindAll(ctx context.Context, institutionName string) ([]*InstitutionDetail, error)
	FindAllPagination(ctx context.Context, offset int, limit int, institutionName string) ([]*InstitutionDetail, error)
	FindByInstitutionName(ctx context.Context, institutionName string) (*InstitutionDetail, error)
	Insert(ctx context.Context, institution *InstitutionDetail) error
}
