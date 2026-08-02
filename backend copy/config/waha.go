package config

type WAHAConfig struct {
	URL               string   `env:"WAHA_URL" envDefault:"http://waha:3000"`
	APIKey            string   `env:"WAHA_API_KEY"`
	WebhookSecret     string   `env:"WAHA_WEBHOOK_SECRET"`
	DashboardUsername string   `env:"WAHA_DASHBOARD_USERNAME"`
	DashboardPassword string   `env:"WAHA_DASHBOARD_PASSWORD"`
	Session           string   `env:"WAHA_SESSION" envDefault:"default"`
	TimeoutSecs       int      `env:"WAHA_TIMEOUT_SECS" envDefault:"10"`
	DelayMillis       int      `env:"WAHA_DELAY_MILLIS" envDefault:"1000"`
	MediaPath         string   `env:"WAHA_MEDIA_PATH" envDefault:"./storage/whatsapp"`
	HookURL           string   `env:"WHATSAPP_HOOK_URL"`
	HookEvents        []string `env:"WHATSAPP_HOOK_EVENTS"`
}
