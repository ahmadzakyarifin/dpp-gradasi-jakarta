package service

import (
	"encoding/json"
	"errors"
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
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/settings/dto"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/settings/model"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/settings/repository"
	"go.uber.org/zap"
)

type SettingsService interface {
	GetSettings() (*dto.SettingsResponse, error)
	UpdateSettings(values map[string]interface{}, updatedBy *uint) (*dto.SettingsResponse, error)
	UploadLogo(file *multipart.FileHeader, updatedBy *uint) (*dto.SettingsResponse, error)
}

type settingsService struct {
	repo       repository.SettingsRepo
	uploadPath string
}

func NewSettingsService(repo repository.SettingsRepo) SettingsService {
	uploadPath := "public/uploads/settings"
	if err := os.MkdirAll(uploadPath, 0o755); err != nil {
		helper.Logger.Error("gagal buat direktori upload settings", zap.Error(err))
	}
	return &settingsService{repo: repo, uploadPath: uploadPath}
}

func (s *settingsService) GetSettings() (*dto.SettingsResponse, error) {
	settings, err := s.repo.Get()
	if err != nil {
		return nil, err
	}
	return toResponse(*settings), nil
}

// UpdateSettings menerima key-value snake_case (sesuai JSON tag) dan menyimpan
// ke database. Key yang tidak dikenal ditolak (400).
func (s *settingsService) UpdateSettings(values map[string]interface{}, updatedBy *uint) (*dto.SettingsResponse, error) {
	errorsMap := make(map[string]string)
	settingsModelType := reflect.TypeOf(model.Settings{})

	// Map key JSON (snake_case) -> field Go (mis. "site_name" -> "SiteName")
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

	for key, val := range values {
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
			for i, v := range missionList {
				str, isString := v.(string)
				if !isString {
					errorsMap[key] = "about_mission harus berupa array of string atau null"
					break
				}
				stringList[i] = str
			}
			if errorsMap[key] != "" {
				continue
			}
			jsonMission, err := json.Marshal(stringList)
			if err != nil {
				return nil, errors.New("Gagal memproses about_mission")
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

		if strings.HasSuffix(key, "_email") && !helper.IsValidEmail(strVal) {
			errorsMap[key] = "Format email tidak valid"
			continue
		}
		if strings.HasSuffix(key, "_url") && !helper.IsValidURL(strVal) {
			errorsMap[key] = "Format URL tidak valid"
			continue
		}

		updates[key] = strVal
	}

	if len(errorsMap) > 0 {
		return nil, fmt.Errorf("Validasi gagal: %v", errorsMap)
	}

	if len(updates) > 0 {
		if updatedBy != nil {
			updates["updated_by"] = *updatedBy
		}
		if err := s.repo.Update(updates); err != nil {
			return nil, err
		}
	}

	updatedSettings, err := s.repo.Get()
	if err != nil {
		return nil, err
	}

	return toResponse(*updatedSettings), nil
}

// UploadLogo menyimpan file logo ke disk lokal dan memperbarui logo_path di DB.
func (s *settingsService) UploadLogo(file *multipart.FileHeader, updatedBy *uint) (*dto.SettingsResponse, error) {
	if file == nil {
		return nil, errors.New("File logo wajib diunggah")
	}

	// Validasi ukuran: maks 2MB
	if file.Size > 2*1024*1024 {
		return nil, errors.New("File logo tidak valid. Maksimal 2MB dengan format PNG, JPG, atau WEBP.")
	}

	// Validasi MIME type
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	header := make([]byte, 512)
	n, err := src.Read(header)
	if err != nil && err != io.EOF {
		return nil, err
	}
	// Reset reader agar bisa dibaca ulang saat copy
	if _, err := src.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	mimeType := http.DetectContentType(header[:n])
	switch mimeType {
	case "image/png", "image/jpeg", "image/webp":
	default:
		return nil, errors.New("File logo tidak valid. Maksimal 2MB dengan format PNG, JPG, atau WEBP.")
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
		return nil, err
	}

	dst := filepath.Join(s.uploadPath, filename)
	out, err := os.Create(dst)
	if err != nil {
		return nil, err
	}
	defer out.Close()

	if _, err := io.Copy(out, src); err != nil {
		return nil, err
	}

	logoPath := "/uploads/settings/" + filename
	if err := s.repo.UpdateLogo(logoPath, updatedBy); err != nil {
		return nil, err
	}

	return s.GetSettings()
}

func toResponse(settings model.Settings) *dto.SettingsResponse {
	var mission []string
	if settings.AboutMission != "" {
		err := json.Unmarshal([]byte(settings.AboutMission), &mission)
		if err != nil {
			helper.Logger.Error("Failed to unmarshal AboutMission", zap.Error(err))
		}
	}
	if mission == nil {
		mission = make([]string, 0)
	}

	return &dto.SettingsResponse{
		ID:                 settings.ID,
		SiteName:           settings.SiteName,
		Tagline:            settings.Tagline,
		LogoPath:           settings.LogoPath,
		ContactEmail:       settings.ContactEmail,
		ContactPhone:       settings.ContactPhone,
		Address:            settings.Address,
		MapsEmbedURL:       settings.MapsEmbedURL,
		FacebookURL:        settings.FacebookURL,
		InstagramURL:       settings.InstagramURL,
		YoutubeURL:         settings.YoutubeURL,
		VideoProfilePath:   settings.VideoProfilePath,
		History:            settings.History,
		AboutTutorial:      settings.AboutTutorial,
		AboutFormationDate: settings.AboutFormationDate,
		AboutNoSK:          settings.AboutNoSK,
		AboutVision:        settings.AboutVision,
		AboutMission:       mission,
		GreetingTitle:      settings.GreetingTitle,
		GreetingSubtitle:   settings.GreetingSubtitle,
		GreetingDate:       settings.GreetingDate,
		GreetingContent:    settings.GreetingContent,
		GreetingImagePath:  settings.GreetingImagePath,
		CreatedAt:          settings.CreatedAt,
		UpdatedAt:          settings.UpdatedAt,
		UpdatedBy:          settings.UpdatedBy,
	}
}
