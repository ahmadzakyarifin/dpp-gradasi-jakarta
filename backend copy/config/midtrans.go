package config

type MidtransConfig struct {
	ServerKey       string `env:"MIDTRANS_SERVER_KEY"`
	ClientKey       string `env:"MIDTRANS_CLIENT_KEY"`
	IsSandbox       bool   `env:"MIDTRANS_IS_SANDBOX" envDefault:"true"`
	TimeoutSecs     int    `env:"MIDTRANS_TIMEOUT_SECS" envDefault:"10"`
	NotificationURL string `env:"MIDTRANS_NOTIFICATION_URL"`
}
