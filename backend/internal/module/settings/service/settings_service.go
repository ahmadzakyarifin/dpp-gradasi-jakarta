package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/helper"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/settings/dto"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/settings/model"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/settings/repository"
	"go.uber.org/zap"
)

type SettingsService interface {
	GetSettings() (*dto.SettingsResponse, error)
	UpdateSettings(values map[string]interface{}) (*dto.SettingsResponse, error)
}

type settingsService struct {
	repo repository.SettingsRepo
}

func NewSettingsService(repo repository.SettingsRepo) SettingsService {
	return &settingsService{repo}
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
func (s *settingsService) UpdateSettings(values map[string]interface{}) (*dto.SettingsResponse, error) {
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
		return nil, errors.New(fmt.Sprintf("Validasi gagal: %v", errorsMap))
	}

	if len(updates) > 0 {
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
		SiteName:           settings.SiteName,
		Tagline:            settings.Tagline,
		LogoURL:            settings.LogoURL,
		ContactEmail:       settings.ContactEmail,
		ContactPhone:       settings.ContactPhone,
		Address:            settings.Address,
		MapsEmbedURL:       settings.MapsEmbedURL,
		FacebookURL:        settings.FacebookURL,
		InstagramURL:       settings.InstagramURL,
		YoutubeURL:         settings.YoutubeURL,
		VideoProfileURL:    settings.VideoProfileURL,
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
		GreetingImageURL:   settings.GreetingImageURL,
		CreatedAt:          settings.CreatedAt,
		UpdatedAt:          settings.UpdatedAt,
	}
}
