package model

import (
	"time"

	"gorm.io/gorm"
)

type ClassTemplateModel struct {
	ID          uint           `gorm:"column:id;primaryKey;autoIncrement"`
	MajorID     *uint          `gorm:"column:major_id"`
	Name        string         `gorm:"column:name"`
	GradeLevel  int            `gorm:"column:grade_level"`
	Description *string        `gorm:"column:description"`
	IsActive    bool           `gorm:"column:is_active;default:true"`
	CreatedAt   time.Time      `gorm:"column:created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index"`

	// Relasi/Join Fields (diisi via Select terpisah, bukan kolom tabel).
	MajorName *string `gorm:"column:major_name;->"`
}

func (ClassTemplateModel) TableName() string {
	return "class_templates"
}
