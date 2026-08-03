package dto

// KegiatanCreateRequest — create (wajib penuh)
type KegiatanCreateRequest struct {
	Title       string `json:"title" binding:"required,min=5,max=300"`
	Category    string `json:"category"`
	EventDate   string `json:"event_date"`
	Location    string `json:"location"`
	Organizer   string `json:"organizer"`
	AuthorID    *uint  `json:"author_id,omitempty"`
	ImagePath   string `json:"image_path"`
	Excerpt     string `json:"excerpt"`
	Content     string `json:"content" binding:"required"`
	IsPublished *bool  `json:"is_published"`
	Tags        string `json:"tags"`    // comma-separated
	GalleryJSON string `json:"gallery"` // JSON string: [{"image_path":"...","caption":"...","sort_order":0}]
}

// KegiatanUpdateRequest — update (partial; field kosong diabaikan)
type KegiatanUpdateRequest struct {
	Title       string `json:"title" binding:"omitempty,min=5,max=300"`
	Category    string `json:"category"`
	EventDate   string `json:"event_date"`
	Location    string `json:"location"`
	Organizer   string `json:"organizer"`
	AuthorID    *uint  `json:"author_id,omitempty"`
	ImagePath   string `json:"image_path"`
	Excerpt     string `json:"excerpt"`
	Content     string `json:"content"`
	IsPublished *bool  `json:"is_published"`
	Tags        string `json:"tags"`    // comma-separated
	GalleryJSON string `json:"gallery"` // JSON string: [{"image_path":"...","caption":"...","sort_order":0}]
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
	ImagePath    string `json:"image_path,omitempty"`
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
	ImagePath   string             `json:"image_path,omitempty"`
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
	ImagePath string `json:"image_path"`
	Caption   string `json:"caption,omitempty"`
	SortOrder int    `json:"sort_order"`
}

type GalleryInput struct {
	ImagePath string `json:"image_path"`
	Caption   string `json:"caption"`
	SortOrder int    `json:"sort_order"`
}

// UploadImageResponse — hasil upload gambar kegiatan (cover/galeri)
type UploadImageResponse struct {
	ImagePath string `json:"image_path"`
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
