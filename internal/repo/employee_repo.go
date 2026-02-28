package repo

import (
	"context"

	"gorm.io/gorm"

	"GODanilich/hitalentGO/internal/model"
)

type EmployeeRepo struct {
	db *gorm.DB
}

func NewEmployeeRepo(db *gorm.DB) *EmployeeRepo {
	return &EmployeeRepo{db: db}
}

func (r *EmployeeRepo) Create(ctx context.Context, emp *model.Employee) error {
	return r.db.WithContext(ctx).Create(emp).Error
}
