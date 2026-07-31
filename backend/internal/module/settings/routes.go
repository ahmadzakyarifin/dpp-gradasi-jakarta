package settings

import (
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/middleware"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/settings/handler"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func RegisterRoutes(api *gin.RouterGroup, handler handler.SettingsHandler, jwtSecret string, redisClient *redis.Client) {
	settings := api.Group("/settings")
	{
		settings.GET("", handler.GetSettings)

		admin := api.Group("/admin/settings")
		admin.Use(middleware.AuthMiddleware(jwtSecret))
		admin.Use(middleware.RateLimitPerUser("settings-admin", 30))
		admin.Use(middleware.RoleMiddleware("super_admin", "admin"))
		{
			admin.POST("", handler.UpdateSettings)
		}
	}
}
