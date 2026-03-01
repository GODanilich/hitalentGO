package service

import (
	"context"

	"GODanilich/hitalentGO/internal/apperr"
	"GODanilich/hitalentGO/internal/dto"
	"GODanilich/hitalentGO/internal/model"
)

// GET /departments/{id}
func (s *DepartmentService) GetDepartmentTree(
	ctx context.Context,
	rootID int64,
	depth int,
	includeEmployees bool,
) (*dto.DepartmentNodeResponse, error) {
	if depth < 1 || depth > maxDepth {
		return nil, apperr.Validation("depth must be between 1 and 5")
	}

	root, err := s.loadDepartment(ctx, rootID, "department not found", "failed to load department")
	if err != nil {
		return nil, err
	}

	allDeps, childrenMap, err := s.loadDepartmentLevels(ctx, root, depth)
	if err != nil {
		return nil, err
	}

	employeesByDept, err := s.loadEmployeesByDepartment(ctx, allDeps, includeEmployees)
	if err != nil {
		return nil, err
	}

	result := buildDepartmentNode(root.ID, includeEmployees, allDeps, childrenMap, employeesByDept)
	return &result, nil
}

func (s *DepartmentService) loadDepartmentLevels(
	ctx context.Context,
	root *model.Department,
	depth int,
) (map[int64]model.Department, map[int64][]int64, error) {
	allDeps := map[int64]model.Department{root.ID: *root}
	childrenMap := map[int64][]int64{}
	currentLevel := []int64{root.ID}

	for level := 1; level <= depth; level++ {
		children, err := s.deps.ListByParentIDs(ctx, currentLevel)
		if err != nil {
			return nil, nil, apperr.Internal("failed to load child departments", err)
		}

		nextLevel := make([]int64, 0, len(children))
		for _, child := range children {
			allDeps[child.ID] = child
			if child.ParentID != nil {
				childrenMap[*child.ParentID] = append(childrenMap[*child.ParentID], child.ID)
			}
			nextLevel = append(nextLevel, child.ID)
		}

		if len(nextLevel) == 0 {
			break
		}
		currentLevel = nextLevel
	}

	return allDeps, childrenMap, nil
}

func (s *DepartmentService) loadEmployeesByDepartment(
	ctx context.Context,
	allDeps map[int64]model.Department,
	includeEmployees bool,
) (map[int64][]dto.EmployeeResponse, error) {
	employeesByDept := map[int64][]dto.EmployeeResponse{}
	if !includeEmployees {
		return employeesByDept, nil
	}

	depIDs := make([]int64, 0, len(allDeps))
	for id := range allDeps {
		depIDs = append(depIDs, id)
	}

	emps, err := s.emps.ListByDepartmentIDs(ctx, depIDs)
	if err != nil {
		return nil, apperr.Internal("failed to load employees", err)
	}

	for _, e := range emps {
		emp := e
		employeesByDept[e.DepartmentID] = append(employeesByDept[e.DepartmentID], *toEmployeeResponse(&emp))
	}

	return employeesByDept, nil
}

func buildDepartmentNode(
	id int64,
	includeEmployees bool,
	allDeps map[int64]model.Department,
	childrenMap map[int64][]int64,
	employeesByDept map[int64][]dto.EmployeeResponse,
) dto.DepartmentNodeResponse {
	dep := allDeps[id]
	node := dto.DepartmentNodeResponse{
		Department: *toDepartmentResponse(&dep),
		Children:   []dto.DepartmentNodeResponse{},
	}

	if includeEmployees {
		node.Employees = employeesByDept[id]
		if node.Employees == nil {
			node.Employees = []dto.EmployeeResponse{}
		}
	}

	for _, childID := range childrenMap[id] {
		node.Children = append(node.Children, buildDepartmentNode(childID, includeEmployees, allDeps, childrenMap, employeesByDept))
	}

	return node
}
