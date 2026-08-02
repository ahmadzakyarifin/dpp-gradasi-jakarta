package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ahmadzakyarifin/schoolpay/backend/config"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/helper"
	jobs "github.com/ahmadzakyarifin/schoolpay/backend/internal/job"
	activitylogdto "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/activitylog/dto"
	activitylogservice "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/activitylog/service"
	authrepo "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/auth/repository"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/notification"
	notifdto "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/notification/dto"
	notifsvc "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/notification/service"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/user/dto"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/user/entity"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/user/mapper"
	userrepo "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/user/repository"
	emailtemplate "github.com/ahmadzakyarifin/schoolpay/backend/internal/template_message/email"
	whatsapptemplate "github.com/ahmadzakyarifin/schoolpay/backend/internal/template_message/whatsapp"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/validator"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type UserService interface {
	GetAll(ctx context.Context, roleID uint) ([]entity.User, error)
	GetPaginated(ctx context.Context, req dto.UserQueryReq) ([]entity.User, int, error)
	Create(ctx context.Context, req dto.UserCreateReq, channel string) (*entity.User, error)
	Update(ctx context.Context, id uint, req dto.UserUpdateReq) (*entity.User, error)
	Delete(ctx context.Context, id uint) error
	ActivateAccount(ctx context.Context, token string, password string) (*entity.User, error)
	SaveRefreshToken(ctx context.Context, userID uint, token string, expiresAt time.Time) error
	ToggleStatus(ctx context.Context, id uint) error
	ResendNotification(ctx context.Context, id uint, channel string) error
	BulkResendNotification(ctx context.Context, ids []uint, channel string) (*BulkResendResult, error)
	ExportExcel(ctx context.Context, search, role, filter, status string) ([]byte, error)
	GetByID(ctx context.Context, id uint) (*entity.User, error)
	BulkDelete(ctx context.Context, ids []uint) error
	Restore(ctx context.Context, id uint) error
	BulkRestore(ctx context.Context, ids []uint) error
	GetDependencyInfo(ctx context.Context, id uint) (map[string]interface{}, error)
	CheckUnique(ctx context.Context, field string, value string, excludeID uint) (bool, error)
	UpdateProfile(ctx context.Context, id uint, name string) (*entity.User, error)
}

type userService struct {
	db       *gorm.DB
	repo     userrepo.UserRepo
	authRepo authrepo.AuthRepo
	job      *jobs.Client
	audit    activitylogservice.ActivityLogService
	notifSvc notifsvc.NotificationService
	cfg      *config.Config
}

func NewUserService(db *gorm.DB, repo userrepo.UserRepo, authRepo authrepo.AuthRepo, job *jobs.Client, audit activitylogservice.ActivityLogService, notifSvc notifsvc.NotificationService, cfg *config.Config) UserService {
	return &userService{
		db:       db,
		repo:     repo,
		authRepo: authRepo,
		job:      job,
		audit:    audit,
		notifSvc: notifSvc,
		cfg:      cfg,
	}
}

func (s *userService) log(ctx context.Context, db *gorm.DB, input *activitylogdto.ActivityLogInput) {
	if s.audit == nil {
		return
	}
	userID, userName, role, ipAddress, userAgent := helper.GetAuditMeta(ctx)
	var uID *uint
	if userID > 0 {
		uID = &userID
	}
	input.ActorID = uID
	input.ActorName = userName
	input.ActorRole = role
	input.IPAddress = ipAddress
	input.UserAgent = userAgent

	_ = s.audit.Log(ctx, db, input)
}

func (s *userService) GetAll(ctx context.Context, roleID uint) ([]entity.User, error) {
	return s.repo.FindAll(ctx, roleID)
}

func (s *userService) GetByID(ctx context.Context, id uint) (*entity.User, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, helper.NewNotFoundError("User tidak ditemukan")
	}
	return user, nil
}

func (s *userService) GetPaginated(ctx context.Context, req dto.UserQueryReq) ([]entity.User, int, error) {
	return s.repo.FindPaginated(ctx, req.Page, req.Limit, req.Search, req.Role, req.Status, req.Sort, req.Relation, req.Trashed)
}

