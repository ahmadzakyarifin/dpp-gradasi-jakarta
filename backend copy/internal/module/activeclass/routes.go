package activeclass

import (
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/middleware"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/activeclass/handler"
	userrepo "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/user/repository"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func RegisterRoutes(api *gin.RouterGroup, h *handler.ActiveClassHandler, jwtSecret string, redisClient *redis.Client, userRepo userrepo.UserRepo) {
	requireAuth := middleware.AuthMiddleware(jwtSecret)

	g := api.Group("/active-classes")
	g.Use(requireAuth)
	{
		g.GET("", middleware.PermissionMiddleware(userRepo, "activeclass.view"), h.GetAll)
		g.POST("", middleware.PermissionMiddleware(userRepo, "activeclass.manage"), h.Create)
		g.GET("/:id", middleware.PermissionMiddleware(userRepo, "activeclass.view"), h.GetByID)
		g.PUT("/:id", middleware.PermissionMiddleware(userRepo, "activeclass.manage"), h.Update)
		g.DELETE("/:id", middleware.PermissionMiddleware(userRepo, "activeclass.manage"), h.Delete)
		g.PATCH("/:id/status", middleware.PermissionMiddleware(userRepo, "activeclass.manage"), h.ToggleStatus)
	}

	// Bulk save kelas aktif per tahun ajaran (hook frontend saveActiveClasses).
	yg := api.Group("/academic-years")
	yg.Use(requireAuth)
	yg.PUT("/:id/active-classes", middleware.PermissionMiddleware(userRepo, "activeclass.manage"), h.BulkUpsertByYear)
}
