package dto

// KegiatanRequest create/update
type KegiatanRequest struct {
	Title       string `json:"title" binding:"required,min=5,max=300"`
	Category    string `json:"category"`
	EventDate   string `json:"event_date"`
	Location    string `json:"location"`
	Organizer   string `json:"organizer"`
	AuthorID    *uint  `json:"author_id,omitempty"`
	ImageURL    string `json:"image_url"`
	Excerpt     string `json:"excerpt"`
	Content     string `json:"content" binding:"required"`
	IsPublished *bool  `json:"is_published"`
	Tags        string `json:"tags"`    // comma-separated
	GalleryJSON string `json:"gallery"` // JSON string: [{"image_url":"...","caption":"...","sort_order":0}]
}

type KegiatanListItem struct {
	ID           uint   `json:"id"`
	Title        string `json:"title"`
	Slug         string `json:"slug"`
	Category     string `json:"category"`
	EventDate    string `json:"event_date,omitempty"`
	Location     string `json:"location,omitempty"`
	Organizer    string `json:"organizer,omitempty"`
	AuthorName   string `json:"author_name,omitempty"`
	ImageURL     string `json:"image_url,omitempty"`
	Excerpt      string `json:"excerpt,omitempty"`
	IsPublished  *bool  `json:"is_published,omitempty"`
	Views        int    `json:"views"`
	GalleryCount int    `json:"gallery_count,omitempty"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at,omitempty"`
}

type KegiatanDetailResponse struct {
	ID          uint               `json:"id"`
	Title       string             `json:"title"`
	Slug        string             `json:"slug"`
	Category    string             `json:"category"`
	EventDate   string             `json:"event_date,omitempty"`
	Location    string             `json:"location,omitempty"`
	Organizer   string             `json:"organizer,omitempty"`
	AuthorName  string             `json:"author_name,omitempty"`
	ImageURL    string             `json:"image_url,omitempty"`
	Excerpt     string             `json:"excerpt,omitempty"`
	Content     string             `json:"content,omitempty"`
	IsPublished bool               `json:"is_published"`
	Views       int                `json:"views"`
	Tags        []string           `json:"tags,omitempty"`
	Gallery     []GalleryImageItem `json:"gallery,omitempty"`
	CreatedAt   string             `json:"created_at"`
}

type GalleryImageItem struct {
	ID        uint   `json:"id"`
	ImageURL  string `json:"image_url"`
	Caption   string `json:"caption,omitempty"`
	SortOrder int    `json:"sort_order"`
}

type GalleryInput struct {
	ImageURL  string `json:"image_url"`
	Caption   string `json:"caption"`
	SortOrder int    `json:"sort_order"`
}

type KegiatanListResponse struct {
	Kegiatan []KegiatanListItem `json:"kegiatan"`
	Meta     PaginationMeta     `json:"meta"`
}

type PaginationMeta struct {
	CurrentPage int `json:"current_page"`
	Limit       int `json:"limit"`
	TotalData   int `json:"total_data"`
	TotalPages  int `json:"total_pages"`
}

type KegiatanQuery struct {
	Page     int    `form:"page"`
	Limit    int    `form:"limit"`
	Search   string `form:"search"`
	Category string `form:"category"`
	Sort     string `form:"sort"`
	Status   string `form:"status"`
}

type BulkRequest struct {
	IDs []uint `json:"ids" binding:"required,min=1"`
}
