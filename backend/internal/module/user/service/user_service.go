package service

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/helper"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/infrastructure"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/user/dto"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/user/model"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/user/repository"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	GetProfile(userID uint) (*dto.UserResponse, error)
	UpdateProfile(userID uint, req *dto.ProfileUpdateRequest) (*dto.UserResponse, string, error)
	VerifyEmail(userID uint, req *dto.VerifyEmailRequest) error
	ChangePassword(userID uint, req *dto.ChangePasswordRequest) error

	// Admin Management
	GetAdmins() ([]dto.UserResponse, error)
	CreateAdmin(req *dto.AdminCreateRequest, appURL string) (*dto.UserResponse, error)
	ResendActivation(adminID uint, targetID uint, appURL string) error
	SetAdminStatus(adminID uint, targetID uint, req *dto.AdminStatusRequest) error
	DeleteAdmin(adminID uint, targetID uint) error
	RestoreAdmin(adminID uint, targetID uint) error
	BulkDeleteAdmin(adminID uint, targetIDs []uint) error
	BulkRestoreAdmin(targetIDs []uint) error
}

type userService struct {
	repo       repository.UserRepo
	redis      *redis.Client
	mailer     *infrastructure.Mailer
	uploadPath string
}

func NewUserService(repo repository.UserRepo, redis *redis.Client, mailer *infrastructure.Mailer) UserService {
	uploadPath := "public/uploads/users"
	os.MkdirAll(uploadPath, 0755)
	return &userService{
		repo:       repo,
		redis:      redis,
		mailer:     mailer,
		uploadPath: uploadPath,
	}
}

func (s *userService) GetProfile(userID uint) (*dto.UserResponse, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, helper.NewServiceError("NOT_FOUND", "User tidak ditemukan", err)
	}
	resp := toResponse(*user)
	return &resp, nil
}

func (s *userService) UpdateProfile(userID uint, req *dto.ProfileUpdateRequest) (*dto.UserResponse, string, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, "", helper.NewServiceError("NOT_FOUND", "User tidak ditemukan", err)
	}

	if req.Photo != nil {
		newPhoto, err := s.handleUpload(req.Photo)
		if err == nil && newPhoto != "" {
			user.PhotoPath = newPhoto
		}
	}

	user.Name = req.Name

	var message string
	// Check if email changed
	if user.Email != req.Email {
		// Check if email is already taken by someone else
		existing, _ := s.repo.FindByEmail(req.Email)
		if existing != nil && existing.ID != userID {
			return nil, "", helper.NewServiceError("EMAIL_TAKEN", "Email sudah terdaftar", nil)
		}

		// Generate token
		token, _ := generateRandomNumber(6)

		// Save to Redis (expire in 1 hour)
		key := fmt.Sprintf("email_verify:%d", userID)
		redisData := map[string]string{
			"email": req.Email,
			"token": token,
		}
		jsonData, _ := json.Marshal(redisData)
		s.redis.Set(context.Background(), key, jsonData, time.Hour)

		// Send Email
		go s.mailer.Send(
			req.Email,
			"Verifikasi Perubahan Email",
			fmt.Sprintf("Kode verifikasi Anda adalah: <b>%s</b>. Kode ini berlaku selama 1 jam.", token),
		)

		message = "Profil berhasil diperbarui. Karena Anda mengubah email, silakan periksa email baru Anda untuk kode verifikasi."
	} else {
		message = "Profil berhasil diperbarui."
	}

	if err := s.repo.Update(user); err != nil {
		return nil, "", helper.NewServiceError("SERVER_ERROR", "Gagal menyimpan profil", err)
	}

	resp := toResponse(*user)
	return &resp, message, nil
}

func (s *userService) VerifyEmail(userID uint, req *dto.VerifyEmailRequest) error {
	key := fmt.Sprintf("email_verify:%d", userID)
	val, err := s.redis.Get(context.Background(), key).Result()
	if err != nil {
		return helper.NewServiceError("INVALID_TOKEN", "Kode verifikasi tidak valid atau kedaluwarsa", err)
	}

	var data map[string]string
	json.Unmarshal([]byte(val), &data)

	if data["token"] != req.Token {
		return helper.NewServiceError("INVALID_TOKEN", "Kode verifikasi salah", nil)
	}

	// Token valid, update email in DB
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return helper.NewServiceError("NOT_FOUND", "User tidak ditemukan", err)
	}

	user.Email = data["email"]
	now := time.Now()
	user.EmailVerifiedAt = &now
	if err := s.repo.Update(user); err != nil {
		return helper.NewServiceError("SERVER_ERROR", "Gagal mengupdate email", err)
	}

	// Delete token
	s.redis.Del(context.Background(), key)
	return nil
}

