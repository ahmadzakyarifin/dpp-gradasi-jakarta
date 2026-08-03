package entity

import "time"

// Pengurus adalah representasi domain dari tabel pengurus.
type Pengurus struct {
	ID           uint
	Name         string
	Role         string
	Department   *string
	Level        string
	Provinsi     *string
	Kabupaten    *string
	ImagePath    string
	FacebookURL  *string
	InstagramURL *string
	LinkedinURL  *string
	Whatsapp     *string
	Periode      string
	SortOrder    int
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}
