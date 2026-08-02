package service

import (
	"context"
	"strconv"

	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/helper"
	activitylogdto "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/activitylog/dto"
	activitylogservice "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/activitylog/service"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/kontak/dto"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/kontak/model"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/kontak/repository"
	"gorm.io/gorm"
)

type KontakService interface {
	GetAll(ctx context.Context, q dto.KontakQuery) (*dto.KontakListResponse, error)
	GetByID(ctx context.Context, id uint) (*dto.KontakDetailResponse, error)
	Submit(ctx context.Context, req *dto.KontakRequest) error
	Delete(ctx context.Context, id uint) error
	Restore(ctx context.Context, id uint) error
	BulkDelete(ctx context.Context, ids []uint) error
	BulkRestore(ctx context.Context, ids []uint) error
}

type kontakService struct {
	db    *gorm.DB
	repo  repository.KontakRepo
	audit activitylogservice.ActivityLogService
}

func NewKontakService(db *gorm.DB, repo repository.KontakRepo, audit activitylogservice.ActivityLogService) KontakService {
	return &kontakService{db: db, repo: repo, audit: audit}
}

func (s *kontakService) log(ctx context.Context, input *activitylogdto.ActivityLogInput) {
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

func (s *kontakService) GetAll(ctx context.Context, q dto.KontakQuery) (*dto.KontakListResponse, error) {
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

func (s *kontakService) GetByID(ctx context.Context, id uint) (*dto.KontakDetailResponse, error) {
	p, err := s.repo.FindByID(id)
	if err != nil {
		return nil, helper.NewServiceError("NOT_FOUND", "Pesan tidak ditemukan.", err)
	}

	// Auto mark as read
	if !p.IsRead {
		if err := s.repo.MarkAsRead(id); err != nil {
			return nil, helper.NewServiceError("SERVER_ERROR", "Gagal menandai pesan sebagai dibaca.", err)
		}
		p.IsRead = true
	}

	resp := toDetail(*p)
	return &resp, nil
}

func (s *kontakService) Submit(ctx context.Context, req *dto.KontakRequest) error {
	p := &model.PesanKontak{
		Nama:   req.Nama,
		Email:  req.Email,
		Subjek: req.Subjek,
		Pesan:  req.Pesan,
	}
	return s.repo.Create(p)
}

func (s *kontakService) Delete(ctx context.Context, id uint) error {
	p, err := s.repo.FindByID(id)
	if err != nil {
		return helper.NewServiceError("NOT_FOUND", "Pesan tidak ditemukan.", err)
	}
	if err := s.repo.Delete(id); err != nil {
		return helper.NewServiceError("SERVER_ERROR", "Gagal menghapus pesan.", err)
	}

	s.log(ctx, &activitylogdto.ActivityLogInput{
		Action:      "kontak.delete",
		EntityType:  "kontak",
		EntityID:    &id,
		EntityLabel: p.Nama,
		Description: "Menghapus pesan kontak: " + p.Nama,
	})

	return nil
}

func (s *kontakService) Restore(ctx context.Context, id uint) error {
	if err := s.repo.Restore(id); err != nil {
		return helper.NewServiceError("SERVER_ERROR", "Gagal memulihkan pesan.", err)
	}

	s.log(ctx, &activitylogdto.ActivityLogInput{
		Action:      "kontak.restore",
		EntityType:  "kontak",
		EntityID:    &id,
		Description: "Memulihkan pesan kontak (ID: " + strconv.FormatUint(uint64(id), 10) + ")",
	})

	return nil
}

func (s *kontakService) BulkDelete(ctx context.Context, ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	if err := s.repo.BulkSoftDelete(ids); err != nil {
		return helper.NewServiceError("SERVER_ERROR", "Gagal menghapus pesan.", err)
	}

	s.log(ctx, &activitylogdto.ActivityLogInput{
		Action:      "kontak.bulk_delete",
		EntityType:  "kontak",
		Description: "Menghapus " + strconv.FormatUint(uint64(len(ids)), 10) + " pesan kontak (soft delete)",
		Metadata: map[string]any{
			"ids": ids,
		},
	})

	return nil
}

func (s *kontakService) BulkRestore(ctx context.Context, ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	if err := s.repo.BulkRestore(ids); err != nil {
		return helper.NewServiceError("SERVER_ERROR", "Gagal memulihkan pesan.", err)
	}

	s.log(ctx, &activitylogdto.ActivityLogInput{
		Action:      "kontak.bulk_restore",
		EntityType:  "kontak",
		Description: "Memulihkan " + strconv.FormatUint(uint64(len(ids)), 10) + " pesan kontak",
		Metadata: map[string]any{
			"ids": ids,
		},
	})

	return nil
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
