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
	authrepo "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/auth/repository"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/user/dto"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/user/entity"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/user/mapper"
	userrepo "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/user/repository"
	emailtemplate "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/template_message/email"
	"gorm.io/gorm"
)

type UserService interface {
	GetAll(ctx context.Context, role string) ([]entity.User, error)
	GetPaginated(ctx context.Context, req dto.UserQueryReq) ([]entity.User, int, error)
	Create(ctx context.Context, req dto.UserCreateReq) (*entity.User, error)
	Update(ctx context.Context, id uint, req dto.UserUpdateReq) (*entity.User, error)
	Delete(ctx context.Context, id uint) error
	ActivateAccount(ctx context.Context, token string, password string) (*entity.User, error)
	SaveRefreshToken(ctx context.Context, userID uint, token string, expiresAt time.Time) error
	ToggleStatus(ctx context.Context, id uint) error
	ResendNotification(ctx context.Context, id uint) error
	BulkResendNotification(ctx context.Context, ids []uint) (*BulkResendResult, error)
	GetByID(ctx context.Context, id uint) (*entity.User, error)
	BulkDelete(ctx context.Context, ids []uint) error
	Restore(ctx context.Context, id uint) error
	BulkRestore(ctx context.Context, ids []uint) error
	GetDependencyInfo(ctx context.Context, id uint) (map[string]interface{}, error)
	CheckUnique(ctx context.Context, field string, value string, excludeID uint) (bool, error)
	UpdateProfile(ctx context.Context, id uint, name string, email string) (*entity.User, error)
	GetProfile(ctx context.Context, userID uint) (*entity.User, error)
	ChangePassword(ctx context.Context, userID uint, oldPassword string, newPassword string) error
	VerifyEmail(ctx context.Context, userID uint, token string) error
}

type userService struct {
	db       *gorm.DB
	repo     userrepo.UserRepo
	authRepo authrepo.AuthRepo
	audit    activitylogservice.ActivityLogService
	mailer   *infrastructure.Mailer
	cfg      *config.Config
}

func NewUserService(db *gorm.DB, repo userrepo.UserRepo, authRepo authrepo.AuthRepo, audit activitylogservice.ActivityLogService, mailer *infrastructure.Mailer, cfg *config.Config) UserService {
	return &userService{
		db:       db,
		repo:     repo,
		authRepo: authRepo,
		audit:    audit,
		mailer:   mailer,
		cfg:      cfg,
	}
}

func (s *userService) log(ctx context.Context, db *gorm.DB, input *activitylogdto.ActivityLogInput) {
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

func (s *userService) GetAll(ctx context.Context, role string) ([]entity.User, error) {
	return s.repo.FindAll(ctx, role)
}

func (s *userService) GetByID(ctx context.Context, id uint) (*entity.User, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, helper.NewNotFoundError("Pengguna tidak ditemukan")
	}
	return user, nil
}

func (s *userService) GetPaginated(ctx context.Context, req dto.UserQueryReq) ([]entity.User, int, error) {
	req.Normalize()
	return s.repo.FindPaginated(ctx, req.Page, req.Limit, req.Search, req.Role, req.Status, req.Sort, req.Trashed)
}

func (s *userService) Create(ctx context.Context, req dto.UserCreateReq) (*entity.User, error) {
	user := mapper.CreateReqToEntity(&req)
	user.Name = strings.TrimSpace(user.Name)
	user.Email = strings.ToLower(strings.TrimSpace(user.Email))
	if user.Status == "" {
		user.Status = "inactive"
	}

	// Validasi email unik
	if existing, _ := s.repo.FindByEmail(ctx, user.Email); existing != nil {
		v := helper.NewValidationError()
		v.Add("email", fmt.Sprintf("Email '%s' sudah terdaftar", user.Email))
		return nil, v
	}

	if err := s.repo.CreateTx(ctx, user); err != nil {
		return nil, err
	}

	// Kirim email aktivasi (langsung, sync via template email)
	if err := s.sendActivationEmail(ctx, user); err != nil {
		// Email gagal bukan berarti user gagal dibuat — log saja.
		_ = err
	}

	s.log(ctx, s.db, &activitylogdto.ActivityLogInput{
		ActorID:     &user.ID,
		ActorName:   user.Name,
		ActorRole:   user.Role,
		Action:      "users.create",
		EntityType:  "users",
		EntityID:    &user.ID,
		EntityLabel: user.Name,
		Description: fmt.Sprintf("Membuat pengguna baru: %s", user.Name),
		Metadata: map[string]any{
			"email": user.Email,
		},
	})

	return user, nil
}

