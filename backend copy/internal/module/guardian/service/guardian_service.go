package service

import (
	"context"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/helper"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/guardian/dto"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/guardian/entity"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/guardian/repository"
)

type GuardianService interface {
	GetAll(ctx context.Context, req dto.GuardianQueryReq) ([]entity.Guardian, int64, error)
	GetByID(ctx context.Context, id uint) (*entity.Guardian, error)
	Create(ctx context.Context, req dto.GuardianCreateReq) (*entity.Guardian, error)
	Update(ctx context.Context, id uint, req dto.GuardianUpdateReq) (*entity.Guardian, error)
	Delete(ctx context.Context, id uint) error
}

type guardianService struct {
	repo repository.GuardianRepo
}

func NewGuardianService(repo repository.GuardianRepo) GuardianService {
	return &guardianService{repo: repo}
}

func (s *guardianService) GetAll(ctx context.Context, req dto.GuardianQueryReq) ([]entity.Guardian, int64, error) {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 || req.Limit > 100 {
		req.Limit = 10
	}
	return s.repo.FindAll(ctx, req.Page, req.Limit, req.Search)
}

func (s *guardianService) GetByID(ctx context.Context, id uint) (*entity.Guardian, error) {
	g, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, helper.NewNotFoundError("wali tidak ditemukan")
	}
	return g, nil
}

func (s *guardianService) Create(ctx context.Context, req dto.GuardianCreateReq) (*entity.Guardian, error) {
	e := &entity.Guardian{
		Name:        req.Name,
		Phone:       req.Phone,
		Email:       req.Email,
		NIK:         req.NIK,
		Education:   req.Education,
		Occupation:  req.Occupation,
		IncomeRange: req.IncomeRange,
	}
	if err := s.repo.Create(ctx, e); err != nil {
		return nil, err
	}
	return e, nil
}

func (s *guardianService) Update(ctx context.Context, id uint, req dto.GuardianUpdateReq) (*entity.Guardian, error) {
	g, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, helper.NewNotFoundError("wali tidak ditemukan")
	}
	g.Name = req.Name
	g.Phone = req.Phone
	g.Email = req.Email
	g.NIK = req.NIK
	g.Education = req.Education
	g.Occupation = req.Occupation
	g.IncomeRange = req.IncomeRange
	if err := s.repo.Update(ctx, g); err != nil {
		return nil, err
	}
	return g, nil
}

func (s *guardianService) Delete(ctx context.Context, id uint) error {
	_, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return helper.NewNotFoundError("wali tidak ditemukan")
	}
	return s.repo.Delete(ctx, id)
}
