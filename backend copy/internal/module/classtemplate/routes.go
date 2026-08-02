package classtemplate

import (
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/middleware"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/classtemplate/handler"
	userrepo "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/user/repository"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func RegisterRoutes(
	api *gin.RouterGroup,
	ctH *handler.Handler,
	jwtSecret string,
	redisClient *redis.Client,
	userRepo userrepo.UserRepo,
) {
	requireAuth := middleware.AuthMiddleware(jwtSecret)

	g := api.Group("/class-template")
	g.Use(requireAuth)
	{
		g.GET("", middleware.PermissionMiddleware(userRepo, "classtemplate.view"), ctH.GetAll)
		g.POST("", middleware.PermissionMiddleware(userRepo, "classtemplate.manage"), ctH.Create)
		g.PUT("/:id", middleware.PermissionMiddleware(userRepo, "classtemplate.manage"), ctH.Update)
		g.DELETE("/:id", middleware.PermissionMiddleware(userRepo, "classtemplate.manage"), ctH.Delete)
		g.PATCH("/:id/status", middleware.PermissionMiddleware(userRepo, "classtemplate.manage"), ctH.ToggleStatus)
		g.PATCH("/:id/restore", middleware.PermissionMiddleware(userRepo, "classtemplate.manage"), ctH.Restore)
		g.POST("/bulk-delete", middleware.PermissionMiddleware(userRepo, "classtemplate.manage"), ctH.BulkDelete)
		g.PATCH("/bulk-restore", middleware.PermissionMiddleware(userRepo, "classtemplate.manage"), ctH.BulkRestore)
		g.GET("/:id/dependency-info", middleware.PermissionMiddleware(userRepo, "classtemplate.view"), ctH.GetDependencyInfo)
		g.GET("/check-unique", middleware.PermissionMiddleware(userRepo, "classtemplate.view"), ctH.CheckUnique)
		g.GET("/suggest-next-name", middleware.PermissionMiddleware(userRepo, "classtemplate.view"), ctH.SuggestNextName)
	}
}
