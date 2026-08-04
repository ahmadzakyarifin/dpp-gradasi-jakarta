package model

import (
	"time"

	"gorm.io/gorm"
)

// UserModel memetakan tabel users.
type UserModel struct {
	ID                 uint           `gorm:"column:id;primaryKey;autoIncrement"`
	Role               string         `gorm:"column:role;not null;type:enum('super_admin','admin');default:admin"`
	Name               string         `gorm:"column:name;not null"`
	Email              string               `gorm:"column:email;unique;not null"`
	Password           string               `gorm:"column:password"`
	PhotoPath          *string              `gorm:"column:photo_path"`
	Status             string               `gorm:"column:status;not null;default:inactive"`
	MustChangePassword bool                 `gorm:"column:must_change_password;not null;default:0"`
	EmailVerifiedAt    *time.Time           `gorm:"column:email_verified_at"`
	LastLoginAt        *time.Time           `gorm:"column:last_login_at"`
	CreatedAt          time.Time            `gorm:"column:created_at;not null;autoCreateTime"`
	UpdatedAt          time.Time            `gorm:"column:updated_at;not null;autoUpdateTime"`
	DeletedAt          gorm.DeletedAt       `gorm:"column:deleted_at;index"`
}

func (UserModel) TableName() string { return "users" }
