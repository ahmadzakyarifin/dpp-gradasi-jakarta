package model

import (
	"time"
)

type PasswordResetTokenModel struct {
	ID        uint       `gorm:"column:id;primaryKey;autoIncrement"`
	UserID    uint       `gorm:"column:user_id;not null"`
	TokenHash string     `gorm:"column:token_hash;not null"`
	ExpiresAt time.Time  `gorm:"column:expires_at;not null"`
	UsedAt    *time.Time `gorm:"column:used_at"`
	CreatedAt time.Time  `gorm:"column:created_at;not null;autoCreateTime"`
}

func (PasswordResetTokenModel) TableName() string { return "password_reset_tokens" }
