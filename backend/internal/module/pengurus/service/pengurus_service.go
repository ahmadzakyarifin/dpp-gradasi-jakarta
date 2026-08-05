package service

import (
	"context"
	"errors"
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
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/pengurus/dto"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/pengurus/mapper"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/pengurus/repository"
	"gorm.io/gorm"
)

type PengurusService interface {
	GetAllPublic(ctx context.Context, query dto.PengurusQuery) (*dto.PengurusListResponse, error)
	GetAllAdmin(ctx context.Context, query dto.PengurusQuery) (*dto.PengurusListResponse, error)
	GetRegions(ctx context.Context) (*dto.RegionsResponse, error)
	GetByID(ctx context.Context, id uint) (*dto.PengurusResponse, error)
	Create(ctx context.Context, req *dto.PengurusRequest) (*dto.PengurusResponse, error)
	Update(ctx context.Context, id uint, req *dto.PengurusRequest) (*dto.PengurusResponse, error)
	Delete(ctx context.Context, id uint) error
	Restore(ctx context.Context, id uint) error
	BulkDelete(ctx context.Context, ids []uint) error
	BulkRestore(ctx context.Context, ids []uint) error
	Reorder(ctx context.Context, ids []uint) error
}

type pengurusService struct {
	db         *gorm.DB
	repo       repository.PengurusRepo
	audit      activitylogservice.ActivityLogService
	uploadPath string
}

func NewPengurusService(db *gorm.DB, repo repository.PengurusRepo, audit activitylogservice.ActivityLogService) PengurusService {
	// Create upload dir if not exists
	uploadPath := "public/uploads/pengurus"
	_ = os.MkdirAll(uploadPath, 0755)
	return &pengurusService{db: db, repo: repo, audit: audit, uploadPath: uploadPath}
}

func (s *pengurusService) log(ctx context.Context, db *gorm.DB, input *activitylogdto.ActivityLogInput) {
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

func (s *pengurusService) GetAllPublic(ctx context.Context, query dto.PengurusQuery) (*dto.PengurusListResponse, error) {
	return s.list(ctx, query, false)
}

func (s *pengurusService) GetAllAdmin(ctx context.Context, query dto.PengurusQuery) (*dto.PengurusListResponse, error) {
	return s.list(ctx, query, true)
}

func (s *pengurusService) list(ctx context.Context, q dto.PengurusQuery, adminMode bool) (*dto.PengurusListResponse, error) {
	results, total, err := s.repo.FindAll(q, adminMode)
	if err != nil {
		return nil, helper.NewServiceError("SERVER_ERROR", "Gagal mengambil data pengurus.", err)
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

	resp := &dto.PengurusListResponse{
		Data: mapper.EntityListToResponse(results),
		Meta: dto.PaginationMeta{
			CurrentPage: page,
			Limit:       limit,
			TotalData:   int(total),
			TotalPages:  totalPages,
		},
	}
	return resp, nil
}

func (s *pengurusService) GetRegions(ctx context.Context) (*dto.RegionsResponse, error) {
	raw, err := s.repo.GetRegions()
	if err != nil {
		return nil, helper.NewServiceError("SERVER_ERROR", "Gagal mengambil data wilayah.", err)
	}

	resp := &dto.RegionsResponse{
		Provinsi:  make([]string, 0),
		Kabupaten: make(map[string][]string),
	}

	provMap := make(map[string]bool)
	kabMap := make(map[string]map[string]bool)

	for _, r := range raw {
		if r.Provinsi != nil && *r.Provinsi != "" {
			prov := *r.Provinsi
			if !provMap[prov] {
				provMap[prov] = true
				resp.Provinsi = append(resp.Provinsi, prov)
			}

			if r.Kabupaten != nil && *r.Kabupaten != "" {
				kab := *r.Kabupaten
				if kabMap[prov] == nil {
					kabMap[prov] = make(map[string]bool)
				}
				if !kabMap[prov][kab] {
					kabMap[prov][kab] = true
					resp.Kabupaten[prov] = append(resp.Kabupaten[prov], kab)
				}
			}
		}
	}
	return resp, nil
}

func (s *pengurusService) GetByID(ctx context.Context, id uint) (*dto.PengurusResponse, error) {
	p, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helper.NewNotFoundError("Pengurus tidak ditemukan")
		}
		return nil, helper.NewServiceError("SERVER_ERROR", "Gagal mengambil data pengurus.", err)
	}
	resp := mapper.EntityToResponse(p)
	return &resp, nil
}

