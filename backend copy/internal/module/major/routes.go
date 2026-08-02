package major

import (
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/middleware"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/major/handler"
	userrepo "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/user/repository"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func RegisterRoutes(
	api *gin.RouterGroup,
	majorH *handler.MajorHandler,
	jwtSecret string,
	redisClient *redis.Client,
	userRepo userrepo.UserRepo,
) {
	requireAuth := middleware.AuthMiddleware(jwtSecret)

	g := api.Group("/major")
	g.Use(requireAuth)
	{
		g.GET("", middleware.PermissionMiddleware(userRepo, "major.view"), majorH.GetAll)
		g.POST("", middleware.PermissionMiddleware(userRepo, "major.manage"), majorH.Create)
		g.PUT("/:id", middleware.PermissionMiddleware(userRepo, "major.manage"), majorH.Update)
		g.DELETE("/:id", middleware.PermissionMiddleware(userRepo, "major.manage"), majorH.Delete)
		g.PATCH("/:id/status", middleware.PermissionMiddleware(userRepo, "major.manage"), majorH.UpdateStatus)
		g.PATCH("/:id/restore", middleware.PermissionMiddleware(userRepo, "major.manage"), majorH.Restore)
		g.POST("/bulk-delete", middleware.PermissionMiddleware(userRepo, "major.manage"), majorH.BulkDelete)
		g.PATCH("/bulk-restore", middleware.PermissionMiddleware(userRepo, "major.manage"), majorH.BulkRestore)
		g.GET("/:id/dependency-info", middleware.PermissionMiddleware(userRepo, "major.view"), majorH.GetDependencyInfo)
	}
}
