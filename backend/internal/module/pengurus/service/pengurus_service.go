package service

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/helper"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/pengurus/dto"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/pengurus/model"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/pengurus/repository"
)

type PengurusService interface {
	GetAllPublic(query dto.PengurusQuery) (*dto.PengurusListResponse, error)
	GetAllAdmin(query dto.PengurusQuery) (*dto.PengurusListResponse, error)
	GetRegions() (*dto.RegionsResponse, error)
	GetByID(id uint) (*dto.PengurusResponse, error)
	Create(req *dto.PengurusRequest) (*dto.PengurusResponse, error)
	Update(id uint, req *dto.PengurusRequest) (*dto.PengurusResponse, error)
	Delete(id uint) error
	Restore(id uint) error
	BulkDelete(ids []uint) error
	BulkRestore(ids []uint) error
}

type pengurusService struct {
	repo       repository.PengurusRepo
	uploadPath string
}

func NewPengurusService(repo repository.PengurusRepo) PengurusService {
	// Create upload dir if not exists
	uploadPath := "public/uploads/pengurus"
	_ = os.MkdirAll(uploadPath, 0755)
	return &pengurusService{repo: repo, uploadPath: uploadPath}
}

func (s *pengurusService) GetAllPublic(query dto.PengurusQuery) (*dto.PengurusListResponse, error) {
	return s.list(query, false)
}

func (s *pengurusService) GetAllAdmin(query dto.PengurusQuery) (*dto.PengurusListResponse, error) {
	return s.list(query, true)
}

func (s *pengurusService) list(q dto.PengurusQuery, adminMode bool) (*dto.PengurusListResponse, error) {
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
		Data: make([]dto.PengurusResponse, 0),
		Meta: dto.PaginationMeta{
			CurrentPage: page,
			Limit:       limit,
			TotalData:   int(total),
			TotalPages:  totalPages,
		},
	}

	for _, p := range results {
		resp.Data = append(resp.Data, toResponse(p))
	}
	return resp, nil
}

func (s *pengurusService) GetRegions() (*dto.RegionsResponse, error) {
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

func (s *pengurusService) GetByID(id uint) (*dto.PengurusResponse, error) {
	p, err := s.repo.FindByID(id)
	if err != nil {
		return nil, helper.NewServiceError("NOT_FOUND", "Pengurus tidak ditemukan.", err)
	}
	resp := toResponse(*p)
	return &resp, nil
}

func (s *pengurusService) Create(req *dto.PengurusRequest) (*dto.PengurusResponse, error) {
	// Validasi wilayah sesuai level (contract: required_if)
	if verr := req.ValidateRegionRules(); len(verr) > 0 {
		return nil, helper.NewServiceError("VALIDATION_ERROR", "validasi gagal", nil)
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

	p := &model.Pengurus{
		Name:         req.Name,
		Role:         req.Role,
		Department:   strPtr(req.Department),
		Level:        req.Level,
		Provinsi:     strPtr(req.Provinsi),
		Kabupaten:    strPtr(req.Kabupaten),
		ImagePath:    imageURL,
		FacebookURL:  strPtr(req.FacebookURL),
		InstagramURL: strPtr(req.InstagramURL),
		LinkedinURL:  strPtr(req.LinkedinURL),
		Whatsapp:     strPtr(req.Whatsapp),
		Periode:      req.Periode,
		SortOrder:    req.SortOrder,
		IsActive:     isActive,
	}

	if err := s.repo.Create(p); err != nil {
		return nil, helper.NewServiceError("SERVER_ERROR", "Gagal membuat data pengurus.", err)
	}

	resp := toResponse(*p)
	return &resp, nil
}

func (s *pengurusService) Update(id uint, req *dto.PengurusRequest) (*dto.PengurusResponse, error) {
	// Validasi wilayah sesuai level (contract: required_if)
	if verr := req.ValidateRegionRules(); len(verr) > 0 {
		return nil, helper.NewServiceError("VALIDATION_ERROR", "validasi gagal", nil)
	}

	p, err := s.repo.FindByID(id)
	if err != nil {
		return nil, helper.NewServiceError("NOT_FOUND", "Pengurus tidak ditemukan.", err)
	}

	if req.Image != nil {
		newImg, err := s.handleUpload(req.Image)
		if err == nil && newImg != "" {
			p.ImagePath = newImg
		}
	}

	p.Name = req.Name
	p.Role = req.Role
	p.Department = strPtr(req.Department)
	p.Level = req.Level
	p.Provinsi = strPtr(req.Provinsi)
	p.Kabupaten = strPtr(req.Kabupaten)
	p.FacebookURL = strPtr(req.FacebookURL)
	p.InstagramURL = strPtr(req.InstagramURL)
	p.LinkedinURL = strPtr(req.LinkedinURL)
	p.Whatsapp = strPtr(req.Whatsapp)
	p.Periode = req.Periode
	p.SortOrder = req.SortOrder

	if req.IsActive != nil {
		p.IsActive = *req.IsActive
	}

	if err := s.repo.Update(p); err != nil {
		return nil, helper.NewServiceError("SERVER_ERROR", "Gagal mengupdate data pengurus.", err)
	}

	resp := toResponse(*p)
	return &resp, nil
}

func (s *pengurusService) Delete(id uint) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return helper.NewServiceError("NOT_FOUND", "Pengurus tidak ditemukan.", err)
	}
	return s.repo.SoftDelete(id)
}

func (s *pengurusService) Restore(id uint) error {
	return s.repo.Restore(id)
}

func (s *pengurusService) BulkDelete(ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	return s.repo.BulkSoftDelete(ids)
}

func (s *pengurusService) BulkRestore(ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	return s.repo.BulkRestore(ids)
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

func toResponse(p model.Pengurus) dto.PengurusResponse {
	resp := dto.PengurusResponse{
		ID:        p.ID,
		Name:      p.Name,
		Role:      p.Role,
		Level:     p.Level,
		ImagePath: p.ImagePath,
		Periode:   p.Periode,
		SortOrder: p.SortOrder,
		IsActive:  p.IsActive,
		CreatedAt: p.CreatedAt.Format(time.RFC3339),
		UpdatedAt: p.UpdatedAt.Format(time.RFC3339),
	}
	if p.Department != nil {
		resp.Department = *p.Department
	}
	if p.Provinsi != nil {
		resp.Provinsi = *p.Provinsi
	}
	if p.Kabupaten != nil {
		resp.Kabupaten = *p.Kabupaten
	}
	if p.FacebookURL != nil {
		resp.FacebookURL = *p.FacebookURL
	}
	if p.InstagramURL != nil {
		resp.InstagramURL = *p.InstagramURL
	}
	if p.LinkedinURL != nil {
		resp.LinkedinURL = *p.LinkedinURL
	}
	if p.Whatsapp != nil {
		resp.Whatsapp = *p.Whatsapp
	}
	return resp
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
