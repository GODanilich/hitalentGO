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

func (r *EmployeeRepo) ListByDepartmentIDs(ctx context.Context, depIDs []int64) ([]model.Employee, error) {
	if len(depIDs) == 0 {
		return []model.Employee{}, nil
	}
	var emps []model.Employee
	err := r.db.WithContext(ctx).
		Where("department_id IN ?", depIDs).
		Order("created_at ASC, id ASC").
		Find(&emps).Error
	return emps, err
}
