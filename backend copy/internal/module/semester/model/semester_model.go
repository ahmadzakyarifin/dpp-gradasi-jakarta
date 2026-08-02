package model

import (
	"time"

	"gorm.io/gorm"
)

type Semester struct {
	ID uint `gorm:"column:id;primaryKey;autoIncrement"`

	AcademicYearID uint           `gorm:"column:academic_year_id"`
	Name           string         `gorm:"column:name"`
	StartDate      time.Time      `gorm:"column:start_date"`
	EndDate        time.Time      `gorm:"column:end_date"`
	IsActive       bool           `gorm:"column:is_active;default:false"`
	CreatedAt      time.Time      `gorm:"column:created_at"`
	UpdatedAt      time.Time      `gorm:"column:updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at;index"`

	ClassMembershipCount int `gorm:"-"`
	BillingRuleCount     int `gorm:"-"`
	InvoiceCount         int `gorm:"-"`
	BatchCount           int `gorm:"-"`
}

func (Semester) TableName() string {
	return "semesters"
}
