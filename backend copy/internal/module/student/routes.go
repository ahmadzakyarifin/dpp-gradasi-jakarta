package student

import (
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/middleware"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/student/handler"
	userrepo "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/user/repository"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func RegisterRoutes(
	api *gin.RouterGroup,
	h *handler.StudentHandler,
	jwtSecret string,
	redisClient *redis.Client,
	userRepo userrepo.UserRepo,
) {
	requireAuth := middleware.AuthMiddleware(jwtSecret)

	g := api.Group("/students")
	g.Use(requireAuth)
	{
		g.GET("", middleware.PermissionMiddleware(userRepo, "student.view"), h.GetAll)
		g.POST("", middleware.PermissionMiddleware(userRepo, "student.create"), h.Create)
		g.PUT("/:id", middleware.PermissionMiddleware(userRepo, "student.update"), h.Update)
		g.DELETE("/:id", middleware.PermissionMiddleware(userRepo, "student.delete"), h.Delete)
		g.PATCH("/:id/restore", middleware.PermissionMiddleware(userRepo, "student.update"), h.Restore)
		g.PATCH("/:id/status", middleware.PermissionMiddleware(userRepo, "student.update"), h.ToggleStatus)
		g.POST("/bulk-delete", middleware.PermissionMiddleware(userRepo, "student.delete"), h.BulkDelete)
		g.PATCH("/bulk-restore", middleware.PermissionMiddleware(userRepo, "student.update"), h.BulkRestore)
		g.GET("/filters", middleware.PermissionMiddleware(userRepo, "student.view"), h.GetFilters)
		g.GET("/export", middleware.PermissionMiddleware(userRepo, "student.export"), h.Export)
		g.POST("/academic-movements", middleware.PermissionMiddleware(userRepo, "student.promote"), h.BulkPromote)
		g.POST("/graduations", middleware.PermissionMiddleware(userRepo, "student.graduate"), h.BulkGraduate)
		g.GET("/:id", middleware.PermissionMiddleware(userRepo, "student.view"), h.GetByID)
		g.GET("/:id/history", middleware.PermissionMiddleware(userRepo, "student.view"), h.GetClassHistory)
		g.GET("/:id/dependency-info", middleware.PermissionMiddleware(userRepo, "student.view"), h.GetDependencyInfo)
		g.GET("/check-unique", middleware.PermissionMiddleware(userRepo, "student.view"), h.CheckUnique)
	}
}
