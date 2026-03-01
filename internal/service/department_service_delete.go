package service

import (
	"context"
	"strings"

	"GODanilich/hitalentGO/internal/apperr"
	"GODanilich/hitalentGO/internal/repo"
)

// DELETE /departments/{id}
func (s *DepartmentService) DeleteDepartment(
	ctx context.Context,
	id int64,
	mode string,
	reassignTo *int64,
) error {
	mode = normalizeDeleteMode(mode)

	if err := s.ensureDepartmentExists(ctx, id, "department not found", "failed to load department"); err != nil {
		return err
	}

	switch mode {
	case "cascade":
		return s.deleteCascade(ctx, id)
	case "reassign":
		return s.deleteWithReassign(ctx, id, reassignTo)
	default:
		return apperr.Validation("mode must be 'cascade' or 'reassign'")
	}
}

func normalizeDeleteMode(mode string) string {
	trimmed := strings.TrimSpace(mode)
	if trimmed == "" {
		return "cascade"
	}
	return trimmed
}

func (s *DepartmentService) deleteCascade(ctx context.Context, id int64) error {
	if err := s.deps.DeleteByID(ctx, id); err != nil {
		return apperr.Internal("failed to delete department", err)
	}
	return nil
}

func (s *DepartmentService) deleteWithReassign(ctx context.Context, id int64, reassignTo *int64) error {
	if reassignTo == nil {
		return apperr.Validation("reassign_to_department_id is required for mode=reassign")
	}
	if *reassignTo <= 0 {
		return apperr.Validation("reassign_to_department_id must be a positive int64")
	}

	if err := s.ensureDepartmentExists(ctx, *reassignTo, "reassign_to_department_id department not found", "failed to load reassign_to department"); err != nil {
		return err
	}

	subtreeIDs, err := s.collectSubtreeIDs(ctx, id)
	if err != nil {
		return apperr.Internal("failed to collect subtree", err)
	}
	if isInsideSubtree(subtreeIDs, *reassignTo) {
		return apperr.Conflict("reassign_to_department_id cannot be inside deleted subtree", nil)
	}

	if err := s.reassignAndDelete(ctx, id, subtreeIDs, *reassignTo); err != nil {
		return apperr.Internal("failed to delete department (reassign)", err)
	}
	return nil
}

func (s *DepartmentService) reassignAndDelete(
	ctx context.Context,
	id int64,
	subtreeIDs []int64,
	reassignTo int64,
) error {
	return s.depsTx(ctx, func(txDeps *repo.DepartmentRepo, txEmps *repo.EmployeeRepo) error {
		if err := txEmps.ReassignDepartments(ctx, subtreeIDs, reassignTo); err != nil {
			return err
		}
		if err := txDeps.DeleteByID(ctx, id); err != nil {
			return err
		}
		return nil
	})
}

func isInsideSubtree(subtreeIDs []int64, candidateID int64) bool {
	for _, id := range subtreeIDs {
		if id == candidateID {
			return true
		}
	}
	return false
}

// DELETE /departments/{id}
func (s *DepartmentService) collectSubtreeIDs(ctx context.Context, rootID int64) ([]int64, error) {
	seen := map[int64]struct{}{rootID: {}}
	level := []int64{rootID}

	for len(level) > 0 {
		children, err := s.deps.ListByParentIDs(ctx, level)
		if err != nil {
			return nil, err
		}

		next := make([]int64, 0, len(children))
		for _, child := range children {
			if _, ok := seen[child.ID]; ok {
				continue
			}
			seen[child.ID] = struct{}{}
			next = append(next, child.ID)
		}
		level = next
	}

	ids := make([]int64, 0, len(seen))
	for id := range seen {
		ids = append(ids, id)
	}
	return ids, nil
}
