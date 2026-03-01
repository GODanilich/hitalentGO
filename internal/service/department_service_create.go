package service

import (
	"context"

	"GODanilich/hitalentGO/internal/apperr"
	"GODanilich/hitalentGO/internal/dto"
	"GODanilich/hitalentGO/internal/model"
)

// POST /departments
func (s *DepartmentService) CreateDepartment(ctx context.Context, req dto.CreateDepartmentRequest) (*dto.DepartmentResponse, error) {
	name, err := validateNameField("name", req.Name)
	if err != nil {
		return nil, err
	}

	if req.ParentID != nil {
		if err := s.ensureDepartmentExists(ctx, *req.ParentID, "parent department not found", "failed to load parent department"); err != nil {
			return nil, err
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

	return toDepartmentResponse(dep), nil
}

// POST /departments/{id}/employees
func (s *DepartmentService) CreateEmployee(ctx context.Context, departmentID int64, req dto.CreateEmployeeRequest) (*dto.EmployeeResponse, error) {
	if err := s.ensureDepartmentExists(ctx, departmentID, "department not found", "failed to load department"); err != nil {
		return nil, err
	}

	fullName, err := validateNameField("full_name", req.FullName)
	if err != nil {
		return nil, err
	}

	position, err := validateNameField("position", req.Position)
	if err != nil {
		return nil, err
	}

	hiredAt, err := parseOptionalDate(req.HiredAt)
	if err != nil {
		return nil, err
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

	return toEmployeeResponse(emp), nil
}
