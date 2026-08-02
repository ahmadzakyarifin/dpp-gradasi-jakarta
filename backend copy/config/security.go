package config

type SecurityConfig struct {
	CaptchaEnabled     bool   `env:"CAPTCHA_ENABLED" envDefault:"false"`
	TurnstileSiteKey   string `env:"TURNSTILE_SITE_KEY"`
	TurnstileSecretKey string `env:"TURNSTILE_SECRET_KEY"`
}
