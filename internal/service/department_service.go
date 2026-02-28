package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgconn"
	"gorm.io/gorm"

	"GODanilich/hitalentGO/internal/apperr"
	"GODanilich/hitalentGO/internal/dto"
	"GODanilich/hitalentGO/internal/model"
	"GODanilich/hitalentGO/internal/repo"
)

const (
	maxNameLen = 200
)

type DepartmentService struct {
	deps *repo.DepartmentRepo
	emps *repo.EmployeeRepo
}

func NewDepartmentService(deps *repo.DepartmentRepo, emps *repo.EmployeeRepo) *DepartmentService {
	return &DepartmentService{deps: deps, emps: emps}
}

// POST /departments
func (s *DepartmentService) CreateDepartment(ctx context.Context, req dto.CreateDepartmentRequest) (*dto.DepartmentResponse, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" || len(name) > maxNameLen {
		return nil, apperr.Validation("name must be 1..200 characters")
	}

	if req.ParentID != nil {
		_, err := s.deps.GetByID(ctx, *req.ParentID)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperr.NotFound("parent department not found")
		}
		if err != nil {
			return nil, apperr.Internal("failed to load parent department", err)
		}
	}

	dep := &model.Department{
		Name:     name,
		ParentID: req.ParentID,
	}

	if err := s.deps.Create(ctx, dep); err != nil {
		if isUniqueViolation(err) {
			return nil, apperr.Conflict("department name already exists for this parent", err)
		}
		return nil, apperr.Internal("failed to create department", err)
	}

	return &dto.DepartmentResponse{
		ID:        dep.ID,
		Name:      dep.Name,
		ParentID:  dep.ParentID,
		CreatedAt: dep.CreatedAt,
	}, nil
}

// POST /departments/{id}/employees
func (s *DepartmentService) CreateEmployee(ctx context.Context, departmentID int64, req dto.CreateEmployeeRequest) (*dto.EmployeeResponse, error) {

	_, err := s.deps.GetByID(ctx, departmentID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperr.NotFound("department not found")
	}
	if err != nil {
		return nil, apperr.Internal("failed to load department", err)
	}

	fullName := strings.TrimSpace(req.FullName)
	if fullName == "" || len(fullName) > maxNameLen {
		return nil, apperr.Validation("full_name must be 1..200 characters")
	}

	position := strings.TrimSpace(req.Position)
	if position == "" || len(position) > maxNameLen {
		return nil, apperr.Validation("position must be 1..200 characters")
	}

	var hiredAt *time.Time
	if strings.TrimSpace(req.HiredAt) != "" {
		t, perr := time.Parse("2006-01-02", req.HiredAt)
		if perr != nil {
			return nil, apperr.Validation("hired_at must be YYYY-MM-DD")
		}
		hiredAt = &t
	}

	emp := &model.Employee{
		DepartmentID: departmentID,
		FullName:     fullName,
		Position:     position,
		HiredAt:      hiredAt,
	}

	if err := s.emps.Create(ctx, emp); err != nil {
		return nil, apperr.Internal("failed to create employee", err)
	}

	return &dto.EmployeeResponse{
		ID:           emp.ID,
		DepartmentID: emp.DepartmentID,
		FullName:     emp.FullName,
		Position:     emp.Position,
		HiredAt:      emp.HiredAt,
		CreatedAt:    emp.CreatedAt,
	}, nil
}

func isUniqueViolation(err error) bool {
	// Postgres unique violation code: 23505
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}
