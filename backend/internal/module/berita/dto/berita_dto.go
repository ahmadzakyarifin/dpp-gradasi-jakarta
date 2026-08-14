package dto

// BeritaCreateRequest — create (wajib title, content, published_date)
type BeritaCreateRequest struct {
	Title         string `json:"title" binding:"required,min=5,max=300"`
	Category      string `json:"category"`
	PublishedDate string `json:"published_date" binding:"required"`
	AuthorID      *uint  `json:"author_id,omitempty"`
	ImagePath     string `json:"image_path"`
	Excerpt       string `json:"excerpt"`
	Content       string `json:"content" binding:"required"`
	Footnote      string `json:"footnote"`
	ImageSource   string `json:"image_source"`
	IsFeatured    *bool  `json:"is_featured"`
	IsPublished   *bool  `json:"is_published"`
	Tags          string `json:"tags"` // comma-separated
}

// BeritaUpdateRequest — update (partial; field kosong diabaikan)
type BeritaUpdateRequest struct {
	Title         string `json:"title" binding:"omitempty,min=5,max=300"`
	Category      string `json:"category"`
	PublishedDate string `json:"published_date"`
	AuthorID      *uint  `json:"author_id,omitempty"`
	ImagePath     string `json:"image_path"`
	Excerpt       string `json:"excerpt"`
	Content       string `json:"content"`
	Footnote      string `json:"footnote"`
	ImageSource   string `json:"image_source"`
	IsFeatured    *bool  `json:"is_featured"`
	IsPublished   *bool  `json:"is_published"`
	Tags          string `json:"tags"` // comma-separated
}

type BeritaListItem struct {
	ID            uint   `json:"id"`
	Title         string `json:"title"`
	Slug          string `json:"slug"`
	Category      string `json:"category"`
	PublishedDate string `json:"published_date"`
	AuthorName    string `json:"author_name,omitempty"`
	ImagePath     string `json:"image_path,omitempty"`
	Excerpt       string `json:"excerpt,omitempty"`
	IsFeatured    bool   `json:"is_featured"`
	IsPublished   *bool  `json:"is_published,omitempty"`
	Views         int    `json:"views"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at,omitempty"`
}

type BeritaDetailResponse struct {
	ID            uint     `json:"id"`
	Title         string   `json:"title"`
	Slug          string   `json:"slug"`
	Category      string   `json:"category"`
	PublishedDate string   `json:"published_date"`
	AuthorName    string   `json:"author_name,omitempty"`
	ImagePath     string   `json:"image_path,omitempty"`
	Excerpt       string   `json:"excerpt,omitempty"`
	Content       string   `json:"content,omitempty"`
	Footnote      string   `json:"footnote,omitempty"`
	ImageSource   string   `json:"image_source,omitempty"`
	IsFeatured    bool     `json:"is_featured"`
	IsPublished   bool     `json:"is_published"`
	Views         int      `json:"views"`
	Tags          []string `json:"tags,omitempty"`
	CreatedAt     string   `json:"created_at"`
}

type BeritaListResponse struct {
	Berita []BeritaListItem `json:"berita"`
	Meta   PaginationMeta   `json:"meta"`
}

type PaginationMeta struct {
	CurrentPage int `json:"current_page"`
	Limit       int `json:"limit"`
	TotalData   int `json:"total_data"`
	TotalPages  int `json:"total_pages"`
}

type BeritaQuery struct {
	Page     int    `form:"page"`
	Limit    int    `form:"limit"`
	Search   string `form:"search"`
	Category string `form:"category"`
	Tag      string `form:"tag"`
	Sort     string `form:"sort"`
	Status   string `form:"status"`
}

type BulkRequest struct {
	IDs []uint `json:"ids" binding:"required,min=1"`
}

// UploadImageResponse — hasil upload gambar berita (cover)
type UploadImageResponse struct {
	ImagePath string `json:"image_path"`
}
