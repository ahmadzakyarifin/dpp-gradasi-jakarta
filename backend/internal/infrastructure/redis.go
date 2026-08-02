package infrastructure

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ahmadzakyarifin/dpp-gradasi/backend/config"
	"github.com/redis/go-redis/v9"
)

// Tuning koneksi Redis — konstanta (bukan env). Nilai aman untuk Upstash & lokal.
const (
	redisPoolSize        = 20
	redisMinIdleConns    = 5
	redisMaxRetries      = 3
	redisDialTimeout     = 5 * time.Second
	redisReadTimeout     = 3 * time.Second
	redisWriteTimeout    = 3 * time.Second
	redisPoolTimeout     = 4 * time.Second
	redisConnMaxIdleTime = 5 * time.Minute
	redisConnMaxLifetime = 30 * time.Minute
)

func ConnectRedis(cfg *config.Config) (*redis.Client, error) {
	addr := fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port)

	options := &redis.Options{
		Addr:            addr,
		Username:        cfg.Redis.Username,
		Password:        cfg.Redis.Pass,
		DB:              cfg.Redis.DB,
		PoolSize:        redisPoolSize,
		MinIdleConns:    redisMinIdleConns,
		MaxRetries:      redisMaxRetries,
		DialTimeout:     redisDialTimeout,
		ReadTimeout:     redisReadTimeout,
		WriteTimeout:    redisWriteTimeout,
		PoolTimeout:     redisPoolTimeout,
		ConnMaxIdleTime: redisConnMaxIdleTime,
		ConnMaxLifetime: redisConnMaxLifetime,
	}

	// Redis managed (Upstash, Redis Cloud, dll) wajib TLS. Auto-detect:
	// host non-lokal => aktifkan TLS. Host lokal (127.0.0.1/localhost) => plain.
	host := strings.ToLower(cfg.Redis.Host)
	if host != "127.0.0.1" && host != "localhost" && !strings.HasSuffix(host, ".local") {
		options.TLSConfig = &tls.Config{MinVersion: tls.VersionTLS12}
	}

	client := redis.NewClient(options)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("gagal konek redis: %w", err)
	}

	log.Println("berhasil terkoneksi ke redis")
	return client, nil
}
