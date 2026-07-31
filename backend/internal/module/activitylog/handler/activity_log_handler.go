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
		helper.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "validasi gagal", validator.Errors(err))
		return
	}

	res, err := h.service.List(c.Request.Context(), &req)
	if err != nil {
		svcErr, ok := err.(*helper.ServiceError)
		if !ok {
			helper.ErrorResponse(c, http.StatusInternalServerError, "SERVER_ERROR", err.Error(), nil)
			return
		}
		helper.ErrorResponse(c, http.StatusBadRequest, svcErr.Code, svcErr.Message, nil)
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
	helper.SuccessResponse(c, http.StatusOK, "ACTIVITY_LOG_LIST", "Daftar log aktivitas berhasil diambil", res, meta)
}

func (h *ActivityLogHandler) Summary(c *gin.Context) {
	var req dto.ActivityLogQueryReq
	if err := c.ShouldBindQuery(&req); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "validasi gagal", validator.Errors(err))
		return
	}

	res, err := h.service.Summary(c.Request.Context(), &req)
	if err != nil {
		svcErr, ok := err.(*helper.ServiceError)
		if !ok {
			helper.ErrorResponse(c, http.StatusInternalServerError, "SERVER_ERROR", err.Error(), nil)
			return
		}
		helper.ErrorResponse(c, http.StatusBadRequest, svcErr.Code, svcErr.Message, nil)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "ACTIVITY_LOG_SUMMARY", "Ringkasan log aktivitas berhasil diambil", res, nil)
}

func (h *ActivityLogHandler) Detail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "ID tidak valid", nil)
		return
	}

	res, err := h.service.Detail(c.Request.Context(), id)
	if err != nil {
		svcErr, ok := err.(*helper.ServiceError)
		if !ok {
			helper.ErrorResponse(c, http.StatusInternalServerError, "SERVER_ERROR", err.Error(), nil)
			return
		}
		helper.ErrorResponse(c, http.StatusBadRequest, svcErr.Code, svcErr.Message, nil)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "ACTIVITY_LOG_DETAIL", "Detail log aktivitas berhasil diambil", gin.H{"data": res}, nil)
}

func (h *ActivityLogHandler) EntityLogs(c *gin.Context) {
	entityType := c.Param("entityType")

	entityID, err := strconv.ParseUint(c.Param("entityID"), 10, 64)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "ID entitas tidak valid", nil)
		return
	}

	res, err := h.service.EntityLogs(c.Request.Context(), entityType, entityID)
	if err != nil {
		svcErr, ok := err.(*helper.ServiceError)
		if !ok {
			helper.ErrorResponse(c, http.StatusInternalServerError, "SERVER_ERROR", err.Error(), nil)
			return
		}
		helper.ErrorResponse(c, http.StatusBadRequest, svcErr.Code, svcErr.Message, nil)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "ENTITY_LOGS_RETRIEVED", "Riwayat aktivitas entitas berhasil diambil", gin.H{"data": res}, nil)
}
