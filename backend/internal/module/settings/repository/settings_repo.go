package repository

import (
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/settings/model"
	"gorm.io/gorm"
)

type SettingsRepo interface {
	Get() (*model.Settings, error)
	Update(settings map[string]interface{}) error
	UpdateLogo(logoPath string, updatedBy *uint) error
}

type settingsRepo struct {
	db *gorm.DB
}

func NewSettingsRepo(db *gorm.DB) SettingsRepo {
	return &settingsRepo{db}
}

func (r *settingsRepo) Get() (*model.Settings, error) {
	var settings model.Settings
	err := r.db.First(&settings).Error
	if err == gorm.ErrRecordNotFound {
		// Seed default if empty
		settings = model.Settings{
			SiteName:     "DPP GRADASI",
			ContactEmail: "admin@gradasi.org",
			ContactPhone: "+628123456789",
		}
		r.db.Create(&settings)
		return &settings, nil
	}
	return &settings, err
}

func (r *settingsRepo) Update(values map[string]interface{}) error {
	if len(values) == 0 {
		return nil
	}
	return r.db.Model(&model.Settings{}).Where("id = ?", 1).Updates(values).Error
}

func (r *settingsRepo) UpdateLogo(logoPath string, updatedBy *uint) error {
	updates := map[string]interface{}{
		"logo_path": logoPath,
	}
	if updatedBy != nil {
		updates["updated_by"] = *updatedBy
	}
	return r.db.Model(&model.Settings{}).Where("id = ?", 1).Updates(updates).Error
}