func userAuditValues(user *entity.User) map[string]interface{} {
	if user == nil {
		return nil
	}
	vals := map[string]interface{}{
		"name":    user.Name,
		"email":   user.Email,
		"phone":   user.Phone,
		"role_id": user.RoleID,
		"status":  user.Status,
	}
	if user.DateOfBirth != nil {
		vals["date_of_birth"] = user.DateOfBirth.Format("2006-01-02")
	}
	if user.CountryCode != nil {
		vals["country_code"] = *user.CountryCode
	}
	return vals
}

func (s *userService) Create(ctx context.Context, req dto.UserCreateReq, channel string) (*entity.User, error) {
	user := mapper.CreateReqToEntity(&req)
	user.Name = strings.Title(strings.ToLower(strings.TrimSpace(user.Name)))
	user.Email = strings.ToLower(strings.TrimSpace(user.Email))
	user.Phone = validator.NormalizePhoneNumber(user.Phone)
	if user.Status == "" {
		user.Status = "inactive"
	}

	v := helper.NewValidationError()
	if err := validator.ValidatePhoneNumber(user.Phone); err != nil {
		v.Add("phone", "Format nomor WhatsApp tidak valid")
	}
	if user.DateOfBirth != nil && user.DateOfBirth.After(time.Now()) {
		v.Add("date_of_birth", "Tanggal lahir tidak boleh di masa depan")
	}
	if user.CountryCode != nil {
		cc := strings.ToUpper(strings.TrimSpace(*user.CountryCode))
		user.CountryCode = &cc
		if !isValidCountryCode(cc) {
			v.Add("country_code", "Kode negara tidak valid")
		}
	}
	if existing, _ := s.repo.FindByEmail(ctx, user.Email); existing != nil {
		v.Add("email", fmt.Sprintf("Email '%s' sudah terdaftar sebagai Pengguna/Wali Murid", user.Email))
	}
	if existing, _ := s.repo.FindByPhone(ctx, user.Phone); existing != nil {
		v.Add("phone", fmt.Sprintf("Nomor WhatsApp '%s' sudah terdaftar sebagai Pengguna/Wali Murid", user.Phone))
	}
	if len(v.Errors) > 0 {
		return nil, v
	}

	if err := s.repo.CreateTx(ctx, user); err != nil {
		return nil, err
	}

	s.log(ctx, s.db, &activitylogdto.ActivityLogInput{
		Action:      "users.create",
		EntityType:  "users",
		EntityID:    &user.ID,
		EntityLabel: user.Name,
		Description: fmt.Sprintf("Membuat pengguna baru: %s", user.Name),
		Metadata: map[string]any{
			"new_values": userAuditValues(user),
		},
	})

	// Otomatis kirim notifikasi aktivasi akun (B2: ikut pilihan channel FE).
	if user.Status == "active" {
		_ = s.enqueueActivation(ctx, user, channel)
	}
	return user, nil
}

func (s *userService) Update(ctx context.Context, id uint, req dto.UserUpdateReq) (*entity.User, error) {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, helper.NewNotFoundError("Pengguna tidak ditemukan")
	}

	mapper.UpdateReqToEntity(&req, existing)
	existing.Name = strings.Title(strings.ToLower(strings.TrimSpace(existing.Name)))
	existing.Email = strings.ToLower(strings.TrimSpace(existing.Email))
	existing.Phone = validator.NormalizePhoneNumber(existing.Phone)

	v := helper.NewValidationError()
	if err := validator.ValidatePhoneNumber(existing.Phone); err != nil {
		v.Add("phone", "Format nomor WhatsApp tidak valid")
	}
	if existing.DateOfBirth != nil && existing.DateOfBirth.After(time.Now()) {
		v.Add("date_of_birth", "Tanggal lahir tidak boleh di masa depan")
	}
	if existing.CountryCode != nil {
		cc := strings.ToUpper(strings.TrimSpace(*existing.CountryCode))
		existing.CountryCode = &cc
		if !isValidCountryCode(cc) {
			v.Add("country_code", "Kode negara tidak valid")
		}
	}
	if other, _ := s.repo.FindByEmail(ctx, existing.Email); other != nil && other.ID != existing.ID {
		v.Add("email", fmt.Sprintf("Email '%s' sudah digunakan oleh Pengguna lain", existing.Email))
	}
	if other, _ := s.repo.FindByPhone(ctx, existing.Phone); other != nil && other.ID != existing.ID {
		v.Add("phone", fmt.Sprintf("Nomor WhatsApp '%s' sudah digunakan oleh Pengguna lain", existing.Phone))
	}
	if len(v.Errors) > 0 {
		return nil, v
	}

	oldVals := userAuditValues(existing)

	if err := s.repo.UpdateTx(ctx, existing); err != nil {
		return nil, err
	}

	s.log(ctx, s.db, &activitylogdto.ActivityLogInput{
		Action:      "users.update",
		EntityType:  "users",
		EntityID:    &existing.ID,
		EntityLabel: existing.Name,
		Description: fmt.Sprintf("Memperbarui pengguna: %s", existing.Name),
		Metadata: map[string]any{
			"old_values": oldVals,
			"new_values": userAuditValues(existing),
		},
	})
	return existing, nil
}

