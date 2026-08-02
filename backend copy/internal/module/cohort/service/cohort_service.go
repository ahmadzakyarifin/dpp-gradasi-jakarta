package service

import (
	"context"
	"fmt"

	"github.com/ahmadzakyarifin/schoolpay/backend/internal/helper"
	activitylogdto "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/activitylog/dto"
	activitylogservice "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/activitylog/service"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/cohort/dto"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/cohort/mapper"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/cohort/repository"
	"gorm.io/gorm"
)

type CohortService interface {
	GetAll(ctx context.Context, req dto.CohortQueryReq) ([]dto.CohortRes, int, error)
	Create(ctx context.Context, req dto.CohortCreateReq) (*dto.CohortRes, error)
	Update(ctx context.Context, id uint, req dto.CohortUpdateReq) (*dto.CohortRes, error)
	Delete(ctx context.Context, id uint) error
	Restore(ctx context.Context, id uint) error
	ToggleStatus(ctx context.Context, id uint) error
	BulkDelete(ctx context.Context, ids []uint) error
	BulkRestore(ctx context.Context, ids []uint) error
	GetDependencyInfo(ctx context.Context, id uint) (map[string]interface{}, error)
	CheckUnique(ctx context.Context, name string, excludeID uint) (bool, error)
}

type cohortService struct {
	db       *gorm.DB
	repo     repository.CohortRepo
	auditSvc activitylogservice.ActivityLogService
}

func NewCohortService(db *gorm.DB, repo repository.CohortRepo, auditSvc activitylogservice.ActivityLogService) CohortService {
	return &cohortService{db: db, repo: repo, auditSvc: auditSvc}
}

func (s *cohortService) log(ctx context.Context, input *activitylogdto.ActivityLogInput) {
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

func (s *cohortService) GetAll(ctx context.Context, req dto.CohortQueryReq) ([]dto.CohortRes, int, error) {
	entities, total, err := s.repo.FindAll(ctx, req.Page, req.Limit, req.Search, req.Status, req.Sort)
	if err != nil {
		return nil, 0, err
	}
	return mapper.EntitiesToRes(entities), total, nil
}

func (s *cohortService) Create(ctx context.Context, req dto.CohortCreateReq) (*dto.CohortRes, error) {
	v := helper.NewValidationError()
	if req.Name == "" {
		v.Add("name", "Nama angkatan wajib diisi")
	}
	if req.EndDate.Before(req.StartDate) {
		v.Add("end_date", "Tanggal akhir harus sama atau setelah tanggal mulai")
	}
	if len(v.Errors) > 0 {
		return nil, v
	}

	exists, err := s.repo.Exists(ctx, req.Name, 0)
	if err != nil {
		return nil, err
	}
	if exists {
		v.Add("name", "Nama angkatan sudah digunakan")
	}
	if len(v.Errors) > 0 {
		return nil, v
	}

	e := mapper.CreateReqToEntity(&req)
	if err := s.repo.Create(ctx, e); err != nil {
		if err.Error() == "nama angkatan sudah digunakan" {
			v := helper.NewValidationError()
			v.Add("name", "Nama angkatan sudah digunakan")
			return nil, v
		}
		return nil, err
	}

	s.log(ctx, &activitylogdto.ActivityLogInput{
		Action:      "cohort.create",
		EntityType:  "cohort",
		EntityID:    &e.ID,
		EntityLabel: e.Name,
		Description: fmt.Sprintf("Membuat angkatan: %s", e.Name),
	})

	res := mapper.EntityToRes(e)
	return &res, nil
}

func (s *cohortService) Update(ctx context.Context, id uint, req dto.CohortUpdateReq) (*dto.CohortRes, error) {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, helper.NewNotFoundError("Angkatan tidak ditemukan")
	}

	v := helper.NewValidationError()
	if req.Name == "" {
		v.Add("name", "Nama angkatan wajib diisi")
	}
	if req.EndDate.Before(req.StartDate) {
		v.Add("end_date", "Tanggal akhir harus sama atau setelah tanggal mulai")
	}
	if len(v.Errors) > 0 {
		return nil, v
	}

	exists, err := s.repo.Exists(ctx, req.Name, existing.ID)
	if err != nil {
		return nil, err
	}
	if exists {
		v.Add("name", "Nama angkatan sudah digunakan")
	}
	if len(v.Errors) > 0 {
		return nil, v
	}

	e := existing
	mapper.UpdateReqToEntity(&req, e)
	if err := s.repo.Update(ctx, e); err != nil {
		if err.Error() == "nama angkatan sudah digunakan" {
			v := helper.NewValidationError()
			v.Add("name", "Nama angkatan sudah digunakan")
			return nil, v
		}
		return nil, err
	}

	s.log(ctx, &activitylogdto.ActivityLogInput{
		Action:      "cohort.update",
		EntityType:  "cohort",
		EntityID:    &e.ID,
		EntityLabel: e.Name,
		Description: fmt.Sprintf("Memperbarui angkatan: %s", e.Name),
	})

	updated, _ := s.repo.FindByID(ctx, id)
	res := mapper.EntityToRes(updated)
	return &res, nil
}

