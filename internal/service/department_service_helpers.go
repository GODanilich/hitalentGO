package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"

	"GODanilich/hitalentGO/internal/apperr"
	"GODanilich/hitalentGO/internal/dto"
	"GODanilich/hitalentGO/internal/model"
)

func validateNameField(fieldName, value string) (string, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" || len(trimmed) > maxNameLen {
		return "", apperr.Validation(fieldName + " must be 1..200 characters")
	}
	return trimmed, nil
}

func parseOptionalDate(raw string) (*time.Time, error) {
	if strings.TrimSpace(raw) == "" {
		return nil, nil
	}

	parsed, err := time.Parse("2006-01-02", raw)
	if err != nil {
		return nil, apperr.Validation("hired_at must be YYYY-MM-DD")
	}
	return &parsed, nil
}

func (s *DepartmentService) ensureDepartmentExists(ctx context.Context, id int64, notFoundMsg, internalMsg string) error {
	_, err := s.deps.GetByID(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return apperr.NotFound(notFoundMsg)
	}
	if err != nil {
		return apperr.Internal(internalMsg, err)
	}
	return nil
}

func toDepartmentResponse(dep *model.Department) *dto.DepartmentResponse {
	return &dto.DepartmentResponse{
		ID:        dep.ID,
		Name:      dep.Name,
		ParentID:  dep.ParentID,
		CreatedAt: dep.CreatedAt,
	}
}

func toEmployeeResponse(emp *model.Employee) *dto.EmployeeResponse {
	return &dto.EmployeeResponse{
		ID:           emp.ID,
		DepartmentID: emp.DepartmentID,
		FullName:     emp.FullName,
		Position:     emp.Position,
		HiredAt:      emp.HiredAt,
		CreatedAt:    emp.CreatedAt,
	}
}

func (s *DepartmentService) loadDepartment(ctx context.Context, id int64, notFoundMsg, internalMsg string) (*model.Department, error) {
	dep, err := s.deps.GetByID(ctx, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperr.NotFound(notFoundMsg)
	}
	if err != nil {
		return nil, apperr.Internal(internalMsg, err)
	}
	return dep, nil
}
