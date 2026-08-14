package kegiatan

import (
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/middleware"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/kegiatan/handler"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(api *gin.RouterGroup, h *handler.KegiatanHandler, jwtSecret string) {
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
	admin.Use(middleware.RoleMiddleware("super_admin", "admin", "editor"))
	admin.Use(middleware.RateLimitPerUser("kegiatan-crud", 30))
	{
		admin.GET("/admin", h.ListAdmin)
		admin.GET("/id/:id", h.GetByID)
		admin.POST("", h.Create)
		admin.POST("/upload-image", h.UploadImage)
		admin.PUT("/categories", h.RenameCategory)
		admin.DELETE("/categories/:name", h.DeleteCategory)
		admin.PUT("/:id", h.Update)
		admin.POST("/bulk-delete", h.BulkDelete)
		admin.POST("/bulk-restore", h.BulkRestore)
		admin.DELETE("/:id", h.Delete)
		admin.POST("/:id/restore", h.Restore)
		admin.DELETE("/gallery/:gallery_id", h.DeleteGalleryImage)
	}
}