func (s *userService) Update(ctx context.Context, id uint, req dto.UserUpdateReq) (*entity.User, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, helper.NewNotFoundError("Pengguna tidak ditemukan")
	}

	// Cek email unik (kecuali dirinya sendiri)
	if strings.ToLower(strings.TrimSpace(req.Email)) != user.Email {
		if other, _ := s.repo.FindByEmail(ctx, strings.ToLower(strings.TrimSpace(req.Email))); other != nil {
			v := helper.NewValidationError()
			v.Add("email", fmt.Sprintf("Email '%s' sudah digunakan oleh pengguna lain", req.Email))
			return nil, v
		}
	}

	oldName := user.Name
	oldEmail := user.Email
	oldStatus := user.Status

	mapper.UpdateReqToEntity(&req, user)
	user.Name = strings.TrimSpace(user.Name)
	user.Email = strings.ToLower(strings.TrimSpace(user.Email))

	if err := s.repo.UpdateTx(ctx, user); err != nil {
		return nil, err
	}

	s.log(ctx, s.db, &activitylogdto.ActivityLogInput{
		ActorID:     &user.ID,
		ActorName:   user.Name,
		ActorRole:   user.Role,
		Action:      "users.update",
		EntityType:  "users",
		EntityID:    &user.ID,
		EntityLabel: user.Name,
		Description: fmt.Sprintf("Memperbarui pengguna: %s", user.Name),
		Metadata: map[string]any{
			"old_name":   oldName,
			"old_email":  oldEmail,
			"old_status": oldStatus,
			"new_name":   user.Name,
			"new_email":  user.Email,
			"new_status": user.Status,
		},
	})

	return user, nil
}

func (s *userService) Delete(ctx context.Context, id uint) error {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return helper.NewNotFoundError("Pengguna tidak ditemukan")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	s.log(ctx, s.db, &activitylogdto.ActivityLogInput{
		ActorID:     &user.ID,
		ActorName:   user.Name,
		ActorRole:   user.Role,
		Action:      "users.delete",
		EntityType:  "users",
		EntityID:    &user.ID,
		EntityLabel: user.Name,
		Description: fmt.Sprintf("Menghapus pengguna (soft delete): %s", user.Name),
		Metadata: map[string]any{
			"email": user.Email,
		},
	})

	return nil
}

func (s *userService) ToggleStatus(ctx context.Context, id uint) error {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return helper.NewNotFoundError("Pengguna tidak ditemukan")
	}

	if err := s.repo.ToggleStatus(ctx, id); err != nil {
		return err
	}

	// Status baru = kebalikan dari sebelumnya
	newStatus := "active"
	if user.Status == "active" {
		newStatus = "inactive"
	}

	s.log(ctx, s.db, &activitylogdto.ActivityLogInput{
		ActorID:     &user.ID,
		ActorName:   user.Name,
		ActorRole:   user.Role,
		Action:      "users.toggle_status",
		EntityType:  "users",
		EntityID:    &user.ID,
		EntityLabel: user.Name,
		Description: fmt.Sprintf("Mengubah status pengguna menjadi %s: %s", newStatus, user.Name),
		Metadata: map[string]any{
			"email":      user.Email,
			"old_status": user.Status,
			"new_status": newStatus,
		},
	})

	return nil
}

