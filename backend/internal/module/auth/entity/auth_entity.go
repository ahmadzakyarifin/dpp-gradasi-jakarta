package entity

import "time"

// UserEntity adalah domain entity user untuk modul auth.
// Murni struct tanpa GORM tag / json tag (pola skill golang_sintax).
type UserEntity struct {
	ID                 uint
	RoleID             uint
	Name               string
	Email              string
	Password           string
	PhotoPath          string
	Status             string
	MustChangePassword bool
	EmailVerifiedAt    *time.Time
	LastLoginAt        *time.Time
	CreatedAt          time.Time
	Role               RoleEntity
}

// RoleEntity adalah domain entity role.
type RoleEntity struct {
	ID          uint
	Name        string
	DisplayName string
}

// RefreshTokenEntity adalah domain entity refresh token.
type RefreshTokenEntity struct {
	ID         uint
	UserID     uint
	TokenHash  string
	DeviceInfo string
	IPAddress  string
	ExpiresAt  time.Time
	CreatedAt  time.Time
}

// PasswordResetTokenEntity adalah domain entity token reset password.
type PasswordResetTokenEntity struct {
	ID        uint
	UserID    uint
	TokenHash string
	ExpiresAt time.Time
	UsedAt    *time.Time
	CreatedAt time.Time
}

// ActivationTokenEntity adalah domain entity token aktivasi akun.
type ActivationTokenEntity struct {
	ID        uint
	UserID    uint
	TokenHash string
	Channel   string
	ExpiresAt time.Time
	UsedAt    *time.Time
	CreatedAt time.Time
}
