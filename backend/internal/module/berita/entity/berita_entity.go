package entity

import "time"

// BeritaTag adalah representasi domain dari tabel berita_tags.
type BeritaTag struct {
	ID       uint
	BeritaID uint
	Tag      string
}

// Berita adalah representasi domain dari tabel `berita`.
type Berita struct {
	ID            uint
	Slug          string
	Title         string
	Category      string
	PublishedDate string
	AuthorID      *uint
	AuthorName    string
	ImagePath     *string
	Excerpt       *string
	Content       *string
	IsFeatured    bool
	IsNew         bool
	IsPublished   bool
	Views         int
	Footnote      *string
	ImageSource   *string
	Tags          []BeritaTag
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time
}
