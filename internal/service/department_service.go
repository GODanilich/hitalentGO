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
	maxDepth   = 5
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

	root, err := s.deps.GetByID(ctx, rootID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, apperr.NotFound("department not found")
	}
	if err != nil {
		return nil, apperr.Internal("failed to load department", err)
	}

	// allDeps: id -> department
	allDeps := map[int64]model.Department{
		root.ID: *root,
	}
	// childrenMap: parentID -> []childIDs
	childrenMap := map[int64][]int64{}

	currentLevel := []int64{root.ID}

	for level := 1; level <= depth; level++ {
		children, err := s.deps.ListByParentIDs(ctx, currentLevel)
		if err != nil {
			return nil, apperr.Internal("failed to load child departments", err)
		}

		nextLevel := make([]int64, 0, len(children))
		for _, ch := range children {
			allDeps[ch.ID] = ch
			if ch.ParentID != nil {
				childrenMap[*ch.ParentID] = append(childrenMap[*ch.ParentID], ch.ID)
			}
			nextLevel = append(nextLevel, ch.ID)
		}

		if len(nextLevel) == 0 {
			break
		}
		currentLevel = nextLevel
	}

	employeesByDept := map[int64][]dto.EmployeeResponse{}
	if includeEmployees {
		depIDs := make([]int64, 0, len(allDeps))
		for id := range allDeps {
			depIDs = append(depIDs, id)
		}

		emps, err := s.emps.ListByDepartmentIDs(ctx, depIDs)
		if err != nil {
			return nil, apperr.Internal("failed to load employees", err)
		}
		for _, e := range emps {
			employeesByDept[e.DepartmentID] = append(employeesByDept[e.DepartmentID], dto.EmployeeResponse{
				ID:           e.ID,
				DepartmentID: e.DepartmentID,
				FullName:     e.FullName,
				Position:     e.Position,
				HiredAt:      e.HiredAt,
				CreatedAt:    e.CreatedAt,
			})
		}
	}

	var buildNode func(id int64) dto.DepartmentNodeResponse
	buildNode = func(id int64) dto.DepartmentNodeResponse {
		dep := allDeps[id]
		node := dto.DepartmentNodeResponse{
			Department: dto.DepartmentResponse{
				ID:        dep.ID,
				Name:      dep.Name,
				ParentID:  dep.ParentID,
				CreatedAt: dep.CreatedAt,
			},
			Children: []dto.DepartmentNodeResponse{},
		}

		if includeEmployees {
			if list, ok := employeesByDept[id]; ok {
				node.Employees = list
			} else {
				node.Employees = []dto.EmployeeResponse{}
			}
		}

		for _, childID := range childrenMap[id] {
			node.Children = append(node.Children, buildNode(childID))
		}
		return node
	}

	result := buildNode(root.ID)
	return &result, nil
}

func isUniqueViolation(err error) bool {
	// Postgres unique violation code: 23505
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}
