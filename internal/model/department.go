package model

import (
	"time"

	"github.com/google/uuid"
)

type Department struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey;column:id"`
	Name      string     `gorm:"column:name;type:varchar(200);not null"`
	ParentID  *uuid.UUID `gorm:"type:uuid;column:parent_id"`
	CreatedAt time.Time  `gorm:"column:created_at;autoCreateTime"`
}

func (Department) TableName() string { return "departments" }
