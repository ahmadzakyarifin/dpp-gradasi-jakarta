package model

import (
	"time"

	"gorm.io/gorm"
)

type Major struct {
	ID        uint      `gorm:"column:id;primaryKey;autoIncrement"`
	Code      *string   `gorm:"column:code"`
	Name      string    `gorm:"column:name"`
	IsActive  bool      `gorm:"column:is_active;default:true"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`

	// Aggregated read-only fields (diisi via Preload/Select terpisah, bukan kolom tabel).
	AcademicYearCount int `gorm:"-"`
	ClassCount        int `gorm:"-"`
	StudentCount      int `gorm:"-"`

	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (Major) TableName() string {
	return "majors"
}
