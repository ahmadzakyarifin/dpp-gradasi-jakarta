package kontak

import (
	"time"

	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/middleware"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/kontak/handler"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(api *gin.RouterGroup, h *handler.KontakHandler, jwtSecret string) {
	// Publik — submit pesan (rate limit per IP: 5 request/menit anti spam)
	kontakLimit := middleware.RateLimitRules("kontak_submit",
		middleware.IP(5, time.Minute),
	)
	// Captcha juga untuk form kontak publik (target spam klasik tanpa login),
	// hanya aktif kalau CAPTCHA_ENABLED=true + secret terisi di backend.
	needCaptcha := middleware.RequiredCaptchaMiddleware(h.Cfg())
	api.POST("/kontak", needCaptcha, kontakLimit, h.Submit)

	// Admin — super_admin, admin
	admin := api.Group("/admin/kontak")
	admin.Use(middleware.AuthMiddleware(jwtSecret))
	admin.Use(middleware.RoleMiddleware("super_admin", "admin"))
	{
		admin.GET("", h.List)
		admin.GET("/:id", h.GetByID)
		admin.POST("/bulk-delete", h.BulkDelete)
		admin.POST("/bulk-restore", h.BulkRestore)
		admin.DELETE("/:id", h.Delete)
		admin.POST("/:id/restore", h.Restore)
	}
}
