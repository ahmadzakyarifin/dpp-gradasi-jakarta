package pengurus

import (
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/middleware"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/pengurus/handler"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func RegisterRoutes(api *gin.RouterGroup, h *handler.PengurusHandler, jwtSecret string, rl *redis.Client) {
	// Publik
	pengurus := api.Group("/pengurus")
	{
		pengurus.GET("", h.ListPublic)
		pengurus.GET("/regions", h.Regions)
	}

	// Admin — super_admin, admin, editor
	admin := api.Group("/admin/pengurus")
	admin.Use(middleware.AuthMiddleware(jwtSecret))
	admin.Use(middleware.RoleMiddleware("super_admin", "admin"))
	admin.Use(middleware.RateLimitPerUser("pengurus-crud", 30))
	{
		admin.GET("", h.ListAdmin)
		admin.POST("", h.Create)
		admin.PUT("/:id", h.Update)
		admin.POST("/bulk-delete", h.BulkDelete)
		admin.POST("/bulk-restore", h.BulkRestore)
		admin.DELETE("/:id", h.Delete)
		admin.POST("/:id/restore", h.Restore)
	}
}
