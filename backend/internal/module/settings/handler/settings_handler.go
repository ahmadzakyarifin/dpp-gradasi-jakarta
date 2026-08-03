package handler

import (
	"net/http"

	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/helper"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/middleware"
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
		helper.HandleServiceError(c, err)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "SETTINGS_RETRIEVED", "Konfigurasi website berhasil diambil", settings, nil)
}

// UpdateSettings menerima JSON object dengan key snake_case.
// Validasi (key dikenal, tipe string/null, format email/url) ditangani service
// via ValidationError → 422. Error lainnya → HandleServiceError.
func (h *settingsHandler) UpdateSettings(c *gin.Context) {
	var requestBody map[string]interface{}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		v := helper.NewValidationError()
		v.Add("body", "Invalid JSON body")
		helper.HandleServiceError(c, v)
		return
	}

	if len(requestBody) == 0 {
		v := helper.NewValidationError()
		v.Add("body", "Request body tidak boleh kosong")
		helper.HandleServiceError(c, v)
		return
	}

	settings, err := h.service.UpdateSettings(c.Request.Context(), requestBody, h.getUpdatedBy(c))
	if err != nil {
		helper.HandleServiceError(c, err)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "SETTINGS_UPDATED", "Konfigurasi website berhasil diperbarui", settings, nil)
}

// UploadLogo menerima multipart file "logo", memvalidasi ukuran & MIME type,
// menyimpan ke public/uploads/settings, dan memperbarui logo_path di settings.
func (h *settingsHandler) UploadLogo(c *gin.Context) {
	file, err := c.FormFile("logo")
	if err != nil {
		v := helper.NewValidationError()
		v.Add("logo", "File logo wajib diunggah.")
		helper.HandleServiceError(c, v)
		return
	}

	settings, err := h.service.UploadLogo(c.Request.Context(), file, h.getUpdatedBy(c))
	if err != nil {
		helper.HandleServiceError(c, err)
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
