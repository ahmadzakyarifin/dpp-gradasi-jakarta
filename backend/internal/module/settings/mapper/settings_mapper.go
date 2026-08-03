package mapper

import (
	"encoding/json"

	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/helper"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/settings/dto"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/settings/entity"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/settings/model"
	"go.uber.org/zap"
)

// ModelToEntity mengonversi model GORM menjadi entity.
func ModelToEntity(m *model.Settings) *entity.Settings {
	if m == nil {
		return nil
	}
	return &entity.Settings{
		ID:                 m.ID,
		SiteName:           m.SiteName,
		Tagline:            m.Tagline,
		LogoPath:           m.LogoPath,
		ContactEmail:       m.ContactEmail,
		ContactPhone:       m.ContactPhone,
		Address:            m.Address,
		MapsEmbedURL:       m.MapsEmbedURL,
		FacebookURL:        m.FacebookURL,
		InstagramURL:       m.InstagramURL,
		YoutubeURL:         m.YoutubeURL,
		VideoProfilePath:   m.VideoProfilePath,
		History:            m.History,
		AboutTutorial:      m.AboutTutorial,
		AboutFormationDate: m.AboutFormationDate,
		AboutNoSK:          m.AboutNoSK,
		AboutVision:        m.AboutVision,
		AboutMission:       m.AboutMission,
		GreetingTitle:      m.GreetingTitle,
		GreetingSubtitle:   m.GreetingSubtitle,
		GreetingDate:       m.GreetingDate,
		GreetingContent:    m.GreetingContent,
		GreetingImagePath:  m.GreetingImagePath,
		CreatedAt:          m.CreatedAt,
		UpdatedAt:          m.UpdatedAt,
		UpdatedBy:          m.UpdatedBy,
	}
}

// EntityToResponse mengonversi entity menjadi response DTO.
// AboutMission (JSON string) di-unmarshal menjadi []string.
func EntityToResponse(e *entity.Settings) *dto.SettingsResponse {
	if e == nil {
		return nil
	}
	var mission []string
	if e.AboutMission != "" {
		if err := json.Unmarshal([]byte(e.AboutMission), &mission); err != nil {
			helper.Logger.Error("Failed to unmarshal AboutMission", zap.Error(err))
		}
	}
	if mission == nil {
		mission = make([]string, 0)
	}

	return &dto.SettingsResponse{
		ID:                 e.ID,
		SiteName:           e.SiteName,
		Tagline:            e.Tagline,
		LogoPath:           e.LogoPath,
		ContactEmail:       e.ContactEmail,
		ContactPhone:       e.ContactPhone,
		Address:            e.Address,
		MapsEmbedURL:       e.MapsEmbedURL,
		FacebookURL:        e.FacebookURL,
		InstagramURL:       e.InstagramURL,
		YoutubeURL:         e.YoutubeURL,
		VideoProfilePath:   e.VideoProfilePath,
		History:            e.History,
		AboutTutorial:      e.AboutTutorial,
		AboutFormationDate: e.AboutFormationDate,
		AboutNoSK:          e.AboutNoSK,
		AboutVision:        e.AboutVision,
		AboutMission:       mission,
		GreetingTitle:      e.GreetingTitle,
		GreetingSubtitle:   e.GreetingSubtitle,
		GreetingDate:       e.GreetingDate,
		GreetingContent:    e.GreetingContent,
		GreetingImagePath:  e.GreetingImagePath,
		CreatedAt:          e.CreatedAt,
		UpdatedAt:          e.UpdatedAt,
		UpdatedBy:          e.UpdatedBy,
	}
}
