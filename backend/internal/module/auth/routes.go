package auth

import (
	"time"

	"github.com/ahmadzakyarifin/dpp-gradasi/backend/config"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/middleware"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/auth/handler"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func RegisterRoutes(
	api *gin.RouterGroup,
	authHandler *handler.AuthHandler,
	jwtSecret string,
	redisClient *redis.Client,
	cfg *config.Config,
) {
	// Public routes (no auth)
	auth := api.Group("/auth")
	{
		auth.POST("/login",
			middleware.RequiredCaptchaMiddleware(cfg),
			middleware.RateLimitRules("auth-login",
				middleware.IPEmail(5, 1*time.Minute), // 5 req/min per IP+Email
			), authHandler.Login)

		auth.POST("/refresh", authHandler.Refresh)

		auth.POST("/forgot-password",
			middleware.RequiredCaptchaMiddleware(cfg),
			middleware.RateLimitRules("auth-forgot",
				middleware.IPEmail(3, 1*time.Minute), // 3 req/min per IP+Email
			), authHandler.ForgotPassword)

		auth.GET("/validate-reset-token", authHandler.ValidateResetToken)

		auth.POST("/reset-password", authHandler.ResetPassword)

		auth.GET("/validate-activation-token", authHandler.ValidateActivationToken)

		auth.POST("/activate-account", authHandler.ActivateAccount)
	}

	// Protected routes (need auth)
	protected := api.Group("/auth")
	protected.Use(middleware.AuthMiddleware(jwtSecret))
	{
		protected.POST("/logout", authHandler.Logout)
		protected.POST("/change-password", authHandler.ChangePassword)
		protected.GET("/me", authHandler.Me)
	}
}
