package service

import (
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/helper"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/sliders/dto"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/sliders/model"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/sliders/repository"
)

type SlidersService interface {
	GetAll(activeOnly bool) (*dto.SliderListResponse, error)
	GetByID(id uint) (*dto.SliderResponse, error)
	Create(req *dto.SliderRequest, createdBy uint) (*dto.SliderResponse, error)
	Update(id uint, req *dto.SliderRequest) (*dto.SliderResponse, error)
	Delete(id uint) error
	Restore(id uint) error
	BulkDelete(ids []uint) error
	BulkRestore(ids []uint) error
	Reorder(ids []uint) error
}

type slidersService struct {
	repo repository.SlidersRepo
}

func NewSlidersService(repo repository.SlidersRepo) SlidersService {
	return &slidersService{repo: repo}
}

func (s *slidersService) GetAll(activeOnly bool) (*dto.SliderListResponse, error) {
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

func (s *slidersService) GetByID(id uint) (*dto.SliderResponse, error) {
	sl, err := s.repo.FindByID(id)
	if err != nil {
		return nil, helper.NewServiceError("NOT_FOUND", "Slider tidak ditemukan.", err)
	}
	resp := toResponse(*sl)
	return &resp, nil
}

func (s *slidersService) Create(req *dto.SliderRequest, createdBy uint) (*dto.SliderResponse, error) {
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

	resp := toResponse(*sl)
	return &resp, nil
}

func (s *slidersService) Update(id uint, req *dto.SliderRequest) (*dto.SliderResponse, error) {
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

	resp := toResponse(*sl)
	return &resp, nil
}

func (s *slidersService) Delete(id uint) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return helper.NewServiceError("NOT_FOUND", "Slider tidak ditemukan.", err)
	}
	return s.repo.Delete(id)
}

func (s *slidersService) Reorder(ids []uint) error {
	if err := s.repo.Reorder(ids); err != nil {
		return helper.NewServiceError("SERVER_ERROR", "Gagal mengubah urutan slider.", err)
	}
	return nil
}

func (s *slidersService) Restore(id uint) error {
	return s.repo.Restore(id)
}

func (s *slidersService) BulkDelete(ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	return s.repo.BulkSoftDelete(ids)
}

func (s *slidersService) BulkRestore(ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	return s.repo.BulkRestore(ids)
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