func (s *userService) BulkDelete(ctx context.Context, ids []uint) error {
	if err := s.repo.BulkDelete(ctx, ids); err != nil {
		return err
	}

	s.log(ctx, s.db, &activitylogdto.ActivityLogInput{
		Action:      "users.bulk_delete",
		EntityType:  "users",
		Description: fmt.Sprintf("Menghapus %d pengguna (soft delete)", len(ids)),
		Metadata: map[string]any{
			"ids": ids,
		},
	})

	return nil
}

func (s *userService) Restore(ctx context.Context, id uint) error {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return helper.NewNotFoundError("Pengguna tidak ditemukan")
	}

	if err := s.repo.Restore(ctx, id); err != nil {
		return err
	}

	s.log(ctx, s.db, &activitylogdto.ActivityLogInput{
		ActorID:     &user.ID,
		ActorName:   user.Name,
		ActorRole:   user.Role,
		Action:      "users.restore",
		EntityType:  "users",
		EntityID:    &user.ID,
		EntityLabel: user.Name,
		Description: fmt.Sprintf("Memulihkan pengguna: %s", user.Name),
		Metadata: map[string]any{
			"email": user.Email,
		},
	})

	return nil
}

func (s *userService) BulkRestore(ctx context.Context, ids []uint) error {
	if err := s.repo.BulkRestore(ctx, ids); err != nil {
		return err
	}

	s.log(ctx, s.db, &activitylogdto.ActivityLogInput{
		Action:      "users.bulk_restore",
		EntityType:  "users",
		Description: fmt.Sprintf("Memulihkan %d pengguna", len(ids)),
		Metadata: map[string]any{
			"ids": ids,
		},
	})

	return nil
}

func (s *userService) sendActivationEmail(ctx context.Context, user *entity.User) error {
	if s.mailer == nil {
		return nil
	}

	rawToken, _, expiresAt, err := helper.GenerateActivationToken(s.cfg.JWT.Secret, 72)
	if err != nil {
		return err
	}

	if err := s.authRepo.SaveAuthToken(ctx, user.ID, rawToken, "activation", expiresAt); err != nil {
		return err
	}

	link := fmt.Sprintf("%s/reset-password?token=%s",
		strings.TrimSuffix(s.cfg.App.FrontendURL, "/"),
		rawToken,
	)

	html, err := emailtemplate.Render("account_activation.html", map[string]any{
		"Name": user.Name,
		"URL":  link,
	})
	if err != nil {
		return fmt.Errorf("gagal merender template email: %w", err)
	}

	return s.mailer.Send(user.Email, "Aktivasi Akun", html)
}

func (s *userService) ResendNotification(ctx context.Context, id uint) error {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return helper.NewNotFoundError("Pengguna tidak ditemukan")
	}

	if err := s.sendActivationEmail(ctx, user); err != nil {
		return err
	}

	s.log(ctx, s.db, &activitylogdto.ActivityLogInput{
		ActorID:     &user.ID,
		ActorName:   user.Name,
		ActorRole:   user.Role,
		Action:      "users.resend_activation",
		EntityType:  "users",
		EntityID:    &user.ID,
		EntityLabel: user.Name,
		Description: fmt.Sprintf("Mengirim ulang email aktivasi: %s", user.Name),
		Metadata: map[string]any{
			"email": user.Email,
		},
	})

	return nil
}

type BulkResendResult struct {
	Total  int      `json:"total"`
	Sent   int      `json:"sent"`
	Failed int      `json:"failed"`
	Errors []string `json:"errors"`
}

func (s *userService) BulkResendNotification(ctx context.Context, ids []uint) (*BulkResendResult, error) {
	result := &BulkResendResult{
		Total:  len(ids),
		Errors: []string{},
	}

	for _, id := range ids {
		user, err := s.repo.FindByID(ctx, id)
		if err != nil || user == nil {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("ID %d: Data tidak ditemukan", id))
			continue
		}
		if userHasPassword(user) {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("%s: Akun sudah aktif dan sudah memiliki password", user.Name))
			continue
		}

		if err := s.sendActivationEmail(ctx, user); err != nil {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("%s: %s", user.Name, err.Error()))
			continue
		}
		result.Sent++
	}

	return result, nil
}

