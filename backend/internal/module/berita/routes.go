package berita

import (
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/middleware"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/berita/handler"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(api *gin.RouterGroup, h *handler.BeritaHandler, jwtSecret string) {
	berita := api.Group("/berita")
	{
		// Publik
		berita.GET("", h.List)
		berita.GET("/categories", h.GetCategories)
		berita.GET("/:slug", h.GetBySlug)
	}

	// Admin — super_admin, admin, admin_berita
	admin := berita.Group("")
	admin.Use(middleware.AuthMiddleware(jwtSecret))
	admin.Use(middleware.RoleMiddleware("super_admin", "admin", "admin_berita"))
	admin.Use(middleware.RateLimitPerUser("berita-crud", 30))
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
		}
}
