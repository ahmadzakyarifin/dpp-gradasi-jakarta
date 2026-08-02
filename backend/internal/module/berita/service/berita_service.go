package service

import (
	"strings"
	"time"

	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/helper"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/berita/dto"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/berita/model"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/berita/repository"
	"github.com/gosimple/slug"
)

type BeritaService interface {
	GetPublished(query dto.BeritaQuery) (*dto.BeritaListResponse, error)
	GetAll(query dto.BeritaQuery) (*dto.BeritaListResponse, error)
	GetBySlug(slug string) (*dto.BeritaDetailResponse, error)
	GetByID(id uint) (*dto.BeritaDetailResponse, error)
	Create(req *dto.BeritaCreateRequest, authorID uint) (*dto.BeritaDetailResponse, error)
	Update(id uint, req *dto.BeritaUpdateRequest) (*dto.BeritaDetailResponse, error)
	Delete(id uint) error
	Restore(id uint) error
	BulkDelete(ids []uint) error
	BulkRestore(ids []uint) error
	GetCategories() ([]string, error)
}

type beritaService struct {
	repo repository.BeritaRepo
}

func NewBeritaService(repo repository.BeritaRepo) BeritaService {
	return &beritaService{repo: repo}
}

func (s *beritaService) GetPublished(query dto.BeritaQuery) (*dto.BeritaListResponse, error) {
	return s.list(true, query)
}

func (s *beritaService) GetAll(query dto.BeritaQuery) (*dto.BeritaListResponse, error) {
	return s.list(false, query)
}

func (s *beritaService) list(publishedOnly bool, q dto.BeritaQuery) (*dto.BeritaListResponse, error) {
	var beritas []model.Berita
	var total int64
	var err error

	if publishedOnly {
		beritas, total, err = s.repo.FindPublished(q)
	} else {
		beritas, total, err = s.repo.FindAll(q)
	}
	if err != nil {
		return nil, helper.NewServiceError("SERVER_ERROR", "Gagal mengambil data berita.", err)
	}

	page := q.Page
	if page <= 0 {
		page = 1
	}
	limit := q.Limit
	if limit <= 0 {
		limit = 10
	}
	totalPages := (int(total) + limit - 1) / limit

	resp := &dto.BeritaListResponse{
		Berita: make([]dto.BeritaListItem, 0),
		Meta: dto.PaginationMeta{
			CurrentPage: page,
			Limit:       limit,
			TotalData:   int(total),
			TotalPages:  totalPages,
		},
	}

	for _, b := range beritas {
		resp.Berita = append(resp.Berita, toListItem(b))
	}

	return resp, nil
}

func (s *beritaService) GetBySlug(slug string) (*dto.BeritaDetailResponse, error) {
	b, err := s.repo.FindBySlug(slug)
	if err != nil {
		return nil, helper.NewServiceError("NOT_FOUND", "Berita tidak ditemukan.", err)
	}

	// Increment views (async-safe enough)
	_ = s.repo.IncrementViews(b.ID)
	b.Views++

	resp := toDetail(*b)
	return &resp, nil
}

func (s *beritaService) GetByID(id uint) (*dto.BeritaDetailResponse, error) {
	b, err := s.repo.FindByID(id)
	if err != nil {
		return nil, helper.NewServiceError("NOT_FOUND", "Berita tidak ditemukan.", err)
	}
	resp := toDetail(*b)
	return &resp, nil
}

func (s *beritaService) Create(req *dto.BeritaCreateRequest, authorID uint) (*dto.BeritaDetailResponse, error) {
	title := strings.TrimSpace(req.Title)
	if title == "" {
		return nil, helper.NewServiceError("VALIDATION_ERROR", "Judul wajib diisi.", nil)
	}
	content := strings.TrimSpace(req.Content)
	if content == "" {
		return nil, helper.NewServiceError("VALIDATION_ERROR", "Konten lengkap wajib diisi.", nil)
	}

	// Validasi judul duplikat — error jelas (bukan auto-suffix)
	dup, err := s.repo.ExistsTitle(title, 0)
	if err != nil {
		return nil, helper.NewServiceError("SERVER_ERROR", "Gagal memvalidasi judul.", err)
	}
	if dup {
		return nil, helper.NewServiceError("DUPLICATE_TITLE", "Judul sudah digunakan, gunakan judul lain.", nil)
	}

	slugStr := slug.Make(title)

	cat := strings.TrimSpace(req.Category)
	if cat == "" {
		cat = "Berita Organisasi"
	}

	var isFeatured, isPublished bool
	if req.IsFeatured != nil {
		isFeatured = *req.IsFeatured
	}
	if req.IsPublished != nil {
		isPublished = *req.IsPublished
	} else {
		isPublished = true
	}

	b := &model.Berita{
		Slug:          slugStr,
		Title:         title,
		Category:      cat,
		PublishedDate: req.PublishedDate,
		ImagePath:     strPtr(req.ImagePath),
		Excerpt:       strPtr(req.Excerpt),
		Content:       strPtr(req.Content),
		IsFeatured:    isFeatured,
		IsPublished:   isPublished,
	}

	if req.AuthorID != nil {
		b.AuthorID = req.AuthorID
	} else {
		b.AuthorID = &authorID
	}

	if err := s.repo.Create(b); err != nil {
		return nil, helper.NewServiceError("SERVER_ERROR", "Gagal membuat berita.", err)
	}

	// Save tags
	if req.Tags != "" {
		tags := parseTags(req.Tags)
		_ = s.repo.SaveTags(b.ID, tags)
	}

	// Reload with tags
	b, _ = s.repo.FindByID(b.ID)

	resp := toDetail(*b)
	return &resp, nil
}

