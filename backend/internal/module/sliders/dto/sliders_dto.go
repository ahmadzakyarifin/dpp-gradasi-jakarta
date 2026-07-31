package dto

// SliderRequest create/update
type SliderRequest struct {
	Title     string  `json:"title" binding:"required"`
	Subtitle  *string `json:"subtitle,omitempty"`
	Tag       *string `json:"tag,omitempty"`
	IsNew     bool    `json:"is_new"`
	EventDate *string `json:"event_date,omitempty"`
	Location  *string `json:"location,omitempty"`
	ImageURL  string  `json:"image_url" binding:"required"`
	LinkURL   *string `json:"link_url,omitempty"`
	SortOrder int     `json:"sort_order"`
	IsActive  bool    `json:"is_active"`
}

// SliderResponse response
type SliderResponse struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	Subtitle  string `json:"subtitle,omitempty"`
	Tag       string `json:"tag,omitempty"`
	IsNew     bool   `json:"is_new"`
	EventDate string `json:"event_date,omitempty"`
	Location  string `json:"location,omitempty"`
	ImageURL  string `json:"image_url"`
	LinkURL   string `json:"link_url,omitempty"`
	SortOrder int    `json:"sort_order"`
	IsActive  bool   `json:"is_active"`
}

// SliderListResponse list
type SliderListResponse struct {
	Sliders []SliderResponse `json:"sliders"`
	Total   int64            `json:"total"`
}

// ReorderRequest reorder
type ReorderRequest struct {
	IDs []uint `json:"ids" binding:"required"`
}

type BulkRequest struct {
	IDs []uint `json:"ids" binding:"required,min=1"`
}
