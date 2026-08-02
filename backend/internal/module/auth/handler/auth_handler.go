package handler

import (
	"context"
	"net/http"

	"github.com/ahmadzakyarifin/dpp-gradasi/backend/config"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/helper"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/middleware"
	activitylogdto "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/activitylog/dto"
	activitylogservice "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/activitylog/service"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/auth/dto"
	authservice "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/auth/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	svc    authservice.AuthService
	cfg    *config.Config
	logSvc activitylogservice.ActivityLogService
}

func NewAuthHandler(svc authservice.AuthService, cfg *config.Config, logSvc activitylogservice.ActivityLogService) *AuthHandler {
	return &AuthHandler{svc: svc, cfg: cfg, logSvc: logSvc}
}

// Login menangani login user
// POST /api/v1/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ValidationErrorResponse(c, []helper.ValidationErrorItem{
			{Field: "email", Tag: "required", Message: "Email dan password wajib diisi."},
		})
		return
	}

	// Captcha Verification if Enabled (Cloudflare Turnstile)
	if h.cfg != nil && h.cfg.Security.CaptchaEnabled {
		if req.CaptchaToken == "" {
			helper.ErrorResponse(c, http.StatusBadRequest, "AUTH_CAPTCHA_REQUIRED", "Silakan selesaikan verifikasi CAPTCHA.", nil)
			return
		}
		ok, err := helper.VerifyTurnstile(
			c.Request.Context(),
			h.cfg.Security.TurnstileSecretKey,
			req.CaptchaToken,
			c.ClientIP(),
			h.cfg.Security.TurnstileVerifyURL,
		)
		if err != nil || !ok {
			helper.ErrorResponse(c, http.StatusBadRequest, "AUTH_CAPTCHA_INVALID", "Verifikasi CAPTCHA gagal. Silakan coba lagi.", nil)
			return
		}
	}

	// Audit context (untuk SYNC audit di service)
	req.IPAddress = c.ClientIP()
	req.UserAgent = c.Request.UserAgent()

	resp, refreshToken, maxAge, svcErr := h.svc.Login(&req)
	if svcErr != nil {
		// Log failed login (async tetap OK untuk event gagal — bukan audit sukses)
		ip := c.ClientIP()
		ua := c.Request.UserAgent()
		go h.logSvc.Log(context.Background(), nil, &activitylogdto.ActivityLogInput{
			Action:      "auth.login_failed",
			EntityType:  "auth",
			Description: "Gagal login dengan email: " + req.Email,
			IPAddress:   ip,
			UserAgent:   ua,
		})

		helper.HandleServiceError(c, svcErr)
		return
	}

	// Set refresh token cookie
	c.SetCookie("refresh_token", refreshToken,
		maxAge,
		h.cfg.Cookie.Path,
		h.cfg.Cookie.Domain,
		h.cfg.Cookie.Secure,
		h.cfg.Cookie.HTTPOnly,
	)

	helper.SuccessResponse(c, http.StatusOK, "AUTH_LOGIN_SUCCESS", "Login berhasil.", resp, nil)
}

// Refresh memperbarui access token
// POST /api/v1/auth/refresh
func (h *AuthHandler) Refresh(c *gin.Context) {
	refreshTokenStr, err := c.Cookie("refresh_token")
	if err != nil {
		refreshTokenStr = ""
	}

	resp, newRefreshToken, maxAge, svcErr := h.svc.Refresh(refreshTokenStr)
	if svcErr != nil {
		helper.HandleServiceError(c, svcErr)
		return
	}

	c.SetCookie("refresh_token", newRefreshToken,
		maxAge,
		h.cfg.Cookie.Path,
		h.cfg.Cookie.Domain,
		h.cfg.Cookie.Secure,
		h.cfg.Cookie.HTTPOnly,
	)

	helper.SuccessResponse(c, http.StatusOK, "AUTH_REFRESH_SUCCESS", "Sesi berhasil diperbarui.", resp, nil)
}

