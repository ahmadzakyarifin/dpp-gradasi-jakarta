package activitylog

import (
	activityloghandler "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/activitylog/handler"
	userrepo "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/user/repository"

	"github.com/ahmadzakyarifin/schoolpay/backend/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// RegisterRoutes mendaftarkan endpoint activity log ke grup /api.
func RegisterRoutes(
	apiGroup *gin.RouterGroup,
	h *activityloghandler.ActivityLogHandler,
	jwtSecret string,
	redisClient *redis.Client,
	userRepo userrepo.UserRepo,
) {
	requireAuth := middleware.AuthMiddleware(jwtSecret)

	g := apiGroup.Group("/activity-logs")
	g.Use(requireAuth)
	{
		g.GET("", middleware.PermissionMiddleware(userRepo, "activitylog.view"), h.List)
		g.GET("/:id", middleware.PermissionMiddleware(userRepo, "activitylog.view"), h.Detail)
		g.GET("/entity/:entityType/:entityID", middleware.PermissionMiddleware(userRepo, "activitylog.view"), h.EntityLogs)
	}
}
