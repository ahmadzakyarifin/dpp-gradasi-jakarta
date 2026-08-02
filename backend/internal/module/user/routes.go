package user

import (
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/middleware"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/user/handler"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes mendaftarkan endpoint manajemen pengguna ke grup /api.
// Self-service endpoints (profile, activate) tidak menggunakan otorisasi role.
func RegisterRoutes(
	apiGroup *gin.RouterGroup,
	h *handler.UserHandler,
	jwtSecret string,
) {
	requireAuth := middleware.AuthMiddleware(jwtSecret)

	// Self-service: profil & aktivasi akun (tanpa otorisasi role)
	users := apiGroup.Group("/users")
	{
		users.PUT("/profile", requireAuth, h.UpdateProfile)
		users.POST("/:id/activate", h.Activate)
	}

	// Admin: manajemen pengguna — hanya super_admin & admin
	admin := apiGroup.Group("/admin/users")
	admin.Use(requireAuth, middleware.RoleMiddleware("super_admin", "admin"))
	{
		admin.GET("", h.GetAll)
		admin.GET("/:id", h.GetByID)
		admin.POST("", h.Create)
		admin.PUT("/:id", h.Update)
		admin.DELETE("/:id", h.Delete)
		admin.POST("/:id/status", h.ToggleStatus)
		admin.POST("/:id/resend-activation", h.ResendNotification)
		admin.POST("/:id/restore", h.Restore)
		admin.POST("/bulk-delete", h.BulkDelete)
		admin.POST("/bulk-restore", h.BulkRestore)
	}
}
