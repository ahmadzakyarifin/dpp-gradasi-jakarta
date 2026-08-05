package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/helper"
	activitylogdto "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/activitylog/dto"
	activitylogservice "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/activitylog/service"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/berita/dto"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/berita/entity"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/berita/mapper"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/berita/repository"
	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

type BeritaService interface {
	GetPublished(ctx context.Context, query dto.BeritaQuery) (*dto.BeritaListResponse, error)
	GetAll(ctx context.Context, query dto.BeritaQuery) (*dto.BeritaListResponse, error)
	GetBySlug(ctx context.Context, slug string) (*dto.BeritaDetailResponse, error)
	GetByID(ctx context.Context, id uint) (*dto.BeritaDetailResponse, error)
	Create(ctx context.Context, req *dto.BeritaCreateRequest, authorID uint) (*dto.BeritaDetailResponse, error)
	Update(ctx context.Context, id uint, req *dto.BeritaUpdateRequest) (*dto.BeritaDetailResponse, error)
	Delete(ctx context.Context, id uint) error
	Restore(ctx context.Context, id uint) error
	BulkDelete(ctx context.Context, ids []uint) error
	BulkRestore(ctx context.Context, ids []uint) error
	GetCategories(ctx context.Context) ([]string, error)
	UploadImage(ctx context.Context, file *multipart.FileHeader) (*dto.UploadImageResponse, error)
}

type beritaService struct {
	db         *gorm.DB
	repo       repository.BeritaRepo
	audit      activitylogservice.ActivityLogService
	uploadPath string
}

func NewBeritaService(db *gorm.DB, repo repository.BeritaRepo, audit activitylogservice.ActivityLogService) BeritaService {
	return &beritaService{
		db:         db,
		repo:       repo,
		audit:      audit,
		uploadPath: "public/uploads/berita",
	}
}

func (s *beritaService) log(ctx context.Context, db *gorm.DB, input *activitylogdto.ActivityLogInput) {
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

	_ = s.audit.Log(ctx, db, input)
}

func (s *beritaService) GetPublished(ctx context.Context, query dto.BeritaQuery) (*dto.BeritaListResponse, error) {
	return s.list(ctx, true, query)
}

func (s *beritaService) GetAll(ctx context.Context, query dto.BeritaQuery) (*dto.BeritaListResponse, error) {
	return s.list(ctx, false, query)
}

func (s *beritaService) list(ctx context.Context, publishedOnly bool, q dto.BeritaQuery) (*dto.BeritaListResponse, error) {
	var beritas []entity.Berita
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
		Berita: mapper.EntityListToItem(beritas),
		Meta: dto.PaginationMeta{
			CurrentPage: page,
			Limit:       limit,
			TotalData:   int(total),
			TotalPages:  totalPages,
		},
	}
	return resp, nil
}

func (s *beritaService) GetBySlug(ctx context.Context, slug string) (*dto.BeritaDetailResponse, error) {
	b, err := s.repo.FindBySlug(slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helper.NewNotFoundError("Berita tidak ditemukan")
		}
		return nil, helper.NewServiceError("SERVER_ERROR", "Gagal mengambil berita.", err)
	}

	// Increment views (async-safe enough)
	_ = s.repo.IncrementViews(b.ID)
	b.Views++

	resp := mapper.EntityToDetail(b)
	return &resp, nil
}

func (s *beritaService) GetByID(ctx context.Context, id uint) (*dto.BeritaDetailResponse, error) {
	b, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helper.NewNotFoundError("Berita tidak ditemukan")
		}
		return nil, helper.NewServiceError("SERVER_ERROR", "Gagal mengambil berita.", err)
	}
	resp := mapper.EntityToDetail(b)
	return &resp, nil
}

