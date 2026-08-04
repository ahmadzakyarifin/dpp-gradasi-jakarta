package model

import (
	"time"

	"gorm.io/gorm"
)

// RoleModel memetakan tabel roles.
// Disalin dari modul role (dipakai join user->role untuk role_name/is_system).
type RoleModel struct {
	ID          uint           `gorm:"column:id;primaryKey;autoIncrement"`
	Name        string         `gorm:"column:name;not null"`
	DisplayName string         `gorm:"column:display_name;not null"`
	IsSystem    bool           `gorm:"column:is_system;not null;default:0"`
	IsActive    bool           `gorm:"column:is_active;not null;default:1"`
	CreatedAt   time.Time      `gorm:"column:created_at;not null;autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;not null;autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

func (RoleModel) TableName() string { return "roles" }
