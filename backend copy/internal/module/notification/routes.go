package notification

import (
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/middleware"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/notification/handler"
	userrepo "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/user/repository"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// RegisterRoutes mendaftarkan endpoint notifikasi (efikasi).
func RegisterRoutes(apiGroup *gin.RouterGroup, h *handler.NotificationHandler, jwtSecret string, redisClient *redis.Client, userRepo userrepo.UserRepo) {
	requireAuth := middleware.AuthMiddleware(jwtSecret)
	notif := apiGroup.Group("/notifications")
	notif.Use(requireAuth)
	{
		notif.GET("", middleware.PermissionMiddleware(userRepo, "notification.view"), h.ListNotifications)
	}
}
