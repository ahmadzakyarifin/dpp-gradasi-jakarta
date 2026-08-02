package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ahmadzakyarifin/schoolpay/backend/internal/helper"
	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
)

func writeRateLimitHeaders(c *gin.Context, name string, info limiter.Context) {
	c.Header("X-RateLimit-"+name+"-Limit", strconv.FormatInt(info.Limit, 10))
	c.Header("X-RateLimit-"+name+"-Remaining", strconv.FormatInt(info.Remaining, 10))
	c.Header("X-RateLimit-"+name+"-Reset", strconv.FormatInt(info.Reset, 10))
}

func abortRateLimitError(c *gin.Context) {
	helper.ErrorResponse(
		c,
		http.StatusInternalServerError,
		"SYSTEM_SECURITY_NOT_CONFIGURED",
		"Gagal memproses rate limit.",
		nil,
	)
	c.Abort()
}

func abortRateLimitUnauthorized(c *gin.Context) {
	helper.ErrorResponse(
		c,
		http.StatusUnauthorized,
		"SYSTEM_SECURITY_NOT_CONFIGURED",
		"Sesi pengguna tidak valid.",
		nil,
	)
	c.Abort()
}

func abortTooManyRequests(c *gin.Context, resetAt int64) {
	retryAfter := retryAfterSeconds(resetAt)

	c.Header("Retry-After", strconv.Itoa(retryAfter))

	helper.RateLimitResponse(
		c,
		"AUTH_RATE_LIMIT_EXCEEDED",
		"Terlalu banyak percobaan. Silakan coba lagi nanti.",
		retryAfter,
	)
	c.Abort()
}

func retryAfterSeconds(resetAt int64) int {
	seconds := int(resetAt - time.Now().Unix())
	if seconds < 1 {
		return 1
	}

	return seconds
}
