package classmembership

import (
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/middleware"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/classmembership/handler"
	userrepo "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/user/repository"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func RegisterRoutes(api *gin.RouterGroup, h *handler.ClassMembershipHandler, jwtSecret string, redisClient *redis.Client, userRepo userrepo.UserRepo) {
	requireAuth := middleware.AuthMiddleware(jwtSecret)

	// Keanggotaan kelas dikelola di bawah active-class agar URL terkait.
	g := api.Group("/active-classes")
	g.Use(requireAuth)
	{
		g.GET("/:id/memberships", middleware.PermissionMiddleware(userRepo, "classmembership.view"), h.GetAll)
		g.POST("/:id/memberships", middleware.PermissionMiddleware(userRepo, "classmembership.manage"), h.Enroll)
	}

	mg := api.Group("/class-memberships")
	mg.Use(requireAuth)
	{
		mg.GET("", middleware.PermissionMiddleware(userRepo, "classmembership.view"), h.GetAll)
		mg.PATCH("/:id/move", middleware.PermissionMiddleware(userRepo, "classmembership.manage"), h.Move)
		mg.PATCH("/:id/status", middleware.PermissionMiddleware(userRepo, "classmembership.manage"), h.SetStatus)
	}
}
