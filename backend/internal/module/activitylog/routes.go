package activitylog

import (
	activityloghandler "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/activitylog/handler"

	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes mendaftarkan endpoint activity log ke grup /api.
func RegisterRoutes(
	apiGroup *gin.RouterGroup,
	h *activityloghandler.ActivityLogHandler,
	jwtSecret string,
) {
	requireAuth := middleware.AuthMiddleware(jwtSecret)

	g := apiGroup.Group("/activity-logs")
	g.Use(requireAuth)
	{
		g.GET("", middleware.RoleMiddleware("super_admin"), h.List)
		g.GET("/:id", middleware.RoleMiddleware("super_admin"), h.Detail)
		g.GET("/entity/:entityType/:entityID", middleware.RoleMiddleware("super_admin"), h.EntityLogs)
	}
}
