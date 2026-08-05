package sliders

import (
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/middleware"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/sliders/handler"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(api *gin.RouterGroup, h *handler.SlidersHandler, jwtSecret string) {
	sliders := api.Group("/sliders")
	{
		// Publik — siapa pun bisa lihat
		sliders.GET("", h.GetAll)
		sliders.GET("/:id", h.GetByID)
	}

	// Protected — butuh auth + role super_admin/admin/editor
	protected := sliders.Group("")
	protected.Use(middleware.AuthMiddleware(jwtSecret))
	protected.Use(middleware.RoleMiddleware("super_admin", "admin"))
	protected.Use(middleware.RateLimitPerUser("sliders-crud", 30))
	{
		protected.GET("/admin", h.GetAllAdmin)
		protected.POST("", h.Create)
		protected.POST("/upload-image", h.UploadImage)
		protected.PUT("/reorder", h.Reorder)
		protected.PUT("/:id", h.Update)
		protected.POST("/bulk-delete", h.BulkDelete)
		protected.POST("/bulk-restore", h.BulkRestore)
		protected.DELETE("/:id", h.Delete)
		protected.POST("/:id/restore", h.Restore)
	}
}
