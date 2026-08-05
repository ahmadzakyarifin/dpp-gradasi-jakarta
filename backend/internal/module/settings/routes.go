package settings

import (
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/middleware"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/settings/handler"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(api *gin.RouterGroup, handler handler.SettingsHandler, jwtSecret string) {
	settings := api.Group("/settings")
	{
		settings.GET("", handler.GetSettings)

		admin := api.Group("/admin/settings")
		admin.Use(middleware.AuthMiddleware(jwtSecret))
		admin.Use(middleware.RateLimitPerUser("settings-admin", 30))
		admin.Use(middleware.RoleMiddleware("super_admin"))
		{
			admin.PUT("", handler.UpdateSettings)
			admin.POST("/logo", handler.UploadLogo)
		}
	}
}
