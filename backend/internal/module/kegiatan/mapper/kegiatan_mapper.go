package mapper

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/kegiatan/dto"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/kegiatan/entity"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/kegiatan/model"
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

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// TagsToEntity mengonversi string tags (comma-separated) menjadi slice entity.KegiatanTag.
func TagsToEntity(tags string) []entity.KegiatanTag {
	parts := ParseTags(tags)
	list := make([]entity.KegiatanTag, 0, len(parts))
	for _, t := range parts {
		list = append(list, entity.KegiatanTag{Tag: t})
	}
	return list
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

// ParseGalleryJSON mengonversi JSON string gallery menjadi slice GalleryInput.
func ParseGalleryJSON(s string) ([]dto.GalleryInput, error) {
	if s == "" {
		return nil, nil
	}
	var items []dto.GalleryInput
	if err := json.Unmarshal([]byte(s), &items); err != nil {
		return nil, err
	}
	return items, nil
}

// GalleryInputsToEntity mengonversi slice GalleryInput menjadi slice entity.KegiatanGallery.
func GalleryInputsToEntity(items []dto.GalleryInput) []entity.KegiatanGallery {
	list := make([]entity.KegiatanGallery, 0, len(items))
	for _, it := range items {
		if it.ImagePath == "" {
			continue
		}
		list = append(list, entity.KegiatanGallery{
			ImagePath: it.ImagePath,
			Caption:   it.Caption,
			SortOrder: it.SortOrder,
		})
	}
	return list
}

// CreateReqToEntity mengonversi request create menjadi entity.
func CreateReqToEntity(req *dto.KegiatanCreateRequest) *entity.Kegiatan {
	if req == nil {
		return nil
	}
	return &entity.Kegiatan{
		Title:     req.Title,
		Category:  req.Category,
		EventDate: req.EventDate,
		Location:  req.Location,
		Organizer: req.Organizer,
		AuthorID:  req.AuthorID,
		ImagePath: strPtr(req.ImagePath),
		Excerpt:   strPtr(req.Excerpt),
		Content:   strPtr(req.Content),
		Tags:      TagsToEntity(req.Tags),
	}
}

// UpdateReqToEntity menerapkan field request (partial update) ke entity existing.
func UpdateReqToEntity(req *dto.KegiatanUpdateRequest, k *entity.Kegiatan) {
	if req == nil || k == nil {
		return
	}
	if req.Title != "" {
		k.Title = req.Title
	}
	if req.Category != "" {
		k.Category = req.Category
	}
	if req.EventDate != "" {
		k.EventDate = req.EventDate
	}
	if req.Location != "" {
		k.Location = req.Location
	}
	if req.Organizer != "" {
		k.Organizer = req.Organizer
	}
	if req.ImagePath != "" {
		k.ImagePath = strPtr(req.ImagePath)
	}
	if req.Excerpt != "" {
		k.Excerpt = strPtr(req.Excerpt)
	}
	if req.Content != "" {
		k.Content = strPtr(req.Content)
	}
	if req.IsPublished != nil {
		k.IsPublished = *req.IsPublished
	}
	if req.AuthorID != nil {
		k.AuthorID = req.AuthorID
	}
	if req.Tags != "" {
		k.Tags = TagsToEntity(req.Tags)
	}
	if req.GalleryJSON != "" {
		if items, err := ParseGalleryJSON(req.GalleryJSON); err == nil {
			k.Gallery = GalleryInputsToEntity(items)
		}
	}
}

// EntityToModel mengonversi entity menjadi model GORM.
func EntityToModel(e *entity.Kegiatan) *model.Kegiatan {
	if e == nil {
		return nil
	}
	tags := make([]model.KegiatanTag, 0, len(e.Tags))
	for _, t := range e.Tags {
		tags = append(tags, model.KegiatanTag{ID: t.ID, KegiatanID: t.KegiatanID, Tag: t.Tag})
	}
	gallery := make([]model.KegiatanGallery, 0, len(e.Gallery))
	for _, g := range e.Gallery {
		gallery = append(gallery, model.KegiatanGallery{
			ID:         g.ID,
			KegiatanID: g.KegiatanID,
			ImagePath:  g.ImagePath,
			Caption:    g.Caption,
			SortOrder:  g.SortOrder,
		})
	}
	return &model.Kegiatan{
		ID:          e.ID,
		Slug:        e.Slug,
		Title:       e.Title,
		Category:    e.Category,
		EventDate:   e.EventDate,
		Location:    e.Location,
		Organizer:   e.Organizer,
		AuthorID:    e.AuthorID,
		AuthorName:  e.AuthorName,
		ImagePath:   e.ImagePath,
		Excerpt:     e.Excerpt,
		Content:     e.Content,
		IsPublished: e.IsPublished,
		Views:       e.Views,
		Tags:        tags,
		Gallery:     gallery,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
		DeletedAt:   toGormDeletedAt(e.DeletedAt),
	}
}

// ModelToEntity mengonversi model GORM menjadi entity.
func ModelToEntity(m *model.Kegiatan) *entity.Kegiatan {
	if m == nil {
		return nil
	}
	tags := make([]entity.KegiatanTag, 0, len(m.Tags))
	for _, t := range m.Tags {
		tags = append(tags, entity.KegiatanTag{ID: t.ID, KegiatanID: t.KegiatanID, Tag: t.Tag})
	}
	gallery := make([]entity.KegiatanGallery, 0, len(m.Gallery))
	for _, g := range m.Gallery {
		gallery = append(gallery, entity.KegiatanGallery{
			ID:         g.ID,
			KegiatanID: g.KegiatanID,
			ImagePath:  g.ImagePath,
			Caption:    g.Caption,
			SortOrder:  g.SortOrder,
		})
	}
	return &entity.Kegiatan{
		ID:          m.ID,
		Slug:        m.Slug,
		Title:       m.Title,
		Category:    m.Category,
		EventDate:   m.EventDate,
		Location:    m.Location,
		Organizer:   m.Organizer,
		AuthorID:    m.AuthorID,
		AuthorName:  m.AuthorName,
		ImagePath:   m.ImagePath,
		Excerpt:     m.Excerpt,
		Content:     m.Content,
		IsPublished: m.IsPublished,
		Views:       m.Views,
		Tags:        tags,
		Gallery:     gallery,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		DeletedAt:   fromGormDeletedAt(m.DeletedAt),
	}
}

// ModelListToEntity mengonversi daftar model menjadi daftar entity.
func ModelListToEntity(models []model.Kegiatan) []entity.Kegiatan {
	list := make([]entity.Kegiatan, 0, len(models))
	for i := range models {
		list = append(list, *ModelToEntity(&models[i]))
	}
	return list
}

// EntityToListItem mengonversi entity menjadi list item response.
func EntityToListItem(e *entity.Kegiatan) dto.KegiatanListItem {
	if e == nil {
		return dto.KegiatanListItem{}
	}
	item := dto.KegiatanListItem{
		ID:           e.ID,
		Title:        e.Title,
		Slug:         e.Slug,
		Category:     e.Category,
		EventDate:    FormatDate(e.EventDate),
		Location:     e.Location,
		Organizer:    e.Organizer,
		AuthorName:   e.AuthorName,
		Views:        e.Views,
		IsPublished:  &e.IsPublished,
		GalleryCount: len(e.Gallery),
		CreatedAt:    e.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:    e.UpdatedAt.Format("2006-01-02T15:04:05Z"),
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
func EntityListToItem(entities []entity.Kegiatan) []dto.KegiatanListItem {
	list := make([]dto.KegiatanListItem, 0, len(entities))
	for i := range entities {
		list = append(list, EntityToListItem(&entities[i]))
	}
	return list
}

// EntityToDetail mengonversi entity menjadi detail response.
func EntityToDetail(e *entity.Kegiatan) dto.KegiatanDetailResponse {
	if e == nil {
		return dto.KegiatanDetailResponse{}
	}
	resp := dto.KegiatanDetailResponse{
		ID:          e.ID,
		Title:       e.Title,
		Slug:        e.Slug,
		Category:    e.Category,
		EventDate:   FormatDate(e.EventDate),
		Location:    e.Location,
		Organizer:   e.Organizer,
		AuthorName:  e.AuthorName,
		IsPublished: e.IsPublished,
		Views:       e.Views,
		CreatedAt:   e.CreatedAt.Format("2006-01-02T15:04:05Z"),
		Tags:        make([]string, 0),
		Gallery:     make([]dto.GalleryImageItem, 0),
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
	for _, g := range e.Gallery {
		resp.Gallery = append(resp.Gallery, dto.GalleryImageItem{
			ID:        g.ID,
			ImagePath: g.ImagePath,
			Caption:   g.Caption,
			SortOrder: g.SortOrder,
		})
	}
	return resp
}

// FormatDate menormalkan format tanggal event_date ke "2006-01-02".
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
	// Return as-is if it's a custom string (e.g., "30 Desember 2026")
	return v
}
