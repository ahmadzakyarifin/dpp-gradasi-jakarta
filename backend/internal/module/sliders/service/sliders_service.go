package service

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/helper"
	activitylogdto "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/activitylog/dto"
	activitylogservice "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/activitylog/service"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/sliders/dto"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/sliders/entity"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/sliders/mapper"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/sliders/repository"
	"gorm.io/gorm"
)

type SlidersService interface {
	GetAll(ctx context.Context, publishedOnly bool) (*dto.SliderListResponse, error)
	GetByID(ctx context.Context, id uint) (*dto.SliderResponse, error)
	Create(ctx context.Context, req *dto.SliderRequest, createdBy uint) (*dto.SliderResponse, error)
	Update(ctx context.Context, id uint, req *dto.SliderRequest) (*dto.SliderResponse, error)
	Delete(ctx context.Context, id uint) error
	Restore(ctx context.Context, id uint) error
	BulkDelete(ctx context.Context, ids []uint) error
	BulkRestore(ctx context.Context, ids []uint) error
	Reorder(ctx context.Context, ids []uint) error
	UploadImage(ctx context.Context, fileHeader *multipart.FileHeader) (*dto.UploadImageResponse, error)
}

type slidersService struct {
	db    *gorm.DB
	repo  repository.SlidersRepo
	audit activitylogservice.ActivityLogService
}

func NewSlidersService(db *gorm.DB, repo repository.SlidersRepo, audit activitylogservice.ActivityLogService) SlidersService {
	return &slidersService{db: db, repo: repo, audit: audit}
}

