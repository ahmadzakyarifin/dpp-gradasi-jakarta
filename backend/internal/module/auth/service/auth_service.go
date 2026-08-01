package auth

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/ahmadzakyarifin/dpp-gradasi/backend/config"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/helper"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/infrastructure"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/auth/dto"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/auth/model"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/auth/repository"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type AuthService interface {
	Login(req *dto.LoginRequest) (*dto.AuthResponse, string, int, error)
	Refresh(refreshTokenStr string) (*dto.RefreshTokenResponse, string, int, error)
	Logout(userID uint, refreshTokenStr string) error
	ForgotPassword(req *dto.ForgotPasswordRequest, appURL string) error
	ValidateResetToken(token string) error
	ResetPassword(req *dto.ResetPasswordRequest) error
	ValidateActivationToken(token string) error
	ActivateAccount(req *dto.ActivateAccountRequest) error
	ChangePassword(userID uint, req *dto.ChangePasswordRequest) error
	GetProfile(userID uint) (*dto.AuthUserResponse, error)
}

type authService struct {
	repo  repository.AuthRepo
	cfg   *config.Config
	mail  *infrastructure.Mailer
	redis *redis.Client
}

func NewAuthService(repo repository.AuthRepo, cfg *config.Config, mail *infrastructure.Mailer, redis *redis.Client) AuthService {
	return &authService{
		repo:  repo,
		cfg:   cfg,
		mail:  mail,
		redis: redis,
	}
}

func (s *authService) Login(req *dto.LoginRequest) (*dto.AuthResponse, string, int, error) {
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", 0, helper.NewServiceError("AUTH_INVALID_CREDENTIALS", "Email atau password salah.", nil)
		}
		return nil, "", 0, helper.NewServiceError("SERVER_ERROR", "Terjadi kesalahan pada server.", err)
	}

	if !helper.CheckPassword(req.Password, user.Password) {
		return nil, "", 0, helper.NewServiceError("AUTH_INVALID_CREDENTIALS", "Email atau password salah.", nil)
	}

	if user.Status == "pending_activation" {
		return nil, "", 0, helper.NewServiceError("AUTH_ACCOUNT_PENDING", "Akun Anda belum diaktifkan. Silakan periksa email Anda untuk mengaktifkan akun.", nil)
	}

	if user.Status == "inactive" {
		return nil, "", 0, helper.NewServiceError("AUTH_ACCOUNT_INACTIVE", "Akun Anda telah dinonaktifkan. Silakan hubungi Super Admin.", nil)
	}

	_ = s.repo.UpdateLastLogin(user.ID)

	accessTTL := s.cfg.JWT.AccessTTLMinutes
	refreshTTL := s.cfg.JWT.RefreshTTLHours
	if req.RememberMe != nil && *req.RememberMe {
		refreshTTL = s.cfg.JWT.RememberMeTTLHours
	}

	accessToken, err := helper.GenerateAccessToken(
		user.ID, user.Email, user.RoleID, user.Name, user.Role.Name,
		s.cfg.JWT.Secret, accessTTL,
	)
	if err != nil {
		return nil, "", 0, helper.NewServiceError("SERVER_ERROR", "Gagal membuat token.", err)
	}

	refreshRaw, refreshExpiry, err := helper.GenerateRefreshToken(refreshTTL)
	if err != nil {
		return nil, "", 0, helper.NewServiceError("SERVER_ERROR", "Gagal membuat refresh token.", err)
	}

	tokenHash := helper.HashToken(refreshRaw, s.cfg.JWT.Secret)
	refreshToken := &model.RefreshToken{
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: refreshExpiry,
	}

	if err := s.repo.SaveRefreshToken(refreshToken); err != nil {
		return nil, "", 0, helper.NewServiceError("SERVER_ERROR", "Gagal menyimpan refresh token.", err)
	}

	maxAge := refreshTTL * 3600

	resp := &dto.AuthResponse{
		AccessToken: accessToken,
		User: dto.AuthUserResponse{
			ID:                 user.ID,
			Name:               user.Name,
			Email:              user.Email,
			Status:             user.Status,
			MustChangePassword: user.MustChangePassword,
			Role: dto.RoleInfo{
				ID:          user.Role.ID,
				Name:        user.Role.Name,
				DisplayName: user.Role.DisplayName,
			},
		},
	}

	return resp, refreshRaw, maxAge, nil
}

