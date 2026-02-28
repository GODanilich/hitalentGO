package model

import (
	"time"

	"github.com/google/uuid"
)

type Employee struct {
	ID           uuid.UUID  `gorm:"type:uuid;primaryKey;column:id"`
	DepartmentID uuid.UUID  `gorm:"type:uuid;column:department_id;not null"`
	FullName     string     `gorm:"column:full_name;type:varchar(200);not null"`
	Position     string     `gorm:"column:position;type:varchar(200);not null"`
	HiredAt      *time.Time `gorm:"column:hired_at;type:date"`
	CreatedAt    time.Time  `gorm:"column:created_at;autoCreateTime"`
}

func (Employee) TableName() string { return "employees" }
