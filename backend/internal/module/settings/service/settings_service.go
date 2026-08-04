package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/helper"
	activitylogdto "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/activitylog/dto"
	activitylogservice "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/activitylog/service"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/settings/dto"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/settings/mapper"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/settings/repository"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type SettingsService interface {
	GetSettings(ctx context.Context) (*dto.SettingsResponse, error)
	UpdateSettings(ctx context.Context, values map[string]interface{}, updatedBy *uint) (*dto.SettingsResponse, error)
	UploadLogo(ctx context.Context, file *multipart.FileHeader, updatedBy *uint) (*dto.SettingsResponse, error)
}

type settingsService struct {
	db         *gorm.DB
	repo       repository.SettingsRepo
	audit      activitylogservice.ActivityLogService
	uploadPath string
}

func NewSettingsService(db *gorm.DB, repo repository.SettingsRepo, audit activitylogservice.ActivityLogService) SettingsService {
	uploadPath := "public/uploads/settings"
	if err := os.MkdirAll(uploadPath, 0o755); err != nil {
		helper.Logger.Error("gagal buat direktori upload settings", zap.Error(err))
	}
	return &settingsService{db: db, repo: repo, audit: audit, uploadPath: uploadPath}
}

func (s *settingsService) log(ctx context.Context, db *gorm.DB, input *activitylogdto.ActivityLogInput) {
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

func (s *settingsService) GetSettings(ctx context.Context) (*dto.SettingsResponse, error) {
	settings, err := s.repo.Get()
	if err != nil {
		return nil, helper.NewServiceError("SERVER_ERROR", "Gagal mengambil konfigurasi website.", err)
	}
	return mapper.EntityToResponse(mapper.ModelToEntity(settings)), nil
}

// UpdateSettings menerima key-value snake_case (sesuai JSON tag) dan menyimpan
// ke database. Key yang tidak dikenal ditolak (VALIDATION_ERROR).
func (s *settingsService) UpdateSettings(ctx context.Context, values map[string]interface{}, updatedBy *uint) (*dto.SettingsResponse, error) {
	errorsMap := make(map[string]string)
	// Map key JSON (snake_case) -> field Go (mis. "site_name" -> "SiteName")
	// Acuan kontrak: dto.SettingsResponse (json tag), bukan model (model murni gorm).
	settingsModelType := reflect.TypeOf(dto.SettingsResponse{})
	jsonToField := make(map[string]string)
	for i := 0; i < settingsModelType.NumField(); i++ {
		field := settingsModelType.Field(i)
		tag := field.Tag.Get("json")
		if tag == "" {
			continue
		}
		jsonName := strings.Split(tag, ",")[0]
		if jsonName != "" && jsonName != "-" {
			jsonToField[jsonName] = field.Name
		}
	}

	updates := make(map[string]interface{})

	// Field yang dikirim balik FE dari hasil GET /settings tapi BUKAN kolom yang
	// boleh diupdate lewat endpoint ini: id/created_at/updated_at/updated_by adalah
	// field meta, sedangkan captcha_enabled/captcha_site_key dihitung dari config
	// backend saat GET (lihat settings_handler.go), bukan kolom di tabel settings.
	// Dilewati diam-diam supaya FE boleh mengirim balik object hasil GET tanpa
	// harus membuang field-field ini satu per satu (sebelumnya captcha_enabled yang
	// bertipe bool selalu ditolak sebagai "Tipe data tidak valid" dan membuat
	// SELURUH proses simpan pengaturan gagal setiap kali).
	ignoredKeys := map[string]bool{
		"id": true, "created_at": true, "updated_at": true, "updated_by": true,
		"captcha_enabled": true, "captcha_site_key": true,
	}

	for key, val := range values {
		if ignoredKeys[key] {
			continue
		}
		fieldName, found := jsonToField[key]
		if !found {
			errorsMap[key] = "Field tidak dikenal"
			continue
		}

		// AboutMission disimpan sebagai JSON string di DB
		if fieldName == "AboutMission" {
			if val == nil {
				updates[key] = ""
				continue
			}
			missionList, isList := val.([]interface{})
			if !isList {
				errorsMap[key] = "about_mission harus berupa array of string atau null"
				continue
			}
			stringList := make([]string, len(missionList))
			valid := true
			for i, v := range missionList {
				str, isString := v.(string)
				if !isString {
					errorsMap[key] = "about_mission harus berupa array of string atau null"
					valid = false
					break
				}
				stringList[i] = str
			}
			if !valid {
				continue
			}
			jsonMission, err := json.Marshal(stringList)
			if err != nil {
				return nil, helper.NewServiceError("SERVER_ERROR", "Gagal memproses about_mission.", err)
			}
			updates[key] = string(jsonMission)
			continue
		}

		// Field biasa: string
		strVal, ok := val.(string)
		if !ok {
			errorsMap[key] = "Tipe data tidak valid. Harus string atau null."
			continue
		}

		// String kosong berarti field opsional ini sengaja dikosongkan oleh admin —
		// valid, jangan divalidasi format (kosong bukan berarti format email/URL salah).
		if strVal != "" {
			if strings.HasSuffix(key, "_email") && !helper.IsValidEmail(strVal) {
				errorsMap[key] = "Format email tidak valid"
				continue
			}
			if strings.HasSuffix(key, "_url") && !helper.IsValidURL(strVal) {
				errorsMap[key] = "Format URL tidak valid"
				continue
			}
		}

		updates[key] = strVal
	}

	if len(errorsMap) > 0 {
		v := helper.NewValidationError()
		for field, msg := range errorsMap {
			v.Add(field, msg)
		}
		return nil, v
	}

	if len(updates) > 0 {
		if updatedBy != nil {
			updates["updated_by"] = *updatedBy
		}
		if err := s.repo.Update(updates); err != nil {
			return nil, helper.NewServiceError("SERVER_ERROR", "Gagal memperbarui konfigurasi website.", err)
		}
	}

	updatedSettings, err := s.repo.Get()
	if err != nil {
		return nil, helper.NewServiceError("SERVER_ERROR", "Gagal mengambil konfigurasi website.", err)
	}

	if len(updates) > 0 {
		keys := make([]string, 0, len(updates))
		for k := range updates {
			keys = append(keys, k)
		}
		s.log(ctx, s.db, &activitylogdto.ActivityLogInput{
			Action:      "settings.update",
			EntityType:  "settings",
			EntityID:    &updatedSettings.ID,
			EntityLabel: "Konfigurasi Website",
			Description: "Memperbarui konfigurasi website",
			Metadata: map[string]any{
				"fields": keys,
			},
		})
	}

	return mapper.EntityToResponse(mapper.ModelToEntity(updatedSettings)), nil
}

// UploadLogo menyimpan file logo ke disk lokal dan memperbarui logo_path di DB.
func (s *settingsService) UploadLogo(ctx context.Context, file *multipart.FileHeader, updatedBy *uint) (*dto.SettingsResponse, error) {
	if file == nil {
		return nil, helper.NewServiceError("VALIDATION_ERROR", "File logo wajib diunggah", nil)
	}

	// Validasi ukuran: maks 2MB
	if file.Size > 2*1024*1024 {
		return nil, helper.NewServiceError("VALIDATION_ERROR", "File logo tidak valid. Maksimal 2MB dengan format PNG, JPG, atau WEBP.", nil)
	}

	// Validasi MIME type
	src, err := file.Open()
	if err != nil {
		return nil, helper.NewServiceError("SERVER_ERROR", "Gagal membaca file logo.", err)
	}
	defer src.Close()

	header := make([]byte, 512)
	n, err := src.Read(header)
	if err != nil && err != io.EOF {
		return nil, helper.NewServiceError("SERVER_ERROR", "Gagal membaca file logo.", err)
	}
	// Reset reader agar bisa dibaca ulang saat copy
	if _, err := src.Seek(0, io.SeekStart); err != nil {
		return nil, helper.NewServiceError("SERVER_ERROR", "Gagal membaca file logo.", err)
	}

	mimeType := http.DetectContentType(header[:n])
	switch mimeType {
	case "image/png", "image/jpeg", "image/webp":
	default:
		return nil, helper.NewServiceError("VALIDATION_ERROR", "File logo tidak valid. Maksimal 2MB dengan format PNG, JPG, atau WEBP.", nil)
	}

	ext := filepath.Ext(file.Filename)
	if ext == "" {
		switch mimeType {
		case "image/png":
			ext = ".png"
		case "image/webp":
			ext = ".webp"
		default:
			ext = ".jpg"
		}
	}
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)

	// Pastikan direktori upload ada
	if err := os.MkdirAll(s.uploadPath, 0o755); err != nil {
		return nil, helper.NewServiceError("SERVER_ERROR", "Gagal menyimpan file logo.", err)
	}

	dst := filepath.Join(s.uploadPath, filename)
	out, err := os.Create(dst)
	if err != nil {
		return nil, helper.NewServiceError("SERVER_ERROR", "Gagal menyimpan file logo.", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, src); err != nil {
		return nil, helper.NewServiceError("SERVER_ERROR", "Gagal menyimpan file logo.", err)
	}

	logoPath := "/uploads/settings/" + filename
	if err := s.repo.UpdateLogo(logoPath, updatedBy); err != nil {
		return nil, helper.NewServiceError("SERVER_ERROR", "Gagal memperbarui logo website.", err)
	}

	settingsID := uint(1)
	s.log(ctx, s.db, &activitylogdto.ActivityLogInput{
		Action:      "settings.update",
		EntityType:  "settings",
		EntityID:    &settingsID,
		EntityLabel: "Logo Website",
		Description: "Mengunggah logo website",
		Metadata: map[string]any{
			"logo_path": logoPath,
		},
	})

	return s.GetSettings(ctx)
}
