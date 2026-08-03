package mapper

import (
	"strings"
	"time"

	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/berita/dto"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/berita/entity"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/berita/model"
	"gorm.io/gorm"
)

func toGormDeletedAt(t *time.Time) gorm.DeletedAt {
	if t == nil {
		return gorm.DeletedAt{}
	}
	return gorm.DeletedAt{Time: *t, Valid: true}
}

func fromGormDeletedAt(d gorm.DeletedAt) *time.Time {
	if !d.Valid {
		return nil
	}
	t := d.Time
	return &t
}

// strPtr mengembalikan pointer string, nil jika kosong.
func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// TagsToEntity mengonversi string tags (comma-separated) menjadi slice entity.BeritaTag.
func TagsToEntity(tags string) []entity.BeritaTag {
	parts := ParseTags(tags)
	list := make([]entity.BeritaTag, 0, len(parts))
	for _, t := range parts {
		list = append(list, entity.BeritaTag{Tag: t})
	}
	return list
}

// CreateReqToEntity mengonversi request create menjadi entity.
func CreateReqToEntity(req *dto.BeritaCreateRequest) *entity.Berita {
	if req == nil {
		return nil
	}
	b := &entity.Berita{
		Title:         req.Title,
		Category:      req.Category,
		PublishedDate: req.PublishedDate,
		AuthorID:      req.AuthorID,
		ImagePath:     strPtr(req.ImagePath),
		Excerpt:       strPtr(req.Excerpt),
		Content:       strPtr(req.Content),
		Tags:          TagsToEntity(req.Tags),
	}
	// is_published & is_featured wajib di-map — default DB is_published=1
	// membuat semua berita otomatis published kalau tidak di-set eksplisit.
	if req.IsPublished != nil {
		b.IsPublished = *req.IsPublished
	} else {
		b.IsPublished = false // default aman: draft sampai admin menerbitkan
	}
	if req.IsFeatured != nil {
		b.IsFeatured = *req.IsFeatured
	}
	return b
}

// UpdateReqToEntity menerapkan field request (partial update) ke entity existing.
func UpdateReqToEntity(req *dto.BeritaUpdateRequest, b *entity.Berita) {
	if req == nil || b == nil {
		return
	}
	if req.Title != "" {
		b.Title = req.Title
	}
	if req.Category != "" {
		b.Category = req.Category
	}
	if req.PublishedDate != "" {
		b.PublishedDate = req.PublishedDate
	}
	if req.ImagePath != "" {
		b.ImagePath = strPtr(req.ImagePath)
	}
	if req.Excerpt != "" {
		b.Excerpt = strPtr(req.Excerpt)
	}
	if req.Content != "" {
		b.Content = strPtr(req.Content)
	}
	if req.IsFeatured != nil {
		b.IsFeatured = *req.IsFeatured
	}
	if req.IsPublished != nil {
		b.IsPublished = *req.IsPublished
	}
	if req.AuthorID != nil {
		b.AuthorID = req.AuthorID
	}
	if req.Tags != "" {
		b.Tags = TagsToEntity(req.Tags)
	}
}

// EntityToModel mengonversi entity menjadi model GORM.
func EntityToModel(e *entity.Berita) *model.Berita {
	if e == nil {
		return nil
	}
	tags := make([]model.BeritaTag, 0, len(e.Tags))
	for _, t := range e.Tags {
		tags = append(tags, model.BeritaTag{ID: t.ID, BeritaID: t.BeritaID, Tag: t.Tag})
	}
	return &model.Berita{
		ID:            e.ID,
		Slug:          e.Slug,
		Title:         e.Title,
		Category:      e.Category,
		PublishedDate: normalizeDate(e.PublishedDate),
		AuthorID:      e.AuthorID,
		AuthorName:    e.AuthorName,
		ImagePath:     e.ImagePath,
		Excerpt:       e.Excerpt,
		Content:       e.Content,
		IsFeatured:    e.IsFeatured,
		IsPublished:   e.IsPublished,
		Views:         e.Views,
		Tags:          tags,
		CreatedAt:     e.CreatedAt,
		UpdatedAt:     e.UpdatedAt,
		DeletedAt:     toGormDeletedAt(e.DeletedAt),
	}
}

