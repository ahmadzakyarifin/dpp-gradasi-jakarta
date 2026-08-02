package service

import (
	"strings"
	"time"

	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/helper"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/kegiatan/dto"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/kegiatan/model"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/kegiatan/repository"
	"github.com/gosimple/slug"
)

type KegiatanService interface {
	GetPublished(q dto.KegiatanQuery) (*dto.KegiatanListResponse, error)
	GetAll(q dto.KegiatanQuery) (*dto.KegiatanListResponse, error)
	GetBySlug(slug string) (*dto.KegiatanDetailResponse, error)
	GetByID(id uint) (*dto.KegiatanDetailResponse, error)
	Create(req *dto.KegiatanCreateRequest, authorID uint) (*dto.KegiatanDetailResponse, error)
	Update(id uint, req *dto.KegiatanUpdateRequest) (*dto.KegiatanDetailResponse, error)
	Delete(id uint) error
	Restore(id uint) error
	BulkDelete(ids []uint) error
	BulkRestore(ids []uint) error
	DeleteGalleryImage(galleryID uint) error
	GetCategories() ([]string, error)
}

type kegiatanService struct {
	repo repository.KegiatanRepo
}

func NewKegiatanService(repo repository.KegiatanRepo) KegiatanService {
	return &kegiatanService{repo: repo}
}

func (s *kegiatanService) GetPublished(q dto.KegiatanQuery) (*dto.KegiatanListResponse, error) {
	return s.list(true, q)
}

func (s *kegiatanService) GetAll(q dto.KegiatanQuery) (*dto.KegiatanListResponse, error) {
	return s.list(false, q)
}

func (s *kegiatanService) list(publishedOnly bool, q dto.KegiatanQuery) (*dto.KegiatanListResponse, error) {
	var items []model.Kegiatan
	var total int64
	var err error

	if publishedOnly {
		items, total, err = s.repo.FindPublished(q)
	} else {
		items, total, err = s.repo.FindAll(q)
	}
	if err != nil {
		return nil, helper.NewServiceError("SERVER_ERROR", "Gagal mengambil data kegiatan.", err)
	}

	page := maxInt(q.Page, 1)
	limit := q.Limit
	if limit <= 0 {
		limit = 10
	}
	totalPages := (int(total) + limit - 1) / limit

	resp := &dto.KegiatanListResponse{
		Kegiatan: make([]dto.KegiatanListItem, 0),
		Meta: dto.PaginationMeta{
			CurrentPage: page,
			Limit:       limit,
			TotalData:   int(total),
			TotalPages:  totalPages,
		},
	}
	for _, k := range items {
		resp.Kegiatan = append(resp.Kegiatan, toListItem(k))
	}
	return resp, nil
}

func (s *kegiatanService) GetBySlug(slug string) (*dto.KegiatanDetailResponse, error) {
	k, err := s.repo.FindBySlug(slug)
	if err != nil {
		return nil, helper.NewServiceError("NOT_FOUND", "Kegiatan tidak ditemukan.", err)
	}
	_ = s.repo.IncrementViews(k.ID)
	k.Views++
	resp := toDetail(*k)
	return &resp, nil
}

func (s *kegiatanService) GetByID(id uint) (*dto.KegiatanDetailResponse, error) {
	k, err := s.repo.FindByID(id)
	if err != nil {
		return nil, helper.NewServiceError("NOT_FOUND", "Kegiatan tidak ditemukan.", err)
	}
	resp := toDetail(*k)
	return &resp, nil
}

func (s *kegiatanService) Create(req *dto.KegiatanCreateRequest, authorID uint) (*dto.KegiatanDetailResponse, error) {
	cat := strings.TrimSpace(req.Category)
	if cat == "" {
		cat = "Kegiatan"
	}

	// Validasi unik slug — jika sudah ada (termasuk di trash), tambah suffix -2, -3, dst
	uniqueSlug, err := s.repo.FindUniqueSlug(slug.Make(req.Title))
	if err != nil {
		return nil, helper.NewServiceError("SERVER_ERROR", "Gagal memvalidasi slug.", err)
	}

	isPub := true
	if req.IsPublished != nil {
		isPub = *req.IsPublished
	}

	k := &model.Kegiatan{
		Slug:        uniqueSlug,
		Title:       req.Title,
		Category:    cat,
		EventDate:   req.EventDate,
		Location:    req.Location,
		Organizer:   req.Organizer,
		ImagePath:   strPtr(req.ImagePath),
		Excerpt:     strPtr(req.Excerpt),
		Content:     strPtr(req.Content),
		IsPublished: isPub,
	}

	if req.AuthorID != nil {
		k.AuthorID = req.AuthorID
	} else {
		k.AuthorID = &authorID
	}

	if err := s.repo.Create(k); err != nil {
		return nil, helper.NewServiceError("SERVER_ERROR", "Gagal membuat kegiatan.", err)
	}

	if req.Tags != "" {
		if err := s.repo.SaveTags(k.ID, parseTags(req.Tags)); err != nil {
			return nil, helper.NewServiceError("SERVER_ERROR", "Gagal menyimpan tag kegiatan.", err)
		}
	}

	if req.GalleryJSON != "" {
		items, err := repository.ParseGalleryJSON(req.GalleryJSON)
		if err == nil && len(items) > 0 {
			if err := s.repo.SaveGallery(k.ID, items); err != nil {
				return nil, helper.NewServiceError("SERVER_ERROR", "Gagal menyimpan galeri kegiatan.", err)
			}
		}
	}

	k, _ = s.repo.FindByID(k.ID)
	resp := toDetail(*k)
	return &resp, nil
}

