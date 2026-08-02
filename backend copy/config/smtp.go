package config

type SMTPConfig struct {
	Host       string `env:"SMTP_HOST" envDefault:"smtp.gmail.com"`
	Port       int    `env:"SMTP_PORT" envDefault:"587"`
	Email      string `env:"SMTP_EMAIL"`
	Pass       string `env:"SMTP_PASS"`
	FromName   string `env:"SMTP_FROM_NAME" envDefault:"SchoolPay"`
	Encryption string `env:"SMTP_ENCRYPTION" envDefault:"tls"`
	TimeoutSec int    `env:"SMTP_TIMEOUT_SECS" envDefault:"10"`
}
