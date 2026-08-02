package entity

import (
	"time"

	roleentity "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/role/entity"
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
	Phone           string
	DateOfBirth     *time.Time
	CountryCode     *string
	PhotoPath       *string
	PasswordHash    string
	Status          string
	EmailVerifiedAt *time.Time
	PhoneVerifiedAt *time.Time
	LastLoginAt     *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       *time.Time
	StudentCount    int
	StudentNames    string
}
