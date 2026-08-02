package config

import (
	"errors"
	"log"

	"github.com/caarlos0/env/v11"
)

// Config merangkum semua konfigurasi dari file-file terpisah
type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Worker   WorkerConfig
	JWT      JWTConfig
	Cookie   CookieConfig
	Security SecurityConfig
	Upload   UploadConfig
	SMTP     SMTPConfig
	WAHA     WAHAConfig
	Midtrans MidtransConfig
	Dev      DevConfig
}

// MustLoad memuat konfigurasi dan akan menghentikan aplikasi (fatal) jika gagal
func MustLoad() *Config {
	cfg := &Config{}

	// env.Parse membaca environment variables berdasarkan struct tag
	if err := env.Parse(cfg); err != nil {
		log.Fatalf("KRITIS: gagal memuat konfigurasi: %v", err)
	}

	// Menjalankan validasi setelah konfigurasi dimuat
	if err := cfg.Validate(); err != nil {
		log.Fatalf("KRITIS: VALIDASI KONFIGURASI GAGAL: %v", err)
	}

	return cfg
}

// Validate memastikan semua pengaturan wajib telah diisi, khususnya untuk production
func (c *Config) Validate() error {
	if c.App.Env == "production" {
		if c.JWT.Secret == "" || len(c.JWT.Secret) < 32 {
			return errors.New("JWT_SECRET harus diisi dan minimal panjangnya 32 karakter di lingkungan production")
		}
		if c.WAHA.WebhookSecret == "" {
			return errors.New("WAHA_WEBHOOK_SECRET harus diisi di lingkungan production untuk keamanan webhook")
		}
		if c.Security.CaptchaEnabled && c.Security.TurnstileSecretKey == "" {
			return errors.New("TURNSTILE_SECRET_KEY harus diisi jika CAPTCHA diaktifkan di lingkungan production")
		}
		if c.Database.Pass == "" {
			return errors.New("DB_PASS tidak boleh kosong di lingkungan production")
		}
	}
	return nil
}
