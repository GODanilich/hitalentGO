package service

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"GODanilich/hitalentGO/internal/apperr"
	"GODanilich/hitalentGO/internal/dto"
	"GODanilich/hitalentGO/internal/model"
)

// PATCH /departments/{id}
func (s *DepartmentService) UpdateDepartment(
	ctx context.Context,
	id int64,
	req dto.PatchDepartmentRequest,
) (*dto.DepartmentResponse, error) {
	dep, err := s.loadDepartment(ctx, id, "department not found", "failed to load department")
	if err != nil {
		return nil, err
	}

	if req.HasName() {
		name, err := validateNameField("name", req.NameTrimmed())
		if err != nil {
			return nil, err
		}
		dep.Name = name
	}

	if req.HasParentID() {
		if err := s.applyParentUpdate(ctx, id, dep, req); err != nil {
			return nil, err
		}
	}

	if err := s.deps.Update(ctx, dep); err != nil {
		if isUniqueViolation(err) {
			return nil, apperr.Conflict("department name already exists for this parent", err)
		}
		return nil, apperr.Internal("failed to update department", err)
	}

	return toDepartmentResponse(dep), nil
}

func (s *DepartmentService) applyParentUpdate(
	ctx context.Context,
	id int64,
	dep *model.Department,
	req dto.PatchDepartmentRequest,
) error {
	if req.ParentID.Value == nil {
		dep.ParentID = nil
		return nil
	}

	newParentID := *req.ParentID.Value
	if newParentID <= 0 {
		return apperr.Validation("parent_id must be a positive int64 or null")
	}
	if newParentID == id {
		return apperr.Conflict("department cannot be parent of itself", nil)
	}

	parent, err := s.deps.GetByID(ctx, newParentID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return apperr.NotFound("parent department not found")
	}
	if err != nil {
		return apperr.Internal("failed to load parent department", err)
	}

	if err := s.ensureNoCycles(ctx, id, parent); err != nil {
		return err
	}

	dep.ParentID = &newParentID
	return nil
}

func (s *DepartmentService) ensureNoCycles(ctx context.Context, departmentID int64, parent *model.Department) error {
	cur := parent
	for cur != nil && cur.ParentID != nil {
		if *cur.ParentID == departmentID {
			return apperr.Conflict("cycle detected: cannot move department into its subtree", nil)
		}

		next, err := s.deps.GetByID(ctx, *cur.ParentID)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			break
		}
		if err != nil {
			return apperr.Internal("failed during cycle check", err)
		}
		cur = next
	}
	return nil
}
