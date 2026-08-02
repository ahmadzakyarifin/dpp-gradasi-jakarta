package role

import (
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/middleware"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/role/handler"
	userrepo "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/user/repository"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// RegisterRoutes mendaftarkan endpoint manajemen role ke grup /api.
// Permission dipasang sesuai mapping docs/rbac-design.md section 4.11:
//   - role.view:   GET (list, by-id, permissions, check-unique, dependency-info)
//   - role.create: POST
//   - role.update: PUT, PATCH (status, restore), PATCH bulk-restore
//   - role.delete: DELETE, POST bulk-delete
func RegisterRoutes(
	apiGroup *gin.RouterGroup,
	h *handler.RoleHandler,
	jwtSecret string,
	redisClient *redis.Client,
	userRepo userrepo.UserRepo,
) {
	requireAuth := middleware.AuthMiddleware(jwtSecret)

	roles := apiGroup.Group("/roles")
	roles.Use(requireAuth)
	{
		// role.view
		roles.GET("", middleware.PermissionMiddleware(userRepo, "role.view"), h.GetAll)
		roles.GET("/permissions", middleware.PermissionMiddleware(userRepo, "role.view"), h.GetPermissions)
		roles.GET("/check-unique", middleware.PermissionMiddleware(userRepo, "role.view"), h.CheckUnique)
		roles.GET("/:id", middleware.PermissionMiddleware(userRepo, "role.view"), h.GetByID)
		roles.GET("/:id/dependency-info", middleware.PermissionMiddleware(userRepo, "role.view"), h.GetDependencyInfo)

		// role.create
		roles.POST("", middleware.PermissionMiddleware(userRepo, "role.create"), h.Create)

		// role.update
		roles.PUT("/:id", middleware.PermissionMiddleware(userRepo, "role.update"), h.Update)
		roles.PATCH("/:id/restore", middleware.PermissionMiddleware(userRepo, "role.update"), h.Restore)
		roles.PATCH("/:id/status", middleware.PermissionMiddleware(userRepo, "role.update"), h.UpdateStatus)
		roles.PATCH("/bulk-restore", middleware.PermissionMiddleware(userRepo, "role.update"), h.BulkRestore)

		// role.delete
		roles.DELETE("/:id", middleware.PermissionMiddleware(userRepo, "role.delete"), h.Delete)
		roles.POST("/bulk-delete", middleware.PermissionMiddleware(userRepo, "role.delete"), h.BulkDelete)
	}
}
