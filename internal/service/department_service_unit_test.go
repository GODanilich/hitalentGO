package service

import (
	"testing"
	"time"

	"GODanilich/hitalentGO/internal/apperr"
	"GODanilich/hitalentGO/internal/dto"
	"GODanilich/hitalentGO/internal/model"
)

func TestValidateNameField(t *testing.T) {
	t.Run("trims spaces", func(t *testing.T) {
		got, err := validateNameField("name", "  Backend  ")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != "Backend" {
			t.Fatalf("expected Backend, got %q", got)
		}
	})

	t.Run("rejects empty", func(t *testing.T) {
		_, err := validateNameField("name", "   ")
		if err == nil {
			t.Fatal("expected validation error")
		}
		ae, ok := apperr.As(err)
		if !ok || ae.Code != apperr.CodeValidation {
			t.Fatalf("expected validation app error, got %T %v", err, err)
		}
	})
}

func TestParseOptionalDate(t *testing.T) {
	v, err := parseOptionalDate("2024-12-31")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v == nil || v.Format("2006-01-02") != "2024-12-31" {
		t.Fatalf("unexpected parsed value: %v", v)
	}

	nilVal, err := parseOptionalDate("   ")
	if err != nil {
		t.Fatalf("unexpected error for empty value: %v", err)
	}
	if nilVal != nil {
		t.Fatalf("expected nil for empty value, got: %v", nilVal)
	}
}

func TestGetDepartmentTreeDepthValidation(t *testing.T) {
	svc := &DepartmentService{}
	_, err := svc.GetDepartmentTree(t.Context(), 1, 0, true)
	if err == nil {
		t.Fatal("expected validation error for depth=0")
	}
	ae, ok := apperr.As(err)
	if !ok || ae.Code != apperr.CodeValidation {
		t.Fatalf("expected validation app error, got %T %v", err, err)
	}
}

func TestBuildDepartmentNode(t *testing.T) {
	now := time.Now()
	rootID := int64(1)
	childID := int64(2)

	allDeps := map[int64]model.Department{
		rootID:  {ID: rootID, Name: "Root", CreatedAt: now},
		childID: {ID: childID, Name: "Child", ParentID: &rootID, CreatedAt: now},
	}
	childrenMap := map[int64][]int64{rootID: {childID}}
	employees := map[int64][]dto.EmployeeResponse{
		rootID: {{ID: 10, DepartmentID: rootID, FullName: "Alice", Position: "Lead", CreatedAt: now}},
	}

	node := buildDepartmentNode(rootID, true, allDeps, childrenMap, employees)
	if node.Department.ID != rootID {
		t.Fatalf("unexpected root id: %d", node.Department.ID)
	}
	if len(node.Children) != 1 || node.Children[0].Department.ID != childID {
		t.Fatalf("unexpected children: %#v", node.Children)
	}
	if len(node.Employees) != 1 || node.Employees[0].FullName != "Alice" {
		t.Fatalf("unexpected employees: %#v", node.Employees)
	}
}

func TestDeleteHelpers(t *testing.T) {
	if normalizeDeleteMode("  ") != "cascade" {
		t.Fatal("empty mode should default to cascade")
	}
	if !isInsideSubtree([]int64{1, 2, 3}, 2) {
		t.Fatal("expected candidate to be detected in subtree")
	}
	if isInsideSubtree([]int64{1, 2, 3}, 5) {
		t.Fatal("did not expect candidate outside subtree")
	}
}
