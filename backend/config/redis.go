package config

type RedisConfig struct {
	Host                string `env:"REDIS_HOST" envDefault:"127.0.0.1"`
	Port                string `env:"REDIS_PORT" envDefault:"6379"`
	Username            string `env:"REDIS_USERNAME" envDefault:""`
	Pass                string `env:"REDIS_PASS"`
	DB                  int    `env:"REDIS_DB" envDefault:"0"`
	PoolSize            int    `env:"REDIS_POOL_SIZE" envDefault:"20"`
	MinIdleConns        int    `env:"REDIS_MIN_IDLE_CONNS" envDefault:"5"`
	MaxRetries          int    `env:"REDIS_MAX_RETRIES" envDefault:"3"`
	DialTimeoutSecs     int    `env:"REDIS_DIAL_TIMEOUT_SECS" envDefault:"5"`
	ReadTimeoutSecs     int    `env:"REDIS_READ_TIMEOUT_SECS" envDefault:"3"`
	WriteTimeoutSecs    int    `env:"REDIS_WRITE_TIMEOUT_SECS" envDefault:"3"`
	PoolTimeoutSecs     int    `env:"REDIS_POOL_TIMEOUT_SECS" envDefault:"4"`
	ConnMaxIdleTimeMins int    `env:"REDIS_CONN_MAX_IDLE_TIME_MINS" envDefault:"5"`
	ConnMaxLifetimeMins int    `env:"REDIS_CONN_MAX_LIFETIME_MINS" envDefault:"30"`
}