func (s *userService) ChangePassword(userID uint, req *dto.ChangePasswordRequest) error {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return helper.NewServiceError("NOT_FOUND", "User tidak ditemukan", err)
	}

	// Verify old password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		return helper.NewServiceError("INVALID_PASSWORD", "Password lama salah", err)
	}

	// Hash new password
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return helper.NewServiceError("SERVER_ERROR", "Gagal mengenkripsi password", err)
	}

	user.Password = string(hashed)
	if err := s.repo.Update(user); err != nil {
		return helper.NewServiceError("SERVER_ERROR", "Gagal menyimpan password baru", err)
	}
	return nil
}

func (s *userService) GetAdmins() ([]dto.UserResponse, error) {
	users, err := s.repo.FindAllAdmins()
	if err != nil {
		return nil, helper.NewServiceError("SERVER_ERROR", "Gagal mengambil data admin", err)
	}

	var res []dto.UserResponse
	for _, u := range users {
		res = append(res, toResponse(u))
	}
	return res, nil
}

func (s *userService) CreateAdmin(req *dto.AdminCreateRequest, appURL string) (*dto.UserResponse, error) {
	existing, _ := s.repo.FindByEmail(req.Email)
	if existing != nil {
		return nil, helper.NewServiceError("EMAIL_TAKEN", "Email sudah terdaftar", nil)
	}

	user := &model.User{
		Name:            req.Name,
		Email:           req.Email,
		RoleID:          req.RoleID,
		Password:        "", // Belum ada password sampai aktivasi
		Status:          "pending_activation",
		EmailVerifiedAt: nil,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, helper.NewServiceError("SERVER_ERROR", "Gagal membuat admin", err)
	}

	// Generate activation token
	rawToken, _, _, err := helper.GenerateResetToken("activation", 24*60)
	if err != nil {
		return nil, helper.NewServiceError("SERVER_ERROR", "Gagal membuat token aktivasi", err)
	}

	// Simpan ke Redis (berlaku 24 jam)
	key := fmt.Sprintf("activation_token:%s", rawToken)
	s.redis.Set(context.Background(), key, user.ID, 24*time.Hour)

	// Kirim email undangan aktivasi secara SYNC (sesuai preferensi: response langsung tahu sukses/gagal)
	activationURL := fmt.Sprintf("%s/activate-account?token=%s&email=%s", appURL, rawToken, req.Email)
	subject := "Undangan Aktivasi Akun Pengelola - DPP GRADASI"
	body := s.buildInvitationEmail(user.Name, activationURL)

	if err := s.mailer.Send(req.Email, subject, body); err != nil {
		// Email gagal → hapus user + token agar tidak ada akun pending tanpa undangan
		_ = s.repo.Delete(user.ID)
		s.redis.Del(context.Background(), key)
		return nil, helper.NewServiceError("MAIL_SEND_FAILED", "Gagal mengirim email undangan. Periksa konfigurasi SMTP dan coba lagi.", err)
	}

	resp := toResponse(*user)
	return &resp, nil
}

func (s *userService) buildInvitationEmail(name, activationURL string) string {
	return `<!DOCTYPE html>
<html><head><meta charset="UTF-8"></head>
<body style="font-family: Arial, sans-serif; background: #f4f4f4; padding: 20px;">
<div style="max-width: 600px; margin: auto; background: white; border-radius: 12px; padding: 40px; border: 1px solid #e2e8f0;">
<h2 style="color: #1e3a8a;">Selamat Datang, ` + name + `!</h2>
<p style="color: #475569; line-height: 1.6;">Anda telah diundang oleh Super Admin untuk menjadi bagian dari pengelola <strong>DPP GRADASI</strong>.</p>
<p style="color: #475569; line-height: 1.6;">Silakan klik tombol di bawah ini untuk mengaktifkan akun Anda dan membuat password pertama Anda. Tautan ini berlaku selama 24 jam.</p>
<div style="text-align: center; margin: 32px 0;">
<a href="` + activationURL + `" style="background: #16a34a; color: white; text-decoration: none; padding: 14px 32px; border-radius: 10px; display: inline-block; font-weight: bold; font-size: 15px;">Aktivasi Akun Saya</a>
</div>
<p style="color: #94a3b8; font-size: 13px;">Jika Anda merasa tidak mengenali undangan ini, Anda dapat mengabaikan email ini.</p>
<hr style="border: none; border-top: 1px solid #e2e8f0; margin: 24px 0;">
<p style="color: #94a3b8; font-size: 12px;">DPP GRADASI — Generasi Digital Indonesia</p>
</div></body></html>`
}

