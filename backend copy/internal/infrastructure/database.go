package infrastructure

import (
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/ahmadzakyarifin/schoolpay/backend/config"
	mysqldriver "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ConnectDB membuka koneksi ke MariaDB dan mengembalikan *gorm.DB yang siap pakai.
func ConnectDB(cfg *config.Config) (*gorm.DB, error) {
	dsn, err := databaseDSN(cfg)
	if err != nil {
		return nil, err
	}

	gormConfig := &gorm.Config{
		Logger: configureQueryLogger(cfg),
	}

	db, err := gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("gagal membuka koneksi database: %w", err)
	}

	// Atur connection pool lewat sql.DB di bawah GORM.
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil sql.DB: %w", err)
	}

	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.Database.ConnMaxLifetimeMins) * time.Minute)
	sqlDB.SetConnMaxIdleTime(time.Duration(cfg.Database.ConnMaxIdleTimeMins) * time.Minute)

	if err := sqlDB.Ping(); err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("database tidak merespons: %w", err)
	}

	log.Println("Database berhasil terkoneksi")

	return db, nil
}

// databaseDSN menyusun DSN MariaDB/MySQL dari config (gaya sama seperti bun).
func databaseDSN(cfg *config.Config) (string, error) {
	locationName := strings.TrimSpace(cfg.App.Timezone)

	location, err := time.LoadLocation(locationName)
	if err != nil {
		return "", fmt.Errorf("timezone %s tidak valid: %w", locationName, err)
	}

	mysqlConfig := mysqldriver.Config{
		User:   cfg.Database.User,
		Passwd: cfg.Database.Pass,

		Net:  "tcp",
		Addr: net.JoinHostPort(cfg.Database.Host, cfg.Database.Port),

		DBName: cfg.Database.Name,

		ParseTime: true,
		Loc:       location,

		Collation: "utf8mb4_unicode_ci",

		AllowNativePasswords: true,

		Timeout:      5 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,

		Params: map[string]string{
			"charset": "utf8mb4",
		},
	}

	return mysqlConfig.FormatDSN(), nil
}

// configureQueryLogger mengembalikan GORM logger.
// Di development: log semua SQL (verbose). Di lainnya: silent.
func configureQueryLogger(cfg *config.Config) logger.Interface {
	if strings.ToLower(cfg.App.Env) != "development" {
		return logger.Default.LogMode(logger.Silent)
	}

	return logger.Default.LogMode(logger.Info)
}