func (s *userService) ensureDeletableUser(ctx context.Context, user *entity.User) error {
	if user == nil {
		return helper.NewNotFoundError("Pengguna tidak ditemukan")
	}
	if user.RoleName == "parent" {
		count, err := s.repo.CountStudentsByParent(ctx, user.ID)
		if err == nil && count > 0 {
			v := helper.NewValidationError()
			v.Add("general", fmt.Sprintf("Akun %s masih terhubung dengan %d siswa aktif", user.Name, count))
			return v
		}
	}
	return nil
}

func (s *userService) Delete(ctx context.Context, id uint) error {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return helper.NewNotFoundError("Pengguna tidak ditemukan")
	}
	if err := s.ensureDeletableUser(ctx, existing); err != nil {
		return err
	}

	err = s.repo.Delete(ctx, id)
	if err == nil {
		s.log(ctx, s.db, &activitylogdto.ActivityLogInput{
			Action:      "users.delete",
			EntityType:  "users",
			EntityID:    &id,
			EntityLabel: existing.Name,
			Description: fmt.Sprintf("Menghapus pengguna: %s", existing.Name),
			Metadata: map[string]any{
				"old_values": map[string]interface{}{"status": "active"},
				"new_values": map[string]interface{}{"status": "deleted"},
			},
		})
	}
	return err
}

func (s *userService) BulkDelete(ctx context.Context, ids []uint) error {
	for _, id := range ids {
		existing, err := s.repo.FindByID(ctx, id)
		if err != nil {
			return err
		}
		if existing == nil {
			return helper.NewNotFoundError("Salah satu pengguna tidak ditemukan")
		}
		if err := s.ensureDeletableUser(ctx, existing); err != nil {
			return err
		}
	}

	err := s.repo.BulkDelete(ctx, ids)
	if err == nil {
		s.log(ctx, s.db, &activitylogdto.ActivityLogInput{
			Action:      "users.delete",
			EntityType:  "users",
			Description: fmt.Sprintf("Menghapus massal pengguna dengan ID: %v", ids),
		})
	}
	return err
}

func (s *userService) Restore(ctx context.Context, id uint) error {
	err := s.repo.Restore(ctx, id)
	if err == nil {
		s.log(ctx, s.db, &activitylogdto.ActivityLogInput{
			Action:      "users.restore",
			EntityType:  "users",
			EntityID:    &id,
			EntityLabel: fmt.Sprintf("User ID %d", id),
			Description: fmt.Sprintf("Memulihkan pengguna ID: %d", id),
		})
	}
	return err
}

func (s *userService) BulkRestore(ctx context.Context, ids []uint) error {
	err := s.repo.BulkRestore(ctx, ids)
	if err == nil {
		s.log(ctx, s.db, &activitylogdto.ActivityLogInput{
			Action:      "users.restore",
			EntityType:  "users",
			Description: fmt.Sprintf("Memulihkan massal pengguna dengan ID: %v", ids),
		})
	}
	return err
}

