package handler

import (
	"net/http"
	"strconv"

	"github.com/ahmadzakyarifin/schoolpay/backend/internal/helper"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/academicyear/dto"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/academicyear/service"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/validator"
	"github.com/gin-gonic/gin"
)

type AcademicYearHandler struct {
	s service.AcademicYearService
}

func NewAcademicYearHandler(s service.AcademicYearService) *AcademicYearHandler {
	return &AcademicYearHandler{s: s}
}

func (h *AcademicYearHandler) GetAll(c *gin.Context) {
	var req dto.AcademicYearQueryReq
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
		"search": req.Search,
		"status": req.Status,
		"sort":   req.Sort,
	}
	meta := helper.GetPaginationMeta(int(total), req.Page, req.Limit, filters)
	helper.SuccessResponse(c, http.StatusOK, "ACADEMIC_YEAR_LIST_RETRIEVED", "berhasil mengambil data tahun ajaran", res, meta)
}

func (h *AcademicYearHandler) Create(c *gin.Context) {
	var req dto.AcademicYearCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "validasi gagal", validator.Errors(err))
		return
	}

	res, err := h.s.Create(c.Request.Context(), req)
	if err != nil {
		helper.HandleServiceError(c, err)
		return
	}

	helper.SuccessResponse(c, http.StatusCreated, "ACADEMIC_YEAR_CREATED", "tahun ajaran berhasil dibuat", res, nil)
}

func (h *AcademicYearHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "ID tidak valid", nil)
		return
	}
	var req dto.AcademicYearUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "validasi gagal", validator.Errors(err))
		return
	}

	re, err := h.s.Update(c.Request.Context(), uint(id), req)
	if err != nil {
		helper.HandleServiceError(c, err)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "ACADEMIC_YEAR_UPDATED", "tahun ajaran berhasil diperbarui", re, nil)
}

func (h *AcademicYearHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "ID tidak valid", nil)
		return
	}
	if err := h.s.Delete(c.Request.Context(), uint(id)); err != nil {
		helper.HandleServiceError(c, err)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "ACADEMIC_YEAR_DELETED", "tahun ajaran berhasil dihapus", nil, nil)
}

func (h *AcademicYearHandler) Restore(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "ID tidak valid", nil)
		return
	}
	if err := h.s.Restore(c.Request.Context(), uint(id)); err != nil {
		helper.HandleServiceError(c, err)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "ACADEMIC_YEAR_RESTORED", "tahun ajaran berhasil dipulihkan", nil, nil)
}

func (h *AcademicYearHandler) ToggleStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "ID tidak valid", nil)
		return
	}
	if err := h.s.ToggleStatus(c.Request.Context(), uint(id)); err != nil {
		helper.HandleServiceError(c, err)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "ACADEMIC_YEAR_STATUS_TOGGLED", "status tahun ajaran berhasil diubah", nil, nil)
}

func (h *AcademicYearHandler) BulkDelete(c *gin.Context) {
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

	helper.SuccessResponse(c, http.StatusOK, "ACADEMIC_YEAR_BULK_DELETED", "tahun ajaran terpilih berhasil dihapus", nil, nil)
}

func (h *AcademicYearHandler) BulkRestore(c *gin.Context) {
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

	helper.SuccessResponse(c, http.StatusOK, "ACADEMIC_YEAR_BULK_RESTORED", "tahun ajaran terpilih berhasil dipulihkan", nil, nil)
}

func (h *AcademicYearHandler) GetDependencyInfo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "ID tidak valid", nil)
		return
	}
	info, err := h.s.GetDependencyInfo(c.Request.Context(), uint(id))
	if err != nil {
		helper.HandleServiceError(c, err)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "ACADEMIC_YEAR_DEPENDENCY_INFO_RETRIEVED", "berhasil mengambil info dependensi", info, nil)
}

func (h *AcademicYearHandler) CheckUnique(c *gin.Context) {
	name := c.Query("name")
	excludeID, _ := strconv.ParseUint(c.Query("exclude_id"), 10, 32)

	if name == "" {
		helper.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "name harus diisi", nil)
		return
	}

	isUnique, err := h.s.CheckUnique(c.Request.Context(), name, uint(excludeID))
	if err != nil {
		helper.HandleServiceError(c, err)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "ACADEMIC_YEAR_UNIQUE_CHECKED", "berhasil mengecek keunikan", map[string]bool{"is_unique": isUnique}, nil)
}
