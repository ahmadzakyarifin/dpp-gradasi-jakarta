package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/helper"
	"github.com/gin-gonic/gin"
)

const (
	ContextUserID    = "user_id"
	ContextRoleID    = "role_id"
	ContextEmail     = "email"
	ContextRoleName  = "role_name"
	ContextUserName  = "user_name"
	ContextIPAddress = "ip_address"
	ContextUserAgent = "user_agent"
)

func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
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
		c.Set(ContextUserName, claims.Name)
		c.Set(ContextRoleName, claims.RoleName)

		c.Set(ContextIPAddress, c.ClientIP())
		c.Set(ContextUserAgent, c.Request.UserAgent())

		ctx := c.Request.Context()
		ctx = context.WithValue(ctx, "user_id", claims.UserID)
		ctx = context.WithValue(ctx, "role_id", claims.RoleID)
		ctx = context.WithValue(ctx, "email", claims.Email)
		ctx = context.WithValue(ctx, "user_name", claims.Name)
		ctx = context.WithValue(ctx, "role_name", claims.RoleName)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

func GetUserID(c *gin.Context) (uint, bool) {
	id, ok := c.Get(ContextUserID)
	if !ok {
		return 0, false
	}
	userID, ok := id.(uint)
	return userID, ok
}

func GetRoleID(c *gin.Context) (uint, bool) {
	id, ok := c.Get(ContextRoleID)
	if !ok {
		return 0, false
	}
	roleID, ok := id.(uint)
	return roleID, ok
}

func GetEmail(c *gin.Context) (string, bool) {
	email, ok := c.Get(ContextEmail)
	if !ok {
		return "", false
	}
	return email.(string), ok
}