func (s *userService) SaveRefreshToken(ctx context.Context, userID uint, token string, expiresAt time.Time) error {
	_, _, _, ip, userAgent := helper.GetAuditMeta(ctx)
	return s.authRepo.SaveRefreshToken(ctx, userID, token, expiresAt, ip, userAgent, "")
}

func (s *userService) ActivateAccount(ctx context.Context, token string, password string) (*entity.User, error) {
	userID, err := s.authRepo.FindAuthToken(ctx, token, "activation")
	if err != nil {
		return nil, errors.New("link aktivasi tidak valid atau sudah kedaluwarsa")
	}

	hashed, err := helper.HashPassword(password)
	if err != nil {
		return nil, errors.New("gagal memproses password")
	}
	_ = s.repo.UpdatePassword(ctx, userID, hashed)

	_ = s.authRepo.DeleteAuthToken(ctx, token, "activation")
	user, err := s.repo.FindByID(ctx, userID)
	if err == nil && user != nil {
		s.log(ctx, s.db, &activitylogdto.ActivityLogInput{
			Action:      "users.activate_account",
			EntityType:  "users",
			EntityID:    &user.ID,
			EntityLabel: user.Name,
			Description: fmt.Sprintf("Mengaktivasi akun pengguna: %s", user.Name),
		})
	}
	return user, err
}

func (s *userService) ToggleStatus(ctx context.Context, id uint) error {
	existing, _ := s.repo.FindByID(ctx, id)
	if existing == nil {
		return helper.NewNotFoundError("Pengguna tidak ditemukan")
	}
	if existing.RoleName == "parent" && existing.Status == "active" {
		count, err := s.repo.CountStudentsByParent(ctx, id)
		if err == nil && count > 0 {
			v := helper.NewValidationError()
			v.Add("general", fmt.Sprintf("Akun %s masih terhubung dengan %d siswa aktif dan tidak dapat dinonaktifkan", existing.Name, count))
			return v
		}
	}
	err := s.repo.ToggleStatus(ctx, id)
	if err == nil {
		s.log(ctx, s.db, &activitylogdto.ActivityLogInput{
			Action:      "users.update",
			EntityType:  "users",
			EntityID:    &id,
			EntityLabel: existing.Name,
			Description: fmt.Sprintf("Mengubah status keaktifan pengguna: %s", existing.Name),
		})
	}
	return err
}

func (s *userService) enqueueActivation(ctx context.Context, user *entity.User, channel string) error {
	if s.notifSvc == nil {
		return nil
	}
	token, err := helper.GenerateRandomToken(32)
	if err != nil {
		return err
	}
	expiresAt := time.Now().Add(72 * time.Hour)
	if err := s.authRepo.SaveAuthToken(ctx, user.ID, token, "activation", expiresAt); err != nil {
		return err
	}

	link := fmt.Sprintf("%s/activate?token=%s",
		strings.TrimSuffix(s.cfg.App.URL, "/"),
		token,
	)

	channel = normalizeActivationChannel(channel)

	switch channel {
	case "whatsapp":
		text, err := whatsapptemplate.Render("account_activation.txt", map[string]any{
			"Name": user.Name,
			"URL":  link,
		})
		if err != nil {
			return fmt.Errorf("gagal merender template whatsapp: %w", err)
		}
		return s.notifSvc.SendWhatsApp(ctx, notifdto.SendWhatsAppRequest{
			To:     user.Phone,
			Text:   text,
			UserID: &user.ID,
		})
	case "email":
		html, err := emailtemplate.Render("account_activation.html", map[string]any{
			"Name": user.Name,
			"URL":  link,
		})
		if err != nil {
			return fmt.Errorf("gagal merender template email: %w", err)
		}
		return s.notifSvc.SendEmail(ctx, notifdto.SendEmailRequest{
			To:      user.Email,
			Subject: notification.EventToSubject(notification.EventAuthAccountActivation),
			HTML:    html,
			UserID:  &user.ID,
		})
	default: // "all" — kirim email + whatsapp
		html, err := emailtemplate.Render("account_activation.html", map[string]any{
			"Name": user.Name,
			"URL":  link,
		})
		if err == nil {
			_ = s.notifSvc.SendEmail(ctx, notifdto.SendEmailRequest{
				To:      user.Email,
				Subject: notification.EventToSubject(notification.EventAuthAccountActivation),
				HTML:    html,
				UserID:  &user.ID,
			})
		}
		text, err := whatsapptemplate.Render("account_activation.txt", map[string]any{
			"Name": user.Name,
			"URL":  link,
		})
		if err == nil {
			_ = s.notifSvc.SendWhatsApp(ctx, notifdto.SendWhatsAppRequest{
				To:     user.Phone,
				Text:   text,
				UserID: &user.ID,
			})
		}
		return nil
	}
}

