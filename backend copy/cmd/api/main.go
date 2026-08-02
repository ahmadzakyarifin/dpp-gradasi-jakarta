package main

import (
	"log"

	"github.com/ahmadzakyarifin/schoolpay/backend/config"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/app"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/infrastructure"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/validator"
)

func main() {
	cfg := config.MustLoad()

	validator.RegisterValidator()

	db, err := infrastructure.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("gagal melakukan koneksi ke database: %v", err)
	}

	redisClient, err := infrastructure.ConnectRedis(cfg)
	if err != nil {
		log.Fatalf("gagal melakukan koneksi ke redis: %v", err)
	}
	log.Println("berhasil terkoneksi ke redis")
	defer redisClient.Close()

	app := app.NewApp(db, cfg, redisClient)

	port := ":" + app.Cfg.App.Port
	log.Printf("server berjalan di %s", port)

	if err := app.Server.Run(port); err != nil {
		log.Fatalf("server crash: %v", err)
	}
}
