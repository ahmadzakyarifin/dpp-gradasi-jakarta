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
	"github.com/redis/go-redis/v9"
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
	settingsHandler    settingshandler.SettingsHandler
	activityLogHandler *activityloghandler.ActivityLogHandler
	ActivityLogSvc     activitylogservice.ActivityLogService
	dashboardHandler   *dashboardhandler.DashboardHandler
}

func NewApp(database *gorm.DB, appConfig *config.Config, redisClient *redis.Client) *App {
	if appConfig.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	if redisClient == nil {
		panic("Redis client is required")
	}

	rl, err := middleware.NewRedisRateLimiter(redisClient)
	if err != nil {
		panic("failed to initialize Redis rate limiter: " + err.Error())
	}
	middleware.SetDefaultRateLimiter(rl)

	routerEngine := gin.Default()

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
	authSvc := authservice.NewAuthService(authRepo, appConfig, mailer, redisClient)
	authHandler := authhandler.NewAuthHandler(authSvc, appConfig, activityLogSvc)

	// Sliders
	slidersRepo := slidersrepo.NewSlidersRepo(database)
	slidersSvc := slidersservice.NewSlidersService(slidersRepo)
	slidersHandler := slidershandler.NewSlidersHandler(slidersSvc, activityLogSvc)

	// Berita
	beritaRepo := beritarepo.NewBeritaRepo(database)
	beritaSvc := beritaservice.NewBeritaService(beritaRepo)
	beritaHandler := beritahandler.NewBeritaHandler(beritaSvc, activityLogSvc)

	// Kegiatan
	kegiatanRepo := kegiatanrepo.NewKegiatanRepo(database)
	kegiatanSvc := kegiatanservice.NewKegiatanService(kegiatanRepo)
	kegiatanHandler := kegiatanhandler.NewKegiatanHandler(kegiatanSvc, activityLogSvc)

	// Kontak
	kontakRepo := kontakrepo.NewKontakRepo(database)
	kontakSvc := kontakservice.NewKontakService(kontakRepo)
	kontakHandler := kontakhandler.NewKontakHandler(kontakSvc, activityLogSvc)

	// Pengurus
	pengurusRepo := pengurusrepo.NewPengurusRepo(database)
	pengurusSvc := pengurusservice.NewPengurusService(pengurusRepo)
	pengurusHandler := pengurushandler.NewPengurusHandler(pengurusSvc, activityLogSvc)

	// User
	userRepo := userrepo.NewUserRepo(database)
	userSvc := userservice.NewUserService(userRepo, redisClient, mailer)
	userHandler := userhandler.NewUserHandler(userSvc, activityLogSvc, appConfig)

	// Settings
	settingsRepo := settingsrepo.NewSettingsRepo(database)
	settingsSvc := settingsservice.NewSettingsService(settingsRepo)
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
		settingsHandler:    settingsHandler,
		activityLogHandler: activityLogHandler,
		ActivityLogSvc:     activityLogSvc,
		dashboardHandler:   dashboardHandler,
	}

	registerRoutes(appInstance, redisClient)

	return appInstance
}
