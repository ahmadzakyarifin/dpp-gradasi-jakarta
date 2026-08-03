package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ahmadzakyarifin/dpp-gradasi/backend/config"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/helper"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/infrastructure"
	activitylogdto "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/activitylog/dto"
	activitylogservice "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/activitylog/service"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/auth/dto"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/auth/mapper"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/auth/repository"
	userentity "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/user/entity"
	emailtemplate "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/template_message/email"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type authService struct {
	r      repository.AuthRepo
	audit  activitylogservice.ActivityLogService
	mailer *infrastructure.Mailer
	cfg    *config.Config
}

func NewAuthService(
	repo repository.AuthRepo,
	audit activitylogservice.ActivityLogService,
	mailer *infrastructure.Mailer,
	cfg *config.Config,
) AuthService {
	return &authService{
		r:      repo,
		audit:  audit,
		mailer: mailer,
		cfg:    cfg,
	}
}

func (s *authService) log(ctx context.Context, db *gorm.DB, input *activitylogdto.ActivityLogInput) {
	if s.audit == nil {
		return
	}
	userID, userName, role, ipAddress, userAgent := helper.GetAuditMeta(ctx)
	if input.ActorID == nil && userID > 0 {
		input.ActorID = &userID
	}
	if input.ActorName == "" {
		input.ActorName = userName
	}
	if input.ActorRole == "" {
		input.ActorRole = role
	}
	if input.IPAddress == "" {
		input.IPAddress = ipAddress
	}
	if input.UserAgent == "" {
		input.UserAgent = userAgent
	}

	_ = s.audit.Log(ctx, db, input)
}

func (s *authService) Login(ctx context.Context, req dto.LoginRequest, ip string, userAgent string, device string) (*dto.LoginResponse, error) {
	user, err := s.r.FindByEmail(ctx, req.Email)
	if err != nil {
		s.log(ctx, s.r.GetDB(), &activitylogdto.ActivityLogInput{
			Action:      "auth.login_failed",
			EntityType:  "auth",
			Description: "Login gagal: email tidak ditemukan",
			IPAddress:   ip,
			UserAgent:   userAgent,
			Metadata: map[string]any{
				"email":  req.Email,
				"reason": "email_not_found",
			},
		})
		return nil, &helper.AuthenticationError{Message: "Email atau password salah.", Code: "AUTH_INVALID_CREDENTIALS"}
	}

	if user.Status != "active" {
		s.log(ctx, s.r.GetDB(), &activitylogdto.ActivityLogInput{
			ActorID:     &user.ID,
			ActorName:   user.Name,
			ActorRole:   user.RoleName,
			Action:      "auth.login_failed",
			EntityType:  "user",
			EntityID:    &user.ID,
			EntityLabel: user.Name,
			Description: "Login gagal: akun tidak aktif",
			IPAddress:   ip,
			UserAgent:   userAgent,
			Metadata: map[string]any{
				"email":  user.Email,
				"status": user.Status,
				"reason": "account_inactive",
			},
		})
		return nil, &helper.AuthenticationError{Message: "Akun Anda belum aktif atau telah dinonaktifkan. Silakan hubungi Admin.", Code: "AUTH_ACCOUNT_INACTIVE"}
	}

	if !helper.CheckPassword(req.Password, user.Password) {
		s.log(ctx, s.r.GetDB(), &activitylogdto.ActivityLogInput{
			ActorID:     &user.ID,
			ActorName:   user.Name,
			ActorRole:   user.RoleName,
			Action:      "auth.login_failed",
			EntityType:  "user",
			EntityID:    &user.ID,
			EntityLabel: user.Name,
			Description: "Login gagal: password salah",
			IPAddress:   ip,
			UserAgent:   userAgent,
			Metadata: map[string]any{
				"email":  user.Email,
				"reason": "wrong_password",
			},
		})
		return nil, &helper.AuthenticationError{Message: "Email atau password salah.", Code: "AUTH_INVALID_CREDENTIALS"}
	}

	accessToken, err := helper.GenerateAccessToken(user.ID, user.Email, user.RoleID, user.Name, user.RoleName, s.cfg.JWT.Secret, s.cfg.JWT.AccessTTLMinutes)
	if err != nil {
		return nil, errors.New("gagal membuat access token")
	}

	ttlRefresh := s.cfg.JWT.RefreshTTLHours
	if req.RememberMe {
		ttlRefresh = s.cfg.JWT.RememberMeTTLHours
	}

	refreshToken, expiry, err := helper.GenerateRefreshToken(ttlRefresh)
	if err != nil {
		return nil, errors.New("gagal membuat refresh token")
	}

	if err := s.r.SaveRefreshToken(ctx, user.ID, refreshToken, expiry, ip, userAgent, device); err != nil {
		return nil, errors.New("gagal menyimpan session")
	}

	s.log(ctx, s.r.GetDB(), &activitylogdto.ActivityLogInput{
		ActorID:     &user.ID,
		ActorName:   user.Name,
		ActorRole:   user.RoleName,
		Action:      "auth.login",
		EntityType:  "user",
		EntityID:    &user.ID,
		EntityLabel: user.Name,
		Description: "User berhasil login",
		IPAddress:   ip,
		UserAgent:   userAgent,
		Metadata: map[string]any{
			"email":       user.Email,
			"remember_me": req.RememberMe,
		},
	})

	return &dto.LoginResponse{
		AccessToken:        accessToken,
		RefreshToken:       refreshToken,
		RefreshTokenExpiry: expiry,
		User:               mapper.UserEntityToAuth(*user),
	}, nil
}

