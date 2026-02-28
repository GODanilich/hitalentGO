package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateEmployeeRequest struct {
	FullName string `json:"full_name"`
	Position string `json:"position"`
	HiredAt  string `json:"hired_at"` // принимается строка YYYY-MM-DD, распарс в сервисе
}

type EmployeeResponse struct {
	ID           uuid.UUID  `json:"id"`
	DepartmentID uuid.UUID  `json:"department_id"`
	FullName     string     `json:"full_name"`
	Position     string     `json:"position"`
	HiredAt      *time.Time `json:"hired_at"`
	CreatedAt    time.Time  `json:"created_at"`
}
