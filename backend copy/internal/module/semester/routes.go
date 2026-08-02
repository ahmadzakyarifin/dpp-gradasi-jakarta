package semester

import (
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/middleware"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/semester/handler"
	userrepo "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/user/repository"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func RegisterRoutes(api *gin.RouterGroup, h *handler.SemesterHandler, jwtSecret string, redisClient *redis.Client, userRepo userrepo.UserRepo) {
	requireAuth := middleware.AuthMiddleware(jwtSecret)

	g := api.Group("/semesters")
	g.Use(requireAuth)
	{
		g.GET("", middleware.PermissionMiddleware(userRepo, "semester.view"), h.GetAll)
		g.POST("", middleware.PermissionMiddleware(userRepo, "semester.manage"), h.Create)
		g.GET("/check-unique", middleware.PermissionMiddleware(userRepo, "semester.view"), h.CheckUnique)
		g.PUT("/:id", middleware.PermissionMiddleware(userRepo, "semester.manage"), h.Update)
		g.DELETE("/:id", middleware.PermissionMiddleware(userRepo, "semester.manage"), h.Delete)
		g.PATCH("/:id/status", middleware.PermissionMiddleware(userRepo, "semester.manage"), h.ToggleStatus)
		g.PATCH("/:id/restore", middleware.PermissionMiddleware(userRepo, "semester.manage"), h.Restore)
		g.POST("/bulk-delete", middleware.PermissionMiddleware(userRepo, "semester.manage"), h.BulkDelete)
		g.PATCH("/bulk-restore", middleware.PermissionMiddleware(userRepo, "semester.manage"), h.BulkRestore)
		g.GET("/:id/dependency-info", middleware.PermissionMiddleware(userRepo, "semester.view"), h.GetDependencyInfo)
	}
}
