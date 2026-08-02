package app

import (
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/academicyear"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/activeclass"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/activitylog"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/auth"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/classmembership"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/classtemplate"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/cohort"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/guardian"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/major"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/notification"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/role"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/semester"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/student"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/user"
	"github.com/redis/go-redis/v9"
)

// registerRoutes mendelegasikan pendaftaran endpoint ke masing-masing modul.
// Setiap modul memegang routes-nya sendiri agar app/routes.go tetap tipis.
// Notification template routes telah dihapus — template hanya dari embedded file.
// Semua modul menerima userRepo untuk PermissionMiddleware (kecuali auth — publik).
func registerRoutes(app *App, redisClient *redis.Client) {
	jwtSecret := app.Cfg.JWT.Secret
	api := app.Server.Group("/api/v1")

	// Auth — publik, tidak perlu PermissionMiddleware
	auth.RegisterRoutes(api, app.authHandler, jwtSecret)

	// Semua admin module — userRepo untuk PermissionMiddleware
	user.RegisterRoutes(api, app.userHandler, jwtSecret, redisClient, app.userRepo)
	role.RegisterRoutes(api, app.roleHandler, jwtSecret, redisClient, app.userRepo)
	activitylog.RegisterRoutes(api, app.activityLogHandler, jwtSecret, redisClient, app.userRepo)
	major.RegisterRoutes(api, app.majorHandler, jwtSecret, redisClient, app.userRepo)
	classtemplate.RegisterRoutes(api, app.classTemplateHandler, jwtSecret, redisClient, app.userRepo)
	cohort.RegisterRoutes(api, app.cohortHandler, jwtSecret, redisClient, app.userRepo)
	academicyear.RegisterRoutes(api, app.academicYearHandler, jwtSecret, redisClient, app.userRepo)
	semester.RegisterRoutes(api, app.semesterHandler, jwtSecret, redisClient, app.userRepo)
	activeclass.RegisterRoutes(api, app.activeClassHandler, jwtSecret, redisClient, app.userRepo)
	classmembership.RegisterRoutes(api, app.classMembershipHandler, jwtSecret, redisClient, app.userRepo)
	student.RegisterRoutes(api, app.studentHandler, jwtSecret, redisClient, app.userRepo)
	notification.RegisterRoutes(api, app.notificationHandler, jwtSecret, redisClient, app.userRepo)
	guardian.RegisterRoutes(api, app.guardianHandler, jwtSecret, redisClient, app.userRepo)
}
