package model

import (
	"time"

	"gorm.io/gorm"
)

type ActiveClass struct {
	ID                  uint           `gorm:"column:id;primaryKey;autoIncrement"`
	AcademicYearID      uint           `gorm:"column:academic_year_id"`
	ClassTemplateID     uint           `gorm:"column:class_template_id"`
	Name                string         `gorm:"column:name"`
	HomeroomNumber      *string        `gorm:"column:homeroom_number"`
	HomeroomTeacherName *string        `gorm:"column:homeroom_teacher_name"`
	Room                *string        `gorm:"column:room"`
	Capacity            int            `gorm:"column:capacity;default:0"`
	StudentCount        int            `gorm:"column:student_count;default:0"`
	IsActive            bool           `gorm:"column:is_active;default:true"`
	CreatedAt           time.Time      `gorm:"column:created_at"`
	UpdatedAt           time.Time      `gorm:"column:updated_at"`
	DeletedAt           gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (ActiveClass) TableName() string {
	return "active_classes"
}