func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*dto.LoginResponse, error) {
	user, _, err := s.r.FindUserByRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, &helper.AuthenticationError{Message: "Sesi telah berakhir. Silakan login kembali.", Code: "AUTH_SESSION_EXPIRED"}
	}

	if user.Status != "active" {
		_ = s.r.DeleteRefreshToken(ctx, refreshToken)
		return nil, &helper.AuthenticationError{Message: "Akun Anda telah dinonaktifkan. Silakan hubungi Admin.", Code: "AUTH_ACCOUNT_INACTIVE"}
	}

	// Rotate refresh token: revoke lama → buat baru
	_ = s.r.DeleteRefreshToken(ctx, refreshToken)

	ttlRefresh := s.cfg.JWT.RefreshTTLHours
	newToken, newExpiry, err := helper.GenerateRefreshToken(ttlRefresh)
	if err != nil {
		return nil, errors.New("gagal membuat refresh token baru")
	}

	if err := s.r.SaveRefreshToken(ctx, user.ID, newToken, newExpiry, "", "", ""); err != nil {
		return nil, errors.New("gagal menyimpan session baru")
	}

	newAccessToken, err := helper.GenerateAccessToken(user.ID, user.Email, user.RoleID, user.Name, user.RoleName, s.cfg.JWT.Secret, s.cfg.JWT.AccessTTLMinutes)
	if err != nil {
		return nil, errors.New("gagal membuat access token baru")
	}

	_, _, _, ipAddress, userAgent := helper.GetAuditMeta(ctx)
	s.log(ctx, s.r.GetDB(), &activitylogdto.ActivityLogInput{
		ActorID:     &user.ID,
		ActorName:   user.Name,
		ActorRole:   user.RoleName,
		Action:      "auth.refresh_token",
		EntityType:  "user",
		EntityID:    &user.ID,
		EntityLabel: user.Name,
		Description: "Access token diperbarui, refresh token di-rotate",
		IPAddress:   ipAddress,
		UserAgent:   userAgent,
		Metadata: map[string]any{
			"email": user.Email,
		},
	})

	return &dto.LoginResponse{
		AccessToken:        newAccessToken,
		RefreshToken:       newToken,
		RefreshTokenExpiry: newExpiry,
		User:               mapper.UserEntityToAuth(*user),
	}, nil
}

func (s *authService) ForgotPassword(ctx context.Context, req dto.ForgotPasswordRequest, ip string, userAgent string) (string, error) {
	email := strings.TrimSpace(req.Email)
	user, err := s.r.FindByEmail(ctx, email)

	if err != nil {
		// Tetap return sukses — jangan bocorkan apakah email terdaftar (security best practice)
		return "", nil
	}

	token := uuid.New().String()

	expiryMinutes := s.cfg.JWT.PasswordResetTTLMinutes
	expiresAt := time.Now().Add(time.Duration(expiryMinutes) * time.Minute)

	if err := s.r.SaveAuthToken(ctx, user.ID, token, repository.TokenResetPassword, expiresAt); err != nil {
		return "", errors.New("gagal memproses permintaan reset password")
	}

	link := fmt.Sprintf("%s/reset-password?token=%s", strings.TrimSuffix(s.cfg.App.URL, "/"), token)
	if s.cfg.App.Env != "production" {
		fmt.Printf("\n[DEBUG] Forgot Password Link for %s: %s\n\n", user.Email, link)
	}

	html, err := emailtemplate.Render("forgot_password.html", map[string]any{
		"Name":    user.Name,
		"URL":     link,
		"Expired": expiryMinutes,
	})
	if err != nil {
		return "", fmt.Errorf("gagal merender template email: %w", err)
	}

	if err := s.mailer.Send(user.Email, "Reset Password", html); err != nil {
		return "", errors.New("gagal mengirim email reset password")
	}

	s.log(ctx, s.r.GetDB(), &activitylogdto.ActivityLogInput{
		ActorID:     &user.ID,
		ActorName:   user.Name,
		ActorRole:   user.RoleName,
		Action:      "auth.forgot_password",
		EntityType:  "user",
		EntityID:    &user.ID,
		EntityLabel: user.Name,
		Description: "Link reset password dibuat dan email dikirim",
		IPAddress:   ip,
		UserAgent:   userAgent,
		Metadata: map[string]any{
			"email":      user.Email,
			"expires_in": fmt.Sprintf("%d menit", expiryMinutes),
		},
	})

	return link, nil
}

