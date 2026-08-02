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

	// Captcha hanya aktif jika CAPTCHA_ENABLED=true + secret key terisi.
	needCaptcha := middleware.RequiredCaptchaMiddleware(authH.Cfg())

	// Rate limit: login & forgot-password
	loginLimit := middleware.RateLimitRules("auth",
		middleware.IPEmail(5, time.Minute),
	)
	forgotLimit := middleware.RateLimitRules("forgot_password",
		middleware.IPEmail(3, time.Minute),
	)

	auth := api.Group("/auth")
	{
		auth.POST("/login", needCaptcha, loginLimit, authH.Login)
		auth.POST("/refresh", authH.Refresh)
		auth.POST("/logout", requireAuth, authH.Logout)
		auth.POST("/forgot-password", needCaptcha, forgotLimit, authH.ForgotPassword)
		auth.POST("/reset-password", authH.ResetPassword)
		auth.GET("/validate-reset-token", authH.ValidateResetToken)
		auth.POST("/change-password", requireAuth, authH.ChangePassword)
	}
}
