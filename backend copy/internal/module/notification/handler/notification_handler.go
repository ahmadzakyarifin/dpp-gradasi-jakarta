package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/ahmadzakyarifin/schoolpay/backend/internal/helper"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/notification/repository"
	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	repo repository.NotificationRepo
}

func NewNotificationHandler(repo repository.NotificationRepo) *NotificationHandler {
	return &NotificationHandler{repo: repo}
}

// ListNotifications menangani GET /notifications
// Query params: channel, status, search, page, pageSize
func (h *NotificationHandler) ListNotifications(c *gin.Context) {
	channel := strings.TrimSpace(c.DefaultQuery("channel", "all"))
	status := strings.TrimSpace(c.DefaultQuery("status", "all"))
	search := strings.TrimSpace(c.Query("search"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	rows, total, err := h.repo.FindPaginated(c.Request.Context(), channel, status, search, page, pageSize)
	if err != nil {
		helper.HandleServiceError(c, err)
		return
	}

	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	items := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		item := map[string]any{
			"id":        row.ID,
			"createdAt": row.CreatedAt.Format("2006-01-02 15:04"),
			"channel":   row.Channel,
			"recipient": row.RecipientName,
			"contact":   row.Destination,
			"title":     row.Subject,
			"message":   row.Message,
			"status":    row.Status,
		}
		errMsg := row.ErrorMessage
		if errMsg == "" {
			errMsg = row.ProviderError
		}
		if errMsg != "" {
			item["error"] = errMsg
		}
		items = append(items, item)
	}

	c.JSON(http.StatusOK, helper.Response{
		Success: true,
		Message: "Data notifikasi berhasil dimuat",
		Data: map[string]any{
			"data":       items,
			"total":      total,
			"page":       page,
			"pageSize":   pageSize,
			"totalPages": totalPages,
		},
	})
}
