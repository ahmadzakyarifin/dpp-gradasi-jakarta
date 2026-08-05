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
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/kegiatan/dto"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/kegiatan/entity"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/kegiatan/mapper"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/kegiatan/repository"
	"github.com/gosimple/slug"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type KegiatanService interface {
	GetPublished(ctx context.Context, q dto.KegiatanQuery) (*dto.KegiatanListResponse, error)
	GetAll(ctx context.Context, q dto.KegiatanQuery) (*dto.KegiatanListResponse, error)
	GetBySlug(ctx context.Context, slug string) (*dto.KegiatanDetailResponse, error)
	GetByID(ctx context.Context, id uint) (*dto.KegiatanDetailResponse, error)
	Create(ctx context.Context, req *dto.KegiatanCreateRequest, authorID uint) (*dto.KegiatanDetailResponse, error)
	Update(ctx context.Context, id uint, req *dto.KegiatanUpdateRequest) (*dto.KegiatanDetailResponse, error)
	UploadImage(ctx context.Context, file *multipart.FileHeader) (*dto.UploadImageResponse, error)
	Delete(ctx context.Context, id uint) error
	Restore(ctx context.Context, id uint) error
	BulkDelete(ctx context.Context, ids []uint) error
	BulkRestore(ctx context.Context, ids []uint) error
	DeleteGalleryImage(ctx context.Context, galleryID uint) error
	GetCategories(ctx context.Context) ([]string, error)
}

type kegiatanService struct {
	db         *gorm.DB
	repo       repository.KegiatanRepo
	audit      activitylogservice.ActivityLogService
	uploadPath string
}

func NewKegiatanService(db *gorm.DB, repo repository.KegiatanRepo, audit activitylogservice.ActivityLogService) KegiatanService {
	uploadPath := "public/uploads/kegiatan"
	if err := os.MkdirAll(uploadPath, 0o755); err != nil {
		helper.Logger.Error("gagal buat direktori upload kegiatan", zap.Error(err))
	}
	return &kegiatanService{db: db, repo: repo, audit: audit, uploadPath: uploadPath}
}

func (s *kegiatanService) log(ctx context.Context, db *gorm.DB, input *activitylogdto.ActivityLogInput) {
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

func (s *kegiatanService) GetPublished(ctx context.Context, q dto.KegiatanQuery) (*dto.KegiatanListResponse, error) {
	return s.list(ctx, true, q)
}

func (s *kegiatanService) GetAll(ctx context.Context, q dto.KegiatanQuery) (*dto.KegiatanListResponse, error) {
	return s.list(ctx, false, q)
}

func (s *kegiatanService) list(ctx context.Context, publishedOnly bool, q dto.KegiatanQuery) (*dto.KegiatanListResponse, error) {
	var (
		total    int64
		entities []entity.Kegiatan
		err      error
	)

	if publishedOnly {
		entities, total, err = s.repo.FindPublished(q)
	} else {
		entities, total, err = s.repo.FindAll(q)
	}
	if err != nil {
		return nil, helper.NewServiceError("SERVER_ERROR", "Gagal mengambil data kegiatan.", err)
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

	resp := &dto.KegiatanListResponse{
		Kegiatan: mapper.EntityListToItem(entities),
		Meta: dto.PaginationMeta{
			CurrentPage: page,
			Limit:       limit,
			TotalData:   int(total),
			TotalPages:  totalPages,
		},
	}
	return resp, nil
}

func (s *kegiatanService) GetBySlug(ctx context.Context, slug string) (*dto.KegiatanDetailResponse, error) {
	k, err := s.repo.FindBySlug(slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helper.NewNotFoundError("Kegiatan tidak ditemukan")
		}
		return nil, helper.NewServiceError("SERVER_ERROR", "Gagal mengambil kegiatan.", err)
	}
	_ = s.repo.IncrementViews(k.ID)
	k.Views++
	resp := mapper.EntityToDetail(k)
	return &resp, nil
}

func (s *kegiatanService) GetByID(ctx context.Context, id uint) (*dto.KegiatanDetailResponse, error) {
	k, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helper.NewNotFoundError("Kegiatan tidak ditemukan")
		}
		return nil, helper.NewServiceError("SERVER_ERROR", "Gagal mengambil kegiatan.", err)
	}
	resp := mapper.EntityToDetail(k)
	return &resp, nil
}

