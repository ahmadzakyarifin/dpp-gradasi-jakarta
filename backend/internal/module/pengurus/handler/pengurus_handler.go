package handler

import (
	"net/http"
	"strconv"

	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/helper"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/pengurus/dto"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/pengurus/service"
	"github.com/gin-gonic/gin"
)

type PengurusHandler struct {
	svc service.PengurusService
}

func NewPengurusHandler(svc service.PengurusService) *PengurusHandler {
	return &PengurusHandler{svc: svc}
}

// GET /api/v1/pengurus
func (h *PengurusHandler) ListPublic(c *gin.Context) {
	var q dto.PengurusQuery
	_ = c.ShouldBindQuery(&q)

	resp, err := h.svc.GetAllPublic(c.Request.Context(), q)
	if err != nil {
		helper.HandleServiceError(c, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "PENGURUS_LIST", "Daftar pengurus berhasil diambil", resp, nil)
}

// GET /api/v1/pengurus/regions
func (h *PengurusHandler) Regions(c *gin.Context) {
	resp, err := h.svc.GetRegions(c.Request.Context())
	if err != nil {
		helper.HandleServiceError(c, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "PENGURUS_REGIONS", "Daftar wilayah berhasil diambil", resp, nil)
}

// GET /api/v1/admin/pengurus
func (h *PengurusHandler) ListAdmin(c *gin.Context) {
	var q dto.PengurusQuery
	_ = c.ShouldBindQuery(&q)

	resp, err := h.svc.GetAllAdmin(c.Request.Context(), q)
	if err != nil {
		helper.HandleServiceError(c, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "PENGURUS_LIST", "Daftar pengurus berhasil diambil", resp, nil)
}

// POST /api/v1/admin/pengurus
func (h *PengurusHandler) Create(c *gin.Context) {
	var req dto.PengurusRequest
	if err := c.ShouldBind(&req); err != nil { // ShouldBind for form-data
		helper.ValidationErrorResponse(c, []helper.ValidationErrorItem{
			{Field: "form", Tag: "invalid", Message: "Data form tidak valid atau tidak lengkap."},
		})
		return
	}

	// Validate required file manually since gin binding file is tricky sometimes
	file, err := c.FormFile("image")
	if err != nil {
		helper.ValidationErrorResponse(c, []helper.ValidationErrorItem{
			{Field: "image", Tag: "required", Message: "Image wajib diunggah."},
		})
		return
	}
	req.Image = file

	resp, err := h.svc.Create(c.Request.Context(), &req)
	if err != nil {
		helper.HandleServiceError(c, err)
		return
	}
	helper.SuccessResponse(c, http.StatusCreated, "PENGURUS_CREATED", "Pengurus berhasil ditambahkan", resp, nil)
}

// PUT /api/v1/admin/pengurus/:id
func (h *PengurusHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "ID pengurus tidak valid.", nil)
		return
	}

	var req dto.PengurusRequest
	if err := c.ShouldBind(&req); err != nil {
		helper.ValidationErrorResponse(c, []helper.ValidationErrorItem{
			{Field: "form", Tag: "invalid", Message: "Data form tidak valid atau tidak lengkap."},
		})
		return
	}

	// Optional file upload
	file, _ := c.FormFile("image")
	if file != nil {
		req.Image = file
	}

	resp, err := h.svc.Update(c.Request.Context(), uint(id), &req)
	if err != nil {
		helper.HandleServiceError(c, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "PENGURUS_UPDATED", "Data pengurus berhasil diperbarui", resp, nil)
}

// DELETE /api/v1/admin/pengurus/:id
func (h *PengurusHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "ID pengurus tidak valid.", nil)
		return
	}
	if err := h.svc.Delete(c.Request.Context(), uint(id)); err != nil {
		helper.HandleServiceError(c, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "PENGURUS_DELETED", "Pengurus berhasil dihapus", nil, nil)
}

// POST /api/v1/admin/pengurus/:id/restore
func (h *PengurusHandler) Restore(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "ID pengurus tidak valid.", nil)
		return
	}
	if err := h.svc.Restore(c.Request.Context(), uint(id)); err != nil {
		helper.HandleServiceError(c, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "PENGURUS_RESTORED", "Pengurus berhasil dipulihkan", nil, nil)
}

// POST /api/v1/admin/pengurus/bulk-delete — admin (bulk soft delete)
func (h *PengurusHandler) BulkDelete(c *gin.Context) {
	var req dto.BulkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ValidationErrorResponse(c, []helper.ValidationErrorItem{
			{Field: "ids", Tag: "required", Message: "IDs wajib diisi minimal 1."},
		})
		return
	}
	if err := h.svc.BulkDelete(c.Request.Context(), req.IDs); err != nil {
		helper.HandleServiceError(c, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "PENGURUS_BULK_DELETED", "Pengurus berhasil dihapus massal", nil, nil)
}

// POST /api/v1/admin/pengurus/bulk-restore — admin (bulk restore)
func (h *PengurusHandler) BulkRestore(c *gin.Context) {
	var req dto.BulkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ValidationErrorResponse(c, []helper.ValidationErrorItem{
			{Field: "ids", Tag: "required", Message: "IDs wajib diisi minimal 1."},
		})
		return
	}
	if err := h.svc.BulkRestore(c.Request.Context(), req.IDs); err != nil {
		helper.HandleServiceError(c, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "PENGURUS_BULK_RESTORED", "Pengurus berhasil dipulihkan massal", nil, nil)
}