// UploadImage menyimpan file gambar cover berita ke public/uploads/berita.
// Mengembalikan path relatif yang bisa langsung dipakai di image_path.
func (s *beritaService) UploadImage(ctx context.Context, file *multipart.FileHeader) (*dto.UploadImageResponse, error) {
	if file == nil {
		return nil, helper.NewServiceError("VALIDATION_ERROR", "File gambar wajib diunggah", nil)
	}

	// Validasi ukuran: maks 5MB
	if file.Size > 5*1024*1024 {
		return nil, helper.NewServiceError("VALIDATION_ERROR", "File gambar tidak valid. Maksimal 5MB dengan format PNG, JPG, atau WEBP.", nil)
	}

	// Validasi MIME type
	src, err := file.Open()
	if err != nil {
		return nil, helper.NewServiceError("SERVER_ERROR", "Gagal membaca file gambar.", err)
	}
	defer src.Close()

	header := make([]byte, 512)
	n, err := src.Read(header)
	if err != nil && err != io.EOF {
		return nil, helper.NewServiceError("SERVER_ERROR", "Gagal membaca file gambar.", err)
	}
	// Reset reader agar bisa dibaca ulang saat copy
	if _, err := src.Seek(0, io.SeekStart); err != nil {
		return nil, helper.NewServiceError("SERVER_ERROR", "Gagal membaca file gambar.", err)
	}

	mimeType := http.DetectContentType(header[:n])
	switch mimeType {
	case "image/png", "image/jpeg", "image/webp":
	default:
		return nil, helper.NewServiceError("VALIDATION_ERROR", "File gambar tidak valid. Maksimal 5MB dengan format PNG, JPG, atau WEBP.", nil)
	}

	ext := filepath.Ext(file.Filename)
	if ext == "" {
		switch mimeType {
		case "image/png":
			ext = ".png"
		case "image/webp":
			ext = ".webp"
		default:
			ext = ".jpg"
		}
	}
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)

	// Pastikan direktori upload ada
	if err := os.MkdirAll(s.uploadPath, 0o755); err != nil {
		return nil, helper.NewServiceError("SERVER_ERROR", "Gagal menyimpan file gambar.", err)
	}

	dst := filepath.Join(s.uploadPath, filename)
	out, err := os.Create(dst)
	if err != nil {
		return nil, helper.NewServiceError("SERVER_ERROR", "Gagal menyimpan file gambar.", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, src); err != nil {
		return nil, helper.NewServiceError("SERVER_ERROR", "Gagal menyimpan file gambar.", err)
	}

	imagePath := "/uploads/berita/" + filename
	return &dto.UploadImageResponse{ImagePath: imagePath}, nil
}

func (s *beritaService) Create(ctx context.Context, req *dto.BeritaCreateRequest, authorID uint) (*dto.BeritaDetailResponse, error) {
	v := helper.NewValidationError()

	title := strings.TrimSpace(req.Title)
	if title == "" {
		v.Add("title", "Judul wajib diisi.")
	}
	content := strings.TrimSpace(req.Content)
	if content == "" {
		v.Add("content", "Konten lengkap wajib diisi.")
	}
	if len(v.Errors) > 0 {
		return nil, v
	}

	slugStr := slug.Make(title)

	// Validasi slug duplikat — error jelas (bukan auto-suffix)
	dup, err := s.repo.ExistsSlug(slugStr, 0)
	if err != nil {
		return nil, helper.NewServiceError("SERVER_ERROR", "Gagal memvalidasi judul.", err)
	}
	if dup {
		return nil, helper.NewServiceError("DUPLICATE_TITLE", "Judul berita sudah terpakai karena menghasilkan slug yang sama dengan berita lain. Gunakan judul lain.", nil)
	}

	cat := strings.TrimSpace(req.Category)
	if cat == "" {
		cat = "Berita Organisasi"
	}

	b := mapper.CreateReqToEntity(req)
	b.Slug = slugStr
	b.Title = title
	b.Category = cat

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
		_ = s.repo.SaveTags(b.ID, mapper.ParseTags(req.Tags))
	}

	s.log(ctx, s.db, &activitylogdto.ActivityLogInput{
		Action:      "berita.create",
		EntityType:  "berita",
		EntityID:    &b.ID,
		EntityLabel: b.Title,
		Description: "Membuat berita baru: " + b.Title,
		Metadata: map[string]any{
			"slug":     b.Slug,
			"category": b.Category,
		},
	})

	// Reload with tags
	created, err := s.repo.FindByID(b.ID)
	if err == nil {
		b = created
	}

	resp := mapper.EntityToDetail(b)
	return &resp, nil
}