// Logout menghapus sesi
// POST /api/v1/auth/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	refreshTokenStr, _ := c.Cookie("refresh_token")

	_ = h.svc.Logout(userID, refreshTokenStr)

	// Hapus cookie
	c.SetCookie("refresh_token", "",
		-1,
		h.cfg.Cookie.Path,
		h.cfg.Cookie.Domain,
		h.cfg.Cookie.Secure,
		h.cfg.Cookie.HTTPOnly,
	)

	helper.SuccessResponse(c, http.StatusOK, "AUTH_LOGOUT_SUCCESS", "Logout berhasil.", nil, nil)

	actorID, actorName, actorRole, ip, ua := helper.GetAuditMeta(c.Request.Context())
	go h.logSvc.Log(context.Background(), nil, &activitylogdto.ActivityLogInput{
		ActorID:     &actorID,
		ActorName:   actorName,
		ActorRole:   actorRole,
		Action:      "auth.logout",
		EntityType:  "auth",
		Description: "Logout berhasil",
		IPAddress:   ip,
		UserAgent:   ua,
	})
}

// ForgotPassword mengirim email reset password
// POST /api/v1/auth/forgot-password
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req dto.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ValidationErrorResponse(c, []helper.ValidationErrorItem{
			{Field: "email", Tag: "required", Message: "Email wajib diisi."},
		})
		return
	}

	// Captcha Verification if Enabled (Cloudflare Turnstile)
	if h.cfg != nil && h.cfg.Security.CaptchaEnabled {
		if req.CaptchaToken == "" {
			helper.ErrorResponse(c, http.StatusBadRequest, "AUTH_CAPTCHA_REQUIRED", "Silakan selesaikan verifikasi CAPTCHA.", nil)
			return
		}
		ok, err := helper.VerifyTurnstile(
			c.Request.Context(),
			h.cfg.Security.TurnstileSecretKey,
			req.CaptchaToken,
			c.ClientIP(),
			h.cfg.Security.TurnstileVerifyURL,
		)
		if err != nil || !ok {
			helper.ErrorResponse(c, http.StatusBadRequest, "AUTH_CAPTCHA_INVALID", "Verifikasi CAPTCHA gagal. Silakan coba lagi.", nil)
			return
		}
	}

	if err := h.svc.ForgotPassword(&req, h.cfg.App.URL); err != nil {
		helper.HandleServiceError(c, err)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "AUTH_FORGOT_PASSWORD_SUCCESS",
		"Jika email terdaftar, tautan reset password akan dikirim ke email tersebut.", nil, nil)

	ip := c.ClientIP()
	ua := c.Request.UserAgent()
	go h.logSvc.Log(context.Background(), nil, &activitylogdto.ActivityLogInput{
		Action:      "auth.forgot_password",
		EntityType:  "auth",
		Description: "Request reset password untuk: " + req.Email,
		IPAddress:   ip,
		UserAgent:   ua,
	})
}

// ValidateResetToken memvalidasi token reset password
// GET /api/v1/auth/validate-reset-token?token=xxx
func (h *AuthHandler) ValidateResetToken(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		helper.ErrorResponse(c, http.StatusBadRequest, "AUTH_TOKEN_INVALID_OR_EXPIRED",
			"Token reset password tidak valid atau telah kedaluwarsa.", nil)
		return
	}

	if err := h.svc.ValidateResetToken(token); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "AUTH_TOKEN_INVALID_OR_EXPIRED",
			"Token reset password tidak valid atau telah kedaluwarsa.", nil)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "AUTH_RESET_TOKEN_VALID", "Token reset valid.", nil, nil)
}

