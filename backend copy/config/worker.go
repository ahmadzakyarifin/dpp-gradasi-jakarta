package config

type WorkerConfig struct {
	Concurrency         int    `env:"WORKER_CONCURRENCY" envDefault:"10"`
	MaxRetry            int    `env:"WORKER_MAX_RETRY" envDefault:"10"`
	TimeoutSecs         int    `env:"WORKER_TIMEOUT_SECS" envDefault:"60"`
	ShutdownTimeoutSecs int    `env:"WORKER_SHUTDOWN_TIMEOUT_SECS" envDefault:"30"`
	StrictPriority      bool   `env:"WORKER_STRICT_PRIORITY" envDefault:"true"`
	RetentionHours      int    `env:"WORKER_RETENTION_HOURS" envDefault:"168"`
	QueueCritical       string `env:"WORKER_QUEUE_CRITICAL" envDefault:"critical"`
	QueueDefault        string `env:"WORKER_QUEUE_DEFAULT" envDefault:"default"`
	QueueEmail          string `env:"WORKER_QUEUE_EMAIL" envDefault:"email"`
	QueueWhatsapp       string `env:"WORKER_QUEUE_WHATSAPP" envDefault:"whatsapp"`
	QueueNotification   string `env:"WORKER_QUEUE_NOTIFICATION" envDefault:"notification"`
	QueueLow            string `env:"WORKER_QUEUE_LOW" envDefault:"low"`
}
