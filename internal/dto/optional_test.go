package dto

import (
	"encoding/json"
	"testing"
)

func TestOptionalInt64Unmarshal(t *testing.T) {
	var v struct {
		ParentID OptionalInt64 `json:"parent_id"`
	}

	if err := json.Unmarshal([]byte(`{"parent_id":null}`), &v); err != nil {
		t.Fatalf("unmarshal null: %v", err)
	}
	if !v.ParentID.Present || v.ParentID.Value != nil {
		t.Fatalf("expected present=true and value=nil, got %+v", v.ParentID)
	}

	if err := json.Unmarshal([]byte(`{"parent_id":123}`), &v); err != nil {
		t.Fatalf("unmarshal number: %v", err)
	}
	if !v.ParentID.Present || v.ParentID.Value == nil || *v.ParentID.Value != 123 {
		t.Fatalf("unexpected parsed value: %+v", v.ParentID)
	}
}
