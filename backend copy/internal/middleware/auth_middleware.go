package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/ahmadzakyarifin/schoolpay/backend/internal/helper"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/user/repository"
	"github.com/gin-gonic/gin"
)

const (
	ContextUserID    = "user_id"
	ContextRoleID    = "role_id"
	ContextEmail     = "email"
	ContextIPAddress = "ip_address"
	ContextUserAgent = "user_agent"
)

func AuthMiddleware(
	jwtSecret string,
) gin.HandlerFunc {

	return func(c *gin.Context) {

		if c.Request.Method == http.MethodOptions {
			c.Next()
			return
		}

		authHeader := strings.TrimSpace(c.GetHeader("Authorization"))
		if authHeader == "" {
			helper.ErrorResponse(c, http.StatusUnauthorized, "AUTH_TOKEN_INVALID_OR_EXPIRED", "Sesi login telah berakhir. Silakan login kembali.", nil)
			c.Abort()
			return
		}

		parts := strings.Fields(authHeader)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			helper.ErrorResponse(c, http.StatusUnauthorized, "AUTH_TOKEN_INVALID_OR_EXPIRED", "Sesi login telah berakhir. Silakan login kembali.", nil)
			c.Abort()
			return
		}

		tokenStr := strings.TrimSpace(parts[1])
		if tokenStr == "" {
			helper.ErrorResponse(c, http.StatusUnauthorized, "AUTH_TOKEN_INVALID_OR_EXPIRED", "Sesi login telah berakhir. Silakan login kembali.", nil)
			c.Abort()
			return
		}

		claims, err := helper.ValidateToken(tokenStr, jwtSecret)
		if err != nil {
			helper.ErrorResponse(c, http.StatusUnauthorized, "AUTH_TOKEN_INVALID_OR_EXPIRED", "Token akses tidak valid atau telah kedaluwarsa.", nil)
			c.Abort()
			return
		}

		if claims.RoleID == 0 {
			helper.ErrorResponse(c, http.StatusForbidden, "AUTH_TOKEN_INVALID_OR_EXPIRED", "Role pengguna tidak valid.", nil)
			c.Abort()
			return
		}

		c.Set(ContextUserID, claims.UserID)
		c.Set(ContextRoleID, claims.RoleID)
		c.Set(ContextEmail, claims.Email)
		c.Set("user_name", claims.Name)
		c.Set("role_name", claims.RoleName)

		c.Set(ContextIPAddress, c.ClientIP())
		c.Set(ContextUserAgent, c.Request.UserAgent())

		// Setup standard context for service layer
		ctx := c.Request.Context()
		ctx = context.WithValue(ctx, "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "role_id", claims.RoleID)
		ctx = context.WithValue(ctx, "email", claims.Email)
		ctx = context.WithValue(ctx, "user_name", claims.Name)
		ctx = context.WithValue(ctx, "role_name", claims.RoleName)
		ctx = context.WithValue(ctx, "ip_address", c.ClientIP())
		ctx = context.WithValue(ctx, "user_agent", c.Request.UserAgent())
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

func PermissionMiddleware(
	userRepo repository.UserRepo,
	requiredPermissions ...string,
) gin.HandlerFunc {

	return func(c *gin.Context) {

		if c.Request.Method == http.MethodOptions {
			c.Next()
			return
		}

		permissions := cleanPermissions(requiredPermissions)
		if len(permissions) == 0 {
			helper.ErrorResponse(c, http.StatusInternalServerError, "SERVER_ERROR", "permission route belum dikonfigurasi", nil)
			c.Abort()
			return
		}

		roleID, ok := GetRoleID(c)
		if !ok || roleID == 0 {
			helper.ErrorResponse(c, http.StatusUnauthorized, "AUTH_TOKEN_INVALID_OR_EXPIRED", "Sesi pengguna tidak valid.", nil)
			c.Abort()
			return
		}

		// Super Admin
		if roleID == 1 {
			c.Next()
			return
		}

		hasPermission, err := userRepo.HasAnyRolePermission(
			c.Request.Context(),
			roleID,
			permissions,
		)

		if err != nil {
			helper.ErrorResponse(c, http.StatusInternalServerError, "SERVER_ERROR", "gagal memeriksa akses", nil)
			c.Abort()
			return
		}

		if !hasPermission {
			helper.ErrorResponse(c, http.StatusForbidden, "FORBIDDEN", "anda tidak memiliki akses ke fitur ini", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}

func cleanPermissions(values []string) []string {

	seen := make(map[string]bool)
	result := make([]string, 0, len(values))

	for _, value := range values {

		permission := strings.ToLower(strings.TrimSpace(value))
		if permission == "" {
			continue
		}

		if seen[permission] {
			continue
		}

		seen[permission] = true
		result = append(result, permission)
	}

	return result
}
