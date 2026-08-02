package helper

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Success    bool   `json:"success"`
	Code       string `json:"code,omitempty"`
	Message    string `json:"message"`
	Data       any    `json:"data,omitempty"`
	Metadata   any    `json:"meta,omitempty"`
	Errors     any    `json:"errors,omitempty"`
	RetryAfter *int   `json:"retry_after,omitempty"`
}

type PaginationMeta struct {
	Total      int            `json:"total"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
	TotalPages int            `json:"total_pages"`
	Filters    map[string]any `json:"filters,omitempty"`
}

func GetPaginationMeta(total int, page int, limit int, filters ...map[string]any) PaginationMeta {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	totalPages := 0
	if total > 0 {
		totalPages = (total + limit - 1) / limit
	}

	var customFilters map[string]any
	if len(filters) > 0 && filters[0] != nil {
		customFilters = filters[0]
	}

	return PaginationMeta{
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
		Filters:    customFilters,
	}
}

func SuccessResponse(
	c *gin.Context,
	httpCode int,
	code string,
	message string,
	data any,
	metadata any,
) {
	c.JSON(httpCode, Response{
		Success:  true,
		Code:     code,
		Message:  message,
		Data:     data,
		Metadata: metadata,
	})
}

func ErrorResponse(
	c *gin.Context,
	httpCode int,
	code string,
	message string,
	errData any,
) {
	c.JSON(httpCode, Response{
		Success: false,
		Code:    code,
		Message: message,
		Errors:  errData,
	})
}

func RateLimitResponse(
	c *gin.Context,
	code string,
	message string,
	retryAfter int,
) {
	c.JSON(http.StatusTooManyRequests, Response{
		Success:    false,
		Code:       code,
		Message:    message,
		RetryAfter: &retryAfter,
	})
}
