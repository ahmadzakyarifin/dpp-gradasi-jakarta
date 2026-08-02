package middleware

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ahmadzakyarifin/schoolpay/backend/config"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/helper"
	"github.com/gin-gonic/gin"
)

const (
	CaptchaVerifiedKey = "captcha_verified"

	turnstileVerifyURL = "https://challenges.cloudflare.com/turnstile/v0/siteverify"
)

type captchaRequest struct {
	TurnstileToken      string `json:"turnstile_token"`
	TurnstileTokenCamel string `json:"turnstileToken"`
	CFTurnstileResponse string `json:"cf-turnstile-response"`
}

type turnstileResponse struct {
	Success    bool     `json:"success"`
	ErrorCodes []string `json:"error-codes"`
}

// RequiredCaptchaMiddleware digunakan untuk endpoint yang wajib CAPTCHA,
// misalnya login, forgot password, dan endpoint sensitif lainnya.
func RequiredCaptchaMiddleware(cfg *config.Config) gin.HandlerFunc {

	if cfg == nil || !cfg.Security.CaptchaEnabled {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	secretKey := strings.TrimSpace(cfg.Security.TurnstileSecretKey)

	return func(c *gin.Context) {

		if c.Request.Method == http.MethodOptions {
			c.Next()
			return
		}

		if secretKey == "" {
			helper.ErrorResponse(
				c,
				http.StatusInternalServerError,
				"SYSTEM_SECURITY_NOT_CONFIGURED",
				"Sistem keamanan belum dikonfigurasi.",
				nil,
			)
			c.Abort()
			return
		}

		var body captchaRequest

		_ = c.ShouldBindBodyWithJSON(&body)
		helper.RestoreBody(c)

		token := captchaTokenFromRequest(body)

		if token == "" {
			helper.ErrorResponse(
				c,
				http.StatusBadRequest,
				"AUTH_CAPTCHA_INVALID",
				"Verifikasi CAPTCHA wajib diselesaikan.",
				nil,
			)
			c.Abort()
			return
		}

		if err := verifyTurnstile(
			client,
			secretKey,
			token,
			c.ClientIP(),
		); err != nil {

			helper.ErrorResponse(
				c,
				http.StatusBadRequest,
				"AUTH_CAPTCHA_INVALID",
				"Verifikasi CAPTCHA tidak valid.",
				nil,
			)

			c.Abort()
			return
		}

		c.Set(CaptchaVerifiedKey, true)

		c.Next()
	}
}

func captchaTokenFromRequest(body captchaRequest) string {

	if token := strings.TrimSpace(body.TurnstileToken); token != "" {
		return token
	}

	if token := strings.TrimSpace(body.TurnstileTokenCamel); token != "" {
		return token
	}

	if token := strings.TrimSpace(body.CFTurnstileResponse); token != "" {
		return token
	}

	return ""
}

func verifyTurnstile(
	client *http.Client,
	secretKey string,
	token string,
	ip string,
) error {

	form := url.Values{}
	form.Set("secret", secretKey)
	form.Set("response", token)

	if strings.TrimSpace(ip) != "" {
		form.Set("remoteip", ip)
	}

	res, err := client.PostForm(
		turnstileVerifyURL,
		form,
	)

	if err != nil {
		return fmt.Errorf("turnstile request failed: %w", err)
	}

	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		return fmt.Errorf("turnstile http status: %d", res.StatusCode)
	}

	var result turnstileResponse

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return fmt.Errorf("turnstile decode failed: %w", err)
	}

	if !result.Success {
		log.Printf("Turnstile verification failed: %v", result.ErrorCodes)
		return fmt.Errorf("turnstile verification failed")
	}

	return nil
}
