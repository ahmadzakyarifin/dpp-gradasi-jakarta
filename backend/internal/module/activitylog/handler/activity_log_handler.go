package handler

import (
	"net/http"
	"strconv"

	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/helper"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/activitylog/dto"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/activitylog/service"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/validator"
	"github.com/gin-gonic/gin"
)

type ActivityLogHandler struct {
	service service.ActivityLogService
}

func NewActivityLogHandler(s service.ActivityLogService) *ActivityLogHandler {
	return &ActivityLogHandler{
		service: s,
	}
}

func (h *ActivityLogHandler) List(c *gin.Context) {
	var req dto.ActivityLogQueryReq
	if err := c.ShouldBindQuery(&req); err != nil {
		helper.ErrorResponse(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "validasi gagal", validator.Errors(err))
		return
	}

	res, err := h.service.List(c.Request.Context(), &req)
	if err != nil {
		helper.HandleServiceError(c, err)
		return
	}

	filters := map[string]any{
		"search": req.Search,
		"action": req.Action,
		"entity": req.Entity,
		"role":   req.Role,
		"risk":   req.Risk,
	}
	meta := helper.GetPaginationMeta(int(res.Pagination.Total), res.Pagination.Page, res.Pagination.Limit, filters)
	helper.SuccessResponse(c, http.StatusOK, "ACTIVITY_LOG_LIST", "berhasil mengambil daftar aktivitas", res, meta)
}

func (h *ActivityLogHandler) Detail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.ErrorResponse(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "ID tidak valid", nil)
		return
	}

	res, err := h.service.Detail(c.Request.Context(), id)
	if err != nil {
		helper.HandleServiceError(c, err)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "ACTIVITY_LOG_DETAIL", "Berhasil mengambil detail activity log.", gin.H{"data": res}, nil)
}

func (h *ActivityLogHandler) EntityLogs(c *gin.Context) {
	entityType := c.Param("entityType")

	entityID, err := strconv.ParseUint(c.Param("entityID"), 10, 64)
	if err != nil {
		helper.ErrorResponse(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "ID entitas tidak valid", nil)
		return
	}

	res, err := h.service.EntityLogs(c.Request.Context(), entityType, entityID)
	if err != nil {
		helper.HandleServiceError(c, err)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "ACTIVITY_LOG_ENTITY_LISTED", "Berhasil mengambil activity log entitas.", gin.H{"data": res}, nil)
}