func (s *authService) Refresh(refreshTokenStr string) (*dto.RefreshTokenResponse, string, int, error) {
	if refreshTokenStr == "" {
		return nil, "", 0, helper.NewServiceError("AUTH_SESSION_EXPIRED", "Sesi telah berakhir. Silakan login kembali.", nil)
	}

	tokenHash := helper.HashToken(refreshTokenStr, s.cfg.JWT.Secret)

	storedToken, err := s.repo.FindRefreshTokenByHash(tokenHash)
	if err != nil {
		return nil, "", 0, helper.NewServiceError("AUTH_SESSION_EXPIRED", "Sesi telah berakhir. Silakan login kembali.", nil)
	}

	user, err := s.repo.FindByID(storedToken.UserID)
	if err != nil {
		return nil, "", 0, helper.NewServiceError("AUTH_SESSION_EXPIRED", "Sesi telah berakhir. Silakan login kembali.", nil)
	}

	if user.Status == "inactive" {
		return nil, "", 0, helper.NewServiceError("AUTH_ACCOUNT_INACTIVE", "Akun Anda telah dinonaktifkan. Silakan hubungi Admin.", nil)
	}

	_ = s.repo.DeleteRefreshToken(storedToken.ID)

	accessTTL := s.cfg.JWT.AccessTTLMinutes
	refreshTTL := s.cfg.JWT.RefreshTTLHours

	accessToken, err := helper.GenerateAccessToken(
		user.ID, user.Email, user.RoleID, user.Name, user.Role.Name,
		s.cfg.JWT.Secret, accessTTL,
	)
	if err != nil {
		return nil, "", 0, helper.NewServiceError("SERVER_ERROR", "Gagal membuat token.", err)
	}

	refreshRaw, refreshExpiry, err := helper.GenerateRefreshToken(refreshTTL)
	if err != nil {
		return nil, "", 0, helper.NewServiceError("SERVER_ERROR", "Gagal membuat refresh token.", err)
	}

	newHash := helper.HashToken(refreshRaw, s.cfg.JWT.Secret)
	newToken := &model.RefreshToken{
		UserID:    user.ID,
		TokenHash: newHash,
		ExpiresAt: refreshExpiry,
	}

	if err := s.repo.SaveRefreshToken(newToken); err != nil {
		return nil, "", 0, helper.NewServiceError("SERVER_ERROR", "Gagal menyimpan refresh token.", err)
	}

	maxAge := refreshTTL * 3600

	resp := &dto.RefreshTokenResponse{
		AccessToken: accessToken,
	}

	return resp, refreshRaw, maxAge, nil
}

func (s *authService) Logout(userID uint, refreshTokenStr string) error {
	if refreshTokenStr != "" {
		tokenHash := helper.HashToken(refreshTokenStr, s.cfg.JWT.Secret)
		storedToken, err := s.repo.FindRefreshTokenByHash(tokenHash)
		if err == nil {
			_ = s.repo.DeleteRefreshToken(storedToken.ID)
		}
	}

	return nil
}

