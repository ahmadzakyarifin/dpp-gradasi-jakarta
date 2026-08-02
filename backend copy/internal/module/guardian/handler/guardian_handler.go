package handler

import (
	"net/http"
	"strconv"

	"github.com/ahmadzakyarifin/schoolpay/backend/internal/helper"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/guardian/dto"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/guardian/service"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/validator"
	"github.com/gin-gonic/gin"
)

type GuardianHandler struct {
	s service.GuardianService
}

func NewGuardianHandler(s service.GuardianService) *GuardianHandler {
	return &GuardianHandler{s: s}
}

func (h *GuardianHandler) GetAll(c *gin.Context) {
	var req dto.GuardianQueryReq
	if err := c.ShouldBindQuery(&req); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "validasi gagal", validator.Errors(err))
		return
	}

	list, total, err := h.s.GetAll(c.Request.Context(), req)
	if err != nil {
		helper.HandleServiceError(c, err)
		return
	}

	meta := helper.GetPaginationMeta(int(total), req.Page, req.Limit)
	helper.SuccessResponse(c, http.StatusOK, "GUARDIAN_LIST", "berhasil mengambil data wali", list, meta)
}

func (h *GuardianHandler) Create(c *gin.Context) {
	var req dto.GuardianCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "validasi gagal", validator.Errors(err))
		return
	}

	res, err := h.s.Create(c.Request.Context(), req)
	if err != nil {
		helper.HandleServiceError(c, err)
		return
	}

	helper.SuccessResponse(c, http.StatusCreated, "GUARDIAN_CREATED", "wali berhasil dibuat", gin.H{"data": res}, nil)
}

func (h *GuardianHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "ID tidak valid", nil)
		return
	}

	var req dto.GuardianUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "validasi gagal", validator.Errors(err))
		return
	}

	res, err := h.s.Update(c.Request.Context(), uint(id), req)
	if err != nil {
		helper.HandleServiceError(c, err)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "GUARDIAN_UPDATED", "wali berhasil diperbarui", gin.H{"data": res}, nil)
}

func (h *GuardianHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "ID tidak valid", nil)
		return
	}

	if err := h.s.Delete(c.Request.Context(), uint(id)); err != nil {
		helper.HandleServiceError(c, err)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "GUARDIAN_DELETED", "wali berhasil dihapus", nil, nil)
}
