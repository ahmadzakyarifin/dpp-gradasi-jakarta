package middleware

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/helper"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/ulule/limiter/v3"
	redisstore "github.com/ulule/limiter/v3/drivers/store/redis"
)

type RateLimiter struct {
	store limiter.Store
}

type preparedRule struct {
	rule    Rule
	limiter *limiter.Limiter
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
	store, err := redisstore.NewStoreWithOptions(
		redisClient,
		limiter.StoreOptions{
			Prefix: "dppgradasi_rate_limit",
		},
	)
	if err != nil {
		return nil, err
	}

	return &RateLimiter{
		store: store,
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

	return defaultRateLimiter.Use(scope, rules...)
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

			info, err := item.limiter.Get(
				c.Request.Context(),
				key,
			)
			if err != nil {
				abortRateLimitError(c)
				return
			}

			writeRateLimitHeaders(
				c,
				item.rule.Name,
				info,
			)

			if info.Reached {
				abortTooManyRequests(
					c,
					info.Reset,
				)
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
			limiter: limiter.New(
				r.store,
				limiter.Rate{
					Period: rule.Period,
					Limit:  rule.Limit,
				},
			),
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

	var body struct {
		Email string `json:"email"`
	}

	defer helper.RestoreBody(c)

	if err := c.ShouldBindBodyWithJSON(&body); err != nil {
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
