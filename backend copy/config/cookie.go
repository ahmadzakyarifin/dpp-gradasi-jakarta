package config

type CookieConfig struct {
	Domain   string `env:"COOKIE_DOMAIN"`
	Path     string `env:"COOKIE_PATH" envDefault:"/"`
	HTTPOnly bool   `env:"COOKIE_HTTP_ONLY" envDefault:"true"`
	Secure   bool   `env:"COOKIE_SECURE" envDefault:"false"`
	SameSite string `env:"COOKIE_SAME_SITE" envDefault:"lax"`
	MaxAge   int    `env:"COOKIE_MAX_AGE" envDefault:"86400"`
}
