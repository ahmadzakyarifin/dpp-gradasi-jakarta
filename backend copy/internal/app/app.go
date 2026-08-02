package app

import (
	"context"

	"github.com/ahmadzakyarifin/schoolpay/backend/config"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/middleware"
	academicyearhandler "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/academicyear/handler"
	academicyearrepo "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/academicyear/repository"
	academicyearsvc "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/academicyear/service"
	activeclasshandler "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/activeclass/handler"
	activeclassrepo "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/activeclass/repository"
	activeclasssvc "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/activeclass/service"
	activityloghandler "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/activitylog/handler"
	authhandler "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/auth/handler"
	classmembershiphandler "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/classmembership/handler"
	classmembershiprepo "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/classmembership/repository"
	classmembershipsvc "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/classmembership/service"
	classtemplatehandler "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/classtemplate/handler"
	classtemplaterepo "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/classtemplate/repository"
	classtemplatesvc "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/classtemplate/service"
	cohorthandler "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/cohort/handler"
	cohortrepo "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/cohort/repository"
	cohortsvc "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/cohort/service"
	guardianhandler "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/guardian/handler"
	guardianrepo "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/guardian/repository"
	guardiansvc "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/guardian/service"
	majorhandler "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/major/handler"
	majorrepo "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/major/repository"
	majorsvc "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/major/service"
	notificationhandler "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/notification/handler"
	rolehandler "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/role/handler"
	semesterhandler "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/semester/handler"
	semesterrepo "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/semester/repository"
	semestersvc "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/semester/service"
	studenthandler "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/student/handler"
	userhandler "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/user/handler"
	userrepo "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/user/repository"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type App struct {
	Server *gin.Engine
	DB     *gorm.DB
	Cfg    config.Config

	userRepo               userrepo.UserRepo
	authHandler            *authhandler.AuthHandler
	userHandler            *userhandler.UserHandler
	roleHandler            *rolehandler.RoleHandler
	activityLogHandler     *activityloghandler.ActivityLogHandler
	majorHandler           *majorhandler.MajorHandler
	classTemplateHandler   *classtemplatehandler.Handler
	cohortHandler          *cohorthandler.CohortHandler
	academicYearHandler    *academicyearhandler.AcademicYearHandler
	semesterHandler        *semesterhandler.SemesterHandler
	activeClassHandler     *activeclasshandler.ActiveClassHandler
	classMembershipHandler *classmembershiphandler.ClassMembershipHandler
	studentHandler         *studenthandler.StudentHandler
	notificationHandler    *notificationhandler.NotificationHandler
	guardianHandler        *guardianhandler.GuardianHandler
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
	configureCloudflareClientIP(routerEngine)
	routerEngine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	routerEngine.Static("/uploads", "./public/uploads")

	appContext := context.Background()
	sf := buildSharedFeatures(database, appConfig, redisClient)

	appInstance := &App{
		Server:   routerEngine,
		DB:       database,
		Cfg:      *appConfig,
		userRepo: sf.UserRepository,
	}

	// Inisialisasi handler tiap modul.
	appInstance.authHandler = authhandler.NewAuthHandler(sf.AuthService, appConfig)
	appInstance.userHandler = userhandler.NewUserHandler(sf.UserService, appConfig)
	appInstance.roleHandler = rolehandler.NewRoleHandler(sf.RoleService)
	appInstance.activityLogHandler = activityloghandler.NewActivityLogHandler(sf.AuditLogService)

	// Major (jurusan) module wiring.
	majorRepo := majorrepo.NewMajorRepo(database)
	majorSvc := majorsvc.NewMajorService(database, majorRepo, sf.AuditLogService)
	appInstance.majorHandler = majorhandler.NewMajorHandler(majorSvc)

	// Academic: class_template submodule.
	classRepo := classtemplaterepo.NewClassTemplateRepo(database)
	classSvc := classtemplatesvc.NewClassTemplateService(database, classRepo, majorRepo, sf.AuditLogService)
	appInstance.classTemplateHandler = classtemplatehandler.NewHandler(classSvc)

	// Cohort (angkatan) module wiring.
	cohortRepo := cohortrepo.NewCohortRepo(database)
	cohortSvc := cohortsvc.NewCohortService(database, cohortRepo, sf.AuditLogService)
	appInstance.cohortHandler = cohorthandler.NewCohortHandler(cohortSvc)

	// Academic Year (tahun ajaran) module wiring.
	ayRepo := academicyearrepo.NewAcademicYearRepo(database)
	aySvc := academicyearsvc.NewAcademicYearService(database, ayRepo, sf.AuditLogService)
	appInstance.academicYearHandler = academicyearhandler.NewAcademicYearHandler(aySvc)

	// Semester module wiring.
	semRepo := semesterrepo.NewSemesterRepo(database)
	semSvc := semestersvc.NewSemesterService(database, semRepo, sf.AuditLogService)
	appInstance.semesterHandler = semesterhandler.NewSemesterHandler(semSvc)

	// Active Class module wiring.
	acRepo := activeclassrepo.NewActiveClassRepo(database)
	acSvc := activeclasssvc.NewActiveClassService(database, acRepo, sf.AuditLogService)
	appInstance.activeClassHandler = activeclasshandler.NewActiveClassHandler(acSvc)

	// Class Membership module wiring.
	cmRepo := classmembershiprepo.NewClassMembershipRepo(database)
	cmSvc := classmembershipsvc.NewClassMembershipService(database, cmRepo, sf.AuditLogService)
	appInstance.classMembershipHandler = classmembershiphandler.NewClassMembershipHandler(cmSvc)

	// Student module wiring.
	appInstance.studentHandler = studenthandler.NewStudentHandler(sf.StudentService)

	// Notification module handler.
	appInstance.notificationHandler = sf.NotificationHandler

	// Guardian (wali) module wiring.
	gRepo := guardianrepo.NewGuardianRepo(database)
	gSvc := guardiansvc.NewGuardianService(gRepo)
	appInstance.guardianHandler = guardianhandler.NewGuardianHandler(gSvc)

	// Jalankan proses non-HTTP: worker asynq (email/WAHA) + cleanup idempotency.
	startBackgroundJobs(appContext, database, appConfig, sf.JobClient)

	registerRoutes(appInstance, redisClient)

	return appInstance
}
