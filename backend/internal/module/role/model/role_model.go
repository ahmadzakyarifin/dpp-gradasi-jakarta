package model

import (
	"time"

	"gorm.io/gorm"
)

// PermissionModel memetakan tabel permissions.
type PermissionModel struct {
	ID          uint      `gorm:"column:id;primaryKey;autoIncrement"`
	Name        string    `gorm:"column:name;not null"`
	DisplayName string    `gorm:"column:display_name;not null"`
	Module      string    `gorm:"column:module;not null"`
	Description string    `gorm:"column:description"`
	CreatedAt   time.Time `gorm:"column:created_at;not null;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"column:updated_at;not null;autoUpdateTime"`
}

func (PermissionModel) TableName() string { return "permissions" }

// RoleModel memetakan tabel roles.
type RoleModel struct {
	ID          uint              `gorm:"column:id;primaryKey;autoIncrement"`
	Name        string            `gorm:"column:name;not null"`
	DisplayName string            `gorm:"column:display_name;not null"`
	IsSystem    bool              `gorm:"column:is_system;not null;default:0"`
	IsActive    bool              `gorm:"column:is_active;not null;default:1"`
	Permissions []PermissionModel `gorm:"many2many:role_permissions;joinForeignKey:role_id;joinReferences:permission_id"`
	CreatedAt   time.Time         `gorm:"column:created_at;not null;autoCreateTime"`
	UpdatedAt   time.Time         `gorm:"column:updated_at;not null;autoUpdateTime"`
	DeletedAt   gorm.DeletedAt    `gorm:"column:deleted_at;index"`
}

func (RoleModel) TableName() string { return "roles" }

// RolePermissionModel memetakan tabel pivot role_permissions.
type RolePermissionModel struct {
	ID           uint      `gorm:"column:id;primaryKey;autoIncrement"`
	RoleID       uint      `gorm:"column:role_id;not null"`
	PermissionID uint      `gorm:"column:permission_id;not null"`
	CreatedAt    time.Time `gorm:"column:created_at;not null;autoCreateTime"`

	Role       *RoleModel       `gorm:"foreignKey:RoleID"`
	Permission *PermissionModel `gorm:"foreignKey:PermissionID"`
}

func (RolePermissionModel) TableName() string { return "role_permissions" }
