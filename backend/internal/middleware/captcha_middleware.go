package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/ahmadzakyarifin/dpp-gradasi/backend/config"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/helper"
	"github.com/gin-gonic/gin"
)

// RequiredCaptchaMiddleware memverifikasi Cloudflare Turnstile CAPTCHA.
// Nonaktif otomatis jika CAPTCHA_ENABLED=false (mode dev).
// Token diambil dari body JSON: captcha_token / turnstile_token / turnstileToken
// atau header: cf-turnstile-response.
func RequiredCaptchaMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		if cfg == nil || !cfg.Security.CaptchaEnabled {
			c.Next()
			return
		}

		token := captchaTokenFromRequest(c)
		if token == "" {
			helper.ErrorResponse(c, http.StatusBadRequest, "AUTH_CAPTCHA_REQUIRED", "Silakan selesaikan verifikasi CAPTCHA.", nil)
			c.Abort()
			return
		}

		ok, err := helper.VerifyTurnstile(
			c.Request.Context(),
			cfg.Security.TurnstileSecretKey,
			token,
			c.ClientIP(),
			cfg.Security.TurnstileVerifyURL,
		)
		if err != nil || !ok {
			helper.ErrorResponse(c, http.StatusBadRequest, "AUTH_CAPTCHA_INVALID", "Verifikasi CAPTCHA gagal. Silakan coba lagi.", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}

// captchaTokenFromRequest mencari token Turnstile di body JSON atau header.
func captchaTokenFromRequest(c *gin.Context) string {
	if v := c.GetHeader("cf-turnstile-response"); v != "" {
		return v
	}

	// Baca body sekali, restore agar handler bisa baca lagi.
	body, err := c.GetRawData()
	if err != nil || len(body) == 0 {
		return ""
	}
	c.Request.Body = io.NopCloser(bytes.NewReader(body))

	var payload struct {
		CaptchaToken    string `json:"captcha_token"`
		TurnstileToken  string `json:"turnstile_token"`
		TurnstileToken2 string `json:"turnstileToken"`
	}
	_ = json.Unmarshal(body, &payload)

	if payload.CaptchaToken != "" {
		return payload.CaptchaToken
	}
	if payload.TurnstileToken != "" {
		return payload.TurnstileToken
	}
	return payload.TurnstileToken2
}
