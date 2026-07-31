package activitylog

import (
	activityloghandler "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/activitylog/handler"

	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// RegisterRoutes mendaftarkan endpoint activity log ke grup /api.
func RegisterRoutes(
	apiGroup *gin.RouterGroup,
	h *activityloghandler.ActivityLogHandler,
	jwtSecret string,
	redisClient *redis.Client,
) {
	requireAuth := middleware.AuthMiddleware(jwtSecret)
	requireSuperAdmin := middleware.RoleMiddleware("super_admin")

	g := apiGroup.Group("/activity-logs")
	g.Use(requireAuth)
	g.Use(requireSuperAdmin)
	{
		g.GET("", h.List)
		g.GET("/summary", h.Summary)
		g.GET("/:id", h.Detail)
		g.GET("/entity/:entityType/:entityID", h.EntityLogs)
	}
}