func (s *slidersService) log(ctx context.Context, db *gorm.DB, input *activitylogdto.ActivityLogInput) {
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

func (s *slidersService) GetAll(ctx context.Context, publishedOnly bool) (*dto.SliderListResponse, error) {
	sliders, err := s.repo.FindAll(publishedOnly)
	if err != nil {
		return nil, helper.NewServiceError("SERVER_ERROR", "Gagal mengambil data slider.", err)
	}

	entities := make([]entity.Slider, 0, len(sliders))
	for i := range sliders {
		entities = append(entities, *mapper.ModelToEntity(&sliders[i]))
	}

	resp := &dto.SliderListResponse{
		Sliders: mapper.EntityListToResponse(entities),
		Total:   int64(len(entities)),
	}
	return resp, nil
}

func (s *slidersService) GetByID(ctx context.Context, id uint) (*dto.SliderResponse, error) {
	sl, err := s.repo.FindByID(id)
	if err != nil {
		return nil, helper.NewNotFoundError("Slider tidak ditemukan")
	}
	resp := mapper.EntityToResponse(mapper.ModelToEntity(sl))
	return &resp, nil
}

func (s *slidersService) Create(ctx context.Context, req *dto.SliderRequest, createdBy uint) (*dto.SliderResponse, error) {
	sl := mapper.CreateReqToEntity(req)
	sl.CreatedBy = &createdBy

	if err := s.repo.Create(mapper.EntityToModel(sl)); err != nil {
		return nil, helper.NewServiceError("SERVER_ERROR", "Gagal membuat slider.", err)
	}

	s.log(ctx, s.db, &activitylogdto.ActivityLogInput{
		Action:      "slider.create",
		EntityType:  "slider",
		EntityID:    &sl.ID,
		EntityLabel: sl.Title,
		Description: "Membuat slider baru: " + sl.Title,
		Metadata: map[string]any{
			"title": sl.Title,
		},
	})

	resp := mapper.EntityToResponse(sl)
	return &resp, nil
}

func (s *slidersService) Update(ctx context.Context, id uint, req *dto.SliderRequest) (*dto.SliderResponse, error) {
	sl, err := s.repo.FindByID(id)
	if err != nil {
		return nil, helper.NewNotFoundError("Slider tidak ditemukan")
	}

	e := mapper.ModelToEntity(sl)
	mapper.UpdateReqToEntity(req, e)

	if err := s.repo.Update(mapper.EntityToModel(e)); err != nil {
		return nil, helper.NewServiceError("SERVER_ERROR", "Gagal mengupdate slider.", err)
	}

	s.log(ctx, s.db, &activitylogdto.ActivityLogInput{
		Action:      "slider.update",
		EntityType:  "slider",
		EntityID:    &e.ID,
		EntityLabel: e.Title,
		Description: "Memperbarui slider: " + e.Title,
		Metadata: map[string]any{
			"title": e.Title,
		},
	})

	resp := mapper.EntityToResponse(e)
	return &resp, nil
}

func (s *slidersService) Delete(ctx context.Context, id uint) error {
	sl, err := s.repo.FindByID(id)
	if err != nil {
		return helper.NewNotFoundError("Slider tidak ditemukan")
	}
	if err := s.repo.Delete(id); err != nil {
		return helper.NewServiceError("SERVER_ERROR", "Gagal menghapus slider.", err)
	}

	s.log(ctx, s.db, &activitylogdto.ActivityLogInput{
		Action:      "slider.delete",
		EntityType:  "slider",
		EntityID:    &id,
		EntityLabel: sl.Title,
		Description: "Menghapus slider (soft delete): " + sl.Title,
	})

	return nil
}

func (s *slidersService) Reorder(ctx context.Context, ids []uint) error {
	if err := s.repo.Reorder(ids); err != nil {
		return helper.NewServiceError("SERVER_ERROR", "Gagal mengubah urutan slider.", err)
	}

	s.log(ctx, s.db, &activitylogdto.ActivityLogInput{
		Action:      "slider.update",
		EntityType:  "slider",
		Description: "Mengubah urutan (reorder) slider",
		Metadata: map[string]any{
			"ids": ids,
		},
	})

	return nil
}

func (s *slidersService) Restore(ctx context.Context, id uint) error {
	if err := s.repo.Restore(id); err != nil {
		return helper.NewServiceError("SERVER_ERROR", "Gagal memulihkan slider.", err)
	}

	s.log(ctx, s.db, &activitylogdto.ActivityLogInput{
		Action:      "slider.restore",
		EntityType:  "slider",
		EntityID:    &id,
		Description: "Memulihkan slider (ID: " + strconv.FormatUint(uint64(id), 10) + ")",
	})

	return nil
}

func (s *slidersService) BulkDelete(ctx context.Context, ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	if err := s.repo.BulkSoftDelete(ids); err != nil {
		return helper.NewServiceError("SERVER_ERROR", "Gagal menghapus slider.", err)
	}

	s.log(ctx, s.db, &activitylogdto.ActivityLogInput{
		Action:      "slider.bulk_delete",
		EntityType:  "slider",
		Description: "Menghapus " + strconv.FormatUint(uint64(len(ids)), 10) + " slider (soft delete)",
		Metadata: map[string]any{
			"ids": ids,
		},
	})

	return nil
}

func (s *slidersService) BulkRestore(ctx context.Context, ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	if err := s.repo.BulkRestore(ids); err != nil {
		return helper.NewServiceError("SERVER_ERROR", "Gagal memulihkan slider.", err)
	}

	s.log(ctx, s.db, &activitylogdto.ActivityLogInput{
		Action:      "slider.bulk_restore",
		EntityType:  "slider",
		Description: "Memulihkan " + strconv.FormatUint(uint64(len(ids)), 10) + " slider",
		Metadata: map[string]any{
			"ids": ids,
		},
	})

	return nil
}

func (s *slidersService) UploadImage(ctx context.Context, fileHeader *multipart.FileHeader) (*dto.UploadImageResponse, error) {
	// Pastikan direktori upload ada
	uploadPath := "public/uploads/img/slider"
	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		return nil, helper.NewServiceError("SERVER_ERROR", "Gagal membuat direktori upload sliders.", err)
	}

	// Generate filename unik
	ext := filepath.Ext(fileHeader.Filename)
	filename := fmt.Sprintf("slider-%d%s", time.Now().UnixNano(), ext)
	dst := filepath.Join(uploadPath, filename)

	// Simpan file
	src, err := fileHeader.Open()
	if err != nil {
		return nil, helper.NewServiceError("SERVER_ERROR", "Gagal membuka file terunggah", err)
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return nil, helper.NewServiceError("SERVER_ERROR", "Gagal membuat file baru di server", err)
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	if err != nil {
		return nil, helper.NewServiceError("SERVER_ERROR", "Gagal menulis file", err)
	}

	imagePath := "/uploads/img/slider/" + filename
	return &dto.UploadImageResponse{ImagePath: imagePath}, nil
}
