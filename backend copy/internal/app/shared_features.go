package app

import (
	"github.com/ahmadzakyarifin/schoolpay/backend/config"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/job"
	activitylogrepo "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/activitylog/repository"
	activitylogservice "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/activitylog/service"
	authrepo "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/auth/repository"
	authservice "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/auth/service"
	majorrepo "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/major/repository"
	majorservice "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/major/service"
	notificationhandler "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/notification/handler"
	notificationrepo "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/notification/repository"
	notificationservice "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/notification/service"
	rolerepo "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/role/repository"
	roleservice "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/role/service"
	studentrepo "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/student/repository"
	studentservice "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/student/service"
	userrepo "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/user/repository"
	userservice "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/user/service"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type sharedFeatures struct {
	UserRepository      userrepo.UserRepo
	AuthRepository      authrepo.AuthRepo
	RoleRepository      rolerepo.RoleRepo
	AuditLogService     activitylogservice.ActivityLogService
	UserService         userservice.UserService
	RoleService         roleservice.RoleService
	AuthService         authservice.AuthService
	MajorService        majorservice.MajorService
	StudentService      studentservice.StudentService
	JobClient           *jobs.Client
	NotificationService notificationservice.NotificationService
	NotificationRepo    notificationrepo.NotificationRepo
	NotificationHandler *notificationhandler.NotificationHandler
}

func buildSharedFeatures(
	database *gorm.DB,
	appConfig *config.Config,
	redisClient *redis.Client,
) sharedFeatures {
	userRepository := userrepo.NewUserRepo(database, redisClient)
	authRepository := authrepo.NewAuthRepo(database, redisClient)
	roleRepository := rolerepo.NewRoleRepo(database, redisClient)
	majorRepo := majorrepo.NewMajorRepo(database)

	// Audit log (activitylog) dipakai oleh banyak service sebagai pencatat aktivitas.
	auditLogService := activitylogservice.NewActivityLogService(
		activitylogrepo.NewActivityLogRepository(database),
	)

	// Client job asynq (email/WAHA) dipakai auth & notification untuk enqueue.
	jobClient := jobs.NewClient(appConfig)

	// Notification repos & service — tanpa template CRUD.
	notifRepo := notificationrepo.NewNotificationRepo(database)
	notifService := notificationservice.NewNotificationService(database, notifRepo, jobClient)

	// Service tiap modul.
	authService := authservice.NewAuthService(authRepository, auditLogService, notifService, appConfig)
	userService := userservice.NewUserService(
		database,
		userRepository,
		authRepository,
		jobClient,
		auditLogService,
		notifService,
		appConfig,
	)
	roleService := roleservice.NewRoleService(database, roleRepository, auditLogService)
	majorSvc := majorservice.NewMajorService(database, majorRepo, auditLogService)

	// Student module.
	studentRepo := studentrepo.NewStudentRepo(database)
	studentSvc := studentservice.NewStudentService(database, studentRepo, auditLogService, appConfig)

	return sharedFeatures{
		UserRepository:      userRepository,
		AuthRepository:      authRepository,
		RoleRepository:      roleRepository,
		AuditLogService:     auditLogService,
		UserService:         userService,
		RoleService:         roleService,
		AuthService:         authService,
		MajorService:        majorSvc,
		StudentService:      studentSvc,
		JobClient:           jobClient,
		NotificationService: notifService,
		NotificationRepo:    notifRepo,
		NotificationHandler: notificationhandler.NewNotificationHandler(notifRepo),
	}
}
