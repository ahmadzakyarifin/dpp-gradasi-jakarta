package model

import (
	"time"

	"gorm.io/gorm"
)

type AcademicYear struct {
	ID        uint           `gorm:"column:id;primaryKey;autoIncrement"`
	Name      string         `gorm:"column:name"`
	StartDate time.Time      `gorm:"column:start_date"`
	EndDate   time.Time      `gorm:"column:end_date"`
	IsActive  bool           `gorm:"column:is_active;default:false"`
	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`

	SemesterCount    int `gorm:"-"`
	ActiveClassCount int `gorm:"-"`
	StudentCount     int `gorm:"-"`
	BillingRuleCount int `gorm:"-"`
	InvoiceCount     int `gorm:"-"`
}

func (AcademicYear) TableName() string {
	return "academic_years"
}
