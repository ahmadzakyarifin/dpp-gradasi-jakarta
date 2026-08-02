package handler

import (
	"net/http"
	"strconv"

	"github.com/ahmadzakyarifin/schoolpay/backend/internal/helper"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/semester/dto"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/semester/service"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/validator"
	"github.com/gin-gonic/gin"
)

type SemesterHandler struct {
	s service.SemesterService
}

func NewSemesterHandler(s service.SemesterService) *SemesterHandler {
	return &SemesterHandler{s: s}
}

func (h *SemesterHandler) GetAll(c *gin.Context) {
	var req dto.SemesterQueryReq
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
		"search":           req.Search,
		"status":           req.Status,
		"academic_year_id": req.AcademicYearID,
		"sort":             req.Sort,
	}
	meta := helper.GetPaginationMeta(int(total), req.Page, req.Limit, filters)
	helper.SuccessResponse(c, http.StatusOK, "SEMESTER_LIST_RETRIEVED", "berhasil mengambil data semester", res, meta)
}

func (h *SemesterHandler) Create(c *gin.Context) {
	var req dto.SemesterCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "validasi gagal", validator.Errors(err))
		return
	}
	res, err := h.s.Create(c.Request.Context(), req)
	if err != nil {
		helper.HandleServiceError(c, err)
		return
	}
	helper.SuccessResponse(c, http.StatusCreated, "SEMESTER_CREATED", "semester berhasil dibuat", res, nil)
}

func (h *SemesterHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "ID tidak valid", nil)
		return
	}
	var req dto.SemesterUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "validasi gagal", validator.Errors(err))
		return
	}
	res, err := h.s.Update(c.Request.Context(), uint(id), req)
	if err != nil {
		helper.HandleServiceError(c, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "SEMESTER_UPDATED", "semester berhasil diperbarui", res, nil)
}

func (h *SemesterHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "ID tidak valid", nil)
		return
	}
	if err := h.s.Delete(c.Request.Context(), uint(id)); err != nil {
		helper.HandleServiceError(c, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "SEMESTER_DELETED", "semester berhasil dihapus", nil, nil)
}

func (h *SemesterHandler) Restore(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "ID tidak valid", nil)
		return
	}
	if err := h.s.Restore(c.Request.Context(), uint(id)); err != nil {
		helper.HandleServiceError(c, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "SEMESTER_RESTORED", "semester berhasil dipulihkan", nil, nil)
}

func (h *SemesterHandler) ToggleStatus(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "ID tidak valid", nil)
		return
	}
	if err := h.s.ToggleStatus(c.Request.Context(), uint(id)); err != nil {
		helper.HandleServiceError(c, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "SEMESTER_STATUS_TOGGLED", "status semester berhasil diubah", nil, nil)
}

func (h *SemesterHandler) BulkDelete(c *gin.Context) {
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
	helper.SuccessResponse(c, http.StatusOK, "SEMESTER_BULK_DELETED", "semester terpilih berhasil dihapus", nil, nil)
}

func (h *SemesterHandler) BulkRestore(c *gin.Context) {
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
	helper.SuccessResponse(c, http.StatusOK, "SEMESTER_BULK_RESTORED", "semester terpilih berhasil dipulihkan", nil, nil)
}

func (h *SemesterHandler) GetDependencyInfo(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "ID tidak valid", nil)
		return
	}
	info, err := h.s.GetDependencyInfo(c.Request.Context(), uint(id))
	if err != nil {
		helper.HandleServiceError(c, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "SEMESTER_DEPENDENCY_INFO", "berhasil mengambil info dependensi", gin.H{"data": info}, nil)
}

func (h *SemesterHandler) CheckUnique(c *gin.Context) {
	ayID, _ := strconv.ParseUint(c.Query("academic_year_id"), 10, 32)
	name := c.Query("name")
	excludeID, _ := strconv.ParseUint(c.Query("exclude_id"), 10, 64)
	isUnique, err := h.s.CheckUnique(c.Request.Context(), uint(ayID), name, uint(excludeID))
	if err != nil {
		helper.HandleServiceError(c, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "SEMESTER_UNIQUE_CHECK", "berhasil mengecek keunikan", gin.H{"is_unique": isUnique}, nil)
}
