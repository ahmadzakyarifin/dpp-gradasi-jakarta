package service

import (
	"context"
	"errors"
	"strconv"

	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/helper"
	activitylogdto "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/activitylog/dto"
	activitylogservice "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/activitylog/service"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/kontak/dto"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/kontak/mapper"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/kontak/repository"
	"gorm.io/gorm"
)

type KontakService interface {
	GetAll(ctx context.Context, q dto.KontakQuery) (*dto.KontakListResponse, error)
	GetByID(ctx context.Context, id uint) (*dto.KontakDetailResponse, error)
	Submit(ctx context.Context, req *dto.KontakRequest) error
	Delete(ctx context.Context, id uint) error
	Restore(ctx context.Context, id uint) error
	// BulkDelete & BulkRestore mengembalikan jumlah pesan yang benar-benar diproses.
	BulkDelete(ctx context.Context, ids []uint) (int64, error)
	BulkRestore(ctx context.Context, ids []uint) (int64, error)
}

type kontakService struct {
	db    *gorm.DB
	repo  repository.KontakRepo
	audit activitylogservice.ActivityLogService
}

func NewKontakService(db *gorm.DB, repo repository.KontakRepo, audit activitylogservice.ActivityLogService) KontakService {
	return &kontakService{db: db, repo: repo, audit: audit}
}

func (s *kontakService) log(ctx context.Context, db *gorm.DB, input *activitylogdto.ActivityLogInput) {
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

func (s *kontakService) GetAll(ctx context.Context, q dto.KontakQuery) (*dto.KontakListResponse, error) {
	items, total, err := s.repo.FindAll(q)
	if err != nil {
		return nil, helper.NewServiceError("SERVER_ERROR", "Gagal mengambil data pesan.", err)
	}

	// Pakai normalisasi yang sama dengan repository supaya meta selalu
	// mencerminkan halaman yang benar-benar di-query.
	page, limit, _ := q.Pagination()

	totalPages := 0
	if total > 0 {
		totalPages = (int(total) + limit - 1) / limit
	}

	resp := &dto.KontakListResponse{
		Kontak: mapper.EntityListToItem(items),
		Meta: dto.PaginationMeta{
			CurrentPage: page,
			Limit:       limit,
			TotalData:   int(total),
			TotalPages:  totalPages,
		},
	}
	return resp, nil
}

func (s *kontakService) GetByID(ctx context.Context, id uint) (*dto.KontakDetailResponse, error) {
	// Unscoped: pesan di Sampah tetap bisa dibuka detailnya (read-only) dari tab Sampah.
	p, err := s.repo.FindAnyByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helper.NewNotFoundError("Pesan tidak ditemukan")
		}
		return nil, helper.NewServiceError("SERVER_ERROR", "Gagal mengambil pesan.", err)
	}

	// Auto mark as read — hanya untuk pesan yang masih di kotak masuk.
	// Pesan di Sampah tidak diubah statusnya saat dilihat.
	if !p.IsRead && p.DeletedAt == nil {
		if err := s.repo.MarkAsRead(id); err != nil {
			return nil, helper.NewServiceError("SERVER_ERROR", "Gagal menandai pesan sebagai dibaca.", err)
		}
		p.IsRead = true
	}

	resp := mapper.EntityToDetail(p)
	return &resp, nil
}

func (s *kontakService) Submit(ctx context.Context, req *dto.KontakRequest) error {
	p := mapper.CreateReqToEntity(req)
	return s.repo.Create(p)
}

func (s *kontakService) Delete(ctx context.Context, id uint) error {
	p, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return helper.NewNotFoundError("Pesan tidak ditemukan")
		}
		return helper.NewServiceError("SERVER_ERROR", "Gagal mengambil pesan.", err)
	}
	if err := s.repo.Delete(id); err != nil {
		return helper.NewServiceError("SERVER_ERROR", "Gagal menghapus pesan.", err)
	}

	s.log(ctx, s.db, &activitylogdto.ActivityLogInput{
		Action:      "kontak.delete",
		EntityType:  "kontak",
		EntityID:    &id,
		EntityLabel: p.Nama,
		Description: "Menghapus pesan kontak: " + p.Nama,
	})

	return nil
}

func (s *kontakService) Restore(ctx context.Context, id uint) error {
	p, err := s.repo.FindAnyByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return helper.NewNotFoundError("Pesan tidak ditemukan")
		}
		return helper.NewServiceError("SERVER_ERROR", "Gagal mengambil pesan.", err)
	}
	if p.DeletedAt == nil {
		return helper.NewServiceError("VALIDATION_ERROR", "Pesan tidak berada di Sampah.", nil)
	}

	if err := s.repo.Restore(id); err != nil {
		return helper.NewServiceError("SERVER_ERROR", "Gagal memulihkan pesan.", err)
	}

	s.log(ctx, s.db, &activitylogdto.ActivityLogInput{
		Action:      "kontak.restore",
		EntityType:  "kontak",
		EntityID:    &id,
		EntityLabel: p.Nama,
		Description: "Memulihkan pesan kontak: " + p.Nama + " (ID: " + strconv.FormatUint(uint64(id), 10) + ")",
	})

	return nil
}

func (s *kontakService) BulkDelete(ctx context.Context, ids []uint) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	deleted, err := s.repo.BulkSoftDelete(ids)
	if err != nil {
		return 0, helper.NewServiceError("SERVER_ERROR", "Gagal menghapus pesan.", err)
	}
	if deleted == 0 {
		return 0, helper.NewServiceError("VALIDATION_ERROR", "Tidak ada pesan yang bisa dihapus.", nil)
	}

	s.log(ctx, s.db, &activitylogdto.ActivityLogInput{
		Action:      "kontak.bulk_delete",
		EntityType:  "kontak",
		Description: "Menghapus " + strconv.FormatInt(deleted, 10) + " pesan kontak (soft delete)",
		Metadata: map[string]any{
			"ids":     ids,
			"deleted": deleted,
		},
	})

	return deleted, nil
}

func (s *kontakService) BulkRestore(ctx context.Context, ids []uint) (int64, error) {
	if len(ids) == 0 {
		return 0, nil
	}
	restored, err := s.repo.BulkRestore(ids)
	if err != nil {
		return 0, helper.NewServiceError("SERVER_ERROR", "Gagal memulihkan pesan.", err)
	}
	if restored == 0 {
		return 0, helper.NewServiceError("VALIDATION_ERROR", "Tidak ada pesan di Sampah yang bisa dipulihkan.", nil)
	}

	s.log(ctx, s.db, &activitylogdto.ActivityLogInput{
		Action:      "kontak.bulk_restore",
		EntityType:  "kontak",
		Description: "Memulihkan " + strconv.FormatInt(restored, 10) + " pesan kontak",
		Metadata: map[string]any{
			"ids":      ids,
			"restored": restored,
		},
	})

	return restored, nil
}
