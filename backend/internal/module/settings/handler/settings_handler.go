package handler

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/helper"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/settings/dto"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/settings/service"
	"github.com/gin-gonic/gin"
)

type SettingsHandler interface {
	GetSettings(c *gin.Context)
	UpdateSettings(c *gin.Context)
}

type settingsHandler struct {
	service service.SettingsService
}

func NewSettingsHandler(service service.SettingsService) SettingsHandler {
	return &settingsHandler{service}
}

func (h *settingsHandler) GetSettings(c *gin.Context) {
	settings, err := h.service.GetSettings()
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "GET_SETTINGS_ERROR", err.Error(), nil)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "SETTINGS_RETRIEVED", "Konfigurasi website berhasil diambil", settings, nil)
}

// UpdateSettings menerima JSON object dengan key snake_case (JSON tag model).
// Validasi: key harus dikenal, tipe string/null, email & url divalidasi format.
func (h *settingsHandler) UpdateSettings(c *gin.Context) {
	var requestBody map[string]interface{}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		valErrors := helper.GenerateValidationMessage(map[string]string{"body": "Invalid JSON body"})
		helper.ValidationErrorResponse(c, valErrors)
		return
	}

	if len(requestBody) == 0 {
		helper.ValidationErrorResponse(c, helper.GenerateValidationMessage(map[string]string{"body": "Request body tidak boleh kosong"}))
		return
	}

	// Kumpulkan key yang valid dari JSON tag model.Settings (lewat SettingsResponse DTO sebagai acuan kontrak)
	validKeys := make(map[string]bool)
	respType := reflect.TypeOf(dto.SettingsResponse{})
	for i := 0; i < respType.NumField(); i++ {
		tag := respType.Field(i).Tag.Get("json")
		jsonName := strings.Split(tag, ",")[0]
		if jsonName != "" && jsonName != "-" {
			validKeys[jsonName] = true
		}
	}

	errorsMap := make(map[string]string)
	for key := range requestBody {
		if !validKeys[key] {
			errorsMap[key] = "Field tidak dikenal"
		}
	}

	if len(errorsMap) > 0 {
		valErrors := helper.GenerateValidationMessage(errorsMap)
		helper.ValidationErrorResponse(c, valErrors)
		return
	}

	settings, err := h.service.UpdateSettings(requestBody)
	if err != nil {
		// Check if it's a validation error from service
		if strings.Contains(err.Error(), "Validasi gagal:") {
			helper.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error(), nil)
		} else {
			helper.ErrorResponse(c, http.StatusInternalServerError, "UPDATE_SETTINGS_ERROR", err.Error(), nil)
		}
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "SETTINGS_UPDATED", "Konfigurasi website berhasil diperbarui", settings, nil)
}