func (s *userService) ResendActivation(adminID uint, targetID uint, appURL string) error {
	target, err := s.repo.FindByID(targetID)
	if err != nil {
		return helper.NewServiceError("NOT_FOUND", "Admin tidak ditemukan", err)
	}
	if target.RoleID == 1 {
		return helper.NewServiceError("FORBIDDEN", "Tidak bisa mengirim undangan ulang ke super_admin", nil)
	}
	if target.Status != "pending_activation" {
		return helper.NewServiceError("BAD_REQUEST", "Hanya akun dengan status Menunggu Aktivasi yang bisa dikirimi undangan ulang", nil)
	}

	rawToken, _, _, err := helper.GenerateResetToken("activation", 24*60)
	if err != nil {
		return helper.NewServiceError("SERVER_ERROR", "Gagal membuat token aktivasi", err)
	}

	key := fmt.Sprintf("activation_token:%s", rawToken)
	s.redis.Set(context.Background(), key, target.ID, 24*time.Hour)

	activationURL := fmt.Sprintf("%s/activate-account?token=%s&email=%s", appURL, rawToken, target.Email)
	subject := "Undangan Aktivasi Akun Pengelola - DPP GRADASI"
	body := s.buildInvitationEmail(target.Name, activationURL)

	if err := s.mailer.Send(target.Email, subject, body); err != nil {
		s.redis.Del(context.Background(), key)
		return helper.NewServiceError("MAIL_SEND_FAILED", "Gagal mengirim email undangan. Periksa konfigurasi SMTP dan coba lagi.", err)
	}
	return nil
}

func (s *userService) SetAdminStatus(adminID uint, targetID uint, req *dto.AdminStatusRequest) error {
	target, err := s.repo.FindByID(targetID)
	if err != nil {
		return helper.NewServiceError("NOT_FOUND", "Admin tidak ditemukan", err)
	}
	if target.RoleID == 1 {
		return helper.NewServiceError("FORBIDDEN", "Tidak bisa mengubah status super_admin", nil)
	}

	if req.Status == "active" && target.Status == "pending_activation" {
		return helper.NewServiceError("BAD_REQUEST", "Akun belum melakukan aktivasi. Kirim ulang undangan aktivasi terlebih dahulu.", nil)
	}

	target.Status = req.Status
	if err := s.repo.Update(target); err != nil {
		return helper.NewServiceError("SERVER_ERROR", "Gagal mengubah status admin", err)
	}
	return nil
}

func (s *userService) DeleteAdmin(adminID uint, targetID uint) error {
	if adminID == targetID {
		return helper.NewServiceError("BAD_REQUEST", "Tidak bisa menghapus akun sendiri", nil)
	}
	user, err := s.repo.FindByID(targetID)
	if err != nil {
		return helper.NewServiceError("NOT_FOUND", "Admin tidak ditemukan", err)
	}
	if user.RoleID == 1 {
		return helper.NewServiceError("FORBIDDEN", "Tidak bisa menghapus sesama super_admin", nil)
	}
	return s.repo.Delete(targetID)
}

func (s *userService) RestoreAdmin(adminID uint, targetID uint) error {
	// Optional: can add check for super_admin
	return s.repo.Restore(targetID)
}

func (s *userService) BulkDeleteAdmin(adminID uint, targetIDs []uint) error {
	// Filter out the admin themselves
	filtered := make([]uint, 0)
	for _, id := range targetIDs {
		if id != adminID {
			filtered = append(filtered, id)
		}
	}
	if len(filtered) == 0 {
		return nil
	}
	return s.repo.BulkSoftDelete(filtered)
}

func (s *userService) BulkRestoreAdmin(targetIDs []uint) error {
	if len(targetIDs) == 0 {
		return nil
	}
	return s.repo.BulkRestore(targetIDs)
}

func (s *userService) handleUpload(file *multipart.FileHeader) (string, error) {
	if file == nil {
		return "", nil
	}
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	dst := filepath.Join(s.uploadPath, filename)

	out, err := os.Create(dst)
	if err != nil {
		return "", err
	}
	defer out.Close()

	if _, err = io.Copy(out, src); err != nil {
		return "", err
	}
	return "/uploads/users/" + filename, nil
}

func toResponse(u model.User) dto.UserResponse {
	return dto.UserResponse{
		ID:        u.ID,
		RoleID:    u.RoleID,
		Name:      u.Name,
		Email:     u.Email,
		PhotoPath: u.PhotoPath,
		Status:    u.Status,
	}
}

func generateRandomNumber(length int) (string, error) {
	const charset = "0123456789"
	b := make([]byte, length)
	for i := range b {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		b[i] = charset[n.Int64()]
	}
	return string(b), nil
}
