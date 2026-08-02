package config

type UploadConfig struct {
	Path              string   `env:"UPLOAD_PATH" envDefault:"./public/uploads"`
	PublicURL         string   `env:"UPLOAD_PUBLIC_URL" envDefault:"/uploads"`
	ImagePath         string   `env:"UPLOAD_IMAGE_PATH" envDefault:"./public/uploads/images"`
	DocumentPath      string   `env:"UPLOAD_DOCUMENT_PATH" envDefault:"./public/uploads/documents"`
	MaxSizeMB         int      `env:"UPLOAD_MAX_SIZE_MB" envDefault:"10"`
	ImageMaxSizeMB    int      `env:"UPLOAD_IMAGE_MAX_SIZE_MB" envDefault:"5"`
	DocumentMaxSizeMB int      `env:"UPLOAD_DOCUMENT_MAX_SIZE_MB" envDefault:"10"`
	AllowedExt        []string `env:"UPLOAD_ALLOWED_EXT" envDefault:"jpg,jpeg,png,pdf"`
}
