package handler

import (
	"net/http"
	"strconv"

	"github.com/ahmadzakyarifin/schoolpay/backend/internal/helper"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/activeclass/dto"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/activeclass/service"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/validator"
	"github.com/gin-gonic/gin"
)

type ActiveClassHandler struct {
	s service.ActiveClassService
}

func NewActiveClassHandler(s service.ActiveClassService) *ActiveClassHandler {
	return &ActiveClassHandler{s: s}
}

func (h *ActiveClassHandler) GetAll(c *gin.Context) {
	var req dto.ActiveClassQueryReq
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
	helper.SuccessResponse(c, http.StatusOK, "ACTIVE_CLASS_LIST_RETRIEVED", "berhasil mengambil data kelas aktif", res, meta)
}

func (h *ActiveClassHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, 400, "VALIDATION_ERROR", "ID tidak valid", nil)
		return
	}
	res, err := h.s.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		helper.HandleServiceError(c, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "ACTIVE_CLASS_DETAIL", "berhasil mengambil detail kelas aktif", gin.H{"data": res}, nil)
}

func (h *ActiveClassHandler) Create(c *gin.Context) {
	var req dto.ActiveClassCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "validasi gagal", validator.Errors(err))
		return
	}
	res, err := h.s.Create(c.Request.Context(), req)
	if err != nil {
		helper.HandleServiceError(c, err)
		return
	}
	helper.SuccessResponse(c, http.StatusCreated, "ACTIVE_CLASS_CREATED", "kelas aktif berhasil dibuat", gin.H{"data": res}, nil)
}

func (h *ActiveClassHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, 400, "VALIDATION_ERROR", "ID tidak valid", nil)
		return
	}
	var req dto.ActiveClassUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "validasi gagal", validator.Errors(err))
		return
	}
	res, err := h.s.Update(c.Request.Context(), uint(id), req)
	if err != nil {
		helper.HandleServiceError(c, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "ACTIVE_CLASS_UPDATED", "kelas aktif berhasil diperbarui", gin.H{"data": res}, nil)
}

func (h *ActiveClassHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, 400, "VALIDATION_ERROR", "ID tidak valid", nil)
		return
	}
	if err := h.s.Delete(c.Request.Context(), uint(id)); err != nil {
		helper.HandleServiceError(c, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "ACTIVE_CLASS_DELETED", "kelas aktif berhasil dihapus", nil, nil)
}

func (h *ActiveClassHandler) ToggleStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, 400, "VALIDATION_ERROR", "ID tidak valid", nil)
		return
	}
	if err := h.s.ToggleStatus(c.Request.Context(), uint(id)); err != nil {
		helper.HandleServiceError(c, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "ACTIVE_CLASS_STATUS_TOGGLED", "status kelas aktif berhasil diubah", nil, nil)
}

func (h *ActiveClassHandler) BulkUpsertByYear(c *gin.Context) {
	ayID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "ID tahun ajaran tidak valid", nil)
		return
	}
	var req struct {
		ActiveClasses []dto.BulkUpsertItem `json:"active_classes" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "validasi gagal", validator.Errors(err))
		return
	}
	res, err := h.s.BulkUpsertByYear(c.Request.Context(), uint(ayID), req.ActiveClasses)
	if err != nil {
		helper.HandleServiceError(c, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "ACTIVE_CLASS_BULK_UPSERTED", "kelas aktif berhasil disimpan", gin.H{"data": res}, nil)
}
