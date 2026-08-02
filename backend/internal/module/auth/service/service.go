package service

import (
	"context"

	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/auth/dto"
)

type AuthService interface {
	Login(ctx context.Context, req dto.LoginRequest, ip string, userAgent string, device string) (*dto.LoginResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*dto.LoginResponse, error)
	ForgotPassword(ctx context.Context, req dto.ForgotPasswordRequest, ip string, userAgent string) (string, error)
	ResetPassword(ctx context.Context, req dto.ResetPasswordRequest, ip string, userAgent string) error
	ValidateResetToken(ctx context.Context, token string) error
	ChangePassword(ctx context.Context, userID uint, req dto.ChangePasswordRequest) error
	Logout(ctx context.Context, refreshToken string) error
}