func (s *beritaService) Update(ctx context.Context, id uint, req *dto.BeritaUpdateRequest) (*dto.BeritaDetailResponse, error) {
	b, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helper.NewNotFoundError("Berita tidak ditemukan")
		}
		return nil, helper.NewServiceError("SERVER_ERROR", "Gagal mengambil berita.", err)
	}

	// Partial update: hanya field yang dikirim yang diubah
	if req.Title != "" {
		title := strings.TrimSpace(req.Title)
		if title == "" {
			v := helper.NewValidationError()
			v.Add("title", "Judul wajib diisi.")
			return nil, v
		}
		slugStr := slug.Make(title)
		// Cek duplikat slug (kecuali id ini sendiri)
		dup, err := s.repo.ExistsSlug(slugStr, id)
		if err != nil {
			return nil, helper.NewServiceError("SERVER_ERROR", "Gagal memvalidasi judul.", err)
		}
		if dup {
			return nil, helper.NewServiceError("DUPLICATE_TITLE", "Judul berita sudah terpakai karena menghasilkan slug yang sama dengan berita lain. Gunakan judul lain.", nil)
		}
	}

	mapper.UpdateReqToEntity(req, b)

	if err := s.repo.Update(b); err != nil {
		return nil, helper.NewServiceError("SERVER_ERROR", "Gagal mengupdate berita.", err)
	}

	// Save tags
	if req.Tags != "" {
		_ = s.repo.SaveTags(b.ID, mapper.ParseTags(req.Tags))
	}

	s.log(ctx, s.db, &activitylogdto.ActivityLogInput{
		Action:      "berita.update",
		EntityType:  "berita",
		EntityID:    &b.ID,
		EntityLabel: b.Title,
		Description: "Memperbarui berita: " + b.Title,
		Metadata: map[string]any{
			"slug": b.Slug,
		},
	})

	updated, err := s.repo.FindByID(b.ID)
	if err == nil {
		b = updated
	}
	resp := mapper.EntityToDetail(b)
	return &resp, nil
}

func (s *beritaService) Delete(ctx context.Context, id uint) error {
	b, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return helper.NewNotFoundError("Berita tidak ditemukan")
		}
		return helper.NewServiceError("SERVER_ERROR", "Gagal mengambil berita.", err)
	}

	if err := s.repo.SoftDelete(id); err != nil {
		return helper.NewServiceError("SERVER_ERROR", "Gagal menghapus berita.", err)
	}

	s.log(ctx, s.db, &activitylogdto.ActivityLogInput{
		Action:      "berita.delete",
		EntityType:  "berita",
		EntityID:    &id,
		EntityLabel: b.Title,
		Description: "Menghapus berita (soft delete): " + b.Title,
	})

	return nil
}

func (s *beritaService) Restore(ctx context.Context, id uint) error {
	if err := s.repo.Restore(id); err != nil {
		return helper.NewServiceError("SERVER_ERROR", "Gagal memulihkan berita.", err)
	}

	s.log(ctx, s.db, &activitylogdto.ActivityLogInput{
		Action:      "berita.restore",
		EntityType:  "berita",
		EntityID:    &id,
		Description: "Memulihkan berita (ID: " + strconv.FormatUint(uint64(id), 10) + ")",
	})

	return nil
}

func (s *beritaService) BulkDelete(ctx context.Context, ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	if err := s.repo.BulkSoftDelete(ids); err != nil {
		return helper.NewServiceError("SERVER_ERROR", "Gagal menghapus berita.", err)
	}

	s.log(ctx, s.db, &activitylogdto.ActivityLogInput{
		Action:      "berita.bulk_delete",
		EntityType:  "berita",
		Description: "Menghapus " + strconv.FormatUint(uint64(len(ids)), 10) + " berita (soft delete)",
		Metadata: map[string]any{
			"ids": ids,
		},
	})

	return nil
}

func (s *beritaService) BulkRestore(ctx context.Context, ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	if err := s.repo.BulkRestore(ids); err != nil {
		return helper.NewServiceError("SERVER_ERROR", "Gagal memulihkan berita.", err)
	}

	s.log(ctx, s.db, &activitylogdto.ActivityLogInput{
		Action:      "berita.bulk_restore",
		EntityType:  "berita",
		Description: "Memulihkan " + strconv.FormatUint(uint64(len(ids)), 10) + " berita",
		Metadata: map[string]any{
			"ids": ids,
		},
	})

	return nil
}

// GetCategories mengembalikan daftar kategori unik untuk dropdown admin.
func (s *beritaService) GetCategories(ctx context.Context) ([]string, error) {
	categories, err := s.repo.GetCategories()
	if err != nil {
		return nil, helper.NewServiceError("SERVER_ERROR", "Gagal mengambil daftar kategori.", err)
	}
	return categories, nil
}
