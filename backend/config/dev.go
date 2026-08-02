package config

// DevConfig berisi konfigurasi mode development / seeding.
// Hanya dipakai oleh cmd/seeder (super admin awal) — bukan runtime API.
type DevConfig struct {
	SuperAdminName     string `env:"SUPERADMIN_NAME" envDefault:"Super Admin"`
	SuperAdminEmail    string `env:"SUPERADMIN_EMAIL" envDefault:"admin@gradasi.org"`
	SuperAdminPassword string `env:"SUPERADMIN_PASSWORD" envDefault:"password123"`
}