// normalizeDate memastikan tanggal berformat YYYY-MM-DD.
// GORM bisa meng-scan kolom DATE ke string dengan format RFC3339
// (mis. "2026-08-03T00:00:00+07:00") — MySQL menolak itu di kolom DATE.
func normalizeDate(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}
	// Ambil 10 karakter pertama (YYYY-MM-DD) kalau format lebih panjang
	if len(s) > 10 {
		if t, err := time.Parse(time.RFC3339, s); err == nil {
			return t.Format("2006-01-02")
		}
		if t, err := time.Parse("2006-01-02 15:04:05", s); err == nil {
			return t.Format("2006-01-02")
		}
		return s[:10]
	}
	return s
}

// ModelToEntity mengonversi model GORM menjadi entity.
func ModelToEntity(m *model.Berita) *entity.Berita {
	if m == nil {
		return nil
	}
	tags := make([]entity.BeritaTag, 0, len(m.Tags))
	for _, t := range m.Tags {
		tags = append(tags, entity.BeritaTag{ID: t.ID, BeritaID: t.BeritaID, Tag: t.Tag})
	}
	return &entity.Berita{
		ID:            m.ID,
		Slug:          m.Slug,
		Title:         m.Title,
		Category:      m.Category,
		PublishedDate: m.PublishedDate,
		AuthorID:      m.AuthorID,
		AuthorName:    m.AuthorName,
		ImagePath:     m.ImagePath,
		Excerpt:       m.Excerpt,
		Content:       m.Content,
		IsFeatured:    m.IsFeatured,
		IsPublished:   m.IsPublished,
		Views:         m.Views,
		Tags:          tags,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
		DeletedAt:     fromGormDeletedAt(m.DeletedAt),
	}
}

// ModelListToEntity mengonversi daftar model menjadi daftar entity.
func ModelListToEntity(models []model.Berita) []entity.Berita {
	list := make([]entity.Berita, 0, len(models))
	for i := range models {
		list = append(list, *ModelToEntity(&models[i]))
	}
	return list
}

// EntityToListItem mengonversi entity menjadi list item response.
func EntityToListItem(e *entity.Berita) dto.BeritaListItem {
	if e == nil {
		return dto.BeritaListItem{}
	}
	item := dto.BeritaListItem{
		ID:            e.ID,
		Title:         e.Title,
		Slug:          e.Slug,
		Category:      e.Category,
		PublishedDate: FormatDate(e.PublishedDate),
		IsFeatured:    e.IsFeatured,
		Views:         e.Views,
		AuthorName:    e.AuthorName,
		CreatedAt:     e.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:     e.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
	if e.ImagePath != nil {
		item.ImagePath = *e.ImagePath
	}
	if e.Excerpt != nil {
		item.Excerpt = *e.Excerpt
	}
	return item
}

// EntityListToItem mengonversi daftar entity menjadi daftar list item response.
func EntityListToItem(entities []entity.Berita) []dto.BeritaListItem {
	list := make([]dto.BeritaListItem, 0, len(entities))
	for i := range entities {
		list = append(list, EntityToListItem(&entities[i]))
	}
	return list
}

// EntityToDetail mengonversi entity menjadi detail response.
func EntityToDetail(e *entity.Berita) dto.BeritaDetailResponse {
	if e == nil {
		return dto.BeritaDetailResponse{}
	}
	resp := dto.BeritaDetailResponse{
		ID:            e.ID,
		Title:         e.Title,
		Slug:          e.Slug,
		Category:      e.Category,
		PublishedDate: FormatDate(e.PublishedDate),
		IsFeatured:    e.IsFeatured,
		IsPublished:   e.IsPublished,
		Views:         e.Views,
		AuthorName:    e.AuthorName,
		CreatedAt:     e.CreatedAt.Format("2006-01-02T15:04:05Z"),
		Tags:          make([]string, 0),
	}
	if e.ImagePath != nil {
		resp.ImagePath = *e.ImagePath
	}
	if e.Excerpt != nil {
		resp.Excerpt = *e.Excerpt
	}
	if e.Content != nil {
		resp.Content = *e.Content
	}
	for _, t := range e.Tags {
		resp.Tags = append(resp.Tags, t.Tag)
	}
	return resp
}

// ParseTags memecah string tag (comma-separated) menjadi slice, trim & buang kosong.
func ParseTags(tags string) []string {
	parts := strings.Split(tags, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}

// FormatDate menormalkan format tanggal published_date ke "2006-01-02".
func FormatDate(v string) string {
	if v == "" {
		return ""
	}
	t, err := time.Parse("2006-01-02T15:04:05Z07:00", v)
	if err == nil {
		return t.Format("2006-01-02")
	}
	if _, err = time.Parse("2006-01-02", v); err == nil {
		return v
	}
	if len(v) >= 10 {
		return v[:10]
	}
	return v
}
