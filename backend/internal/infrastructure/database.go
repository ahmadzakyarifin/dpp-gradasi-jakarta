package infrastructure

import (
	"fmt"
	"log"
	"time"

	"github.com/ahmadzakyarifin/dpp-gradasi/backend/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Tuning koneksi DB — konstanta (bukan env). Nilai aman untuk dev & prod kecil.
const (
	dbMaxOpenConns        = 100
	dbMaxIdleConns        = 20
	dbConnMaxLifetimeMins = 30
	dbConnMaxIdleTimeMins = 10
)

func ConnectDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.User,
		cfg.Database.Pass,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	)

	logLevel := logger.Silent
	if cfg.App.Env == "development" {
		logLevel = logger.Info
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("gagal konek database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("gagal get sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(dbMaxOpenConns)
	sqlDB.SetMaxIdleConns(dbMaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(dbConnMaxLifetimeMins) * time.Minute)
	sqlDB.SetConnMaxIdleTime(time.Duration(dbConnMaxIdleTimeMins) * time.Minute)

	log.Println("berhasil terkoneksi ke database")
	return db, nil
}
