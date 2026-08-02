package middleware

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type RateLimiter struct {
	limiter *fixedWindowLimiter
}

type preparedRule struct {
	rule Rule
}

type requestInfo struct {
	scope   string
	path    string
	ip      string
	email   string
	userID  string
	hasUser bool
}

var defaultRateLimiter *RateLimiter

func NewRedisRateLimiter(redisClient *redis.Client) (*RateLimiter, error) {
	return &RateLimiter{
		limiter: newFixedWindowLimiter(redisClient),
	}, nil
}

func SetDefaultRateLimiter(rl *RateLimiter) {
	if rl != nil {
		defaultRateLimiter = rl
	}
}

func RateLimitRules(scope string, rules ...Rule) gin.HandlerFunc {
	if defaultRateLimiter == nil {
		panic("rate limiter belum diinisialisasi")
	}

	handler := defaultRateLimiter.Use(scope, rules...)
	return handler
}

func RateLimiterMiddleware(scope string, rules ...Rule) gin.HandlerFunc {
	return RateLimitRules(scope, rules...)
}

func RateLimitPerUser(scope string, maxRequests int64) gin.HandlerFunc {
	return RateLimitRules(
		scope,
		User(maxRequests, time.Minute),
	)
}

func (r *RateLimiter) Use(scope string, rules ...Rule) gin.HandlerFunc {
	scope = normalizeScope(scope)

	preparedRules := r.prepareRules(rules)

	return func(c *gin.Context) {

		req := newRequestInfo(c, scope, rules)

		for _, item := range preparedRules {

			if item.rule.Kind == KindUser && !req.hasUser {
				abortRateLimitUnauthorized(c)
				return
			}

			key, ok := makeRateLimitKey(item.rule, req)
			if !ok {
				continue
			}

			blocked, retryAfter, _, err := r.limiter.check(
				c,
				key,
				item.rule.Limit,
				item.rule.Period,
			)
			if err != nil {
				// Fail-open: Redis error jangan blokir request (availability > security).
				// Log warning supaya kebaca kalau ini terjadi di production.
				log.Printf("[rate-limit][warn] redis error utk scope=%s kind=%s: %v — fail-open, request diteruskan", scope, item.rule.Kind, err)
				continue
			}

			if blocked {
				abortTooManyRequests(c, time.Now().Add(time.Duration(retryAfter)*time.Second).Unix())
				return
			}
		}

		c.Next()
	}
}

func (r *RateLimiter) prepareRules(rules []Rule) []preparedRule {
	prepared := make([]preparedRule, 0, len(rules))

	for _, rule := range rules {
		rule = normalizeRule(rule)
		prepared = append(prepared, preparedRule{
			rule: rule,
		})
	}

	return prepared
}

func normalizeRule(rule Rule) Rule {

	if strings.TrimSpace(rule.Name) == "" {
		rule.Name = string(rule.Kind)
	}

	rule.Name = normalizeScope(rule.Name)

	if rule.Limit <= 0 {
		rule.Limit = 1
	}

	if rule.Period <= 0 {
		rule.Period = time.Minute
	}

	return rule
}

func newRequestInfo(
	c *gin.Context,
	scope string,
	rules []Rule,
) requestInfo {

	userID, hasUser := getUserID(c)

	return requestInfo{
		scope:   scope,
		path:    rateLimitPath(c),
		ip:      c.ClientIP(),
		email:   emailFromBody(c, rules),
		userID:  userID,
		hasUser: hasUser,
	}
}

func makeRateLimitKey(
	rule Rule,
	req requestInfo,
) (string, bool) {

	prefix := fmt.Sprintf(
		"rate:%s:%s:%s",
		req.scope,
		req.path,
		rule.Name,
	)

	switch rule.Kind {

	case KindIP:
		return prefix + ":ip:" + hash(req.ip), true

	case KindEmail:
		if req.email == "" {
			return "", false
		}

		return prefix + ":email:" + hash(req.email), true

	case KindIPEmail:
		if req.email == "" {
			return "", false
		}

		return prefix + ":ip_email:" + hash(req.ip+":"+req.email), true

	case KindUser:
		if !req.hasUser {
			return "", false
		}

		return prefix + ":user:" + hash(req.userID), true

	default:
		return "", false
	}
}

func rateLimitPath(c *gin.Context) string {

	path := c.FullPath()
	if path == "" {
		path = c.Request.URL.Path
	}

	path = strings.Trim(path, "/")
	path = strings.ReplaceAll(path, "/", ":")

	if path == "" {
		return "unknown"
	}

	return path
}

func emailFromBody(
	c *gin.Context,
	rules []Rule,
) string {

	if !needsEmail(rules) {
		return ""
	}

	// Baca body lalu RESTORE agar handler tetap bisa bind normal.
	raw, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return ""
	}
	c.Request.Body = io.NopCloser(bytes.NewReader(raw))

	var body struct {
		Email string `json:"email"`
	}
	if err := json.Unmarshal(raw, &body); err != nil {
		return ""
	}

	return strings.ToLower(
		strings.TrimSpace(body.Email),
	)
}

func needsEmail(rules []Rule) bool {

	for _, rule := range rules {

		if rule.Kind == KindEmail ||
			rule.Kind == KindIPEmail {
			return true
		}
	}

	return false
}

func getUserID(c *gin.Context) (string, bool) {

	id, ok := GetUserID(c)
	if !ok {
		return "", false
	}

	return strconv.FormatUint(
		uint64(id),
		10,
	), true
}

func normalizeScope(scope string) string {

	scope = strings.ToLower(
		strings.TrimSpace(scope),
	)

	if scope == "" {
		return "global"
	}

	scope = strings.ReplaceAll(scope, " ", "_")
	scope = strings.ReplaceAll(scope, ":", "_")
	scope = strings.ReplaceAll(scope, "/", "_")

	return scope
}

func hash(value string) string {
	value = strings.ToLower(
		strings.TrimSpace(value),
	)

	if value == "" {
		value = "unknown"
	}

	sum := sha256.Sum256([]byte(value))

	return hex.EncodeToString(sum[:])[:32]
}
