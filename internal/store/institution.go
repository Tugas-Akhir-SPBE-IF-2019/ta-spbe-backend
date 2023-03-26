package store

import "context"

type InstitutionDetail struct {
	Name     string
	Category string
}

type Institution interface {
	FindAll(ctx context.Context, institutionName string) ([]*InstitutionDetail, error)
	FindAllPagination(ctx context.Context, offset int, limit int, institutionName string) ([]*InstitutionDetail, error)
}
