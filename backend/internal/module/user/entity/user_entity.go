package entity

import (
	"time"

	roleentity "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/role/entity"
)

type User struct {
	ID              uint
	RoleID          uint
	IsSystem        bool
	RoleName        string
	RoleDisplayName string
	Role            *roleentity.Role
	Name            string
	Email           string
	PhotoPath       *string
	PasswordHash    string
	Status          string
	EmailVerifiedAt *time.Time
	LastLoginAt     *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       *time.Time
}