func (s *pengurusService) Create(ctx context.Context, req *dto.PengurusRequest) (*dto.PengurusResponse, error) {
	// Validasi wilayah sesuai level (contract: required_if)
	if verr := req.ValidateRegionRules(); len(verr) > 0 {
		v := helper.NewValidationError()
		for field, msg := range verr {
			v.Add(field, msg)
		}
		return nil, v
	}

	imageURL, err := s.handleUpload(req.Image)
	if err != nil {
		return nil, helper.NewServiceError("UPLOAD_FAILED", "Gagal mengunggah foto.", err)
	}

	if req.Image == nil && imageURL == "" {
		return nil, helper.NewServiceError("VALIDATION_ERROR", "Foto (image) wajib diunggah saat menambah pengurus", nil)
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	p := mapper.CreateReqToEntity(req, imageURL, isActive)

	if err := s.repo.Create(p); err != nil {
		return nil, helper.NewServiceError("SERVER_ERROR", "Gagal membuat data pengurus.", err)
	}

	s.log(ctx, s.db, &activitylogdto.ActivityLogInput{
		Action:      "pengurus.create",
		EntityType:  "pengurus",
		EntityID:    &p.ID,
		EntityLabel: p.Name,
		Description: "Menambahkan pengurus baru: " + p.Name,
		Metadata: map[string]any{
			"role":    p.Role,
			"periode": p.Periode,
		},
	})

	resp := mapper.EntityToResponse(p)
	return &resp, nil
}

func (s *pengurusService) Update(ctx context.Context, id uint, req *dto.PengurusRequest) (*dto.PengurusResponse, error) {
	// Validasi wilayah sesuai level (contract: required_if)
	if verr := req.ValidateRegionRules(); len(verr) > 0 {
		v := helper.NewValidationError()
		for field, msg := range verr {
			v.Add(field, msg)
		}
		return nil, v
	}

	p, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helper.NewNotFoundError("Pengurus tidak ditemukan")
		}
		return nil, helper.NewServiceError("SERVER_ERROR", "Gagal mengambil data pengurus.", err)
	}

	if req.Image != nil {
		newImg, err := s.handleUpload(req.Image)
		if err == nil && newImg != "" {
			p.ImagePath = newImg
		}
	}

	mapper.ApplyUpdate(p, req)

	if err := s.repo.Update(p); err != nil {
		return nil, helper.NewServiceError("SERVER_ERROR", "Gagal mengupdate data pengurus.", err)
	}

	s.log(ctx, s.db, &activitylogdto.ActivityLogInput{
		Action:      "pengurus.update",
		EntityType:  "pengurus",
		EntityID:    &p.ID,
		EntityLabel: p.Name,
		Description: "Memperbarui pengurus: " + p.Name,
		Metadata: map[string]any{
			"role":    p.Role,
			"periode": p.Periode,
		},
	})

	resp := mapper.EntityToResponse(p)
	return &resp, nil
}

func (s *pengurusService) Delete(ctx context.Context, id uint) error {
	p, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return helper.NewNotFoundError("Pengurus tidak ditemukan")
		}
		return helper.NewServiceError("SERVER_ERROR", "Gagal mengambil data pengurus.", err)
	}

	if err := s.repo.SoftDelete(id); err != nil {
		return helper.NewServiceError("SERVER_ERROR", "Gagal menghapus data pengurus.", err)
	}

	s.log(ctx, s.db, &activitylogdto.ActivityLogInput{
		Action:      "pengurus.delete",
		EntityType:  "pengurus",
		EntityID:    &id,
		EntityLabel: p.Name,
		Description: "Menghapus pengurus (soft delete): " + p.Name,
	})

	return nil
}

func (s *pengurusService) Restore(ctx context.Context, id uint) error {
	if err := s.repo.Restore(id); err != nil {
		return helper.NewServiceError("SERVER_ERROR", "Gagal memulihkan data pengurus.", err)
	}

	s.log(ctx, s.db, &activitylogdto.ActivityLogInput{
		Action:      "pengurus.restore",
		EntityType:  "pengurus",
		EntityID:    &id,
		Description: "Memulihkan pengurus (ID: " + strconv.FormatUint(uint64(id), 10) + ")",
	})

	return nil
}

func (s *pengurusService) BulkDelete(ctx context.Context, ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	if err := s.repo.BulkSoftDelete(ids); err != nil {
		return helper.NewServiceError("SERVER_ERROR", "Gagal menghapus data pengurus.", err)
	}

	s.log(ctx, s.db, &activitylogdto.ActivityLogInput{
		Action:      "pengurus.bulk_delete",
		EntityType:  "pengurus",
		Description: "Menghapus " + strconv.FormatUint(uint64(len(ids)), 10) + " pengurus (soft delete)",
		Metadata: map[string]any{
			"ids": ids,
		},
	})

	return nil
}

func (s *pengurusService) BulkRestore(ctx context.Context, ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	if err := s.repo.BulkRestore(ids); err != nil {
		return helper.NewServiceError("SERVER_ERROR", "Gagal memulihkan data pengurus.", err)
	}

	s.log(ctx, s.db, &activitylogdto.ActivityLogInput{
		Action:      "pengurus.bulk_restore",
		EntityType:  "pengurus",
		Description: "Memulihkan " + strconv.FormatUint(uint64(len(ids)), 10) + " pengurus",
		Metadata: map[string]any{
			"ids": ids,
		},
	})

	return nil
}

func (s *pengurusService) Reorder(ctx context.Context, ids []uint) error {
	if err := s.repo.Reorder(ids); err != nil {
		return helper.NewServiceError("SERVER_ERROR", "Gagal mengubah urutan pengurus.", err)
	}

	s.log(ctx, s.db, &activitylogdto.ActivityLogInput{
		Action:      "pengurus.update",
		EntityType:  "pengurus",
		Description: "Mengubah urutan (reorder) pengurus",
		Metadata: map[string]any{
			"ids": ids,
		},
	})

	return nil
}

func (s *pengurusService) handleUpload(file *multipart.FileHeader) (string, error) {
	if file == nil {
		return "", nil
	}

	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)

	// Pastikan direktori upload ada
	if err := os.MkdirAll(s.uploadPath, 0o755); err != nil {
		return "", err
	}

	dst := filepath.Join(s.uploadPath, filename)

	out, err := os.Create(dst)
	if err != nil {
		return "", err
	}
	defer out.Close()

	if _, err = io.Copy(out, src); err != nil {
		return "", err
	}

	// URL path to be saved in DB
	return "/uploads/pengurus/" + filename, nil
}
