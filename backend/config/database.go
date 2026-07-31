package config

type DatabaseConfig struct {
	Driver              string `env:"DB_DRIVER" envDefault:"mysql"`
	Host                string `env:"DB_HOST" envDefault:"127.0.0.1"`
	Port                string `env:"DB_PORT" envDefault:"3306"`
	Name                string `env:"DB_NAME" envDefault:"dpp_gradasi"`
	User                string `env:"DB_USER"`
	Pass                string `env:"DB_PASS"`
	MaxOpenConns        int    `env:"DB_MAX_OPEN_CONNS" envDefault:"100"`
	MaxIdleConns        int    `env:"DB_MAX_IDLE_CONNS" envDefault:"20"`
	ConnMaxLifetimeMins int    `env:"DB_CONN_MAX_LIFETIME_MINS" envDefault:"30"`
	ConnMaxIdleTimeMins int    `env:"DB_CONN_MAX_IDLE_TIME_MINS" envDefault:"10"`
}
