package dto

import "time"

type CreateEmployeeRequest struct {
	FullName string `json:"full_name"`
	Position string `json:"position"`
	HiredAt  string `json:"hired_at"` // принимается строка YYYY-MM-DD, распарс в сервисе
}

type EmployeeResponse struct {
	ID           int64      `json:"id"`
	DepartmentID int64      `json:"department_id"`
	FullName     string     `json:"full_name"`
	Position     string     `json:"position"`
	HiredAt      *time.Time `json:"hired_at"`
	CreatedAt    time.Time  `json:"created_at"`
}
