package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/helper"
	"github.com/gin-gonic/gin"
)

type RuleKind string

const (
	KindIP      RuleKind = "ip"
	KindEmail   RuleKind = "email"
	KindIPEmail RuleKind = "ip_email"
	KindUser    RuleKind = "user"
)

type Rule struct {
	Kind   RuleKind
	Name   string
	Limit  int64
	Period time.Duration
}

func IP(limit int64, period time.Duration) Rule {
	return Rule{Kind: KindIP, Limit: limit, Period: period}
}

func Email(limit int64, period time.Duration) Rule {
	return Rule{Kind: KindEmail, Limit: limit, Period: period}
}

func IPEmail(limit int64, period time.Duration) Rule {
	return Rule{Kind: KindIPEmail, Limit: limit, Period: period}
}

func User(limit int64, period time.Duration) Rule {
	return Rule{Kind: KindUser, Limit: limit, Period: period}
}

func abortRateLimitUnauthorized(c *gin.Context) {
	helper.ErrorResponse(c, http.StatusUnauthorized, "AUTH_TOKEN_INVALID_OR_EXPIRED", "Sesi login telah berakhir. Silakan login kembali.", nil)
	c.Abort()
}

func abortTooManyRequests(c *gin.Context, resetAt int64) {
	retryAfter := int(resetAt - time.Now().Unix())
	if retryAfter < 1 {
		retryAfter = 1
	}
	c.Header("Retry-After", strconv.Itoa(retryAfter))
	helper.RateLimitResponse(c, "AUTH_RATE_LIMIT_EXCEEDED", "Terlalu banyak percobaan. Silakan coba lagi nanti.", retryAfter)
	c.Abort()
}
