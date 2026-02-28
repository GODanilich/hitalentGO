package model

import "time"

type Department struct {
	ID        int64     `gorm:"primaryKey;column:id"`
	Name      string    `gorm:"column:name;type:varchar(200);not null"`
	ParentID  *int64    `gorm:"column:parent_id"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (Department) TableName() string { return "departments" }
