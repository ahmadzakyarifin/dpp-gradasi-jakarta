package handler

import (
	"net/http"
	"strconv"

	"github.com/ahmadzakyarifin/dpp-gradasi/backend/config"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/helper"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/kontak/dto"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/kontak/service"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/validator"
	"github.com/gin-gonic/gin"
)

type KontakHandler struct {
	svc service.KontakService
	cfg *config.Config
}

func NewKontakHandler(svc service.KontakService, cfg *config.Config) *KontakHandler {
	return &KontakHandler{svc: svc, cfg: cfg}
}

func (h *KontakHandler) Cfg() *config.Config {
	return h.cfg
}

// POST /api/v1/kontak — publik submit
func (h *KontakHandler) Submit(c *gin.Context) {
	var req dto.KontakRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "validasi gagal", validator.Errors(err))
		return
	}
	if err := h.svc.Submit(c.Request.Context(), &req); err != nil {
		helper.HandleServiceError(c, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "KONTAK_SUBMITTED", "Pesan Anda berhasil dikirim. Terima kasih!", nil, nil)
}

// GET /api/v1/kontak — admin list
func (h *KontakHandler) List(c *gin.Context) {
	var q dto.KontakQuery
	_ = c.ShouldBindQuery(&q)
	resp, err := h.svc.GetAll(c.Request.Context(), q)
	if err != nil {
		helper.HandleServiceError(c, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "KONTAK_LIST", "Daftar pesan masuk berhasil diambil.", resp, nil)
}

// GET /api/v1/kontak/:id — admin detail + mark as read
func (h *KontakHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "ID tidak valid", nil)
		return
	}
	resp, err := h.svc.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		helper.HandleServiceError(c, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "KONTAK_DETAIL", "Detail pesan berhasil diambil.", resp, nil)
}

// DELETE /api/v1/kontak/:id — admin
func (h *KontakHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "ID tidak valid", nil)
		return
	}
	if err := h.svc.Delete(c.Request.Context(), uint(id)); err != nil {
		helper.HandleServiceError(c, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "KONTAK_DELETED", "Pesan berhasil dihapus.", nil, nil)
}

// POST /api/v1/admin/kontak/:id/restore — admin
func (h *KontakHandler) Restore(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "ID tidak valid", nil)
		return
	}
	if err := h.svc.Restore(c.Request.Context(), uint(id)); err != nil {
		helper.HandleServiceError(c, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "KONTAK_RESTORED", "Pesan berhasil dipulihkan.", nil, nil)
}

// POST /api/v1/admin/kontak/bulk-delete — admin
func (h *KontakHandler) BulkDelete(c *gin.Context) {
	var req dto.BulkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "validasi gagal", validator.Errors(err))
		return
	}
	affected, err := h.svc.BulkDelete(c.Request.Context(), req.IDs)
	if err != nil {
		helper.HandleServiceError(c, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "KONTAK_BULK_DELETED", "Pesan berhasil dihapus massal.", dto.BulkResponse{Affected: affected}, nil)
}

// POST /api/v1/admin/kontak/bulk-restore — admin
func (h *KontakHandler) BulkRestore(c *gin.Context) {
	var req dto.BulkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "validasi gagal", validator.Errors(err))
		return
	}
	affected, err := h.svc.BulkRestore(c.Request.Context(), req.IDs)
	if err != nil {
		helper.HandleServiceError(c, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "KONTAK_BULK_RESTORED", "Pesan berhasil dipulihkan massal.", dto.BulkResponse{Affected: affected}, nil)
}