func (s *userService) ResendNotification(ctx context.Context, id uint, channel string) error {
	channel = normalizeActivationChannel(channel)
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return helper.NewNotFoundError("Pengguna tidak ditemukan")
	}
	if user.Status != "active" {
		v := helper.NewValidationError()
		v.Add("general", fmt.Sprintf("Gagal: Akun %s sedang Non-Aktif", user.Name))
		return v
	}
	if userHasPassword(user) {
		v := helper.NewValidationError()
		v.Add("general", fmt.Sprintf("Gagal: Akun %s sudah aktif dan sudah memiliki password. Gunakan fitur reset password jika diperlukan", user.Name))
		return v
	}

	s.log(ctx, s.db, &activitylogdto.ActivityLogInput{
		Action:      "users.send_activation",
		EntityType:  "users",
		EntityID:    &id,
		EntityLabel: user.Name,
		Description: fmt.Sprintf("Mengirim ulang notifikasi aktivasi melalui %s", channel),
	})

	return s.enqueueActivation(ctx, user, channel)
}

type BulkResendResult struct {
	Total  int      `json:"total"`
	Sent   int      `json:"sent"`
	Failed int      `json:"failed"`
	Errors []string `json:"errors"`
}

func (s *userService) BulkResendNotification(ctx context.Context, ids []uint, channel string) (*BulkResendResult, error) {
	channel = normalizeActivationChannel(channel)
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
		if user.RoleName == "parent" {
			count, err := s.repo.CountStudentsByParent(ctx, user.ID)
			if err == nil && count == 0 {
				result.Failed++
				result.Errors = append(result.Errors, fmt.Sprintf("%s: Tidak memiliki data siswa yang terhubung", user.Name))
				continue
			}
		}
		if user.Status != "active" {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("%s: Akun sedang Non-Aktif", user.Name))
			continue
		}
		if userHasPassword(user) {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("%s: Akun sudah aktif dan sudah memiliki password", user.Name))
			continue
		}

		result.Sent++
		s.log(ctx, s.db, &activitylogdto.ActivityLogInput{
			Action:      "users.send_activation",
			EntityType:  "users",
			EntityID:    &id,
			EntityLabel: user.Name,
			Description: fmt.Sprintf("Mengirim massal notifikasi aktivasi melalui %s", channel),
		})
		if err := s.enqueueActivation(ctx, user, channel); err != nil {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("%s: %s", user.Name, err.Error()))
		}
	}

	return result, nil
}

func userHasPassword(user *entity.User) bool {
	return user != nil && strings.TrimSpace(user.PasswordHash) != ""
}

func isValidCountryCode(code string) bool {
	if len(code) != 2 {
		return false
	}
	for _, c := range code {
		if c < 'A' || c > 'Z' {
			return false
		}
	}
	return true
}

func normalizeActivationChannel(channel string) string {
	switch strings.ToLower(strings.TrimSpace(channel)) {
	case "email":
		return "email"
	case "whatsapp", "wa":
		return "whatsapp"
	default:
		return ""
	}
}

