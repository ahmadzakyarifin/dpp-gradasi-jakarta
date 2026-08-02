package service

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/helper"
	activitylogdto "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/activitylog/dto"
	activitylogservice "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/activitylog/service"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/kegiatan/dto"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/kegiatan/model"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/kegiatan/repository"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

type KegiatanService interface {
	GetPublished(ctx context.Context, q dto.KegiatanQuery) (*dto.KegiatanListResponse, error)
	GetAll(ctx context.Context, q dto.KegiatanQuery) (*dto.KegiatanListResponse, error)
	GetBySlug(ctx context.Context, slug string) (*dto.KegiatanDetailResponse, error)
	GetByID(ctx context.Context, id uint) (*dto.KegiatanDetailResponse, error)
	Create(ctx context.Context, req *dto.KegiatanCreateRequest, authorID uint) (*dto.KegiatanDetailResponse, error)
	Update(ctx context.Context, id uint, req *dto.KegiatanUpdateRequest) (*dto.KegiatanDetailResponse, error)
	Delete(ctx context.Context, id uint) error
	Restore(ctx context.Context, id uint) error
	BulkDelete(ctx context.Context, ids []uint) error
	BulkRestore(ctx context.Context, ids []uint) error
	DeleteGalleryImage(ctx context.Context, galleryID uint) error
	GetCategories(ctx context.Context) ([]string, error)
}

type kegiatanService struct {
	db    *gorm.DB
	repo  repository.KegiatanRepo
	audit activitylogservice.ActivityLogService
}

func NewKegiatanService(db *gorm.DB, repo repository.KegiatanRepo, audit activitylogservice.ActivityLogService) KegiatanService {
	return &kegiatanService{db: db, repo: repo, audit: audit}
}

func (s *kegiatanService) log(ctx context.Context, input *activitylogdto.ActivityLogInput) {
	if s.audit == nil {
		return
	}
	userID, userName, role, ipAddress, userAgent := helper.GetAuditMeta(ctx)
	if input.ActorID == nil && userID > 0 {
		input.ActorID = &userID
	}
	if input.ActorName == "" {
		input.ActorName = userName
	}
	if input.ActorRole == "" {
		input.ActorRole = role
	}
	if input.IPAddress == "" {
		input.IPAddress = ipAddress
	}
	if input.UserAgent == "" {
		input.UserAgent = userAgent
	}

	_ = s.audit.Log(ctx, s.db, input)
}

func (s *kegiatanService) GetPublished(ctx context.Context, q dto.KegiatanQuery) (*dto.KegiatanListResponse, error) {
	return s.list(ctx, true, q)
}

func (s *kegiatanService) GetAll(ctx context.Context, q dto.KegiatanQuery) (*dto.KegiatanListResponse, error) {
	return s.list(ctx, false, q)
}

