package dto

// KontakRequest submit pesan dari publik
type KontakRequest struct {
	Nama   string `json:"nama" binding:"required,min=3,max=100"`
	Email  string `json:"email" binding:"required,email,max=150"`
	Subjek string `json:"subjek" binding:"required,max=200"`
	Pesan  string `json:"pesan" binding:"required,min=10"`
}

type KontakListItem struct {
	ID        uint   `json:"id"`
	Nama      string `json:"nama"`
	Email     string `json:"email"`
	Subjek    string `json:"subjek"`
	Pesan     string `json:"pesan,omitempty"`
	IsRead    bool   `json:"is_read"`
	CreatedAt string `json:"created_at"`
}

type KontakDetailResponse struct {
	ID           uint   `json:"id"`
	Nama         string `json:"nama"`
	Email        string `json:"email"`
	Subjek       string `json:"subjek"`
	Pesan        string `json:"pesan"`
	IsRead       bool   `json:"is_read"`
	ResponseNote string `json:"response_note,omitempty"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type KontakListResponse struct {
	Kontak []KontakListItem `json:"kontak"`
	Meta   PaginationMeta   `json:"meta"`
}

type PaginationMeta struct {
	CurrentPage int `json:"current_page"`
	Limit       int `json:"limit"`
	TotalData   int `json:"total_data"`
	TotalPages  int `json:"total_pages"`
}

type KontakQuery struct {
	Page   int    `form:"page"`
	Limit  int    `form:"limit"`
	Search string `form:"search"`
	Status string `form:"status"` // unread | read | all
	Sort   string `form:"sort"`   // newest | oldest
}

type BulkRequest struct {
	IDs []uint `json:"ids" binding:"required,min=1"`
}
