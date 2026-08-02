package role

import (
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/middleware"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/role/handler"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	apiGroup *gin.RouterGroup,
	h *handler.RoleHandler,
	jwtSecret string,
) {
	requireAuth := middleware.AuthMiddleware(jwtSecret)

	roles := apiGroup.Group("/roles")
	roles.Use(requireAuth)
	{
		// role.view
		roles.GET("", middleware.RoleMiddleware("super_admin", "admin"), h.GetAll)
		roles.GET("/permissions", middleware.RoleMiddleware("super_admin", "admin"), h.GetPermissions)
		roles.GET("/check-unique", middleware.RoleMiddleware("super_admin", "admin"), h.CheckUnique)
		roles.GET("/:id", middleware.RoleMiddleware("super_admin", "admin"), h.GetByID)
		roles.GET("/:id/dependency-info", middleware.RoleMiddleware("super_admin", "admin"), h.GetDependencyInfo)

		// role.create
		roles.POST("", middleware.RoleMiddleware("super_admin", "admin"), h.Create)

		// role.update
		roles.PUT("/:id", middleware.RoleMiddleware("super_admin", "admin"), h.Update)
		roles.PATCH("/:id/restore", middleware.RoleMiddleware("super_admin", "admin"), h.Restore)
		roles.PATCH("/:id/status", middleware.RoleMiddleware("super_admin", "admin"), h.UpdateStatus)
		roles.PATCH("/bulk-restore", middleware.RoleMiddleware("super_admin", "admin"), h.BulkRestore)

		// role.delete
		roles.DELETE("/:id", middleware.RoleMiddleware("super_admin", "admin"), h.Delete)
		roles.POST("/bulk-delete", middleware.RoleMiddleware("super_admin", "admin"), h.BulkDelete)
	}
}
