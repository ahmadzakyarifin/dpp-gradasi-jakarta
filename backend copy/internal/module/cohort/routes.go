package cohort

import (
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/middleware"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/cohort/handler"
	userrepo "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/user/repository"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func RegisterRoutes(api *gin.RouterGroup, cohortHandler *handler.CohortHandler, jwtSecret string, redisClient *redis.Client, userRepo userrepo.UserRepo) {
	requireAuth := middleware.AuthMiddleware(jwtSecret)

	g := api.Group("/cohorts")
	g.Use(requireAuth)
	{
		g.GET("", middleware.PermissionMiddleware(userRepo, "cohort.view"), cohortHandler.GetAll)
		g.POST("", middleware.PermissionMiddleware(userRepo, "cohort.manage"), cohortHandler.Create)
		g.GET("/check-unique", middleware.PermissionMiddleware(userRepo, "cohort.view"), cohortHandler.CheckUnique)
		g.PUT("/:id", middleware.PermissionMiddleware(userRepo, "cohort.manage"), cohortHandler.Update)
		g.DELETE("/:id", middleware.PermissionMiddleware(userRepo, "cohort.manage"), cohortHandler.Delete)
		g.PATCH("/:id/status", middleware.PermissionMiddleware(userRepo, "cohort.manage"), cohortHandler.ToggleStatus)
		g.PATCH("/:id/restore", middleware.PermissionMiddleware(userRepo, "cohort.manage"), cohortHandler.Restore)
		g.POST("/bulk-delete", middleware.PermissionMiddleware(userRepo, "cohort.manage"), cohortHandler.BulkDelete)
		g.PATCH("/bulk-restore", middleware.PermissionMiddleware(userRepo, "cohort.manage"), cohortHandler.BulkRestore)
		g.GET("/:id/dependency-info", middleware.PermissionMiddleware(userRepo, "cohort.view"), cohortHandler.GetDependencyInfo)
	}
}
