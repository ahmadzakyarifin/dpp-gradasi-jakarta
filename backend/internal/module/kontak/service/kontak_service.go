package service

import (
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/helper"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/kontak/dto"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/kontak/model"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/kontak/repository"
)

type KontakService interface {
	GetAll(q dto.KontakQuery) (*dto.KontakListResponse, error)
	GetByID(id uint) (*dto.KontakDetailResponse, error)
	Submit(req *dto.KontakRequest) error
	Delete(id uint) error
	Restore(id uint) error
	BulkDelete(ids []uint) error
	BulkRestore(ids []uint) error
}

type kontakService struct {
	repo repository.KontakRepo
}

func NewKontakService(repo repository.KontakRepo) KontakService {
	return &kontakService{repo: repo}
}

func (s *kontakService) GetAll(q dto.KontakQuery) (*dto.KontakListResponse, error) {
	items, total, err := s.repo.FindAll(q)
	if err != nil {
		return nil, helper.NewServiceError("SERVER_ERROR", "Gagal mengambil data pesan.", err)
	}

	page := maxInt(q.Page, 1)
	limit := maxInt(q.Limit, 10)
	totalPages := (int(total) + limit - 1) / limit

	resp := &dto.KontakListResponse{
		Kontak: make([]dto.KontakListItem, 0),
		Meta: dto.PaginationMeta{
			CurrentPage: page,
			Limit:       limit,
			TotalData:   int(total),
			TotalPages:  totalPages,
		},
	}
	for _, p := range items {
		resp.Kontak = append(resp.Kontak, toListItem(p))
	}
	return resp, nil
}

func (s *kontakService) GetByID(id uint) (*dto.KontakDetailResponse, error) {
	p, err := s.repo.FindByID(id)
	if err != nil {
		return nil, helper.NewServiceError("NOT_FOUND", "Pesan tidak ditemukan.", err)
	}

	// Auto mark as read
	if !p.IsRead {
		s.repo.MarkAsRead(id)
		p.IsRead = true
	}

	resp := toDetail(*p)
	return &resp, nil
}

func (s *kontakService) Submit(req *dto.KontakRequest) error {
	p := &model.PesanKontak{
		Nama:   req.Nama,
		Email:  req.Email,
		Subjek: req.Subjek,
		Pesan:  req.Pesan,
	}
	return s.repo.Create(p)
}

func (s *kontakService) Delete(id uint) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return helper.NewServiceError("NOT_FOUND", "Pesan tidak ditemukan.", err)
	}
	return s.repo.Delete(id)
}

func (s *kontakService) Restore(id uint) error {
	return s.repo.Restore(id)
}

func (s *kontakService) BulkDelete(ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	return s.repo.BulkSoftDelete(ids)
}

func (s *kontakService) BulkRestore(ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	return s.repo.BulkRestore(ids)
}

func toListItem(p model.PesanKontak) dto.KontakListItem {
	return dto.KontakListItem{
		ID:        p.ID,
		Nama:      p.Nama,
		Email:     p.Email,
		Subjek:    p.Subjek,
		IsRead:    p.IsRead,
		CreatedAt: p.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func toDetail(p model.PesanKontak) dto.KontakDetailResponse {
	resp := dto.KontakDetailResponse{
		ID:        p.ID,
		Nama:      p.Nama,
		Email:     p.Email,
		Subjek:    p.Subjek,
		Pesan:     p.Pesan,
		IsRead:    p.IsRead,
		CreatedAt: p.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: p.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
	if p.ResponseNote != nil {
		resp.ResponseNote = *p.ResponseNote
	}
	return resp
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
