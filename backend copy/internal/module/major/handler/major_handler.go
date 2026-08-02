package handler

import (
	"net/http"
	"strconv"

	"github.com/ahmadzakyarifin/schoolpay/backend/internal/helper"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/major/dto"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/major/mapper"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/major/service"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/validator"
	"github.com/gin-gonic/gin"
)

type MajorHandler struct {
	s service.MajorService
}

func NewMajorHandler(s service.MajorService) *MajorHandler {
	return &MajorHandler{s: s}
}

func (h *MajorHandler) Create(c *gin.Context) {
	var req dto.MajorCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "validasi gagal", validator.Errors(err))
		return
	}

	entity := mapper.CreateReqToEntity(&req)
	if err := h.s.Create(c.Request.Context(), entity); err != nil {
		helper.HandleServiceError(c, err)
		return
	}

	helper.SuccessResponse(c, http.StatusCreated, "MAJOR_CREATED", "jurusan berhasil dibuat", mapper.EntityToResponse(entity), nil)
}

func (h *MajorHandler) GetAll(c *gin.Context) {
	var req dto.MajorQueryReq
	if err := c.ShouldBindQuery(&req); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "validasi gagal", validator.Errors(err))
		return
	}

	req.Normalize()

	list, total, err := h.s.GetAll(c.Request.Context(), req)
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
	helper.SuccessResponse(c, http.StatusOK, "MAJOR_LIST", "berhasil mengambil data jurusan", mapper.EntityListToResponse(list), meta)
}

func (h *MajorHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "ID tidak valid", nil)
		return
	}
	var req dto.MajorUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "validasi gagal", validator.Errors(err))
		return
	}

	updatedEntity, err := h.s.Update(c.Request.Context(), uint(id), &req)
	if err != nil {
		helper.HandleServiceError(c, err)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "MAJOR_UPDATED", "jurusan berhasil diperbarui", mapper.EntityToResponse(updatedEntity), nil)
}

func (h *MajorHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "ID tidak valid", nil)
		return
	}
	if err := h.s.Delete(c.Request.Context(), uint(id)); err != nil {
		helper.HandleServiceError(c, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "MAJOR_DELETED", "jurusan berhasil dihapus", nil, nil)
}

func (h *MajorHandler) Restore(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "ID tidak valid", nil)
		return
	}
	if err := h.s.Restore(c.Request.Context(), uint(id)); err != nil {
		helper.HandleServiceError(c, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "MAJOR_RESTORED", "jurusan berhasil dipulihkan", nil, nil)
}

func (h *MajorHandler) UpdateStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "ID tidak valid", nil)
		return
	}
	var req dto.MajorStatusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "validasi gagal", validator.Errors(err))
		return
	}

	if err := h.s.UpdateStatus(c.Request.Context(), uint(id), &req); err != nil {
		helper.HandleServiceError(c, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "MAJOR_STATUS_UPDATED", "status jurusan berhasil diubah", nil, nil)
}

func (h *MajorHandler) BulkDelete(c *gin.Context) {
	var req dto.BulkDeleteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "validasi gagal", validator.Errors(err))
		return
	}

	if err := h.s.BulkDelete(c.Request.Context(), req.IDs); err != nil {
		helper.HandleServiceError(c, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "MAJOR_BULK_DELETED", "data berhasil dihapus", nil, nil)
}

func (h *MajorHandler) BulkRestore(c *gin.Context) {
	var req dto.BulkDeleteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "validasi gagal", validator.Errors(err))
		return
	}

	if err := h.s.BulkRestore(c.Request.Context(), req.IDs); err != nil {
		helper.HandleServiceError(c, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "MAJOR_BULK_RESTORED", "data berhasil dipulihkan", nil, nil)
}

func (h *MajorHandler) GetDependencyInfo(c *gin.Context) {
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
	helper.SuccessResponse(c, http.StatusOK, "MAJOR_DEPENDENCY_INFO", "berhasil mengambil info dependensi", info, nil)
}
