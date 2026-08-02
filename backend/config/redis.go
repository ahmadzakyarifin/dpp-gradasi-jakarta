package config

type RedisConfig struct {
	Host     string `env:"REDIS_HOST" envDefault:"127.0.0.1"`
	Port     string `env:"REDIS_PORT" envDefault:"6379"`
	Username string `env:"REDIS_USERNAME" envDefault:""`
	Pass     string `env:"REDIS_PASS"`
	DB       int    `env:"REDIS_DB" envDefault:"0"`
}
