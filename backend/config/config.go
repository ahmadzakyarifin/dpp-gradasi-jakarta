package config

import (
	"errors"
	"log"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Cookie   CookieConfig
	Security SecurityConfig
	SMTP     SMTPConfig
	Dev      DevConfig
}

func MustLoad() *Config {
	cfg := &Config{}

	if err := env.Parse(cfg); err != nil {
		log.Fatalf("KRITIS: gagal memuat konfigurasi: %v", err)
	}

	if err := cfg.Validate(); err != nil {
		log.Fatalf("KRITIS: VALIDASI KONFIGURASI GAGAL: %v", err)
	}

	return cfg
}

func (c *Config) Validate() error {
	if c.App.Env == "production" {
		if c.JWT.Secret == "" || len(c.JWT.Secret) < 32 {
			return errors.New("JWT_SECRET harus diisi dan minimal panjangnya 32 karakter di lingkungan production")
		}
		if c.Database.Pass == "" {
			return errors.New("DB_PASS tidak boleh kosong di lingkungan production")
		}
	}
	return nil
}
