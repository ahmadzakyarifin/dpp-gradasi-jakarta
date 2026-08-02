package config

type DatabaseConfig struct {
	Host string `env:"DB_HOST" envDefault:"127.0.0.1"`
	Port string `env:"DB_PORT" envDefault:"3306"`
	Name string `env:"DB_NAME" envDefault:"dpp_gradasi"`
	User string `env:"DB_USER"`
	Pass string `env:"DB_PASS"`
}
