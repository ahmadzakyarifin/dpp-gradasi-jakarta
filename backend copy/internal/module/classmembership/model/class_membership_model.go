package model

import (
	"time"

	"gorm.io/gorm"
)

type ClassMembership struct {
	ID               uint           `gorm:"column:id;primaryKey;autoIncrement"`
	StudentID        uint           `gorm:"column:student_id"`
	ActiveClassID    uint           `gorm:"column:active_class_id"`
	AcademicYearID   uint           `gorm:"column:academic_year_id"`
	SemesterID       *uint          `gorm:"column:semester_id"`
	AttendanceNumber *int           `gorm:"column:attendance_number"`
	StartDate        *time.Time     `gorm:"column:start_date"`
	EndDate          *time.Time     `gorm:"column:end_date"`
	Status           string         `gorm:"column:status;default:active"`
	Note             *string        `gorm:"column:note"`
	CreatedAt        time.Time      `gorm:"column:created_at"`
	UpdatedAt        time.Time      `gorm:"column:updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (ClassMembership) TableName() string {
	return "class_memberships"
}
