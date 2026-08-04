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
	// DeletedAt hanya terisi untuk pesan di Sampah — dipakai tabel tab Sampah
	// yang mengurutkan dan menampilkan kapan pesan dibuang.
	DeletedAt string `json:"deleted_at,omitempty"`
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
	// DeletedAt terisi kalau pesan ada di Sampah — supaya UI tidak perlu menebak
	// status dari tab yang sedang dibuka.
	DeletedAt string `json:"deleted_at,omitempty"`
}

// BulkResponse melaporkan jumlah baris yang benar-benar terpengaruh, supaya UI
// tidak mengklaim angka dari jumlah yang dipilih (bisa beda kalau ada yang sudah
// dihapus/dipulihkan lewat sesi lain).
type BulkResponse struct {
	Affected int64 `json:"affected"`
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
	Tab    string `form:"tab"`    // active (kotak masuk) | trash (sampah)
	Status string `form:"status"` // unread | read | all | trashed (alias tab=trash)
	Sort   string `form:"sort"`   // newest | oldest
}

// IsTrashed menentukan apakah query meminta isi Sampah.
// Menerima dua bentuk yang dipakai frontend/modul lain: tab=trash atau status=trashed.
func (q KontakQuery) IsTrashed() bool {
	return q.Tab == "trash" || q.Status == "trashed"
}

// DefaultLimit & MaxLimit membatasi ukuran halaman list pesan.
const (
	DefaultLimit = 10
	MaxLimit     = 100
)

// Pagination menormalkan page & limit. Dipakai bersama oleh repository dan service
// supaya offset yang di-query dan meta yang dikembalikan tidak pernah berbeda.
func (q KontakQuery) Pagination() (page, limit, offset int) {
	page = q.Page
	if page <= 0 {
		page = 1
	}
	limit = q.Limit
	if limit <= 0 {
		limit = DefaultLimit
	}
	if limit > MaxLimit {
		limit = MaxLimit
	}
	return page, limit, (page - 1) * limit
}

type BulkRequest struct {
	IDs []uint `json:"ids" binding:"required,min=1"`
}
