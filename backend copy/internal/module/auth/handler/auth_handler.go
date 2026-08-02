package handler

import (
	"net/http"
	"time"

	"github.com/ahmadzakyarifin/schoolpay/backend/config"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/helper"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/auth/dto"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/auth/service"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/validator"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	s   service.AuthService
	cfg *config.Config
}

func NewAuthHandler(authService service.AuthService, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		s:   authService,
		cfg: cfg,
	}
}

func (h *AuthHandler) Cfg() *config.Config {
	return h.cfg
}

func (h *AuthHandler) secureCookie() bool {
	return h.cfg != nil && h.cfg.App.Env == "production"
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "Validasi gagal.", validator.Errors(err))
		return
	}

	userAgent := c.GetHeader("User-Agent")
	ipAddress := c.ClientIP()

	res, err := h.s.Login(c.Request.Context(), req, ipAddress, userAgent, "")
	if err != nil {
		helper.HandleServiceError(c, err)
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("refresh_token", res.RefreshToken, int(time.Until(res.RefreshTokenExpiry).Seconds()), "/", "", h.secureCookie(), true)

	helper.SuccessResponse(c, http.StatusOK, "AUTH_LOGIN_SUCCESS", "Login berhasil.", gin.H{
		"access_token": res.AccessToken,
		"user":         res.User,
	}, nil)
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	cookieToken, err := c.Cookie("refresh_token")
	if err != nil {
		helper.ErrorResponse(c, http.StatusUnauthorized, "AUTH_SESSION_EXPIRED", "Sesi telah berakhir. Silakan login kembali.", nil)
		return
	}

	res, err := h.s.RefreshToken(c.Request.Context(), cookieToken)
	if err != nil {
		helper.HandleServiceError(c, err)
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("refresh_token", res.RefreshToken, int(time.Until(res.RefreshTokenExpiry).Seconds()), "/", "", h.secureCookie(), true)

	helper.SuccessResponse(c, http.StatusOK, "AUTH_REFRESH_SUCCESS", "Sesi berhasil diperbarui.", gin.H{
		"access_token": res.AccessToken,
	}, nil)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	cookieToken, err := c.Cookie("refresh_token")
	if err == nil {
		_ = h.s.Logout(c.Request.Context(), cookieToken)
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("refresh_token", "", -1, "/", "", h.secureCookie(), true)

	helper.SuccessResponse(c, http.StatusOK, "AUTH_LOGOUT_SUCCESS", "Logout berhasil.", nil, nil)
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req dto.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "validasi gagal", validator.Errors(err))
		return
	}

	ipAddress := c.ClientIP()
	userAgent := c.Request.UserAgent()
	link, err := h.s.ForgotPassword(c.Request.Context(), req, ipAddress, userAgent)
	if err != nil {
		helper.HandleServiceError(c, err)
		return
	}

	var data map[string]any
	if h.cfg.App.Env != "production" && link != "" {
		data = map[string]any{
			"debug_link": link,
		}
	}

	helper.SuccessResponse(c, http.StatusOK, "AUTH_FORGOT_PASSWORD_SUCCESS", "Jika email terdaftar, tautan reset password akan dikirim ke email tersebut.", data, nil)
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req dto.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "validasi gagal", validator.Errors(err))
		return
	}

	ipAddress := c.ClientIP()
	userAgent := c.Request.UserAgent()

	if err := h.s.ResetPassword(c.Request.Context(), req, ipAddress, userAgent); err != nil {
		helper.HandleServiceError(c, err)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "AUTH_RESET_PASSWORD_SUCCESS", "Password berhasil direset.", nil, nil)
}

func (h *AuthHandler) ValidateResetToken(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		helper.ErrorResponse(c, http.StatusBadRequest, "AUTH_TOKEN_INVALID_OR_EXPIRED", "Token reset password tidak valid atau telah kedaluwarsa.", nil)
		return
	}

	if err := h.s.ValidateResetToken(c.Request.Context(), token); err != nil {
		helper.HandleServiceError(c, err)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "AUTH_RESET_TOKEN_VALID", "Token reset valid.", nil, nil)
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusUnprocessableEntity, "VALIDATION_ERROR", "Validasi gagal.", validator.Errors(err))
		return
	}

	value, exists := c.Get("user_id")
	if !exists {
		helper.ErrorResponse(c, http.StatusUnauthorized, "AUTH_TOKEN_INVALID_OR_EXPIRED", "Sesi login telah berakhir. Silakan login kembali.", nil)
		return
	}

	userID, ok := value.(uint)
	if !ok {
		helper.ErrorResponse(c, http.StatusUnauthorized, "AUTH_TOKEN_INVALID_OR_EXPIRED", "Sesi login telah berakhir. Silakan login kembali.", nil)
		return
	}

	if err := h.s.ChangePassword(c.Request.Context(), userID, req); err != nil {
		helper.HandleServiceError(c, err)
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("refresh_token", "", -1, "/", "", h.secureCookie(), true)

	helper.SuccessResponse(c, http.StatusOK, "AUTH_CHANGE_PASSWORD_SUCCESS", "Password berhasil diperbarui, silakan login ulang.", nil, nil)
}
