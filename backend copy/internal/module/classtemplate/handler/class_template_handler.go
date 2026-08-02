package handler

import (
	"net/http"
	"strconv"

	"github.com/ahmadzakyarifin/schoolpay/backend/internal/helper"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/classtemplate/dto"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/classtemplate/service"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/validator"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	s service.ClassTemplateService
}

func NewHandler(s service.ClassTemplateService) *Handler {
	return &Handler{s: s}
}

func (h *Handler) GetAll(c *gin.Context) {
	var req dto.ClassTemplateQueryReq
	if err := c.ShouldBindQuery(&req); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "validasi gagal", validator.Errors(err))
		return
	}

	res, total, err := h.s.GetAll(c.Request.Context(), req)
	if err != nil {
		helper.HandleServiceError(c, err)
		return
	}

	filters := map[string]any{
		"search":      req.Search,
		"status":      req.Status,
		"major_id":    req.MajorID,
		"grade_level": req.GradeLevel,
		"sort":        req.Sort,
	}
	meta := helper.GetPaginationMeta(int(total), req.Page, req.Limit, filters)
	helper.SuccessResponse(c, http.StatusOK, "CLASS_TEMPLATE_LIST", "berhasil mengambil data template kelas", res, meta)
}

func (h *Handler) Create(c *gin.Context) {
	var req dto.ClassTemplateCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "validasi gagal", validator.Errors(err))
		return
	}

	res, err := h.s.Create(c.Request.Context(), req)
	if err != nil {
		helper.HandleServiceError(c, err)
		return
	}

	helper.SuccessResponse(c, http.StatusCreated, "CLASS_TEMPLATE_CREATED", "template kelas berhasil dibuat", gin.H{"data": res}, nil)
}

func (h *Handler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "ID tidak valid", nil)
		return
	}

	var req dto.ClassTemplateUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "validasi gagal", validator.Errors(err))
		return
	}

	res, err := h.s.Update(c.Request.Context(), uint(id), req)
	if err != nil {
		helper.HandleServiceError(c, err)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "CLASS_TEMPLATE_UPDATED", "template kelas berhasil diperbarui", gin.H{"data": res}, nil)
}

func (h *Handler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "ID tidak valid", nil)
		return
	}

	if err := h.s.Delete(c.Request.Context(), uint(id)); err != nil {
		helper.HandleServiceError(c, err)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "CLASS_TEMPLATE_DELETED", "template kelas berhasil dihapus", nil, nil)
}

func (h *Handler) Restore(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "ID tidak valid", nil)
		return
	}

	if err := h.s.Restore(c.Request.Context(), uint(id)); err != nil {
		helper.HandleServiceError(c, err)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "CLASS_TEMPLATE_RESTORED", "template kelas berhasil dipulihkan", nil, nil)
}

func (h *Handler) ToggleStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "ID tidak valid", nil)
		return
	}

	if err := h.s.ToggleStatus(c.Request.Context(), uint(id)); err != nil {
		helper.HandleServiceError(c, err)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "CLASS_TEMPLATE_STATUS_UPDATED", "status template kelas berhasil diubah", nil, nil)
}

func (h *Handler) BulkDelete(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "validasi gagal", validator.Errors(err))
		return
	}

	if err := h.s.BulkDelete(c.Request.Context(), req.IDs); err != nil {
		helper.HandleServiceError(c, err)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "CLASS_TEMPLATE_BULK_DELETED", "template kelas terpilih berhasil dihapus", nil, nil)
}

func (h *Handler) BulkRestore(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "validasi gagal", validator.Errors(err))
		return
	}

	if err := h.s.BulkRestore(c.Request.Context(), req.IDs); err != nil {
		helper.HandleServiceError(c, err)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "CLASS_TEMPLATE_BULK_RESTORED", "template kelas terpilih berhasil dipulihkan", nil, nil)
}

func (h *Handler) GetDependencyInfo(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "ID tidak valid", nil)
		return
	}

	info, err := h.s.GetDependencyInfo(c.Request.Context(), uint(id))
	if err != nil {
		helper.HandleServiceError(c, err)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "CLASS_TEMPLATE_DEPENDENCY_INFO", "berhasil mengambil info dependensi", gin.H{"data": info}, nil)
}

func (h *Handler) CheckUnique(c *gin.Context) {
	name := c.Query("name")
	majorIDStr := c.Query("major_id")
	excludeIDStr := c.Query("exclude_id")

	var majorID *uint
	if majorIDStr != "" && majorIDStr != "0" {
		parsed, _ := strconv.ParseUint(majorIDStr, 10, 32)
		val := uint(parsed)
		majorID = &val
	}

	var excludeID uint
	if excludeIDStr != "" && excludeIDStr != "0" {
		parsed, _ := strconv.ParseUint(excludeIDStr, 10, 32)
		excludeID = uint(parsed)
	}

	isUnique, err := h.s.CheckUnique(c.Request.Context(), name, majorID, excludeID)
	if err != nil {
		helper.HandleServiceError(c, err)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "CLASS_TEMPLATE_UNIQUE_CHECK", "berhasil mengecek keunikan", gin.H{"is_unique": isUnique}, nil)
}

func (h *Handler) SuggestNextName(c *gin.Context) {
	name := c.Query("name")
	suggested, err := h.s.SuggestNextName(c.Request.Context(), name)
	if err != nil {
		helper.HandleServiceError(c, err)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "CLASS_TEMPLATE_SUGGEST_NEXT_NAME", "saran nama berhasil", gin.H{"suggested_name": suggested}, nil)
}
