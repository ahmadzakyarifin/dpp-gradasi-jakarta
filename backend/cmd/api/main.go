package main

import (
	"log"

	"github.com/ahmadzakyarifin/dpp-gradasi/backend/config"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/app"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/infrastructure"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/validator"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load("../.env")
	_ = godotenv.Load(".env")

	cfg := config.MustLoad()

	// Fail-closed: CAPTCHA aktif wajib punya secret key — jangan biarkan
	// captcha "menyala" tanpa verifikasi sungguhan (bypass tersembunyi).
	if cfg.Security.CaptchaEnabled && cfg.Security.TurnstileSecretKey == "" {
		log.Fatal("CAPTCHA_ENABLED=true tetapi CAPTCHA_SECRET_KEY kosong. Set CAPTCHA_SECRET_KEY atau matikan CAPTCHA_ENABLED.")
	}

	validator.RegisterValidator()

	db, err := infrastructure.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("gagal melakukan koneksi ke database: %v", err)
	}

	app := app.NewApp(db, cfg)

	port := ":" + app.Cfg.App.Port
	log.Printf("server berjalan di %s", port)

	if err := app.Server.Run(port); err != nil {
		log.Fatalf("server crash: %v", err)
	}
}
