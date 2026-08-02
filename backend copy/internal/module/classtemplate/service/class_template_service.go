package service

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/ahmadzakyarifin/schoolpay/backend/internal/helper"
	activitylogdto "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/activitylog/dto"
	activitylogservice "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/activitylog/service"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/classtemplate/dto"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/classtemplate/mapper"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/classtemplate/repository"
	majorrepo "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/major/repository"
	"gorm.io/gorm"
)

type ClassTemplateService interface {
	GetAll(ctx context.Context, req dto.ClassTemplateQueryReq) ([]dto.ClassTemplateRes, int, error)
	Create(ctx context.Context, req dto.ClassTemplateCreateReq) (*dto.ClassTemplateRes, error)
	Update(ctx context.Context, id uint, req dto.ClassTemplateUpdateReq) (*dto.ClassTemplateRes, error)
	Delete(ctx context.Context, id uint) error
	Restore(ctx context.Context, id uint) error
	ToggleStatus(ctx context.Context, id uint) error
	BulkDelete(ctx context.Context, ids []uint) error
	BulkRestore(ctx context.Context, ids []uint) error
	GetDependencyInfo(ctx context.Context, id uint) (map[string]interface{}, error)
	CheckUnique(ctx context.Context, name string, majorID *uint, excludeID uint) (bool, error)
	SuggestNextName(ctx context.Context, baseName string) (string, error)
}

type service struct {
	db        *gorm.DB
	repo      repository.ClassTemplateRepo
	majorRepo majorrepo.MajorRepo
	auditSvc  activitylogservice.ActivityLogService
}

func NewClassTemplateService(db *gorm.DB, repo repository.ClassTemplateRepo, majorRepo majorrepo.MajorRepo, auditSvc activitylogservice.ActivityLogService) ClassTemplateService {
	return &service{
		db:        db,
		repo:      repo,
		majorRepo: majorRepo,
		auditSvc:  auditSvc,
	}
}

func (s *service) log(ctx context.Context, input *activitylogdto.ActivityLogInput) {
	if s.auditSvc == nil {
		return
	}
	userID, userName, role, ipAddress, userAgent := helper.GetAuditMeta(ctx)
	var uID *uint
	if userID > 0 {
		uID = &userID
	}
	input.ActorID = uID
	input.ActorName = userName
	input.ActorRole = role
	input.IPAddress = ipAddress
	input.UserAgent = userAgent

	_ = s.auditSvc.Log(ctx, s.db, input)
}

func (s *service) GetAll(ctx context.Context, req dto.ClassTemplateQueryReq) ([]dto.ClassTemplateRes, int, error) {
	entities, total, err := s.repo.FindAll(ctx, req.Page, req.Limit, req.Search, req.Status, req.MajorID, req.GradeLevel, req.Sort)
	if err != nil {
		return nil, 0, err
	}

	res := make([]dto.ClassTemplateRes, 0, len(entities))
	for i := range entities {
		res = append(res, mapper.EntityToResponse(&entities[i]))
	}
	return res, total, nil
}

func (s *service) validateMajor(ctx context.Context, majorID *uint) error {
	if majorID == nil || *majorID == 0 {
		return nil
	}
	existing, err := s.majorRepo.FindByID(ctx, *majorID)
	if err != nil {
		return err
	}
	if existing == nil {
		return helper.NewNotFoundError("Jurusan tidak ditemukan")
	}
	if !existing.IsActive {
		v := helper.NewValidationError()
		v.Add("major_id", "Jurusan sedang tidak aktif")
		return v
	}
	return nil
}

func (s *service) Create(ctx context.Context, req dto.ClassTemplateCreateReq) (*dto.ClassTemplateRes, error) {
	if err := s.validateMajor(ctx, req.MajorID); err != nil {
		return nil, err
	}

	req.Name = strings.TrimSpace(req.Name)
	v := helper.NewValidationError()
	if req.Name == "" {
		v.Add("name", "Nama kelas wajib diisi")
	}
	if req.GradeLevel < 1 || req.GradeLevel > 12 {
		v.Add("grade_level", "Tingkat kelas harus antara 1-12")
	}

	exists, err := s.repo.Exists(ctx, req.Name, req.MajorID, 0)
	if err != nil {
		return nil, err
	}
	if exists {
		v.Add("name", "Nama kelas sudah digunakan pada jurusan ini")
	}

	if len(v.Errors) > 0 {
		return nil, v
	}

	e := mapper.CreateReqToEntity(&req)
	if err := s.repo.Create(ctx, e); err != nil {
		return nil, err
	}

	s.log(ctx, &activitylogdto.ActivityLogInput{
		Action:      "class_template.create",
		EntityType:  "class_template",
		EntityID:    &e.ID,
		EntityLabel: e.Name,
		Description: fmt.Sprintf("Membuat template kelas: %s", e.Name),
	})

	res := mapper.EntityToResponse(e)
	return &res, nil
}