func (s *authService) ForgotPassword(req *dto.ForgotPasswordRequest, appURL string) error {
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return nil
	}

	if user.Status == "inactive" {
		subject := "Pemberitahuan Akun DPP GRADASI"
		body := s.buildInactiveAccountEmail(user.Name)
		if err := s.mail.Send(req.Email, subject, body); err != nil {
			return helper.NewServiceError("EMAIL_SEND_FAILED", "Gagal mengirim email. Silakan coba lagi.", err)
		}

		if s.cfg.App.Env == "development" {
			log.Printf("[DEV] Email akun nonaktif akan dikirim ke %s", req.Email)
		}
		return nil
	}

	if user.Status == "pending_activation" {
		subject := "Pemberitahuan Akun DPP GRADASI"
		body := s.buildPendingAccountEmail(user.Name)
		if err := s.mail.Send(req.Email, subject, body); err != nil {
			return helper.NewServiceError("EMAIL_SEND_FAILED", "Gagal mengirim email. Silakan coba lagi.", err)
		}

		if s.cfg.App.Env == "development" {
			log.Printf("[DEV] Email akun belum aktif akan dikirim ke %s", req.Email)
		}
		return nil
	}

	rawToken, tokenHash, expiry, err := helper.GenerateResetToken(
		s.cfg.JWT.Secret,
		s.cfg.JWT.PasswordResetTTLMinutes,
	)
	if err != nil {
		return helper.NewServiceError("SERVER_ERROR", "Gagal membuat token reset.", err)
	}

	resetToken := &model.PasswordResetToken{
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: expiry,
	}

	if err := s.repo.SavePasswordResetToken(resetToken); err != nil {
		return helper.NewServiceError("SERVER_ERROR", "Gagal menyimpan token reset.", err)
	}

	resetURL := appURL + "/reset-password?token=" + rawToken + "&email=" + req.Email
	subject := "Reset Password - DPP GRADASI"
	body := s.buildResetPasswordEmail(user.Name, resetURL)

	if err := s.mail.Send(req.Email, subject, body); err != nil {
		return helper.NewServiceError("EMAIL_SEND_FAILED", "Gagal mengirim email. Silakan coba lagi.", err)
	}

	if s.cfg.App.Env == "development" {
		log.Printf("[DEV] Forgot password link untuk %s: %s", req.Email, resetURL)
		if s.mail != nil {
			log.Printf("[DEV] Email akan dikirim ke %s via %s", req.Email, s.mail.DSN())
		}
	}

	return nil
}

func (s *authService) ValidateResetToken(rawToken string) error {
	tokenHash := helper.HashToken(rawToken, s.cfg.JWT.Secret)
	_, err := s.repo.FindPasswordResetTokenByHash(tokenHash)
	if err != nil {
		return helper.NewServiceError("AUTH_TOKEN_INVALID_OR_EXPIRED", "Token reset password tidak valid atau telah kedaluwarsa.", nil)
	}
	return nil
}

func (s *authService) ResetPassword(req *dto.ResetPasswordRequest) error {
	tokenHash := helper.HashToken(req.Token, s.cfg.JWT.Secret)
	storedToken, err := s.repo.FindPasswordResetTokenByHash(tokenHash)
	if err != nil {
		return helper.NewServiceError("AUTH_TOKEN_INVALID_OR_EXPIRED", "Token reset password tidak valid atau telah kedaluwarsa.", nil)
	}

	hashedPassword, err := helper.HashPassword(req.Password)
	if err != nil {
		return helper.NewServiceError("SERVER_ERROR", "Gagal mengenkripsi password.", err)
	}

	if err := s.repo.SetUserPassword(storedToken.UserID, hashedPassword); err != nil {
		return helper.NewServiceError("SERVER_ERROR", "Gagal menyimpan password baru.", err)
	}

	_ = s.repo.MarkResetTokenUsed(storedToken.ID)
	_ = s.repo.DeleteUserRefreshTokens(storedToken.UserID)

	return nil
}

