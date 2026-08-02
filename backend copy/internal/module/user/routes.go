package user

import (
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/middleware"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/user/handler"
	userrepo "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/user/repository"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// RegisterRoutes mendaftarkan endpoint manajemen pengguna ke grup /api.
// Self-service endpoints (profile, activate) tidak menggunakan PermissionMiddleware.
func RegisterRoutes(
	apiGroup *gin.RouterGroup,
	h *handler.UserHandler,
	jwtSecret string,
	redisClient *redis.Client,
	userRepo userrepo.UserRepo,
) {
	requireAuth := middleware.AuthMiddleware(jwtSecret)

	users := apiGroup.Group("/users")
	users.Use(requireAuth)
	{
		// Admin: manajemen pengguna — diproteksi permission
		users.GET("", middleware.PermissionMiddleware(userRepo, "user.view"), h.GetAll)
		users.GET("/:id", middleware.PermissionMiddleware(userRepo, "user.view"), h.GetByID)
		users.POST("", middleware.PermissionMiddleware(userRepo, "user.create"), h.Create)
		users.PUT("/:id", middleware.PermissionMiddleware(userRepo, "user.update"), h.Update)
		users.DELETE("/:id", middleware.PermissionMiddleware(userRepo, "user.delete"), h.Delete)
		users.POST("/:id/toggle-status", middleware.PermissionMiddleware(userRepo, "user.update"), h.ToggleStatus)
		users.POST("/:id/resend-notification", middleware.PermissionMiddleware(userRepo, "user.update"), h.ResendNotification)
		users.POST("/bulk/resend-notification", middleware.PermissionMiddleware(userRepo, "user.update"), h.BulkResendNotification)
		users.POST("/bulk/delete", middleware.PermissionMiddleware(userRepo, "user.delete"), h.BulkDelete)
		users.POST("/:id/restore", middleware.PermissionMiddleware(userRepo, "user.update"), h.Restore)
		users.POST("/bulk/restore", middleware.PermissionMiddleware(userRepo, "user.update"), h.BulkRestore)
		users.GET("/export", middleware.PermissionMiddleware(userRepo, "user.export"), h.Export)
		users.GET("/:id/dependency-info", middleware.PermissionMiddleware(userRepo, "user.view"), h.GetDependencyInfo)
		users.GET("/check-unique", middleware.PermissionMiddleware(userRepo, "user.view"), h.CheckUnique)

		// Self-service — tanpa permission
		users.PUT("/profile", h.UpdateProfile)
		// POST /:id/activate = aktivasi via token, tidak pakai permission
		users.POST("/:id/activate", h.Activate)
	}
}
