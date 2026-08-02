package config

type DevUser struct {
	Name     string `env:"NAME"`
	Email    string `env:"EMAIL"`
	Password string `env:"PASSWORD"`
	Phone    string `env:"PHONE"`
}

type DevConfig struct {
	Admin  DevUser `envPrefix:"DEV_ADMIN_"`
	Parent DevUser `envPrefix:"DEV_PARENT_"`
}
