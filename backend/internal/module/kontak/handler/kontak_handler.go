package handler

import (
	"net/http"
	"strconv"

	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/helper"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/kontak/dto"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/kontak/service"
	"github.com/gin-gonic/gin"
)

type KontakHandler struct {
	svc service.KontakService
}

func NewKontakHandler(svc service.KontakService) *KontakHandler {
	return &KontakHandler{svc: svc}
}

// POST /api/v1/kontak — publik submit
func (h *KontakHandler) Submit(c *gin.Context) {
	var req dto.KontakRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ValidationErrorResponse(c, []helper.ValidationErrorItem{
			{Field: "nama", Tag: "required", Message: "Nama, email, subjek, dan pesan wajib diisi."},
		})
		return
	}
	if err := h.svc.Submit(c.Request.Context(), &req); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "SERVER_ERROR", "Gagal mengirim pesan.", nil)
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
		svcErr, _ := err.(*helper.ServiceError)
		helper.ErrorResponse(c, http.StatusInternalServerError, svcErr.Code, svcErr.Message, nil)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "KONTAK_LIST", "Daftar pesan masuk berhasil diambil.", resp, nil)
}

// GET /api/v1/kontak/:id — admin detail + mark as read
func (h *KontakHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "ID tidak valid.", nil)
		return
	}
	resp, err := h.svc.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		svcErr, _ := err.(*helper.ServiceError)
		code := http.StatusInternalServerError
		if svcErr.Code == "NOT_FOUND" {
			code = http.StatusNotFound
		}
		helper.ErrorResponse(c, code, svcErr.Code, svcErr.Message, nil)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "KONTAK_DETAIL", "Detail pesan berhasil diambil.", resp, nil)
}

// DELETE /api/v1/kontak/:id — admin
func (h *KontakHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "ID tidak valid.", nil)
		return
	}
	if err := h.svc.Delete(c.Request.Context(), uint(id)); err != nil {
		svcErr, _ := err.(*helper.ServiceError)
		code := http.StatusInternalServerError
		if svcErr.Code == "NOT_FOUND" {
			code = http.StatusNotFound
		}
		helper.ErrorResponse(c, code, svcErr.Code, svcErr.Message, nil)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "KONTAK_DELETED", "Pesan berhasil dihapus.", nil, nil)
}

// POST /api/v1/admin/kontak/:id/restore — admin
func (h *KontakHandler) Restore(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "ID tidak valid.", nil)
		return
	}
	if err := h.svc.Restore(c.Request.Context(), uint(id)); err != nil {
		svcErr, _ := err.(*helper.ServiceError)
		helper.ErrorResponse(c, http.StatusInternalServerError, svcErr.Code, svcErr.Message, nil)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "KONTAK_RESTORED", "Pesan berhasil dipulihkan.", nil, nil)
}

// POST /api/v1/admin/kontak/bulk-delete — admin
func (h *KontakHandler) BulkDelete(c *gin.Context) {
	var req dto.BulkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ValidationErrorResponse(c, []helper.ValidationErrorItem{
			{Field: "ids", Tag: "required", Message: "IDs wajib diisi minimal 1."},
		})
		return
	}
	if err := h.svc.BulkDelete(c.Request.Context(), req.IDs); err != nil {
		svcErr, _ := err.(*helper.ServiceError)
		helper.ErrorResponse(c, http.StatusInternalServerError, svcErr.Code, svcErr.Message, nil)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "KONTAK_BULK_DELETED", "Pesan berhasil dihapus massal.", nil, nil)
}

// POST /api/v1/admin/kontak/bulk-restore — admin
func (h *KontakHandler) BulkRestore(c *gin.Context) {
	var req dto.BulkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ValidationErrorResponse(c, []helper.ValidationErrorItem{
			{Field: "ids", Tag: "required", Message: "IDs wajib diisi minimal 1."},
		})
		return
	}
	if err := h.svc.BulkRestore(c.Request.Context(), req.IDs); err != nil {
		svcErr, _ := err.(*helper.ServiceError)
		helper.ErrorResponse(c, http.StatusInternalServerError, svcErr.Code, svcErr.Message, nil)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "KONTAK_BULK_RESTORED", "Pesan berhasil dipulihkan massal.", nil, nil)
}
