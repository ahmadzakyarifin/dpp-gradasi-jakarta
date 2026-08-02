package service

import (
	"context"
	"strconv"

	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/helper"
	activitylogdto "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/activitylog/dto"
	activitylogservice "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/activitylog/service"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/sliders/dto"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/sliders/model"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/sliders/repository"
	"gorm.io/gorm"
)

type SlidersService interface {
	GetAll(ctx context.Context, activeOnly bool) (*dto.SliderListResponse, error)
	GetByID(ctx context.Context, id uint) (*dto.SliderResponse, error)
	Create(ctx context.Context, req *dto.SliderRequest, createdBy uint) (*dto.SliderResponse, error)
	Update(ctx context.Context, id uint, req *dto.SliderRequest) (*dto.SliderResponse, error)
	Delete(ctx context.Context, id uint) error
	Restore(ctx context.Context, id uint) error
	BulkDelete(ctx context.Context, ids []uint) error
	BulkRestore(ctx context.Context, ids []uint) error
	Reorder(ctx context.Context, ids []uint) error
}

type slidersService struct {
	db    *gorm.DB
	repo  repository.SlidersRepo
	audit activitylogservice.ActivityLogService
}

func NewSlidersService(db *gorm.DB, repo repository.SlidersRepo, audit activitylogservice.ActivityLogService) SlidersService {
	return &slidersService{db: db, repo: repo, audit: audit}
}

func (s *slidersService) log(ctx context.Context, input *activitylogdto.ActivityLogInput) {
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

func (s *slidersService) GetAll(ctx context.Context, activeOnly bool) (*dto.SliderListResponse, error) {
	sliders, err := s.repo.FindAll(activeOnly)
	if err != nil {
		return nil, helper.NewServiceError("SERVER_ERROR", "Gagal mengambil data slider.", err)
	}

	resp := &dto.SliderListResponse{}
	resp.Total = int64(len(sliders))
	for _, sl := range sliders {
		resp.Sliders = append(resp.Sliders, toResponse(sl))
	}
	return resp, nil
}

func (s *slidersService) GetByID(ctx context.Context, id uint) (*dto.SliderResponse, error) {
	sl, err := s.repo.FindByID(id)
	if err != nil {
		return nil, helper.NewServiceError("NOT_FOUND", "Slider tidak ditemukan.", err)
	}
	resp := toResponse(*sl)
	return &resp, nil
}

func (s *slidersService) Create(ctx context.Context, req *dto.SliderRequest, createdBy uint) (*dto.SliderResponse, error) {
	sl := &model.Slider{
		Title:     req.Title,
		Subtitle:  req.Subtitle,
		Tag:       req.Tag,
		IsNew:     req.IsNew,
		EventDate: req.EventDate,
		Location:  req.Location,
		ImagePath: req.ImagePath,
		LinkURL:   req.LinkURL,
		SortOrder: req.SortOrder,
		IsActive:  req.IsActive,
		CreatedBy: &createdBy,
	}

	if err := s.repo.Create(sl); err != nil {
		return nil, helper.NewServiceError("SERVER_ERROR", "Gagal membuat slider.", err)
	}

	s.log(ctx, &activitylogdto.ActivityLogInput{
		Action:      "slider.create",
		EntityType:  "slider",
		EntityID:    &sl.ID,
		EntityLabel: sl.Title,
		Description: "Membuat slider baru: " + sl.Title,
		Metadata: map[string]any{
			"title": sl.Title,
		},
	})

	resp := toResponse(*sl)
	return &resp, nil
}

func (s *slidersService) Update(ctx context.Context, id uint, req *dto.SliderRequest) (*dto.SliderResponse, error) {
	sl, err := s.repo.FindByID(id)
	if err != nil {
		return nil, helper.NewServiceError("NOT_FOUND", "Slider tidak ditemukan.", err)
	}

	sl.Title = req.Title
	sl.Subtitle = req.Subtitle
	sl.Tag = req.Tag
	sl.IsNew = req.IsNew
	sl.EventDate = req.EventDate
	sl.Location = req.Location
	sl.ImagePath = req.ImagePath
	sl.LinkURL = req.LinkURL
	sl.SortOrder = req.SortOrder
	sl.IsActive = req.IsActive

	if err := s.repo.Update(sl); err != nil {
		return nil, helper.NewServiceError("SERVER_ERROR", "Gagal mengupdate slider.", err)
	}

	s.log(ctx, &activitylogdto.ActivityLogInput{
		Action:      "slider.update",
		EntityType:  "slider",
		EntityID:    &sl.ID,
		EntityLabel: sl.Title,
		Description: "Memperbarui slider: " + sl.Title,
		Metadata: map[string]any{
			"title": sl.Title,
		},
	})

	resp := toResponse(*sl)
	return &resp, nil
}

func (s *slidersService) Delete(ctx context.Context, id uint) error {
	sl, err := s.repo.FindByID(id)
	if err != nil {
		return helper.NewServiceError("NOT_FOUND", "Slider tidak ditemukan.", err)
	}
	if err := s.repo.Delete(id); err != nil {
		return helper.NewServiceError("SERVER_ERROR", "Gagal menghapus slider.", err)
	}

	s.log(ctx, &activitylogdto.ActivityLogInput{
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

	s.log(ctx, &activitylogdto.ActivityLogInput{
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

	s.log(ctx, &activitylogdto.ActivityLogInput{
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

	s.log(ctx, &activitylogdto.ActivityLogInput{
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

	s.log(ctx, &activitylogdto.ActivityLogInput{
		Action:      "slider.bulk_restore",
		EntityType:  "slider",
		Description: "Memulihkan " + strconv.FormatUint(uint64(len(ids)), 10) + " slider",
		Metadata: map[string]any{
			"ids": ids,
		},
	})

	return nil
}

func toResponse(sl model.Slider) dto.SliderResponse {
	r := dto.SliderResponse{
		ID:        sl.ID,
		Title:     sl.Title,
		IsNew:     sl.IsNew,
		ImagePath: sl.ImagePath,
		SortOrder: sl.SortOrder,
		IsActive:  sl.IsActive,
	}
	if sl.Subtitle != nil {
		r.Subtitle = *sl.Subtitle
	}
	if sl.Tag != nil {
		r.Tag = *sl.Tag
	}
	if sl.EventDate != nil {
		r.EventDate = *sl.EventDate
	}
	if sl.Location != nil {
		r.Location = *sl.Location
	}
	if sl.LinkURL != nil {
		r.LinkURL = *sl.LinkURL
	}
	return r
}
