package model

import (
	"time"

	"gorm.io/gorm"
)

type Cohort struct {
	ID          uint           `gorm:"column:id;primaryKey;autoIncrement"`
	Name        string         `gorm:"column:name"`
	StartDate   time.Time      `gorm:"column:start_date"`
	EndDate     time.Time      `gorm:"column:end_date"`
	Description *string        `gorm:"column:description"`
	IsActive    bool           `gorm:"column:is_active;default:true"`
	CreatedAt   time.Time      `gorm:"column:created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index"`

	// Aggregate counts diisi via sub-query / Preload, bukan kolom tabel.
	StudentCount int `gorm:"-"`
}

func (Cohort) TableName() string {
	return "cohorts"
}
