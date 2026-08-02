package guardian

import (
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/middleware"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/guardian/handler"
	userrepo "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/user/repository"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func RegisterRoutes(api *gin.RouterGroup, h *handler.GuardianHandler, jwtSecret string, redisClient *redis.Client, userRepo userrepo.UserRepo) {
	requireAuth := middleware.AuthMiddleware(jwtSecret)

	g := api.Group("/guardians")
	g.Use(requireAuth)
	{
		g.GET("", middleware.PermissionMiddleware(userRepo, "guardian.view"), h.GetAll)
		g.POST("", middleware.PermissionMiddleware(userRepo, "guardian.manage"), h.Create)
		g.PUT("/:id", middleware.PermissionMiddleware(userRepo, "guardian.manage"), h.Update)
		g.DELETE("/:id", middleware.PermissionMiddleware(userRepo, "guardian.manage"), h.Delete)
	}
}
