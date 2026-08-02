package config

type CookieConfig struct {
	Domain   string `env:"COOKIE_DOMAIN" envDefault:""`
	Secure   bool   `env:"COOKIE_SECURE" envDefault:"false"`
	HTTPOnly bool   `env:"COOKIE_HTTP_ONLY" envDefault:"true"`
	Path     string `env:"COOKIE_PATH" envDefault:"/api/v1/auth"`
}
