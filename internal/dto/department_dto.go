package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateDepartmentRequest struct {
	Name     string     `json:"name"`
	ParentID *uuid.UUID `json:"parent_id"`
}

type DepartmentResponse struct {
	ID        uuid.UUID  `json:"id"`
	Name      string     `json:"name"`
	ParentID  *uuid.UUID `json:"parent_id"`
	CreatedAt time.Time  `json:"created_at"`
}