func (s *kegiatanService) Create(ctx context.Context, req *dto.KegiatanCreateRequest, authorID uint) (*dto.KegiatanDetailResponse, error) {
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

	cat := strings.TrimSpace(req.Category)
	if cat == "" {
		cat = "Kegiatan"
	}

	// Validasi unik slug — jika sudah ada (termasuk di trash), tambah suffix -2, -3, dst
	uniqueSlug, err := s.repo.FindUniqueSlug(slug.Make(title))
	if err != nil {
		return nil, helper.NewServiceError("SERVER_ERROR", "Gagal memvalidasi slug.", err)
	}

	isPub := true
	if req.IsPublished != nil {
		isPub = *req.IsPublished
	}

	k := mapper.CreateReqToEntity(req)
	k.Slug = uniqueSlug
	k.Title = title
	k.Category = cat
	k.IsPublished = isPub

	if req.AuthorID != nil {
		k.AuthorID = req.AuthorID
	} else {
		k.AuthorID = &authorID
	}

	if err := s.repo.Create(k); err != nil {
		return nil, helper.NewServiceError("SERVER_ERROR", "Gagal membuat kegiatan.", err)
	}

	if req.Tags != "" {
		if err := s.repo.SaveTags(k.ID, mapper.ParseTags(req.Tags)); err != nil {
			return nil, helper.NewServiceError("SERVER_ERROR", "Gagal menyimpan tag kegiatan.", err)
		}
	}

	if req.GalleryJSON != "" {
		items, err := mapper.ParseGalleryJSON(req.GalleryJSON)
		if err == nil && len(items) > 0 {
			if err := s.repo.SaveGallery(k.ID, items); err != nil {
				return nil, helper.NewServiceError("SERVER_ERROR", "Gagal menyimpan galeri kegiatan.", err)
			}
		}
	}

	s.log(ctx, s.db, &activitylogdto.ActivityLogInput{
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

	created, err := s.repo.FindByID(k.ID)
	if err == nil {
		k = created
	}
	resp := mapper.EntityToDetail(k)
	return &resp, nil
}

func (s *kegiatanService) Update(ctx context.Context, id uint, req *dto.KegiatanUpdateRequest) (*dto.KegiatanDetailResponse, error) {
	k, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helper.NewNotFoundError("Kegiatan tidak ditemukan")
		}
		return nil, helper.NewServiceError("SERVER_ERROR", "Gagal mengambil kegiatan.", err)
	}

	// Partial update: hanya field yang dikirim yang diubah
	if req.Title != "" {
		title := strings.TrimSpace(req.Title)
		if title == "" {
			v := helper.NewValidationError()
			v.Add("title", "Judul wajib diisi.")
			return nil, v
		}
		// Cek duplikat judul (kecuali id ini sendiri)
		slugStr := slug.Make(title)
		// Cek duplikat slug (kecuali id ini sendiri)
		dup, err := s.repo.ExistsSlug(slugStr, id)
		if err != nil {
			return nil, helper.NewServiceError("SERVER_ERROR", "Gagal memvalidasi judul.", err)
		}
		if dup {
			return nil, helper.NewServiceError("DUPLICATE_TITLE", "Judul kegiatan sudah terpakai karena menghasilkan slug yang sama dengan kegiatan lain. Gunakan judul lain.", nil)
		}
	}

	mapper.UpdateReqToEntity(req, k)

	if err := s.repo.Update(k); err != nil {
		return nil, helper.NewServiceError("SERVER_ERROR", "Gagal mengupdate kegiatan.", err)
	}

	if req.Tags != "" {
		if err := s.repo.SaveTags(k.ID, mapper.ParseTags(req.Tags)); err != nil {
			return nil, helper.NewServiceError("SERVER_ERROR", "Gagal menyimpan tag kegiatan.", err)
		}
	}
	if req.GalleryJSON != "" {
		items, err := mapper.ParseGalleryJSON(req.GalleryJSON)
		if err == nil && len(items) > 0 {
			if err := s.repo.SaveGallery(k.ID, items); err != nil {
				return nil, helper.NewServiceError("SERVER_ERROR", "Gagal menyimpan galeri kegiatan.", err)
			}
		}
	}

	s.log(ctx, s.db, &activitylogdto.ActivityLogInput{
		Action:      "kegiatan.update",
		EntityType:  "kegiatan",
		EntityID:    &k.ID,
		EntityLabel: k.Title,
		Description: "Memperbarui kegiatan: " + k.Title,
		Metadata: map[string]any{
			"slug": k.Slug,
		},
	})

	updated, err := s.repo.FindByID(k.ID)
	if err == nil {
		k = updated
	}
	resp := mapper.EntityToDetail(k)
	return &resp, nil
}

// UploadImage menyimpan file gambar (cover/galeri kegiatan) ke public/uploads/kegiatan.
// Mengembalikan path relatif yang bisa langsung dipakai di image_path / gallery.
func (s *kegiatanService) UploadImage(ctx context.Context, file *multipart.FileHeader) (*dto.UploadImageResponse, error) {
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

	imagePath := "/uploads/kegiatan/" + filename
	return &dto.UploadImageResponse{ImagePath: imagePath}, nil
}

func (s *kegiatanService) Delete(ctx context.Context, id uint) error {
	k, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return helper.NewNotFoundError("Kegiatan tidak ditemukan")
		}
		return helper.NewServiceError("SERVER_ERROR", "Gagal mengambil kegiatan.", err)
	}
	if err := s.repo.SoftDelete(id); err != nil {
		return helper.NewServiceError("SERVER_ERROR", "Gagal menghapus kegiatan.", err)
	}

	s.log(ctx, s.db, &activitylogdto.ActivityLogInput{
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

	s.log(ctx, s.db, &activitylogdto.ActivityLogInput{
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

	s.log(ctx, s.db, &activitylogdto.ActivityLogInput{
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

	s.log(ctx, s.db, &activitylogdto.ActivityLogInput{
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
	if err := s.repo.DeleteGalleryImage(galleryID); err != nil {
		return helper.NewServiceError("SERVER_ERROR", "Gagal menghapus gambar galeri.", err)
	}
	return nil
}

// GetCategories mengembalikan daftar kategori unik untuk dropdown admin.
func (s *kegiatanService) GetCategories(ctx context.Context) ([]string, error) {
	categories, err := s.repo.GetCategories()
	if err != nil {
		return nil, helper.NewServiceError("SERVER_ERROR", "Gagal mengambil daftar kategori.", err)
	}
	return categories, nil
}
