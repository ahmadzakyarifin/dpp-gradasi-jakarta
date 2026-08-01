package service

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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
	GetAdmins(q dto.ListUsersQuery) (*dto.UserListResponse, error)
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
	_ = os.MkdirAll(uploadPath, 0755)
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
		go func() {
			_ = s.mailer.Send(
				req.Email,
				"Verifikasi Perubahan Email",
				fmt.Sprintf("Kode verifikasi Anda adalah: <b>%s</b>. Kode ini berlaku selama 1 jam.", token),
			)
		}()

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
	if err := json.Unmarshal([]byte(val), &data); err != nil {
		return helper.NewServiceError("INVALID_TOKEN", "Kode verifikasi tidak valid", err)
	}

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
	user.MustChangePassword = false // Password sudah diganti, flag reset
	if err := s.repo.Update(user); err != nil {
		return helper.NewServiceError("SERVER_ERROR", "Gagal menyimpan password baru", err)
	}
	return nil
}

func (s *userService) GetAdmins(q dto.ListUsersQuery) (*dto.UserListResponse, error) {
	users, total, err := s.repo.FindAllAdmins(q)
	if err != nil {
		return nil, helper.NewServiceError("SERVER_ERROR", "Gagal mengambil data admin", err)
	}

	page, limit := q.Page, q.Limit
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	items := make([]dto.UserResponse, 0, len(users))
	for _, u := range users {
		items = append(items, toResponse(u))
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))
	return &dto.UserListResponse{
		Items: items,
		Pagination: dto.Pagination{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}

func (s *userService) CreateAdmin(req *dto.AdminCreateRequest, appURL string) (*dto.UserResponse, error) {
	existing, _ := s.repo.FindByEmail(req.Email)
	if existing != nil {
		return nil, helper.NewServiceError("EMAIL_TAKEN", "Email sudah terdaftar", nil)
	}

	// Generate password default (acak 10 karakter)
	defaultPassword, err := generateRandomPassword(10)
	if err != nil {
		return nil, helper.NewServiceError("SERVER_ERROR", "Gagal membuat password default", err)
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(defaultPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, helper.NewServiceError("SERVER_ERROR", "Gagal mengenkripsi password", err)
	}

	user := &model.User{
		Name:               req.Name,
		Email:              req.Email,
		RoleID:             req.RoleID,
		Password:           string(hashed),
		Status:             "active",
		MustChangePassword: true, // Wajib ganti password pada login pertama
		EmailVerifiedAt:    nil,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, helper.NewServiceError("SERVER_ERROR", "Gagal membuat admin", err)
	}

	// Kirim email kredensial (email + password default) secara SYNC
	subject := "Akun Pengelola DPP GRADASI — Kredensial Login"
	body := s.buildCredentialsEmail(user.Name, user.Email, defaultPassword)

	if err := s.mailer.Send(user.Email, subject, body); err != nil {
		// Email gagal → hapus user agar tidak ada akun tanpa kredensial
		_ = s.repo.Delete(user.ID)
		return nil, helper.NewServiceError("MAIL_SEND_FAILED", "Gagal mengirim email kredensial. Periksa konfigurasi SMTP dan coba lagi.", err)
	}

	resp := toResponse(*user)
	return &resp, nil
}

func (s *userService) buildCredentialsEmail(name, email, password string) string {
	return `<!DOCTYPE html>
<html><head><meta charset="UTF-8"></head>
<body style="font-family: Arial, sans-serif; background: #f4f4f4; padding: 20px;">
<div style="max-width: 600px; margin: auto; background: white; border-radius: 12px; padding: 40px; border: 1px solid #e2e8f0;">
<h2 style="color: #1e3a8a;">Selamat Datang, ` + name + `!</h2>
<p style="color: #475569; line-height: 1.6;">Anda telah didaftarkan sebagai pengelola <strong>DPP GRADASI</strong>. Berikut kredensial login Anda:</p>
<div style="background: #f8fafc; border: 1px solid #e2e8f0; border-radius: 10px; padding: 20px; margin: 24px 0;">
<p style="margin: 8px 0;"><strong style="color: #1e3a8a;">Email:</strong> <span style="color: #334155;">` + email + `</span></p>
<p style="margin: 8px 0;"><strong style="color: #1e3a8a;">Password:</strong> <span style="font-family: monospace; background: #f1f5f9; padding: 2px 8px; border-radius: 6px; color: #dc2626;">` + password + `</span></p>
</div>
<p style="color: #475569; line-height: 1.6;">Silakan login di <strong>` + `, lalu <strong>segera ganti password Anda</strong> pada halaman profil.</p>
<p style="color: #dc2626; font-size: 13px; font-weight: bold;">Demi keamanan, Anda akan diminta mengganti password ini pada login pertama.</p>
<p style="color: #94a3b8; font-size: 13px;">Jika Anda merasa tidak mengenali email ini, Anda dapat mengabaikannya.</p>
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
		return helper.NewServiceError("FORBIDDEN", "Tidak bisa mengirim kredensial ulang ke super_admin", nil)
	}

	// Generate password default baru
	defaultPassword, err := generateRandomPassword(10)
	if err != nil {
		return helper.NewServiceError("SERVER_ERROR", "Gagal membuat password default", err)
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(defaultPassword), bcrypt.DefaultCost)
	if err != nil {
		return helper.NewServiceError("SERVER_ERROR", "Gagal mengenkripsi password", err)
	}

	target.Password = string(hashed)
	target.MustChangePassword = true
	if err := s.repo.Update(target); err != nil {
		return helper.NewServiceError("SERVER_ERROR", "Gagal memperbarui kredensial admin", err)
	}

	subject := "Kredensial Login DPP GRADASI — Perbaruan Password"
	body := s.buildCredentialsEmail(target.Name, target.Email, defaultPassword)

	if err := s.mailer.Send(target.Email, subject, body); err != nil {
		return helper.NewServiceError("MAIL_SEND_FAILED", "Gagal mengirim email kredensial. Periksa konfigurasi SMTP dan coba lagi.", err)
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
		return helper.NewServiceError("FORBIDDEN", "Tidak bisa menghapus super_admin", nil)
	}
	return s.repo.Delete(targetID)
}

func (s *userService) RestoreAdmin(adminID uint, targetID uint) error {
	return s.repo.Restore(targetID)
}

func (s *userService) BulkDeleteAdmin(adminID uint, targetIDs []uint) error {
	// Filter out admin sendiri + semua super_admin (role 1)
	filtered := make([]uint, 0)
	for _, id := range targetIDs {
		if id == adminID {
			continue
		}
		user, err := s.repo.FindByID(id)
		if err != nil {
			continue // skip yang tidak ditemukan
		}
		if user.RoleID == 1 {
			continue // super_admin tidak bisa dihapus
		}
		filtered = append(filtered, id)
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

	// Limit ukuran 2MB
	const maxSize = 2 << 20 // 2MB
	if file.Size > maxSize {
		return "", fmt.Errorf("file terlalu besar (maks 2MB)")
	}

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// Detect MIME dari content (bukan ekstensi)
	buffer := make([]byte, 512)
	n, _ := src.Read(buffer)
	if n == 0 {
		return "", fmt.Errorf("file kosong")
	}
	mimeType := http.DetectContentType(buffer[:n])
	if !strings.HasPrefix(mimeType, "image/") {
		return "", fmt.Errorf("file harus berupa gambar (image/*)")
	}

	// Reset reader
	if _, err := src.Seek(0, io.SeekStart); err != nil {
		return "", err
	}

	ext := filepath.Ext(file.Filename)
	if ext == "" {
		switch mimeType {
		case "image/png":
			ext = ".png"
		case "image/jpeg":
			ext = ".jpg"
		default:
			ext = ".png"
		}
	}
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
	roleName := ""
	if u.Role.Name != "" {
		roleName = u.Role.Name
	}
	return dto.UserResponse{
		ID:                 u.ID,
		RoleID:             u.RoleID,
		RoleName:           roleName,
		Name:               u.Name,
		Email:              u.Email,
		PhotoPath:          u.PhotoPath,
		Status:             u.Status,
		MustChangePassword: u.MustChangePassword,
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

// generateRandomPassword membuat password acak alfanumerik (aman untuk dikirim via email).
func generateRandomPassword(length int) (string, error) {
	const charset = "ABCDEFGHJKLMNPQRSTUVWXYZabcdefghjkmnpqrstuvwxyz23456789" // tanpa karakter ambigu (0,1,O,I,l)
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
