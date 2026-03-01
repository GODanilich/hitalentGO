package dto

import "strings"

type PatchDepartmentRequest struct {
	Name     *string       `json:"name"`
	ParentID OptionalInt64 `json:"parent_id"`
}

func (r PatchDepartmentRequest) HasName() bool {
	return r.Name != nil
}

func (r PatchDepartmentRequest) NameTrimmed() string {
	if r.Name == nil {
		return ""
	}
	return strings.TrimSpace(*r.Name)
}

func (r PatchDepartmentRequest) HasParentID() bool {
	return r.ParentID.Present
}