func userHasPassword(user *entity.User) bool {
	return user != nil && strings.TrimSpace(user.Password) != ""
}

// ActivateAccount mengaktifkan akun via token aktivasi.
func (s *userService) ActivateAccount(ctx context.Context, token string, password string) (*entity.User, error) {
	userID, err := s.authRepo.FindAuthToken(ctx, token, "activation")
	if err != nil {
		return nil, &helper.AuthenticationError{Message: "Token aktivasi tidak valid atau telah kedaluwarsa.", Code: "AUTH_TOKEN_INVALID_OR_EXPIRED"}
	}

	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, helper.NewNotFoundError("Pengguna tidak ditemukan")
	}

	hashed, err := helper.HashPassword(password)
	if err != nil {
		return nil, errors.New("gagal memproses password baru")
	}

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.repo.UpdatePassword(ctx, userID, hashed); err != nil {
			return err
		}
		if err := s.repo.Activate(ctx, userID); err != nil {
			return err
		}
		return s.authRepo.DeleteAuthToken(ctx, token, "activation")
	})
	if err != nil {
		return nil, errors.New("gagal mengaktifkan akun")
	}

	s.log(ctx, s.db, &activitylogdto.ActivityLogInput{
		ActorID:     &user.ID,
		ActorName:   user.Name,
		ActorRole:   user.Role,
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

func (s *userService) SaveRefreshToken(ctx context.Context, userID uint, token string, expiresAt time.Time) error {
	return s.authRepo.SaveRefreshToken(ctx, userID, token, expiresAt, "", "", "")
}

func (s *userService) GetDependencyInfo(ctx context.Context, id uint) (map[string]interface{}, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, helper.NewNotFoundError("Pengguna tidak ditemukan")
	}

	return map[string]interface{}{
		"user_id":   user.ID,
		"name":      user.Name,
		"email":     user.Email,
		"has_child": false,
	}, nil
}

func (s *userService) CheckUnique(ctx context.Context, field string, value string, excludeID uint) (bool, error) {
	field = strings.ToLower(strings.TrimSpace(field))
	value = strings.TrimSpace(value)

	switch field {
	case "email":
		value = strings.ToLower(value)
		existing, err := s.repo.FindByEmail(ctx, value)
		if err != nil {
			return true, nil
		}
		if existing != nil && existing.ID != excludeID {
			return false, nil
		}
		return true, nil
	default:
		return true, nil
	}
}

func (s *userService) UpdateProfile(ctx context.Context, id uint, name string, email string) (*entity.User, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, helper.NewNotFoundError("Pengguna tidak ditemukan")
	}

	oldName := user.Name
	oldEmail := user.Email
	user.Name = strings.TrimSpace(name)

	emailChanged := false
	newEmail := strings.ToLower(strings.TrimSpace(email))
	if newEmail != "" && newEmail != user.Email {
		// Check email uniqueness
		if existing, _ := s.repo.FindByEmail(ctx, newEmail); existing != nil {
			v := helper.NewValidationError()
			v.Add("email", fmt.Sprintf("Email '%s' sudah terdaftar", newEmail))
			return nil, v
		}
		emailChanged = true
		user.Email = newEmail
		user.EmailVerifiedAt = nil // Mark as unverified until verified
	}

	if err := s.repo.UpdateTx(ctx, user); err != nil {
		return nil, err
	}

	if emailChanged {
		// Trigger activation/verification email send
		_ = s.sendActivationEmail(ctx, user)
	}

	s.log(ctx, s.db, &activitylogdto.ActivityLogInput{
		ActorID:     &user.ID,
		ActorName:   user.Name,
		ActorRole:   user.Role,
		Action:      "users.update_profile",
		EntityType:  "users",
		EntityID:    &user.ID,
		EntityLabel: user.Name,
		Description: fmt.Sprintf("Memperbarui profil: %s", user.Name),
		Metadata: map[string]any{
			"old_name":      oldName,
			"new_name":      user.Name,
			"old_email":     oldEmail,
			"new_email":     user.Email,
			"email_changed": emailChanged,
		},
	})

	return user, nil
}