func (s *kegiatanService) list(ctx context.Context, publishedOnly bool, q dto.KegiatanQuery) (*dto.KegiatanListResponse, error) {
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

func (s *kegiatanService) GetBySlug(ctx context.Context, slug string) (*dto.KegiatanDetailResponse, error) {
	k, err := s.repo.FindBySlug(slug)
	if err != nil {
		return nil, helper.NewServiceError("NOT_FOUND", "Kegiatan tidak ditemukan.", err)
	}
	_ = s.repo.IncrementViews(k.ID)
	k.Views++
	resp := toDetail(*k)
	return &resp, nil
}

func (s *kegiatanService) GetByID(ctx context.Context, id uint) (*dto.KegiatanDetailResponse, error) {
	k, err := s.repo.FindByID(id)
	if err != nil {
		return nil, helper.NewServiceError("NOT_FOUND", "Kegiatan tidak ditemukan.", err)
	}
	resp := toDetail(*k)
	return &resp, nil
}

func (s *kegiatanService) Create(ctx context.Context, req *dto.KegiatanCreateRequest, authorID uint) (*dto.KegiatanDetailResponse, error) {
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

	s.log(ctx, &activitylogdto.ActivityLogInput{
		Action:      "kegiatan.create",
		EntityType:  "kegiatan",
		EntityID:    &k.ID,
		EntityLabel: k.Title,
		Description: "Membuat kegiatan baru: " + k.Title,
		Metadata: map[string]any{
			"slug":     k.Slug,
			"category": k.Category,
		},
	})

	k, _ = s.repo.FindByID(k.ID)
	resp := toDetail(*k)
	return &resp, nil
}

func (s *kegiatanService) Update(ctx context.Context, id uint, req *dto.KegiatanUpdateRequest) (*dto.KegiatanDetailResponse, error) {
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

	s.log(ctx, &activitylogdto.ActivityLogInput{
		Action:      "kegiatan.update",
		EntityType:  "kegiatan",
		EntityID:    &k.ID,
		EntityLabel: k.Title,
		Description: "Memperbarui kegiatan: " + k.Title,
		Metadata: map[string]any{
			"slug": k.Slug,
		},
	})

	k, _ = s.repo.FindByID(k.ID)
	resp := toDetail(*k)
	return &resp, nil
}

func (s *kegiatanService) Delete(ctx context.Context, id uint) error {
	k, err := s.repo.FindByID(id)
	if err != nil {
		return helper.NewServiceError("NOT_FOUND", "Kegiatan tidak ditemukan.", err)
	}
	if err := s.repo.SoftDelete(id); err != nil {
		return helper.NewServiceError("SERVER_ERROR", "Gagal menghapus kegiatan.", err)
	}

	s.log(ctx, &activitylogdto.ActivityLogInput{
		Action:      "kegiatan.delete",
		EntityType:  "kegiatan",
		EntityID:    &id,
		EntityLabel: k.Title,
		Description: "Menghapus kegiatan (soft delete): " + k.Title,
	})

	return nil
}

func (s *kegiatanService) Restore(ctx context.Context, id uint) error {
	if err := s.repo.Restore(id); err != nil {
		return helper.NewServiceError("SERVER_ERROR", "Gagal memulihkan kegiatan.", err)
	}

	s.log(ctx, &activitylogdto.ActivityLogInput{
		Action:      "kegiatan.restore",
		EntityType:  "kegiatan",
		EntityID:    &id,
		Description: "Memulihkan kegiatan (ID: " + strconv.FormatUint(uint64(id), 10) + ")",
	})

	return nil
}

func (s *kegiatanService) BulkDelete(ctx context.Context, ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	if err := s.repo.BulkSoftDelete(ids); err != nil {
		return helper.NewServiceError("SERVER_ERROR", "Gagal menghapus kegiatan.", err)
	}

	s.log(ctx, &activitylogdto.ActivityLogInput{
		Action:      "kegiatan.bulk_delete",
		EntityType:  "kegiatan",
		Description: "Menghapus " + strconv.FormatUint(uint64(len(ids)), 10) + " kegiatan (soft delete)",
		Metadata: map[string]any{
			"ids": ids,
		},
	})

	return nil
}

func (s *kegiatanService) BulkRestore(ctx context.Context, ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	if err := s.repo.BulkRestore(ids); err != nil {
		return helper.NewServiceError("SERVER_ERROR", "Gagal memulihkan kegiatan.", err)
	}

	s.log(ctx, &activitylogdto.ActivityLogInput{
		Action:      "kegiatan.bulk_restore",
		EntityType:  "kegiatan",
		Description: "Memulihkan " + strconv.FormatUint(uint64(len(ids)), 10) + " kegiatan",
		Metadata: map[string]any{
			"ids": ids,
		},
	})

	return nil
}

func (s *kegiatanService) DeleteGalleryImage(ctx context.Context, galleryID uint) error {
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
func (s *kegiatanService) GetCategories(ctx context.Context) ([]string, error) {
	return s.repo.GetCategories()
}
