package handler

import (
	"net/http"
	"strconv"

	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/helper"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/middleware"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/sliders/dto"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/sliders/service"
	"github.com/gin-gonic/gin"
)

type SlidersHandler struct {
	svc service.SlidersService
}

func NewSlidersHandler(svc service.SlidersService) *SlidersHandler {
	return &SlidersHandler{svc: svc}
}

// GET /api/v1/sliders — publik
func (h *SlidersHandler) GetAll(c *gin.Context) {
	// Publik SELALU hanya melihat slider aktif
	resp, err := h.svc.GetAll(c.Request.Context(), true)
	if err != nil {
		svcErr, _ := err.(*helper.ServiceError)
		helper.ErrorResponse(c, http.StatusInternalServerError, svcErr.Code, svcErr.Message, nil)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "SLIDERS_FETCHED", "Data slider berhasil diambil.", resp, nil)
}

// GET /api/v1/sliders/admin — admin
func (h *SlidersHandler) GetAllAdmin(c *gin.Context) {
	activeOnly := c.DefaultQuery("active", "true") == "true"
	resp, err := h.svc.GetAll(c.Request.Context(), activeOnly)
	if err != nil {
		svcErr, _ := err.(*helper.ServiceError)
		helper.ErrorResponse(c, http.StatusInternalServerError, svcErr.Code, svcErr.Message, nil)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "SLIDERS_FETCHED", "Data slider berhasil diambil.", resp, nil)
}

// GET /api/v1/sliders/:id — publik
func (h *SlidersHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "ID slider tidak valid.", nil)
		return
	}
	resp, err := h.svc.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		svcErr, _ := err.(*helper.ServiceError)
		httpCode := http.StatusInternalServerError
		if svcErr.Code == "NOT_FOUND" {
			httpCode = http.StatusNotFound
		}
		helper.ErrorResponse(c, httpCode, svcErr.Code, svcErr.Message, nil)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "SLIDER_FETCHED", "Data slider berhasil diambil.", resp, nil)
}

// POST /api/v1/sliders — admin
func (h *SlidersHandler) Create(c *gin.Context) {
	var req dto.SliderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ValidationErrorResponse(c, []helper.ValidationErrorItem{
			{Field: "title", Tag: "required", Message: "Title dan image_path wajib diisi."},
		})
		return
	}
	userID, _ := middleware.GetUserID(c)
	resp, err := h.svc.Create(c.Request.Context(), &req, userID)
	if err != nil {
		svcErr, _ := err.(*helper.ServiceError)
		helper.ErrorResponse(c, http.StatusInternalServerError, svcErr.Code, svcErr.Message, nil)
		return
	}
	helper.SuccessResponse(c, http.StatusCreated, "SLIDER_CREATED", "Slider berhasil dibuat.", resp, nil)
}

// PUT /api/v1/sliders/:id — admin
func (h *SlidersHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "ID slider tidak valid.", nil)
		return
	}
	var req dto.SliderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ValidationErrorResponse(c, []helper.ValidationErrorItem{
			{Field: "title", Tag: "required", Message: "Title dan image_path wajib diisi."},
		})
		return
	}
	resp, err := h.svc.Update(c.Request.Context(), uint(id), &req)
	if err != nil {
		svcErr, _ := err.(*helper.ServiceError)
		httpCode := http.StatusInternalServerError
		if svcErr.Code == "NOT_FOUND" {
			httpCode = http.StatusNotFound
		}
		helper.ErrorResponse(c, httpCode, svcErr.Code, svcErr.Message, nil)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "SLIDER_UPDATED", "Slider berhasil diupdate.", resp, nil)
}

// DELETE /api/v1/sliders/:id — admin
func (h *SlidersHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "ID slider tidak valid.", nil)
		return
	}
	if err := h.svc.Delete(c.Request.Context(), uint(id)); err != nil {
		svcErr, _ := err.(*helper.ServiceError)
		httpCode := http.StatusInternalServerError
		if svcErr.Code == "NOT_FOUND" {
			httpCode = http.StatusNotFound
		}
		helper.ErrorResponse(c, httpCode, svcErr.Code, svcErr.Message, nil)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "SLIDER_DELETED", "Slider berhasil dihapus.", nil, nil)
}

// PUT /api/v1/sliders/reorder — admin
func (h *SlidersHandler) Reorder(c *gin.Context) {
	var req dto.ReorderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ValidationErrorResponse(c, []helper.ValidationErrorItem{
			{Field: "ids", Tag: "required", Message: "Daftar ID wajib dikirim."},
		})
		return
	}

	if err := h.svc.Reorder(c.Request.Context(), req.IDs); err != nil {
		svcErr, _ := err.(*helper.ServiceError)
		helper.ErrorResponse(c, http.StatusInternalServerError, svcErr.Code, svcErr.Message, nil)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "SLIDERS_REORDERED", "Urutan slider berhasil diperbarui.", nil, nil)
}

// POST /api/v1/sliders/:id/restore — admin
func (h *SlidersHandler) Restore(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "INVALID_ID", "ID slider tidak valid.", nil)
		return
	}
	if err := h.svc.Restore(c.Request.Context(), uint(id)); err != nil {
		svcErr, _ := err.(*helper.ServiceError)
		helper.ErrorResponse(c, http.StatusInternalServerError, svcErr.Code, svcErr.Message, nil)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "SLIDER_RESTORED", "Slider berhasil dipulihkan.", nil, nil)
}

// POST /api/v1/sliders/bulk-delete — admin (bulk soft delete)
func (h *SlidersHandler) BulkDelete(c *gin.Context) {
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
	helper.SuccessResponse(c, http.StatusOK, "SLIDER_BULK_DELETED", "Slider berhasil dihapus massal.", nil, nil)
}

// POST /api/v1/sliders/bulk-restore — admin (bulk restore)
func (h *SlidersHandler) BulkRestore(c *gin.Context) {
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
	helper.SuccessResponse(c, http.StatusOK, "SLIDER_BULK_RESTORED", "Slider berhasil dipulihkan massal.", nil, nil)
}