// ResetPassword mereset password dengan token
// POST /api/v1/auth/reset-password
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req dto.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ValidationErrorResponse(c, []helper.ValidationErrorItem{
			{Field: "password", Tag: "required", Message: "Password dan konfirmasi wajib diisi."},
		})
		return
	}

	if err := h.svc.ResetPassword(&req); err != nil {
		helper.HandleServiceError(c, err)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "AUTH_RESET_PASSWORD_SUCCESS", "Password berhasil direset. Silakan login.", nil, nil)

	ip := c.ClientIP()
	ua := c.Request.UserAgent()
	go h.logSvc.Log(context.Background(), nil, &activitylogdto.ActivityLogInput{
		Action:      "auth.reset_password",
		EntityType:  "auth",
		Description: "Password berhasil direset",
		IPAddress:   ip,
		UserAgent:   ua,
	})
}

// ValidateActivationToken memvalidasi token aktivasi akun
// GET /api/v1/auth/validate-activation-token?token=xxx
func (h *AuthHandler) ValidateActivationToken(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		helper.ErrorResponse(c, http.StatusBadRequest, "AUTH_TOKEN_INVALID_OR_EXPIRED",
			"Token aktivasi tidak valid atau telah kedaluwarsa.", nil)
		return
	}

	if err := h.svc.ValidateActivationToken(token); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "AUTH_TOKEN_INVALID_OR_EXPIRED",
			"Token aktivasi tidak valid atau telah kedaluwarsa.", nil)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "AUTH_ACTIVATION_TOKEN_VALID", "Token aktivasi valid.", nil, nil)
}

// ActivateAccount mengaktifkan akun dan membuat password pertama
// POST /api/v1/auth/activate-account
func (h *AuthHandler) ActivateAccount(c *gin.Context) {
	var req dto.ActivateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ValidationErrorResponse(c, []helper.ValidationErrorItem{
			{Field: "password", Tag: "required", Message: "Password dan konfirmasi password wajib diisi (minimal 6 karakter)."},
		})
		return
	}

	if err := h.svc.ActivateAccount(&req); err != nil {
		helper.HandleServiceError(c, err)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "AUTH_ACCOUNT_ACTIVATED", "Akun berhasil diaktifkan. Silakan login.", nil, nil)

	ip := c.ClientIP()
	ua := c.Request.UserAgent()
	go h.logSvc.Log(context.Background(), nil, &activitylogdto.ActivityLogInput{
		Action:      "auth.activate_account",
		EntityType:  "auth",
		Description: "Akun berhasil diaktifkan dengan password baru",
		IPAddress:   ip,
		UserAgent:   ua,
	})
}

// ChangePassword mengganti password saat sudah login
// POST /api/v1/auth/change-password
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ValidationErrorResponse(c, []helper.ValidationErrorItem{
			{Field: "current_password", Tag: "required", Message: "Password saat ini wajib diisi."},
		})
		return
	}

	if err := h.svc.ChangePassword(userID, &req); err != nil {
		helper.HandleServiceError(c, err)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "AUTH_CHANGE_PASSWORD_SUCCESS",
		"Password berhasil diperbarui, silakan login ulang.", nil, nil)

	actorID, actorName, actorRole, ip, ua := helper.GetAuditMeta(c.Request.Context())
	go h.logSvc.Log(context.Background(), nil, &activitylogdto.ActivityLogInput{
		ActorID:     &actorID,
		ActorName:   actorName,
		ActorRole:   actorRole,
		Action:      "user.change_password",
		EntityType:  "user",
		EntityID:    &userID,
		Description: "Admin mengganti password",
		IPAddress:   ip,
		UserAgent:   ua,
	})
}

// Me mengambil profile user yang sedang login
// GET /api/v1/auth/me
func (h *AuthHandler) Me(c *gin.Context) {
	userID, _ := middleware.GetUserID(c)

	profile, err := h.svc.GetProfile(userID)
	if err != nil {
		helper.HandleServiceError(c, err)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "AUTH_PROFILE", "Profil berhasil diambil.", profile, nil)
}