func (s *cohortService) Delete(ctx context.Context, id uint) error {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return helper.NewNotFoundError("Angkatan tidak ditemukan")
	}

	students, _ := s.repo.CountStudents(ctx, existing.ID)
	if students > 0 {
		v := helper.NewValidationError()
		v.Add("general", fmt.Sprintf("tidak dapat menghapus angkatan karena masih digunakan di %d siswa", students))
		return v
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	s.log(ctx, &activitylogdto.ActivityLogInput{
		Action:      "cohort.delete",
		EntityType:  "cohort",
		EntityID:    &existing.ID,
		EntityLabel: existing.Name,
		Description: fmt.Sprintf("Menghapus angkatan: %s", existing.Name),
	})
	return nil
}

func (s *cohortService) Restore(ctx context.Context, id uint) error {
	existing, err := s.repo.FindByIDUnscoped(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return helper.NewNotFoundError("Angkatan tidak ditemukan")
	}

	if err := s.repo.Restore(ctx, id); err != nil {
		return err
	}

	s.log(ctx, &activitylogdto.ActivityLogInput{
		Action:      "cohort.restore",
		EntityType:  "cohort",
		EntityID:    &existing.ID,
		EntityLabel: existing.Name,
		Description: fmt.Sprintf("Memulihkan angkatan: %s", existing.Name),
	})
	return nil
}

func (s *cohortService) ToggleStatus(ctx context.Context, id uint) error {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return helper.NewNotFoundError("Angkatan tidak ditemukan")
	}

	if err := s.repo.ToggleStatus(ctx, id); err != nil {
		return err
	}

	s.log(ctx, &activitylogdto.ActivityLogInput{
		Action:      "cohort.update",
		EntityType:  "cohort",
		EntityID:    &existing.ID,
		EntityLabel: existing.Name,
		Description: fmt.Sprintf("Mengubah status aktif/nonaktif angkatan: %s", existing.Name),
	})
	return nil
}

func (s *cohortService) BulkDelete(ctx context.Context, ids []uint) error {
	for _, id := range ids {
		students, _ := s.repo.CountStudents(ctx, id)
		if students > 0 {
			v := helper.NewValidationError()
			v.Add("general", fmt.Sprintf("tidak dapat menghapus angkatan karena masih digunakan di siswa"))
			return v
		}
	}

	err := s.repo.BulkDelete(ctx, ids)
	if err != nil {
		if err.Error() == "tidak ada data yang dapat dihapus" {
			v := helper.NewValidationError()
			v.Add("general", "tidak ada data angkatan aktif yang dapat dihapus")
			return v
		}
		return err
	}

	nameMap, _ := s.repo.FindNamesByIDs(ctx, ids)
	for _, id := range ids {
		cohortID := id
		name := nameMap[id]
		s.log(ctx, &activitylogdto.ActivityLogInput{
			Action:      "cohort.delete",
			EntityType:  "cohort",
			EntityID:    &cohortID,
			EntityLabel: name,
			Description: fmt.Sprintf("Menghapus angkatan (bulk): %s", name),
		})
	}
	return nil
}

func (s *cohortService) BulkRestore(ctx context.Context, ids []uint) error {
	err := s.repo.BulkRestore(ctx, ids)
	if err != nil {
		if err.Error() == "tidak ada data yang perlu dipulihkan" {
			v := helper.NewValidationError()
			v.Add("general", "tidak ada data angkatan terhapus yang dapat dipulihkan")
			return v
		}
		return err
	}

	nameMap, _ := s.repo.FindNamesByIDs(ctx, ids)
	for _, id := range ids {
		cohortID := id
		name := nameMap[id]
		s.log(ctx, &activitylogdto.ActivityLogInput{
			Action:      "cohort.restore",
			EntityType:  "cohort",
			EntityID:    &cohortID,
			EntityLabel: name,
			Description: fmt.Sprintf("Memulihkan angkatan (bulk): %s", name),
		})
	}
	return nil
}

func (s *cohortService) GetDependencyInfo(ctx context.Context, id uint) (map[string]interface{}, error) {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, fmt.Errorf("Not found")
	}
	students, err := s.repo.CountStudents(ctx, existing.ID)
	if err != nil {
		return nil, err
	}
	hasDep := students > 0
	return map[string]interface{}{
		"student_count":    students,
		"can_delete":       !hasDep,
		"has_dependencies": hasDep,
		"message":          fmt.Sprintf("%d siswa menggunakan angkatan ini", students),
	}, nil
}

func (s *cohortService) CheckUnique(ctx context.Context, name string, excludeID uint) (bool, error) {
	exists, err := s.repo.Exists(ctx, name, excludeID)
	return !exists, err
}
