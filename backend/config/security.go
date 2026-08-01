package config

// SecurityConfig berisi konfigurasi keamanan, termasuk Cloudflare Turnstile.
type SecurityConfig struct {
	CaptchaEnabled bool `env:"CAPTCHA_ENABLED" envDefault:"false"`
	// SiteKey publik (aman untuk frontend). Nama env mengikuti konvensi
	// CAPTCHA_* (bukan TURNSTILE_*) — konsisten dengan .env.
	TurnstileSiteKey string `env:"CAPTCHA_SITE_KEY"`
	// SecretKey rahasia, hanya untuk verifikasi server-side.
	TurnstileSecretKey string `env:"CAPTCHA_SECRET_KEY"`
	// BaseURL verifikasi Turnstile (dev: http://localhost:3000, prod: https://challenges.cloudflare.com).
	TurnstileVerifyURL string `env:"CAPTCHA_VERIFY_URL"`
}
