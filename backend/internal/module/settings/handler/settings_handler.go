package handler

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/helper"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/middleware"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/settings/dto"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/settings/service"
	"github.com/gin-gonic/gin"
)

type SettingsHandler interface {
	GetSettings(c *gin.Context)
	UpdateSettings(c *gin.Context)
	UploadLogo(c *gin.Context)
}

type settingsHandler struct {
	service service.SettingsService
}

func NewSettingsHandler(service service.SettingsService) SettingsHandler {
	return &settingsHandler{service}
}

func (h *settingsHandler) GetSettings(c *gin.Context) {
	settings, err := h.service.GetSettings(c.Request.Context())
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

	settings, err := h.service.UpdateSettings(c.Request.Context(), requestBody, h.getUpdatedBy(c))
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

// UploadLogo menerima multipart file "logo", memvalidasi ukuran & MIME type,
// menyimpan ke public/uploads/settings, dan memperbarui logo_path di settings.
func (h *settingsHandler) UploadLogo(c *gin.Context) {
	file, err := c.FormFile("logo")
	if err != nil {
		helper.ValidationErrorResponse(c, []helper.ValidationErrorItem{
			{Field: "logo", Tag: "required", Message: "File logo wajib diunggah."},
		})
		return
	}

	settings, err := h.service.UploadLogo(c.Request.Context(), file, h.getUpdatedBy(c))
	if err != nil {
		msg := err.Error()
		if strings.Contains(msg, "tidak valid") || strings.Contains(msg, "wajib") {
			helper.ErrorResponse(c, http.StatusBadRequest, "INVALID_LOGO", msg, nil)
			return
		}
		helper.ErrorResponse(c, http.StatusInternalServerError, "UPLOAD_LOGO_ERROR", msg, nil)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "SETTINGS_LOGO_UPLOADED", "Logo berhasil diunggah", settings, nil)
}

// getUpdatedBy membaca user_id dari context (di-set AuthMiddleware).
func (h *settingsHandler) getUpdatedBy(c *gin.Context) *uint {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return nil
	}
	return &userID
}
