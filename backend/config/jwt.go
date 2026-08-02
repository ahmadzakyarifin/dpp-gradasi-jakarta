package config

type JWTConfig struct {
	Secret                  string `env:"JWT_SECRET"`
	AccessTTLMinutes        int    `env:"JWT_ACCESS_TTL_MINS" envDefault:"15"`
	RefreshTTLHours         int    `env:"JWT_REFRESH_TTL_HOURS" envDefault:"72"`
	RememberMeTTLHours      int    `env:"JWT_REMEMBER_ME_TTL_HOURS" envDefault:"720"`
	PasswordResetTTLMinutes int    `env:"PASSWORD_RESET_TTL_MINS" envDefault:"15"`
}