func (s *authService) ValidateActivationToken(token string) error {
	key := fmt.Sprintf("activation_token:%s", token)
	val, err := s.redis.Get(context.Background(), key).Result()
	if err == nil && val != "" {
		return nil
	}

	tokenHash := helper.HashToken(token, s.cfg.JWT.Secret)
	_, err = s.repo.FindActivationTokenByHash(tokenHash)
	if err != nil {
		return helper.NewServiceError("AUTH_TOKEN_INVALID_OR_EXPIRED", "Token aktivasi tidak valid atau telah kedaluwarsa.", nil)
	}
	return nil
}

func (s *authService) ActivateAccount(req *dto.ActivateAccountRequest) error {
	var userID uint
	key := fmt.Sprintf("activation_token:%s", req.Token)
	val, err := s.redis.Get(context.Background(), key).Result()
	if err == nil && val != "" {
		if _, scanErr := fmt.Sscanf(val, "%d", &userID); scanErr != nil {
			return helper.NewServiceError("AUTH_TOKEN_INVALID_OR_EXPIRED", "Token aktivasi tidak valid atau telah kedaluwarsa.", nil)
		}
	}

	if userID == 0 {
		tokenHash := helper.HashToken(req.Token, s.cfg.JWT.Secret)
		storedToken, err := s.repo.FindActivationTokenByHash(tokenHash)
		if err != nil {
			return helper.NewServiceError("AUTH_TOKEN_INVALID_OR_EXPIRED", "Token aktivasi tidak valid atau telah kedaluwarsa.", nil)
		}
		userID = storedToken.UserID
		_ = s.repo.MarkActivationTokenUsed(storedToken.ID)
	} else {
		s.redis.Del(context.Background(), key)
	}

	hashedPassword, err := helper.HashPassword(req.Password)
	if err != nil {
		return helper.NewServiceError("SERVER_ERROR", "Gagal mengenkripsi password.", err)
	}

	if err := s.repo.SetUserPassword(userID, hashedPassword); err != nil {
		return helper.NewServiceError("SERVER_ERROR", "Gagal menyimpan password baru.", err)
	}

	if err := s.repo.ActivateUser(userID); err != nil {
		return helper.NewServiceError("SERVER_ERROR", "Gagal mengaktifkan akun.", err)
	}

	return nil
}

func (s *authService) ChangePassword(userID uint, req *dto.ChangePasswordRequest) error {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return helper.NewServiceError("NOT_FOUND", "User tidak ditemukan.", err)
	}

	if !helper.CheckPassword(req.CurrentPassword, user.Password) {
		return helper.NewServiceError("AUTH_INVALID_CREDENTIALS", "Password saat ini salah.", nil)
	}

	hashedPassword, err := helper.HashPassword(req.Password)
	if err != nil {
		return helper.NewServiceError("SERVER_ERROR", "Gagal mengenkripsi password.", err)
	}

	if err := s.repo.SetUserPassword(userID, hashedPassword); err != nil {
		return helper.NewServiceError("SERVER_ERROR", "Gagal menyimpan password baru.", err)
	}

	_ = s.repo.DeleteUserRefreshTokens(userID)

	return nil
}

func (s *authService) GetProfile(userID uint) (*dto.AuthUserResponse, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, helper.NewServiceError("NOT_FOUND", "User tidak ditemukan.", err)
	}

	return &dto.AuthUserResponse{
		ID:                 user.ID,
		Name:               user.Name,
		Email:              user.Email,
		PhotoPath:          user.PhotoPath,
		Status:             user.Status,
		MustChangePassword: user.MustChangePassword,
		Role: dto.RoleInfo{
			ID:          user.Role.ID,
			Name:        user.Role.Name,
			DisplayName: user.Role.DisplayName,
		},
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}, nil
}

