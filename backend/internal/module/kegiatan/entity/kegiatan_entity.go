package entity

import "time"

// KegiatanTag adalah representasi domain dari tabel kegiatan_tags.
type KegiatanTag struct {
	ID         uint
	KegiatanID uint
	Tag        string
}

// KegiatanGallery adalah representasi domain dari tabel kegiatan_gallery.
type KegiatanGallery struct {
	ID         uint
	KegiatanID uint
	ImagePath  string
	Caption    string
	SortOrder  int
}

// Kegiatan adalah representasi domain dari tabel kegiatan.
type Kegiatan struct {
	ID          uint
	Slug        string
	Title       string
	Category    string
	EventDate   string
	Location    string
	Organizer   string
	AuthorID    *uint
	AuthorName  string
	ImagePath   *string
	Excerpt     *string
	Content     *string
	IsPublished bool
	Views       int
	Tags        []KegiatanTag
	Gallery     []KegiatanGallery
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time
}
