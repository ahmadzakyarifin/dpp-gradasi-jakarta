package infrastructure

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ahmadzakyarifin/dpp-gradasi/backend/config"
	"github.com/redis/go-redis/v9"
)

func ConnectRedis(cfg *config.Config) (*redis.Client, error) {
	addr := fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port)

	client := redis.NewClient(&redis.Options{
		Addr:            addr,
		Username:        cfg.Redis.Username,
		Password:        cfg.Redis.Pass,
		DB:              cfg.Redis.DB,
		PoolSize:        cfg.Redis.PoolSize,
		MinIdleConns:    cfg.Redis.MinIdleConns,
		MaxRetries:      cfg.Redis.MaxRetries,
		DialTimeout:     time.Duration(cfg.Redis.DialTimeoutSecs) * time.Second,
		ReadTimeout:     time.Duration(cfg.Redis.ReadTimeoutSecs) * time.Second,
		WriteTimeout:    time.Duration(cfg.Redis.WriteTimeoutSecs) * time.Second,
		PoolTimeout:     time.Duration(cfg.Redis.PoolTimeoutSecs) * time.Second,
		ConnMaxIdleTime: time.Duration(cfg.Redis.ConnMaxIdleTimeMins) * time.Minute,
		ConnMaxLifetime: time.Duration(cfg.Redis.ConnMaxLifetimeMins) * time.Minute,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("gagal konek redis: %w", err)
	}

	log.Println("berhasil terkoneksi ke redis")
	return client, nil
}
