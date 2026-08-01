package middleware

import (
	"errors"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

var errUnexpectedRedisReply = errors.New("rate limiter: unexpected redis reply")

// fixedWindowLimiter adalah implementasi fixed-window rate limit di atas Redis
// (INCR + PEXPIRE atomic via Lua). Menggantikan ulule/limiter yang punya bug
// TTL tidak pernah di-set (key tanpa expiry → retry_after selalu 1 & blokir permanen).
type fixedWindowLimiter struct {
	client *redis.Client
	prefix string
}

func newFixedWindowLimiter(client *redis.Client) *fixedWindowLimiter {
	return &fixedWindowLimiter{
		client: client,
		prefix: "dppgradasi_rate_limit",
	}
}

// check menaikkan counter; jika melebihi limit → return blocked=true + retryAfter.
func (l *fixedWindowLimiter) check(c *gin.Context, key string, limit int64, window time.Duration) (blocked bool, retryAfter int, remaining int64, err error) {
	ctx := c.Request.Context()
	fullKey := l.prefix + ":" + key

	res, err := incrFixedWindow(ctx, l.client, fullKey, window)
	if err != nil {
		return false, 0, 0, err
	}

	remaining = limit - res.Count
	if remaining < 0 {
		remaining = 0
	}

	// TTL tersisa = sisa waktu window (retry_after akurat)
	retryAfter = int(res.TTL / time.Second)
	if retryAfter < 1 {
		retryAfter = 1
	}

	c.Header("X-RateLimit-Limit", strconv.FormatInt(limit, 10))
	c.Header("X-RateLimit-Remaining", strconv.FormatInt(remaining, 10))
	c.Header("X-RateLimit-Reset", strconv.FormatInt(time.Now().Unix()+int64(retryAfter), 10))

	if res.Count > limit {
		return true, retryAfter, remaining, nil
	}

	return false, retryAfter, remaining, nil
}
