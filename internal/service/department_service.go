package service

import (
	"context"
	"errors"

	"github.com/jackc/pgconn"
	"gorm.io/gorm"

	"GODanilich/hitalentGO/internal/repo"
)

const (
	maxNameLen = 200
	maxDepth   = 5
)

type DepartmentService struct {
	deps *repo.DepartmentRepo
	emps *repo.EmployeeRepo
}

func NewDepartmentService(deps *repo.DepartmentRepo, emps *repo.EmployeeRepo) *DepartmentService {
	return &DepartmentService{deps: deps, emps: emps}
}

func (s *DepartmentService) depsTx(
	ctx context.Context,
	fn func(txDeps *repo.DepartmentRepo, txEmps *repo.EmployeeRepo) error,
) error {
	base := s.deps.DB().WithContext(ctx)

	return base.Transaction(func(tx *gorm.DB) error {
		txDeps := s.deps.WithDB(tx)
		txEmps := s.emps.WithDB(tx)
		return fn(txDeps, txEmps)
	})
}

func isUniqueViolation(err error) bool {
	// Postgres unique violation code: 23505
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}
