package model

import "time"

type RefreshToken struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID     uint      `gorm:"not null;index" json:"user_id"`
	TokenHash  string    `gorm:"uniqueIndex;size:255;not null" json:"-"`
	DeviceInfo string    `gorm:"size:255" json:"device_info,omitempty"`
	IPAddress  string    `gorm:"size:45" json:"ip_address,omitempty"`
	ExpiresAt  time.Time `gorm:"not null" json:"expires_at"`
	CreatedAt  time.Time `json:"created_at"`
}

func (RefreshToken) TableName() string {
	return "refresh_tokens"
}

type PasswordResetToken struct {
	ID        uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint       `gorm:"not null;index" json:"user_id"`
	TokenHash string     `gorm:"uniqueIndex;size:255;not null" json:"-"`
	ExpiresAt time.Time  `gorm:"not null" json:"expires_at"`
	UsedAt    *time.Time `json:"used_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

func (PasswordResetToken) TableName() string {
	return "password_reset_tokens"
}

type ActivationToken struct {
	ID        uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint       `gorm:"not null;index" json:"user_id"`
	TokenHash string     `gorm:"uniqueIndex;size:255;not null" json:"-"`
	Channel   string     `gorm:"type:enum('email','whatsapp','all');default:'email'" json:"channel"`
	ExpiresAt time.Time  `gorm:"not null" json:"expires_at"`
	UsedAt    *time.Time `json:"used_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

func (ActivationToken) TableName() string {
	return "activation_tokens"
}
