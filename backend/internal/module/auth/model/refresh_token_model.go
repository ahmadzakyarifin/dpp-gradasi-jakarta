package model

import (
	"time"

	usermodel "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/user/model"
)

type RefreshTokenModel struct {
	ID     uint                 `gorm:"column:id;primaryKey;autoIncrement"`
	UserID uint                 `gorm:"column:user_id;not null"`
	User   *usermodel.UserModel `gorm:"foreignKey:UserID;references:ID"`

	TokenHash string `gorm:"column:token_hash;not null"`

	DeviceName *string `gorm:"column:device_name"`
	IPAddress  *string `gorm:"column:ip_address"`
	UserAgent  *string `gorm:"column:user_agent"`

	ExpiresAt time.Time  `gorm:"column:expires_at;not null"`
	RevokedAt *time.Time `gorm:"column:revoked_at"`

	CreatedAt time.Time `gorm:"column:created_at;not null;autoCreateTime"`
}

func (RefreshTokenModel) TableName() string { return "refresh_tokens" }