// GetProfile mengembalikan profil user berdasarkan ID dari token login.
func (s *userService) GetProfile(ctx context.Context, userID uint) (*entity.User, error) {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, helper.NewNotFoundError("Pengguna tidak ditemukan")
	}
	return user, nil
}

// ChangePassword mengganti password akun sendiri. Wajib old_password benar.
// Setelah sukses: must_change_password di-reset false + semua refresh token di-revoke.
func (s *userService) ChangePassword(ctx context.Context, userID uint, oldPassword string, newPassword string) error {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil || user == nil {
		return helper.NewNotFoundError("Pengguna tidak ditemukan")
	}

	if !helper.CheckPassword(oldPassword, user.Password) {
		return &helper.AuthenticationError{Message: "Password lama salah", Code: "INVALID_PASSWORD"}
	}

	hashed, err := helper.HashPassword(newPassword)
	if err != nil {
		return errors.New("gagal memproses password baru")
	}

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.repo.UpdatePassword(ctx, userID, hashed); err != nil {
			return err
		}
		// Reset flag must_change_password + revoke semua refresh token (logout semua perangkat)
		if err := tx.Model(&entity.User{}).Where("id = ?", userID).Update("must_change_password", false).Error; err != nil {
			return err
		}
		return s.authRepo.DeleteAllUserRefreshTokens(ctx, userID)
	})
	if err != nil {
		return errors.New("gagal mengubah password")
	}

	s.log(ctx, s.db, &activitylogdto.ActivityLogInput{
		ActorID:     &userID,
		ActorName:   user.Name,
		ActorRole:   user.Role,
		Action:      "users.change_password",
		EntityType:  "users",
		EntityID:    &userID,
		EntityLabel: user.Name,
		Description: fmt.Sprintf("Mengubah password akun: %s", user.Name),
	})

	return nil
}

// VerifyEmail memverifikasi token aktivasi email untuk user yang sedang login.
// Menggunakan tabel activation_tokens (project Redis-free), bukan Redis.
func (s *userService) VerifyEmail(ctx context.Context, userID uint, token string) error {
	tokenUserID, err := s.authRepo.FindAuthToken(ctx, token, "activation")
	if err != nil {
		return &helper.AuthenticationError{Message: "Token verifikasi tidak valid atau telah kedaluwarsa.", Code: "AUTH_TOKEN_INVALID_OR_EXPIRED"}
	}
	if tokenUserID != userID {
		return &helper.AuthenticationError{Message: "Token verifikasi tidak valid atau telah kedaluwarsa.", Code: "AUTH_TOKEN_INVALID_OR_EXPIRED"}
	}

	now := time.Now()
	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&entity.User{}).Where("id = ?", userID).Update("email_verified_at", now).Error; err != nil {
			return err
		}
		return s.authRepo.DeleteAuthToken(ctx, token, "activation")
	}); err != nil {
		return errors.New("gagal memverifikasi email")
	}

	user, _ := s.repo.FindByID(ctx, userID)
	if user != nil {
		s.log(ctx, s.db, &activitylogdto.ActivityLogInput{
			ActorID:     &userID,
			ActorName:   user.Name,
			ActorRole:   user.Role,
			Action:      "users.verify_email",
			EntityType:  "users",
			EntityID:    &userID,
			EntityLabel: user.Name,
			Description: fmt.Sprintf("Email diverifikasi: %s", user.Name),
		})
	}

	return nil
}
