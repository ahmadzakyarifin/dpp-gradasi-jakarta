package config

type AppConfig struct {
	Name           string `env:"APP_NAME" envDefault:"DPP Gradasi"`
	Env            string `env:"APP_ENV" envDefault:"development"`
	URL            string `env:"APP_URL" envDefault:"http://localhost:8080"`
	Port           string `env:"PORT" envDefault:"8080"`
	// ALLOWED_ORIGINS — whitelist origin CORS, comma-separated. Default dev: Vite.
	AllowedOrigins string `env:"ALLOWED_ORIGINS" envDefault:"http://localhost:5173"`
	FrontendURL    string `env:"FRONTEND_URL" envDefault:"http://localhost:5173"`

	// SEO / Open Graph — dipakai endpoint /api/v1/seo/* untuk generate meta tags.
	SEOSiteName           string `env:"APP_SEO_SITE_NAME" envDefault:"Gradasi Generasi Digital"`
	SEODefaultTitle       string `env:"APP_SEO_DEFAULT_TITLE" envDefault:"Gradasi Generasi Digital Indonesia"`
	SEODefaultDescription string `env:"APP_SEO_DEFAULT_DESCRIPTION" envDefault:"Portal berita dan kegiatan DPP Gradasi Jakarta"`
	SEODefaultImage       string `env:"APP_SEO_DEFAULT_IMAGE" envDefault:"/assets/default-og-image.jpg"`
}
