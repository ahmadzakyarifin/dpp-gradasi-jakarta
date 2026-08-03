package entity

import "time"

// Slider adalah representasi domain dari tabel sliders.
type Slider struct {
	ID        uint
	Title     string
	Subtitle  *string
	Tag       *string
	IsNew     bool
	EventDate *string
	Location  *string
	ImagePath string
	LinkURL   *string
	SortOrder int
	IsActive  bool
	CreatedBy *uint
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
