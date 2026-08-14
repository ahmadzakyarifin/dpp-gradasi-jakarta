package dto

import (
	"mime/multipart"
	"strings"
)

type PengurusRequest struct {
	Name         string                `form:"name" binding:"required,min=3,max=150"`
	Role         string                `form:"role" binding:"required,max=200"`
	Department   string                `form:"department" binding:"max=100"`
	Level        string                `form:"level" binding:"required"`
	Provinsi     string                `form:"provinsi" binding:"max=100"`
	Kabupaten    string                `form:"kabupaten" binding:"max=100"`
	FacebookURL  string                `form:"facebook_url" binding:"max=500"`
	InstagramURL string                `form:"instagram_url" binding:"max=500"`
	LinkedinURL  string                `form:"linkedin_url" binding:"max=500"`
	Whatsapp     string                `form:"whatsapp" binding:"max=20"`
	Email        string                `form:"email" binding:"omitempty,email,max=150"`
	Pekerjaan    string                `form:"pekerjaan" binding:"max=150"`
	Bio          string                `form:"bio"`
	Pendidikan   string                `form:"pendidikan"`
	Sertifikasi  string                `form:"sertifikasi"`
	Periode      string                `form:"periode" binding:"required,max=50"`
	SortOrder    int                   `form:"sort_order"`
	IsActive     *bool                 `form:"is_active"`
	Image        *multipart.FileHeader `form:"image"`
	CV           *multipart.FileHeader `form:"cv"`
}

// ValidateRegionRules: validasi required_if level — provinsi wajib utk dpd/dpc,
// kabupaten wajib utk dpc. Dipanggil dari service (bukan binding tag) agar pesan
// error konsisten dengan contract.
func (r *PengurusRequest) ValidateRegionRules() map[string]string {
	errorsMap := make(map[string]string)
	// Validate Level enum
	validLevels := map[string]bool{
		"Ketua Umum":        true,
		"Pengurus Pusat":    true,
		"Pengurus Provinsi": true,
		"Pengurus Kab/Kota": true,
	}
	if !validLevels[r.Level] {
		errorsMap["level"] = "Level tidak valid"
	}

	switch r.Level {
	case "Pengurus Provinsi", "Pengurus Kab/Kota":
		if strings.TrimSpace(r.Provinsi) == "" {
			errorsMap["provinsi"] = "Provinsi wajib diisi untuk level Pengurus Provinsi/Kab/Kota"
		}
	}
	if r.Level == "Pengurus Kab/Kota" {
		if strings.TrimSpace(r.Kabupaten) == "" {
			errorsMap["kabupaten"] = "Kabupaten wajib diisi untuk level Pengurus Kab/Kota"
		}
	}
	return errorsMap
}

type PengurusResponse struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	Role         string `json:"role"`
	Department   string `json:"department,omitempty"`
	Level        string `json:"level"`
	Provinsi     string `json:"provinsi,omitempty"`
	Kabupaten    string `json:"kabupaten,omitempty"`
	ImagePath    string `json:"image_path"`
	CVPath       string `json:"cv_path,omitempty"`
	FacebookURL  string `json:"facebook_url,omitempty"`
	InstagramURL string `json:"instagram_url,omitempty"`
	LinkedinURL  string `json:"linkedin_url,omitempty"`
	Whatsapp     string `json:"whatsapp,omitempty"`
	Email        string `json:"email,omitempty"`
	Pekerjaan    string `json:"pekerjaan,omitempty"`
	Bio          string `json:"bio,omitempty"`
	Pendidikan   string `json:"pendidikan,omitempty"`
	Sertifikasi  string `json:"sertifikasi,omitempty"`
	Periode      string `json:"periode"`
	SortOrder    int    `json:"sort_order"`
	IsActive     bool   `json:"is_active"`
	CreatedAt    string `json:"created_at,omitempty"`
	UpdatedAt    string `json:"updated_at,omitempty"`
}

type PengurusListResponse struct {
	Data []PengurusResponse `json:"data"`
	Meta PaginationMeta     `json:"meta"`
}

type PengurusQuery struct {
	Level     string `form:"level"`
	Provinsi  string `form:"provinsi"`
	Kabupaten string `form:"kabupaten"`
	Search    string `form:"search"`
	Status    string `form:"status"` // active, inactive, all
	Page      int    `form:"page"`
	Limit     int    `form:"limit"`
	Sort      string `form:"sort"`
	Trashed   bool   `form:"trashed"`
}

type RegionsResponse struct {
	Provinsi  []string            `json:"provinsi"`
	Kabupaten map[string][]string `json:"kabupaten"`
}

type PaginationMeta struct {
	CurrentPage int `json:"current_page"`
	Limit       int `json:"limit"`
	TotalData   int `json:"total_data"`
	TotalPages  int `json:"total_pages"`
}

type BulkRequest struct {
	IDs []uint `json:"ids" binding:"required,min=1"`
}
