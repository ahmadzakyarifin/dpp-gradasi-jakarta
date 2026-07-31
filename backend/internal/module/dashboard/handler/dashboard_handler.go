package handler

import (
	"net/http"

	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/helper"
	dashservice "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/dashboard/service"
	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	service dashservice.DashboardService
}

func NewDashboardHandler(s dashservice.DashboardService) *DashboardHandler {
	return &DashboardHandler{service: s}
}

func (h *DashboardHandler) Summary(c *gin.Context) {
	res, err := h.service.GetSummary(c.Request.Context())
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "SERVER_ERROR", err.Error(), nil)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "DASHBOARD_SUMMARY", "Ringkasan dashboard berhasil diambil", res, nil)
}