func (s *beritaService) Update(id uint, req *dto.BeritaUpdateRequest) (*dto.BeritaDetailResponse, error) {
	b, err := s.repo.FindByID(id)
	if err != nil {
		return nil, helper.NewServiceError("NOT_FOUND", "Berita tidak ditemukan.", err)
	}

	// Partial update: hanya field yang dikirim yang diubah
	if req.Title != "" {
		title := strings.TrimSpace(req.Title)
		if title == "" {
			return nil, helper.NewServiceError("VALIDATION_ERROR", "Judul wajib diisi.", nil)
		}
		// Cek duplikat judul (kecuali id ini sendiri)
		dup, err := s.repo.ExistsTitle(title, id)
		if err != nil {
			return nil, helper.NewServiceError("SERVER_ERROR", "Gagal memvalidasi judul.", err)
		}
		if dup {
			return nil, helper.NewServiceError("DUPLICATE_TITLE", "Judul sudah digunakan, gunakan judul lain.", nil)
		}
		b.Title = title
	}

	cat := strings.TrimSpace(req.Category)
	if cat != "" {
		b.Category = cat
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

	if err := s.repo.Update(b); err != nil {
		return nil, helper.NewServiceError("SERVER_ERROR", "Gagal mengupdate berita.", err)
	}

	// Save tags
	if req.Tags != "" {
		tags := parseTags(req.Tags)
		_ = s.repo.SaveTags(b.ID, tags)
	}

	b, _ = s.repo.FindByID(b.ID)
	resp := toDetail(*b)
	return &resp, nil
}

func (s *beritaService) Delete(id uint) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return helper.NewServiceError("NOT_FOUND", "Berita tidak ditemukan.", err)
	}
	return s.repo.SoftDelete(id)
}

func (s *beritaService) Restore(id uint) error {
	return s.repo.Restore(id)
}

func (s *beritaService) BulkDelete(ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	return s.repo.BulkSoftDelete(ids)
}

func (s *beritaService) BulkRestore(ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	return s.repo.BulkRestore(ids)
}

func toListItem(b model.Berita) dto.BeritaListItem {
	item := dto.BeritaListItem{
		ID:            b.ID,
		Title:         b.Title,
		Slug:          b.Slug,
		Category:      b.Category,
		PublishedDate: formatDate(b.PublishedDate),
		IsFeatured:    b.IsFeatured,
		Views:         b.Views,
		AuthorName:    b.AuthorName,
		CreatedAt:     b.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:     b.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
	if b.ImagePath != nil {
		item.ImagePath = *b.ImagePath
	}
	if b.Excerpt != nil {
		item.Excerpt = *b.Excerpt
	}
	return item
}

func toDetail(b model.Berita) dto.BeritaDetailResponse {
	resp := dto.BeritaDetailResponse{
		ID:            b.ID,
		Title:         b.Title,
		Slug:          b.Slug,
		Category:      b.Category,
		PublishedDate: formatDate(b.PublishedDate),
		IsFeatured:    b.IsFeatured,
		IsPublished:   b.IsPublished,
		Views:         b.Views,
		AuthorName:    b.AuthorName,
		CreatedAt:     b.CreatedAt.Format("2006-01-02T15:04:05Z"),
		Tags:          make([]string, 0),
	}
	if b.ImagePath != nil {
		resp.ImagePath = *b.ImagePath
	}
	if b.Excerpt != nil {
		resp.Excerpt = *b.Excerpt
	}
	if b.Content != nil {
		resp.Content = *b.Content
	}
	for _, t := range b.Tags {
		resp.Tags = append(resp.Tags, t.Tag)
	}
	return resp
}

func parseTags(tags string) []string {
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

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func formatDate(v string) string {
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
	// Return as-is if can't parse (e.g. already date string)
	if len(v) >= 10 {
		return v[:10]
	}
	return v
}

// GetCategories mengembalikan daftar kategori unik untuk dropdown admin.
func (s *beritaService) GetCategories() ([]string, error) {
	return s.repo.GetCategories()
}
