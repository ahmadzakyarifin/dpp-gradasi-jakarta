package auth

import (
	"time"

	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/middleware"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/auth/handler"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	api *gin.RouterGroup,
	authH *handler.AuthHandler,
	jwtSecret string,
) {
	requireAuth := middleware.AuthMiddleware(jwtSecret)

	// Rate limit: login & forgot-password
	loginLimit := middleware.RateLimitRules("auth",
		middleware.IPEmail(5, time.Minute),
	)
	forgotLimit := middleware.RateLimitRules("forgot_password",
		middleware.IPEmail(3, time.Minute),
	)

	// Rate limit endpoint token sensitif (reset & aktivasi) — per IP,
	// batas ketat karena token sekali pakai rawan brute-force/enumerasi.
	resetSubmitLimit := middleware.RateLimitRules("reset_password_submit",
		middleware.IP(10, 15*time.Minute),
	)
	validateResetLimit := middleware.RateLimitRules("validate_reset_token",
		middleware.IP(30, 15*time.Minute),
	)
	activateLimit := middleware.RateLimitRules("activate_account",
		middleware.IP(10, 15*time.Minute),
	)
	validateActivateLimit := middleware.RateLimitRules("validate_activation_token",
		middleware.IP(30, 15*time.Minute),
	)

	auth := api.Group("/auth")
	{
		auth.POST("/login", loginLimit, authH.Login)
		auth.POST("/refresh", authH.Refresh)
		auth.POST("/logout", requireAuth, authH.Logout)
		auth.POST("/forgot-password", forgotLimit, authH.ForgotPassword)
		auth.POST("/reset-password", resetSubmitLimit, authH.ResetPassword)
		auth.GET("/validate-reset-token", validateResetLimit, authH.ValidateResetToken)
		auth.POST("/change-password", requireAuth, authH.ChangePassword)
		auth.GET("/me", requireAuth, authH.Me)
		auth.POST("/activate-account", activateLimit, authH.ActivateAccount)
		auth.GET("/validate-activation-token", validateActivateLimit, authH.ValidateActivationToken)
	}
}
