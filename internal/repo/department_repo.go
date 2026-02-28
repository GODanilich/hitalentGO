package repo

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"GODanilich/hitalentGO/internal/model"
)

type DepartmentRepo struct {
	db *gorm.DB
}

func NewDepartmentRepo(db *gorm.DB) *DepartmentRepo {
	return &DepartmentRepo{db: db}
}

func (r *DepartmentRepo) Create(ctx context.Context, dep *model.Department) error {
	return r.db.WithContext(ctx).Create(dep).Error
}

func (r *DepartmentRepo) GetByID(ctx context.Context, id int64) (*model.Department, error) {
	var dep model.Department
	err := r.db.WithContext(ctx).First(&dep, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, gorm.ErrRecordNotFound
	}
	if err != nil {
		return nil, err
	}
	return &dep, nil
}