func (s *authService) buildResetPasswordEmail(name, resetURL string) string {
	return `<!DOCTYPE html>
<html><head><meta charset="UTF-8"></head>
<body style="font-family: Arial, sans-serif; background: #f4f4f4; padding: 20px;">
<div style="max-width: 600px; margin: auto; background: white; border-radius: 12px; padding: 40px; border: 1px solid #e2e8f0;">
<h2 style="color: #1e3a8a;">Halo ` + name + `,</h2>
<p style="color: #475569; line-height: 1.6;">Kami menerima permintaan reset password untuk akun DPP GRADASI Anda.</p>
<p style="color: #475569; line-height: 1.6;">Klik tombol di bawah untuk mereset password Anda. Tautan ini berlaku 15 menit.</p>
<div style="text-align: center; margin: 32px 0;">
<a href="` + resetURL + `" style="background: #2563eb; color: white; text-decoration: none; padding: 14px 32px; border-radius: 10px; display: inline-block; font-weight: bold; font-size: 15px;">Reset Password</a>
</div>
<p style="color: #94a3b8; font-size: 13px;">Jika Anda tidak meminta reset password, abaikan email ini.</p>
<hr style="border: none; border-top: 1px solid #e2e8f0; margin: 24px 0;">
<p style="color: #94a3b8; font-size: 12px;">DPP GRADASI — Generasi Digital Indonesia</p>
</div></body></html>`
}

func (s *authService) buildInactiveAccountEmail(name string) string {
	return `<!DOCTYPE html>
<html><head><meta charset="UTF-8"></head>
<body style="font-family: Arial, sans-serif; background: #f4f4f4; padding: 20px;">
<div style="max-width: 600px; margin: auto; background: white; border-radius: 12px; padding: 40px; border: 1px solid #e2e8f0;">
<h2 style="color: #1e3a8a;">Halo ` + name + `,</h2>
<p style="color: #475569; line-height: 1.6;">Anda baru saja melakukan permintaan untuk mereset password. Namun, kami mendeteksi bahwa akun Anda saat ini berstatus <strong>Nonaktif</strong>.</p>
<p style="color: #475569; line-height: 1.6;">Selama akun dalam status nonaktif, proses reset password tidak dapat dilanjutkan. Silakan hubungi <strong>Administrator Sistem</strong> atau rekan Anda yang memiliki akses Admin untuk mengaktifkan kembali akun Anda.</p>
<p style="color: #94a3b8; font-size: 13px; margin-top: 32px;">Jika Anda tidak merasa melakukan permintaan ini, silakan abaikan email ini.</p>
<hr style="border: none; border-top: 1px solid #e2e8f0; margin: 24px 0;">
<p style="color: #94a3b8; font-size: 12px;">DPP GRADASI — Generasi Digital Indonesia</p>
</div></body></html>`
}

func (s *authService) buildPendingAccountEmail(name string) string {
	return `<!DOCTYPE html>
<html><head><meta charset="UTF-8"></head>
<body style="font-family: Arial, sans-serif; background: #f4f4f4; padding: 20px;">
<div style="max-width: 600px; margin: auto; background: white; border-radius: 12px; padding: 40px; border: 1px solid #e2e8f0;">
<h2 style="color: #1e3a8a;">Halo ` + name + `,</h2>
<p style="color: #475569; line-height: 1.6;">Anda baru saja melakukan permintaan untuk mereset password. Namun, kami mendeteksi bahwa akun Anda belum <strong>diaktifkan</strong>.</p>
<p style="color: #475569; line-height: 1.6;">Sebelum dapat menggunakan fitur reset password, akun Anda harus diaktifkan terlebih dahulu. Silakan hubungi <strong>Administrator Sistem</strong> atau rekan Anda yang memiliki akses Admin untuk mengirimkan kredensial aktivasi.</p>
<p style="color: #94a3b8; font-size: 13px; margin-top: 32px;">Jika Anda tidak merasa melakukan permintaan ini, silakan abaikan email ini.</p>
<hr style="border: none; border-top: 1px solid #e2e8f0; margin: 24px 0;">
<p style="color: #94a3b8; font-size: 12px;">DPP GRADASI — Generasi Digital Indonesia</p>
</div></body></html>`
}