func (s *kegiatanService) Update(id uint, req *dto.KegiatanUpdateRequest) (*dto.KegiatanDetailResponse, error) {
	k, err := s.repo.FindByID(id)
	if err != nil {
		return nil, helper.NewServiceError("NOT_FOUND", "Kegiatan tidak ditemukan.", err)
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
		k.Title = title
	}

	cat := strings.TrimSpace(req.Category)
	if cat != "" {
		k.Category = cat
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

	if err := s.repo.Update(k); err != nil {
		return nil, helper.NewServiceError("SERVER_ERROR", "Gagal mengupdate kegiatan.", err)
	}

	if req.Tags != "" {
		if err := s.repo.SaveTags(k.ID, parseTags(req.Tags)); err != nil {
			return nil, helper.NewServiceError("SERVER_ERROR", "Gagal menyimpan tag kegiatan.", err)
		}
	}
	if req.GalleryJSON != "" {
		items, err := repository.ParseGalleryJSON(req.GalleryJSON)
		if err == nil && len(items) > 0 {
			if err := s.repo.SaveGallery(k.ID, items); err != nil {
				return nil, helper.NewServiceError("SERVER_ERROR", "Gagal menyimpan galeri kegiatan.", err)
			}
		}
	}

	k, _ = s.repo.FindByID(k.ID)
	resp := toDetail(*k)
	return &resp, nil
}

func (s *kegiatanService) Delete(id uint) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return helper.NewServiceError("NOT_FOUND", "Kegiatan tidak ditemukan.", err)
	}
	return s.repo.SoftDelete(id)
}

func (s *kegiatanService) Restore(id uint) error {
	return s.repo.Restore(id)
}

func (s *kegiatanService) BulkDelete(ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	return s.repo.BulkSoftDelete(ids)
}

func (s *kegiatanService) BulkRestore(ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	return s.repo.BulkRestore(ids)
}

func (s *kegiatanService) DeleteGalleryImage(galleryID uint) error {
	return s.repo.DeleteGalleryImage(galleryID)
}

func toListItem(k model.Kegiatan) dto.KegiatanListItem {
	item := dto.KegiatanListItem{
		ID:           k.ID,
		Title:        k.Title,
		Slug:         k.Slug,
		Category:     k.Category,
		EventDate:    formatDate(k.EventDate),
		Location:     k.Location,
		Organizer:    k.Organizer,
		AuthorName:   k.AuthorName,
		Views:        k.Views,
		IsPublished:  &k.IsPublished,
		GalleryCount: len(k.Gallery),
		CreatedAt:    k.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:    k.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
	if k.ImagePath != nil {
		item.ImagePath = *k.ImagePath
	}
	if k.Excerpt != nil {
		item.Excerpt = *k.Excerpt
	}
	return item
}

func toDetail(k model.Kegiatan) dto.KegiatanDetailResponse {
	resp := dto.KegiatanDetailResponse{
		ID:          k.ID,
		Title:       k.Title,
		Slug:        k.Slug,
		Category:    k.Category,
		EventDate:   formatDate(k.EventDate),
		Location:    k.Location,
		Organizer:   k.Organizer,
		AuthorName:  k.AuthorName,
		IsPublished: k.IsPublished,
		Views:       k.Views,
		CreatedAt:   k.CreatedAt.Format("2006-01-02T15:04:05Z"),
		Tags:        make([]string, 0),
		Gallery:     make([]dto.GalleryImageItem, 0),
	}
	if k.ImagePath != nil {
		resp.ImagePath = *k.ImagePath
	}
	if k.Excerpt != nil {
		resp.Excerpt = *k.Excerpt
	}
	if k.Content != nil {
		resp.Content = *k.Content
	}
	for _, t := range k.Tags {
		resp.Tags = append(resp.Tags, t.Tag)
	}
	for _, g := range k.Gallery {
		resp.Gallery = append(resp.Gallery, dto.GalleryImageItem{
			ID:        g.ID,
			ImagePath: g.ImagePath,
			Caption:   g.Caption,
			SortOrder: g.SortOrder,
		})
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
	if len(v) >= 10 {
		return v[:10]
	}
	return v
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// GetCategories mengembalikan daftar kategori unik untuk dropdown admin.
func (s *kegiatanService) GetCategories() ([]string, error) {
	return s.repo.GetCategories()
}
