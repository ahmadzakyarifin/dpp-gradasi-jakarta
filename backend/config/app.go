package config

type AppConfig struct {
	Name           string   `env:"APP_NAME" envDefault:"DPP Gradasi"`
	Env            string   `env:"APP_ENV" envDefault:"development"`
	URL            string   `env:"APP_URL" envDefault:"http://localhost:8080"`
	Timezone       string   `env:"APP_TIMEZONE" envDefault:"Asia/Jakarta"`
	Port           string   `env:"PORT" envDefault:"8080"`
	AllowedOrigins []string `env:"ALLOWED_ORIGINS" envDefault:"*"`
}
