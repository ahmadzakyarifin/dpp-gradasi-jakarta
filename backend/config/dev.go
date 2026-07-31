package config

type DevConfig struct {
	SeedOnStartup      bool   `env:"DEV_SEED_ON_STARTUP" envDefault:"true"`
	SeedAdmin          bool   `env:"DEV_SEED_ADMIN" envDefault:"true"`
	SuperAdminName     string `env:"SUPERADMIN_NAME" envDefault:"Super Admin"`
	SuperAdminEmail    string `env:"SUPERADMIN_EMAIL" envDefault:"admin@gradasi.org"`
	SuperAdminPassword string `env:"SUPERADMIN_PASSWORD" envDefault:"password123"`
}
