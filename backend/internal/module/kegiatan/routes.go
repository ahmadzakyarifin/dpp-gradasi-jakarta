package kegiatan

import (
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/middleware"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/kegiatan/handler"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func RegisterRoutes(api *gin.RouterGroup, h *handler.KegiatanHandler, jwtSecret string, rl *redis.Client) {
	kgt := api.Group("/kegiatan")
	{
		// Publik
		kgt.GET("", h.List)
		kgt.GET("/categories", h.GetCategories)
		kgt.GET("/:slug", h.GetBySlug)
	}

	// Admin — super_admin, admin, editor
	admin := kgt.Group("")
	admin.Use(middleware.AuthMiddleware(jwtSecret))
	admin.Use(middleware.RoleMiddleware("super_admin", "admin"))
	admin.Use(middleware.RateLimitPerUser("kegiatan-crud", 30))
	{
		admin.GET("/admin", h.ListAdmin)
		admin.GET("/id/:id", h.GetByID)
		admin.POST("", h.Create)
		admin.PUT("/:id", h.Update)
		admin.POST("/bulk-delete", h.BulkDelete)
		admin.POST("/bulk-restore", h.BulkRestore)
		admin.DELETE("/:id", h.Delete)
		admin.POST("/:id/restore", h.Restore)
		admin.DELETE("/gallery/:gallery_id", h.DeleteGalleryImage)
	}
}