func (s *service) Update(ctx context.Context, id uint, req dto.ClassTemplateUpdateReq) (*dto.ClassTemplateRes, error) {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, helper.NewNotFoundError("Template kelas tidak ditemukan")
	}

	if err := s.validateMajor(ctx, req.MajorID); err != nil {
		return nil, err
	}

	req.Name = strings.TrimSpace(req.Name)
	v := helper.NewValidationError()
	if req.Name == "" {
		v.Add("name", "Nama kelas wajib diisi")
	}
	if req.GradeLevel < 1 || req.GradeLevel > 12 {
		v.Add("grade_level", "Tingkat kelas harus antara 1-12")
	}

	exists, err := s.repo.Exists(ctx, req.Name, req.MajorID, existing.ID)
	if err != nil {
		return nil, err
	}
	if exists {
		v.Add("name", "Nama kelas sudah digunakan pada jurusan ini")
	}

	if len(v.Errors) > 0 {
		return nil, v
	}

	e := existing
	mapper.UpdateReqToEntity(&req, e)
	if err := s.repo.Update(ctx, e); err != nil {
		return nil, err
	}

	s.log(ctx, &activitylogdto.ActivityLogInput{
		Action:      "class_template.update",
		EntityType:  "class_template",
		EntityID:    &e.ID,
		EntityLabel: e.Name,
		Description: fmt.Sprintf("Memperbarui template kelas: %s", e.Name),
	})

	updated, _ := s.repo.FindByID(ctx, id)
	res := mapper.EntityToResponse(updated)
	return &res, nil
}

func (s *service) Delete(ctx context.Context, id uint) error {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return helper.NewNotFoundError("Template kelas tidak ditemukan")
	}

	count, err := s.repo.CountActiveClasses(ctx, existing.ID)
	if err != nil {
		return err
	}
	if count > 0 {
		v := helper.NewValidationError()
		v.Add("general", fmt.Sprintf("tidak dapat menghapus kelas karena masih digunakan di %d kelas aktif", count))
		return v
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	s.log(ctx, &activitylogdto.ActivityLogInput{
		Action:      "class_template.delete",
		EntityType:  "class_template",
		EntityID:    &existing.ID,
		EntityLabel: existing.Name,
		Description: fmt.Sprintf("Menghapus template kelas: %s", existing.Name),
	})
	return nil
}

func (s *service) Restore(ctx context.Context, id uint) error {
	if err := s.repo.Restore(ctx, id); err != nil {
		return err
	}

	existing, err := s.repo.FindByID(ctx, id)
	if err != nil || existing == nil {
		// Jika gagal mengambil data untuk log, abaikan saja karena restore sudah berhasil
		return nil
	}

	s.log(ctx, &activitylogdto.ActivityLogInput{
		Action:      "class_template.restore",
		EntityType:  "class_template",
		EntityID:    &existing.ID,
		EntityLabel: existing.Name,
		Description: fmt.Sprintf("Memulihkan template kelas: %s", existing.Name),
	})
	return nil
}

func (s *service) ToggleStatus(ctx context.Context, id uint) error {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return helper.NewNotFoundError("Template kelas tidak ditemukan")
	}

	if err := s.repo.ToggleStatus(ctx, id); err != nil {
		return err
	}

	s.log(ctx, &activitylogdto.ActivityLogInput{
		Action:      "class_template.update",
		EntityType:  "class_template",
		EntityID:    &existing.ID,
		EntityLabel: existing.Name,
		Description: fmt.Sprintf("Mengubah status aktif/nonaktif template kelas: %s", existing.Name),
	})
	return nil
}

func (s *service) BulkDelete(ctx context.Context, ids []uint) error {
	if err := s.repo.BulkDelete(ctx, ids); err != nil {
		return err
	}
	for _, id := range ids {
		s.log(ctx, &activitylogdto.ActivityLogInput{
			Action:      "class_template.delete",
			EntityType:  "class_template",
			EntityID:    &id,
			EntityLabel: strconv.FormatUint(uint64(id), 10),
			Description: "Menghapus massal template kelas: " + strconv.FormatUint(uint64(id), 10),
		})
	}
	return nil
}

func (s *service) BulkRestore(ctx context.Context, ids []uint) error {
	if err := s.repo.BulkRestore(ctx, ids); err != nil {
		return err
	}
	for _, id := range ids {
		s.log(ctx, &activitylogdto.ActivityLogInput{
			Action:      "class_template.restore",
			EntityType:  "class_template",
			EntityID:    &id,
			EntityLabel: strconv.FormatUint(uint64(id), 10),
			Description: "Memulihkan massal template kelas: " + strconv.FormatUint(uint64(id), 10),
		})
	}
	return nil
}

func (s *service) GetDependencyInfo(ctx context.Context, id uint) (map[string]interface{}, error) {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, helper.NewNotFoundError("Template kelas tidak ditemukan")
	}
	count, err := s.repo.CountActiveClasses(ctx, existing.ID)
	if err != nil {
		return nil, err
	}
	message := ""
	if count > 0 {
		message = fmt.Sprintf("%d kelas aktif di tahun ajaran", count)
	}
	return map[string]interface{}{
		"has_dependencies": count > 0,
		"message":          message,
		"counts": map[string]int{
			"active_classes": count,
		},
	}, nil
}

func (s *service) SuggestNextName(ctx context.Context, baseName string) (string, error) {
	baseName = strings.TrimSpace(baseName)
	if baseName == "" {
		return "", nil
	}

	re := regexp.MustCompile(`(\d+)(\D*)$`)
	matches := re.FindStringSubmatchIndex(baseName)

	suggestedName := baseName + " 1"

	if matches != nil {
		numStart := matches[2]
		numEnd := matches[3]
		suffixStart := matches[4]
		suffixEnd := matches[5]

		numStr := baseName[numStart:numEnd]
		num, _ := strconv.Atoi(numStr)

		newNumStr := strconv.Itoa(num + 1)
		suffix := baseName[suffixStart:suffixEnd]

		suggestedName = baseName[:numStart] + newNumStr + suffix
	}

	return suggestedName, nil
}

func (s *service) CheckUnique(ctx context.Context, name string, majorID *uint, excludeID uint) (bool, error) {
	exists, err := s.repo.Exists(ctx, name, majorID, excludeID)
	if err != nil {
		return false, err
	}
	return !exists, nil
}
