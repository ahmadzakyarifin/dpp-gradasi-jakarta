package repository

import (
	"context"
	"time"

	userentity "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/user/entity"
	"gorm.io/gorm"
)

type AuthRepo interface {
	FindByEmail(ctx context.Context, email string) (*userentity.User, error)
	FindUserByID(ctx context.Context, id uint) (*userentity.User, error)
	UpdatePassword(ctx context.Context, userID uint, hashedPassword string) error

	SaveRefreshToken(ctx context.Context, userID uint, token string, expiresAt time.Time, ip string, userAgent string, device string) error
	FindUserByRefreshToken(ctx context.Context, token string) (*userentity.User, time.Time, error)
	DeleteRefreshToken(ctx context.Context, token string) error
	DeleteAllUserRefreshTokens(ctx context.Context, userID uint) error

	SaveAuthToken(ctx context.Context, userID uint, token string, tokenType string, expiresAt time.Time) error
	FindAuthToken(ctx context.Context, token string, tokenType string) (uint, error)
	DeleteAuthToken(ctx context.Context, token string, tokenType string) error
	GetDB() *gorm.DB
	WithTx(tx *gorm.DB) AuthRepo
}
