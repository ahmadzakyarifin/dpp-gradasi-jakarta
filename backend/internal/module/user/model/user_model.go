package model

import (
	"time"

	rolemodel "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/role/model"
	"gorm.io/gorm"
)

type User struct {
	ID              uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	RoleID          uint           `gorm:"not null;index" json:"role_id"`
	Name            string         `gorm:"size:100;not null" json:"name"`
	Email           string         `gorm:"uniqueIndex;size:150;not null" json:"email"`
	Password        string         `gorm:"size:255" json:"-"`
	PhotoPath       string         `gorm:"size:500" json:"photo_path,omitempty"`
	Status          string         `gorm:"type:enum('active','inactive','pending_activation');default:'inactive'" json:"status"`
	EmailVerifiedAt *time.Time     `json:"email_verified_at,omitempty"`
	LastLoginAt     *time.Time     `json:"last_login_at,omitempty"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`

	// Relasi
	Role rolemodel.Role `gorm:"foreignKey:RoleID" json:"role,omitempty"`
}

func (User) TableName() string {
	return "users"
}