func (s *authService) ResetPassword(ctx context.Context, req dto.ResetPasswordRequest, ip string, userAgent string) error {
	userID, err := s.r.FindAuthToken(ctx, req.Token, repository.TokenResetPassword)
	if err != nil {
		s.log(ctx, s.r.GetDB(), &activitylogdto.ActivityLogInput{
			Action:      "auth.reset_password_failed",
			EntityType:  "auth",
			Description: "Reset password gagal: token tidak valid atau sudah kadaluarsa",
			IPAddress:   ip,
			UserAgent:   userAgent,
			Metadata: map[string]any{
				"reason": "invalid_or_expired_token",
			},
		})
		return &helper.AuthenticationError{Message: "Token reset password tidak valid atau telah kedaluwarsa.", Code: "AUTH_TOKEN_INVALID_OR_EXPIRED"}
	}

	hashed, err := helper.HashPassword(req.Password)
	if err != nil {
		s.log(ctx, s.r.GetDB(), &activitylogdto.ActivityLogInput{
			ActorID:     &userID,
			Action:      "auth.reset_password_failed",
			EntityType:  "user",
			EntityID:    &userID,
			Description: "Reset password gagal: gagal hash password",
			IPAddress:   ip,
			UserAgent:   userAgent,
			Metadata: map[string]any{
				"reason": "password_hash_failed",
			},
		})
		return errors.New("gagal memproses password baru")
	}

	err = s.r.GetDB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		repoTx := s.r.WithTx(tx)
		if err := repoTx.UpdatePassword(ctx, userID, hashed); err != nil {
			return err
		}
		_ = repoTx.DeleteAuthToken(ctx, req.Token, repository.TokenResetPassword)
		_ = repoTx.DeleteAllUserRefreshTokens(ctx, userID)
		return nil
	})

	if err != nil {
		s.log(ctx, s.r.GetDB(), &activitylogdto.ActivityLogInput{
			ActorID:     &userID,
			Action:      "auth.reset_password_failed",
			EntityType:  "user",
			EntityID:    &userID,
			Description: "Reset password gagal: error saat update database",
			IPAddress:   ip,
			UserAgent:   userAgent,
			Metadata: map[string]any{
				"reason": "database_update_failed",
			},
		})
		return errors.New("gagal memperbarui password")
	}

	s.log(ctx, s.r.GetDB(), &activitylogdto.ActivityLogInput{
		ActorID:     &userID,
		Action:      "auth.reset_password",
		EntityType:  "user",
		EntityID:    &userID,
		Description: "Password berhasil direset",
		IPAddress:   ip,
		UserAgent:   userAgent,
	})

	return nil
}

func (s *authService) ValidateResetToken(ctx context.Context, token string) error {
	_, err := s.r.FindAuthToken(ctx, token, repository.TokenResetPassword)
	if err != nil {
		return &helper.AuthenticationError{Message: "Token reset password tidak valid atau telah kedaluwarsa.", Code: "AUTH_TOKEN_INVALID_OR_EXPIRED"}
	}
	return nil
}

func (s *authService) ChangePassword(ctx context.Context, userID uint, req dto.ChangePasswordRequest) error {
	user, err := s.r.FindUserByID(ctx, userID)
	if err != nil || user == nil {
		return helper.NewNotFoundError("user tidak ditemukan")
	}

	if !helper.CheckPassword(req.CurrentPassword, user.Password) {
		s.log(ctx, s.r.GetDB(), &activitylogdto.ActivityLogInput{
			ActorID:     &user.ID,
			ActorName:   user.Name,
			ActorRole:   user.RoleName,
			Action:      "auth.change_password_failed",
			EntityType:  "user",
			EntityID:    &user.ID,
			EntityLabel: user.Name,
			Description: "Ganti password gagal: password saat ini tidak cocok",
			Metadata: map[string]any{
				"email":  user.Email,
				"reason": "incorrect_current_password",
			},
		})
		return &helper.AuthenticationError{Message: "Password saat ini tidak cocok.", Code: "AUTH_INVALID_CREDENTIALS"}
	}

	hashed, err := helper.HashPassword(req.Password)
	if err != nil {
		s.log(ctx, s.r.GetDB(), &activitylogdto.ActivityLogInput{
			ActorID:     &user.ID,
			ActorName:   user.Name,
			ActorRole:   user.RoleName,
			Action:      "auth.change_password_failed",
			EntityType:  "user",
			EntityID:    &user.ID,
			EntityLabel: user.Name,
			Description: "Ganti password gagal: gagal hash password baru",
			Metadata: map[string]any{
				"email":  user.Email,
				"reason": "password_hash_failed",
			},
		})
		return errors.New("gagal memproses password baru")
	}

	err = s.r.GetDB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		repoTx := s.r.WithTx(tx)
		if err := repoTx.UpdatePassword(ctx, userID, hashed); err != nil {
			return err
		}
		return repoTx.DeleteAllUserRefreshTokens(ctx, userID)
	})

	if err != nil {
		s.log(ctx, s.r.GetDB(), &activitylogdto.ActivityLogInput{
			ActorID:     &user.ID,
			ActorName:   user.Name,
			ActorRole:   user.RoleName,
			Action:      "auth.change_password_failed",
			EntityType:  "user",
			EntityID:    &user.ID,
			EntityLabel: user.Name,
			Description: "Ganti password gagal: error saat update database",
			Metadata: map[string]any{
				"email":  user.Email,
				"reason": "database_update_failed",
			},
		})
		return errors.New("gagal memperbarui password")
	}

	return nil
}

