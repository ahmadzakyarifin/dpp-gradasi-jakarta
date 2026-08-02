package handler

import (
	"net/http"
	"strconv"

	"github.com/ahmadzakyarifin/schoolpay/backend/internal/helper"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/classmembership/dto"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/classmembership/service"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/validator"
	"github.com/gin-gonic/gin"
)

type ClassMembershipHandler struct {
	s service.ClassMembershipService
}

func NewClassMembershipHandler(s service.ClassMembershipService) *ClassMembershipHandler {
	return &ClassMembershipHandler{s: s}
}

func (h *ClassMembershipHandler) GetAll(c *gin.Context) {
	var req dto.ClassMembershipQueryReq
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
		"active_class_id": req.ActiveClassID,
		"student_id":      req.StudentID,
		"status":          req.Status,
		"sort":            req.Sort,
	}
	meta := helper.GetPaginationMeta(int(total), req.Page, req.Limit, filters)
	helper.SuccessResponse(c, http.StatusOK, "CLASS_MEMBERSHIP_LIST_RETRIEVED", "berhasil mengambil data", res, meta)
}

func (h *ClassMembershipHandler) Enroll(c *gin.Context) {
	var req dto.EnrollReq
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "validasi gagal", validator.Errors(err))
		return
	}
	res, err := h.s.Enroll(c.Request.Context(), req)
	if err != nil {
		helper.HandleServiceError(c, err)
		return
	}
	helper.SuccessResponse(c, http.StatusCreated, "CLASS_MEMBERSHIP_ENROLLED", "siswa berhasil didaftarkan ke kelas", gin.H{"data": res}, nil)
}

func (h *ClassMembershipHandler) Move(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		helper.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "ID tidak valid", nil)
		return
	}
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "ID tidak valid", nil)
		return
	}
	var req dto.MoveReq
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "validasi gagal", validator.Errors(err))
		return
	}
	res, err := h.s.Move(c.Request.Context(), uint(id), req)
	if err != nil {
		helper.HandleServiceError(c, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "CLASS_MEMBERSHIP_MOVED", "siswa berhasil dipindah kelas", gin.H{"data": res}, nil)
}

func (h *ClassMembershipHandler) SetStatus(c *gin.Context) {
	idStr := c.Param("id")
	if idStr == "" {
		helper.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "ID tidak valid", nil)
		return
	}
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "ID tidak valid", nil)
		return
	}
	var req dto.SetStatusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "validasi gagal", validator.Errors(err))
		return
	}
	res, err := h.s.SetStatus(c.Request.Context(), uint(id), req)
	if err != nil {
		helper.HandleServiceError(c, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "CLASS_MEMBERSHIP_STATUS_UPDATED", "status keanggotaan kelas berhasil diubah", gin.H{"data": res}, nil)
}
