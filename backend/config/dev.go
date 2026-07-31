package config

type DevConfig struct {
	SeedOnStartup bool `env:"DEV_SEED_ON_STARTUP" envDefault:"false"`
	SeedAdmin     bool `env:"DEV_SEED_ADMIN" envDefault:"true"`
}