func (s *userService) ExportExcel(ctx context.Context, search, role, filter, status string) ([]byte, error) {
	s.log(ctx, s.db, &activitylogdto.ActivityLogInput{
		Action:      "users.export_excel",
		EntityType:  "users",
		Description: "Melakukan ekspor daftar pengguna ke Excel",
	})
	users, _, _ := s.repo.FindPaginated(ctx, 1, 1000000, search, role, status, "", "", false)
	f := excelize.NewFile()
	sheet := "Daftar Pengguna"
	f.SetSheetName("Sheet1", sheet)

	headers := []string{"NAMA LENGKAP", "EMAIL", "NOMOR WHATSAPP", "ROLE ID", "STATUS"}
	colWidths := []float64{28, 28, 20, 12, 12}

	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true},
		Alignment: &excelize.Alignment{Horizontal: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})

	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
		f.SetCellStyle(sheet, cell, cell, headerStyle)
		colName, _ := excelize.ColumnNumberToName(i + 1)
		if i < len(colWidths) {
			f.SetColWidth(sheet, colName, colName, colWidths[i])
		} else {
			f.SetColWidth(sheet, colName, colName, 20)
		}
	}

	bodyStyle, _ := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Alignment: &excelize.Alignment{Vertical: "center"},
	})

	for i, u := range users {
		row := i + 2
		activeStatus := "Aktif"
		if u.Status != "active" {
			activeStatus = "Non-Aktif"
		}

		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), u.Name)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), u.Email)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), u.Phone)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), u.RoleID)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), activeStatus)

		f.SetCellStyle(sheet, fmt.Sprintf("A%d", row), fmt.Sprintf("E%d", row), bodyStyle)
	}

	buf, _ := f.WriteToBuffer()
	return buf.Bytes(), nil
}

func FormatActivationStudentList(raw string) string {
	parts := strings.Split(raw, "||")
	lines := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		if split := strings.SplitN(part, "::", 2); len(split) == 2 {
			part = split[1]
		}
		lines = append(lines, fmt.Sprintf("• *%s*", part))
	}
	if len(lines) == 0 {
		return "• *Data siswa terhubung*"
	}
	return strings.Join(lines, "\n")
}

func (s *userService) GetDependencyInfo(ctx context.Context, id uint) (map[string]interface{}, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, helper.NewNotFoundError("Pengguna tidak ditemukan")
	}

	var messages []string
	counts := map[string]int{}

	if user.RoleName == "parent" {
		count, err := s.repo.CountStudentsByParent(ctx, id)
		if err == nil && count > 0 {
			messages = append(messages, fmt.Sprintf("terhubung dengan %d siswa aktif", count))
			counts["active_students"] = count
		}
	}

	return map[string]interface{}{
		"has_dependencies": len(messages) > 0,
		"message":          strings.Join(messages, " dan "),
		"counts":           counts,
	}, nil
}

func (s *userService) CheckUnique(ctx context.Context, field string, value string, excludeID uint) (bool, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return true, nil
	}

	switch field {
	case "email":
		value = strings.ToLower(value)
		if existing, _ := s.repo.FindByEmail(ctx, value); existing != nil && existing.ID != excludeID {
			return false, nil
		}
	case "phone_number":
		normalized := validator.NormalizePhoneNumber(value)
		if existing, _ := s.repo.FindByPhone(ctx, normalized); existing != nil && existing.ID != excludeID {
			return false, nil
		}
	case "nik":
		// NIK sudah dipindah ke tabel guardians.
		return true, nil
	}

	return true, nil
}

func (s *userService) UpdateProfile(ctx context.Context, id uint, name string) (*entity.User, error) {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, helper.NewNotFoundError("Pengguna tidak ditemukan")
	}

	name = strings.TrimSpace(name)
	if len(name) < 2 {
		return nil, fmt.Errorf("nama minimal 2 karakter")
	}

	oldVals := map[string]any{"name": existing.Name}
	existing.Name = name

	if err := s.repo.UpdateTx(ctx, existing); err != nil {
		return nil, err
	}

	s.log(ctx, s.db, &activitylogdto.ActivityLogInput{
		Action:      "users.update_profile",
		EntityType:  "users",
		EntityID:    &existing.ID,
		EntityLabel: existing.Name,
		Description: fmt.Sprintf("Pengguna %s memperbarui profil", existing.Name),
		Metadata: map[string]any{
			"old_values": oldVals,
			"new_values": map[string]any{"name": existing.Name},
		},
	})

	return existing, nil
}
