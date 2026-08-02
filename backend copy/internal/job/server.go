package jobs

import (
	"github.com/ahmadzakyarifin/schoolpay/backend/config"
	"github.com/hibiken/asynq"
)

func NewServer(cfg *config.Config) *asynq.Server {
	return asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     cfg.Redis.Host + ":" + cfg.Redis.Port,
			Password: cfg.Redis.Pass,
			DB:       cfg.Redis.DB,
		},
		asynq.Config{
			Concurrency: cfg.Worker.Concurrency,

			Queues: map[string]int{
				QueueCritical: 10,
				QueueEmail:    8,
				QueueWhatsApp: 7,
				QueueDefault:  5,
				QueueLow:      1,
			},
		},
	)
}
