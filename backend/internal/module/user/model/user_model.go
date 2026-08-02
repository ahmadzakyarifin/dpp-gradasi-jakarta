package model

import (
	"time"

	rolemodel "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/role/model"
	"gorm.io/gorm"
)

// UserModel memetakan tabel users.
type UserModel struct {
	ID              uint                 `gorm:"column:id;primaryKey;autoIncrement"`
	RoleID          uint                 `gorm:"column:role_id;not null"`
	Role            *rolemodel.RoleModel `gorm:"foreignKey:RoleID;references:ID"`
	Name            string               `gorm:"column:name;not null"`
	Email           string               `gorm:"column:email;unique;not null"`
	Phone           string               `gorm:"column:phone;unique;not null"`
	DateOfBirth     *time.Time           `gorm:"column:date_of_birth"`
	CountryCode     *string              `gorm:"column:country_code"`
	PhotoPath       *string              `gorm:"column:photo_path"`
	PasswordHash    string               `gorm:"column:password_hash"`
	Status          string               `gorm:"column:status;not null;default:inactive"`
	EmailVerifiedAt *time.Time           `gorm:"column:email_verified_at"`
	PhoneVerifiedAt *time.Time           `gorm:"column:phone_verified_at"`
	LastLoginAt     *time.Time           `gorm:"column:last_login_at"`
	CreatedAt       time.Time            `gorm:"column:created_at;not null;autoCreateTime"`
	UpdatedAt       time.Time            `gorm:"column:updated_at;not null;autoUpdateTime"`
	DeletedAt       gorm.DeletedAt       `gorm:"column:deleted_at;index"`

	// Kolom hasil subquery (read-only, bukan kolom tabel asli).
	StudentCount int    `gorm:"->;column:student_count"`
	StudentNames string `gorm:"->;column:student_names"`
}

func (UserModel) TableName() string { return "users" }
