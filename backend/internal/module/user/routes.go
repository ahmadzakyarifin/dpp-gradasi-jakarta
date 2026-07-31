package user

import (
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/middleware"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/user/handler"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func RegisterRoutes(api *gin.RouterGroup, h *handler.UserHandler, jwtSecret string, rl *redis.Client) {
	// Profil Sendiri (bisa diakses semua admin)
	profile := api.Group("/profile")
	profile.Use(middleware.AuthMiddleware(jwtSecret))
	{
		profile.GET("", h.GetProfile)
		profile.PUT("", h.UpdateProfile)
		profile.PUT("/password", h.ChangePassword)
		profile.POST("/verify-email", h.VerifyEmail)
	}

	// Manajemen Admin Lain (hanya super_admin yang bisa)
	admin := api.Group("/admin/users")
	admin.Use(middleware.AuthMiddleware(jwtSecret))
	admin.Use(middleware.RoleMiddleware("super_admin")) // Hanya role super_admin (ID 1) yang diizinkan
	{
		admin.GET("", h.GetAdmins)
		admin.POST("", h.CreateAdmin)
		admin.POST("/bulk-delete", h.BulkDeleteAdmin)
		admin.POST("/bulk-restore", h.BulkRestoreAdmin)
		admin.DELETE("/:id", h.DeleteAdmin)
		admin.POST("/:id/restore", h.RestoreAdmin)
		admin.POST("/:id/resend-activation", h.ResendActivation)
		admin.PUT("/:id/status", h.SetAdminStatus)
	}
}
