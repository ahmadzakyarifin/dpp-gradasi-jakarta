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
	"github.com/redis/go-redis/v9"
)

func registerRoutes(app *App, redisClient *redis.Client) {
	api := app.Server.Group("/api/v1")

	auth.RegisterRoutes(api, app.authHandler, app.Cfg.JWT.Secret, redisClient, &app.Cfg)
	sliders.RegisterRoutes(api, app.slidersHandler, app.Cfg.JWT.Secret, redisClient)
	berita.RegisterRoutes(api, app.beritaHandler, app.Cfg.JWT.Secret, redisClient)
	kegiatan.RegisterRoutes(api, app.kegiatanHandler, app.Cfg.JWT.Secret, redisClient)
	kontak.RegisterRoutes(api, app.kontakHandler, app.Cfg.JWT.Secret, redisClient)
	pengurus.RegisterRoutes(api, app.pengurusHandler, app.Cfg.JWT.Secret, redisClient)
	user.RegisterRoutes(api, app.userHandler, app.Cfg.JWT.Secret, redisClient)
	settings.RegisterRoutes(api, app.settingsHandler, app.Cfg.JWT.Secret, redisClient)
	activitylog.RegisterRoutes(api, app.activityLogHandler, app.Cfg.JWT.Secret, redisClient)
	dashboard.RegisterRoutes(api, app.dashboardHandler, app.Cfg.JWT.Secret, redisClient)
}
