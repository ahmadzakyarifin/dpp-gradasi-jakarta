package app

import (
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/activitylog"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/auth"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/berita"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/dashboard"
	kegiatan "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/kegiatan"
	kontak "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/kontak"
	pengurus "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/pengurus"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/settings"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/sliders"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/user"
)

func registerRoutes(app *App) {
	api := app.Server.Group("/api/v1")

	auth.RegisterRoutes(api, app.authHandler, app.Cfg.JWT.Secret)
	sliders.RegisterRoutes(api, app.slidersHandler, app.Cfg.JWT.Secret)
	berita.RegisterRoutes(api, app.beritaHandler, app.Cfg.JWT.Secret)
	kegiatan.RegisterRoutes(api, app.kegiatanHandler, app.Cfg.JWT.Secret)
	kontak.RegisterRoutes(api, app.kontakHandler, app.Cfg.JWT.Secret)
	pengurus.RegisterRoutes(api, app.pengurusHandler, app.Cfg.JWT.Secret)
	user.RegisterRoutes(api, app.userHandler, app.Cfg.JWT.Secret)
	settings.RegisterRoutes(api, app.settingsHandler, app.Cfg.JWT.Secret)
	activitylog.RegisterRoutes(api, app.activityLogHandler, app.Cfg.JWT.Secret)
	dashboard.RegisterRoutes(api, app.dashboardHandler, app.Cfg.JWT.Secret)

	// SEO Open Graph Preview routes — group terpisah TANPA rate limiter & Turnstile.
	// Bot crawler social media (WA, FB, Twitter) sering fetch berulang kali.
	seoGroup := app.Server.Group("/api/v1/seo")
	seoGroup.GET("/berita/:slug", app.beritaHandler.SeoBySlug)
	seoGroup.HEAD("/berita/:slug", app.beritaHandler.SeoBySlug)
	seoGroup.GET("/kegiatan/:slug", app.kegiatanHandler.SeoBySlug)
	seoGroup.HEAD("/kegiatan/:slug", app.kegiatanHandler.SeoBySlug)
}
