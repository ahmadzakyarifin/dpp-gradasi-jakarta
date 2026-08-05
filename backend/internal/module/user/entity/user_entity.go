package entity

import (
	"time"
)

type User struct {
	ID                 uint
	Role               string // super_admin | admin (enum users.role)
	IsSystem           bool
	Name               string
	Email              string
	EmailPending       *string
	PhotoPath          *string
	Password           string
	Status             string
	MustChangePassword bool
	EmailVerifiedAt    *time.Time
	LastLoginAt        *time.Time
	CreatedAt          time.Time
	UpdatedAt          time.Time
	DeletedAt          *time.Time
}