func (s *authService) Logout(ctx context.Context, refreshToken string) error {
	if refreshToken == "" {
		return nil
	}
	_ = s.r.DeleteRefreshToken(ctx, refreshToken)

	user, _, err := s.r.FindUserByRefreshToken(ctx, refreshToken)
	if err == nil && user != nil {
		s.log(ctx, s.r.GetDB(), &activitylogdto.ActivityLogInput{
			ActorID:     &user.ID,
			ActorName:   user.Name,
			ActorRole:   user.RoleName,
			Action:      "auth.logout",
			EntityType:  "user",
			EntityID:    &user.ID,
			EntityLabel: user.Name,
			Description: "User logout",
		})
	}

	return nil
}

// SaveRefreshToken menyimpan refresh token untuk user.
func (s *authService) SaveRefreshToken(ctx context.Context, userID uint, token string, expiresAt time.Time) error {
	return s.r.SaveRefreshToken(ctx, userID, token, expiresAt, "", "", "")
}

// Me mengambil data user + role dari token login.
func (s *authService) Me(ctx context.Context, userID uint) (*userentity.User, error) {
	user, err := s.r.FindUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// ValidateActivationToken memvalidasi token aktivasi akun.
func (s *authService) ValidateActivationToken(ctx context.Context, token string) error {
	if token == "" {
		return &helper.AuthenticationError{Message: "Token aktivasi tidak valid atau telah kedaluwarsa.", Code: "AUTH_TOKEN_INVALID_OR_EXPIRED"}
	}

	_, err := s.r.FindAuthToken(ctx, token, "activation")
	if err != nil {
		return &helper.AuthenticationError{Message: "Token aktivasi tidak valid atau telah kedaluwarsa.", Code: "AUTH_TOKEN_INVALID_OR_EXPIRED"}
	}

	return nil
}

// ActivateAccount mengaktifkan akun via token aktivasi (path kontrak /auth/activate-account).
func (s *authService) ActivateAccount(ctx context.Context, token string, password string) (*userentity.User, error) {
	userID, err := s.r.FindAuthToken(ctx, token, "activation")
	if err != nil {
		return nil, &helper.AuthenticationError{Message: "Token aktivasi tidak valid atau telah kedaluwarsa.", Code: "AUTH_TOKEN_INVALID_OR_EXPIRED"}
	}

	user, err := s.r.FindUserByID(ctx, userID)
	if err != nil {
		return nil, helper.NewNotFoundError("Pengguna tidak ditemukan")
	}

	hashed, err := helper.HashPassword(password)
	if err != nil {
		return nil, errors.New("gagal memproses password baru")
	}

	err = s.r.GetDB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.r.WithTx(tx).UpdatePassword(ctx, userID, hashed); err != nil {
			return err
		}
		if err := tx.Model(&userentity.User{}).Where("id = ?", userID).Updates(map[string]any{
			"status":               "active",
			"email_verified_at":    time.Now(),
			"must_change_password": false,
		}).Error; err != nil {
			return err
		}
		return s.r.WithTx(tx).DeleteAuthToken(ctx, token, "activation")
	})
	if err != nil {
		return nil, errors.New("gagal mengaktifkan akun")
	}

	s.log(ctx, s.r.GetDB(), &activitylogdto.ActivityLogInput{
		ActorID:     &user.ID,
		ActorName:   user.Name,
		ActorRole:   user.RoleName,
		Action:      "users.activate",
		EntityType:  "users",
		EntityID:    &user.ID,
		EntityLabel: user.Name,
		Description: fmt.Sprintf("Akun diaktifkan: %s", user.Name),
		Metadata: map[string]any{
			"email": user.Email,
		},
	})

	return user, nil
}
