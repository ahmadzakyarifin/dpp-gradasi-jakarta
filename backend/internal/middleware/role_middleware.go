package middleware

import (
	"net/http"
	"strings"

	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/helper"
	"github.com/gin-gonic/gin"
)

// RoleMiddleware memeriksa apakah role user termasuk dalam allowedRoles.
// Ini menggantikan PermissionMiddleware — cek role_id langsung.
// allowedRoles: slice of role names (strings), misal []string{"super_admin", "admin"}
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodOptions {
			c.Next()
			return
		}

		roleNameRaw, exists := c.Get(ContextRoleName)
		if !exists {
			helper.ErrorResponse(c, http.StatusForbidden, "FORBIDDEN", "Anda tidak memiliki akses ke sumber daya ini.", nil)
			c.Abort()
			return
		}

		roleName, ok := roleNameRaw.(string)
		if !ok {
			helper.ErrorResponse(c, http.StatusForbidden, "FORBIDDEN", "Role pengguna tidak valid.", nil)
			c.Abort()
			return
		}

		for _, allowed := range allowedRoles {
			if strings.EqualFold(roleName, allowed) {
				c.Next()
				return
			}
		}

		helper.ErrorResponse(c, http.StatusForbidden, "FORBIDDEN", "Anda tidak memiliki izin untuk melakukan aksi ini.", nil)
		c.Abort()
	}
}
