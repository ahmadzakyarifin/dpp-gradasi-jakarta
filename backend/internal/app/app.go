package app

import (
	"time"

	"github.com/ahmadzakyarifin/dpp-gradasi/backend/config"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/infrastructure"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/middleware"
	authhandler "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/auth/handler"
	authrepo "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/auth/repository"
	authservice "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/auth/service"
	beritahandler "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/berita/handler"
	beritarepo "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/berita/repository"
	beritaservice "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/berita/service"
	kegiatanhandler "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/kegiatan/handler"
	kegiatanrepo "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/kegiatan/repository"
	kegiatanservice "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/kegiatan/service"
	kontakhandler "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/kontak/handler"
	kontakrepo "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/kontak/repository"
	kontakservice "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/kontak/service"
	pengurushandler "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/pengurus/handler"
	pengurusrepo "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/pengurus/repository"
	pengurusservice "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/pengurus/service"
	rolehandler "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/role/handler"
	rolerepo "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/role/repository"
	roleservice "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/role/service"
	settingshandler "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/settings/handler"
	settingsrepo "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/settings/repository"
	settingsservice "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/settings/service"
	slidershandler "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/sliders/handler"
	slidersrepo "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/sliders/repository"
	slidersservice "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/sliders/service"
	userhandler "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/user/handler"
	userrepo "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/user/repository"
	userservice "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/user/service"

	activityloghandler "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/activitylog/handler"
	activitylogrepo "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/activitylog/repository"
	activitylogservice "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/activitylog/service"
	dashboardhandler "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/dashboard/handler"
	dashboardservice "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/dashboard/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type App struct {
	Server             *gin.Engine
	DB                 *gorm.DB
	Cfg                config.Config
	Mailer             *infrastructure.Mailer
	authHandler        *authhandler.AuthHandler
	slidersHandler     *slidershandler.SlidersHandler
	beritaHandler      *beritahandler.BeritaHandler
	kegiatanHandler    *kegiatanhandler.KegiatanHandler
	kontakHandler      *kontakhandler.KontakHandler
	pengurusHandler    *pengurushandler.PengurusHandler
	userHandler        *userhandler.UserHandler
	roleHandler        *rolehandler.RoleHandler
	settingsHandler    settingshandler.SettingsHandler
	activityLogHandler *activityloghandler.ActivityLogHandler
	ActivityLogSvc     activitylogservice.ActivityLogService
	dashboardHandler   *dashboardhandler.DashboardHandler
}

func NewApp(database *gorm.DB, appConfig *config.Config) *App {
	if appConfig.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Rate limiter in-memory (gaya DPP murni, tanpa Redis).
	rl := middleware.NewRateLimiter()
	middleware.SetDefaultRateLimiter(rl)

	routerEngine := gin.Default()

	// Client IP akurat di belakang Cloudflare (penting utk rate limiter + audit IP)
	if appConfig.App.Env == "production" {
		configureCloudflareClientIP(routerEngine)
	}

	// Enable CORS
	routerEngine.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	routerEngine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Serve uploaded files (logo, image pengurus, foto user, dll)
	routerEngine.Static("/uploads", "./public/uploads")

	mailer := infrastructure.NewMailer(appConfig)

	// Activity Log
	activityLogRepo := activitylogrepo.NewActivityLogRepository(database)
	activityLogSvc := activitylogservice.NewActivityLogService(activityLogRepo)
	activityLogHandler := activityloghandler.NewActivityLogHandler(activityLogSvc)

	// Auth
	authRepo := authrepo.NewAuthRepo(database)
	authSvc := authservice.NewAuthService(authRepo, activityLogSvc, mailer, appConfig)
	authHandler := authhandler.NewAuthHandler(authSvc, appConfig)

	// Sliders
	slidersRepo := slidersrepo.NewSlidersRepo(database)
	slidersSvc := slidersservice.NewSlidersService(database, slidersRepo, activityLogSvc)
	slidersHandler := slidershandler.NewSlidersHandler(slidersSvc)

	// Berita
	beritaRepo := beritarepo.NewBeritaRepo(database)
	beritaSvc := beritaservice.NewBeritaService(database, beritaRepo, activityLogSvc)
	beritaHandler := beritahandler.NewBeritaHandler(beritaSvc)

	// Kegiatan
	kegiatanRepo := kegiatanrepo.NewKegiatanRepo(database)
	kegiatanSvc := kegiatanservice.NewKegiatanService(database, kegiatanRepo, activityLogSvc)
	kegiatanHandler := kegiatanhandler.NewKegiatanHandler(kegiatanSvc)

	// Kontak
	kontakRepo := kontakrepo.NewKontakRepo(database)
	kontakSvc := kontakservice.NewKontakService(database, kontakRepo, activityLogSvc)
	kontakHandler := kontakhandler.NewKontakHandler(kontakSvc)

	// Pengurus
	pengurusRepo := pengurusrepo.NewPengurusRepo(database)
	pengurusSvc := pengurusservice.NewPengurusService(database, pengurusRepo, activityLogSvc)
	pengurusHandler := pengurushandler.NewPengurusHandler(pengurusSvc)

	// User
	userRepo := userrepo.NewUserRepo(database)
	userSvc := userservice.NewUserService(database, userRepo, authRepo, activityLogSvc, mailer, appConfig)
	userHandler := userhandler.NewUserHandler(userSvc, appConfig)

	// Role
	roleRepo := rolerepo.NewRoleRepo(database)
	roleSvc := roleservice.NewRoleService(database, roleRepo, activityLogSvc)
	roleHandler := rolehandler.NewRoleHandler(roleSvc)

	// Settings
	settingsRepo := settingsrepo.NewSettingsRepo(database)
	settingsSvc := settingsservice.NewSettingsService(database, settingsRepo, activityLogSvc)
	settingsHandler := settingshandler.NewSettingsHandler(settingsSvc)

	// Dashboard
	dashboardSvc := dashboardservice.NewDashboardService(database, activityLogSvc)
	dashboardHandler := dashboardhandler.NewDashboardHandler(dashboardSvc)

	appInstance := &App{
		Server:             routerEngine,
		DB:                 database,
		Cfg:                *appConfig,
		Mailer:             mailer,
		authHandler:        authHandler,
		slidersHandler:     slidersHandler,
		beritaHandler:      beritaHandler,
		kegiatanHandler:    kegiatanHandler,
		kontakHandler:      kontakHandler,
		pengurusHandler:    pengurusHandler,
		userHandler:        userHandler,
		roleHandler:        roleHandler,
		settingsHandler:    settingsHandler,
		activityLogHandler: activityLogHandler,
		ActivityLogSvc:     activityLogSvc,
		dashboardHandler:   dashboardHandler,
	}

	registerRoutes(appInstance)

	return appInstance
}
