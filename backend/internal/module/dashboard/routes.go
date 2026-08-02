package dashboard

import (
	dashhandler "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/dashboard/handler"

	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	apiGroup *gin.RouterGroup,
	h *dashhandler.DashboardHandler,
	jwtSecret string,
) {
	g := apiGroup.Group("/admin/dashboard")
	g.Use(middleware.AuthMiddleware(jwtSecret))
	g.Use(middleware.RoleMiddleware("super_admin"))
	{
		g.GET("/summary", h.Summary)
	}
}
