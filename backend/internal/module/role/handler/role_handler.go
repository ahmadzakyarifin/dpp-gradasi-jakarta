package handler

import (
	"net/http"
	"strconv"

	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/helper"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/role/dto"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/role/mapper"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/role/service"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/validator"
	"github.com/gin-gonic/gin"
)

type RoleHandler struct {
	s service.RoleService
}

func NewRoleHandler(s service.RoleService) *RoleHandler {
	return &RoleHandler{s: s}
}

func (h *RoleHandler) GetAll(c *gin.Context) {
	var req dto.RoleQueryReq
	if err := c.ShouldBindQuery(&req); err != nil {
		helper.ErrorResponse(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "validasi gagal", validator.Errors(err))
		return
	}

	req.Normalize()

	list, total, err := h.s.GetAllRoles(c.Request.Context(), req)
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
	helper.SuccessResponse(c, http.StatusOK, "ROLE_LIST", "berhasil mengambil data role", mapper.EntityListToResponse(list), meta)
}

func (h *RoleHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helper.ErrorResponse(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "ID tidak valid", nil)
		return
	}

	role, err := h.s.GetRoleByID(c.Request.Context(), uint(id))
	if err != nil {
		helper.HandleServiceError(c, err)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "ROLE_DETAIL", "berhasil mengambil data role", mapper.EntityToResponse(role), nil)
}

func (h *RoleHandler) Create(c *gin.Context) {
	var req dto.RoleCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "validasi gagal", validator.Errors(err))
		return
	}

	created, err := h.s.CreateRole(c.Request.Context(), req)
	if err != nil {
		helper.HandleServiceError(c, err)
		return
	}

	helper.SuccessResponse(c, http.StatusCreated, "ROLE_CREATED", "role berhasil dibuat", mapper.EntityToResponse(created), nil)
}

func (h *RoleHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helper.ErrorResponse(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "ID tidak valid", nil)
		return
	}

	var req dto.RoleUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "validasi gagal", validator.Errors(err))
		return
	}

	updated, err := h.s.UpdateRole(c.Request.Context(), uint(id), req)
	if err != nil {
		helper.HandleServiceError(c, err)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "ROLE_UPDATED", "role berhasil diperbarui", mapper.EntityToResponse(updated), nil)
}

func (h *RoleHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helper.ErrorResponse(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "ID tidak valid", nil)
		return
	}

	if err := h.s.DeleteRole(c.Request.Context(), uint(id)); err != nil {
		helper.HandleServiceError(c, err)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "ROLE_DELETED", "role berhasil dihapus", nil, nil)
}

func (h *RoleHandler) Restore(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helper.ErrorResponse(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "ID tidak valid", nil)
		return
	}

	if err := h.s.RestoreRole(c.Request.Context(), uint(id)); err != nil {
		helper.HandleServiceError(c, err)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "ROLE_RESTORED", "role berhasil dipulihkan", nil, nil)
}

func (h *RoleHandler) UpdateStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helper.ErrorResponse(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "ID tidak valid", nil)
		return
	}

	var req dto.RoleStatusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "validasi gagal", validator.Errors(err))
		return
	}

	if err := h.s.UpdateStatus(c.Request.Context(), uint(id), req); err != nil {
		helper.HandleServiceError(c, err)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "ROLE_STATUS_UPDATED", "status role berhasil diubah", nil, nil)
}

func (h *RoleHandler) BulkDelete(c *gin.Context) {
	var req dto.BulkDeleteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "validasi gagal", validator.Errors(err))
		return
	}

	if err := h.s.BulkDelete(c.Request.Context(), req.IDs); err != nil {
		helper.HandleServiceError(c, err)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "ROLE_BULK_DELETED", "role terpilih berhasil dihapus", nil, nil)
}

func (h *RoleHandler) BulkRestore(c *gin.Context) {
	var req dto.BulkDeleteReq
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "validasi gagal", validator.Errors(err))
		return
	}

	if err := h.s.BulkRestore(c.Request.Context(), req.IDs); err != nil {
		helper.HandleServiceError(c, err)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "ROLE_BULK_RESTORED", "role terpilih berhasil dipulihkan", nil, nil)
}

func (h *RoleHandler) GetDependencyInfo(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		helper.ErrorResponse(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "ID tidak valid", nil)
		return
	}

	info, err := h.s.GetDependencyInfo(c.Request.Context(), uint(id))
	if err != nil {
		helper.HandleServiceError(c, err)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "ROLE_DEPENDENCY_INFO", "berhasil mengambil info dependensi", info, nil)
}

func (h *RoleHandler) CheckUnique(c *gin.Context) {
	field := c.Query("field")
	value := c.Query("value")
	excludeID, _ := strconv.ParseUint(c.Query("exclude_id"), 10, 32)

	if field == "" || value == "" {
		helper.ErrorResponse(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "field dan value harus diisi", nil)
		return
	}

	exists, err := h.s.CheckUnique(c.Request.Context(), field, value, uint(excludeID))
	if err != nil {
		helper.HandleServiceError(c, err)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "ROLE_UNIQUE_CHECK", "berhasil mengecek keunikan", gin.H{"is_unique": !exists}, nil)
}
