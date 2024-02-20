package gormx

import (
	"time"

	"gorm.io/gorm"
)

type Model struct {
	ID        uint64         `gorm:"column:id;autoIncrement;primaryKey" json:"id"`
	CreatedAt time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}
