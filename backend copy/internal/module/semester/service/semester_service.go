package service

import (
	"context"
	"fmt"

	"github.com/ahmadzakyarifin/schoolpay/backend/internal/helper"
	activitylogdto "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/activitylog/dto"
	activitylogservice "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/activitylog/service"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/semester/dto"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/semester/entity"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/semester/mapper"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/semester/repository"
	"gorm.io/gorm"
)

type SemesterService interface {
	GetAll(ctx context.Context, req dto.SemesterQueryReq) ([]dto.SemesterRes, int, error)
	Create(ctx context.Context, req dto.SemesterCreateReq) (*dto.SemesterRes, error)
	Update(ctx context.Context, id uint, req dto.SemesterUpdateReq) (*dto.SemesterRes, error)
	Delete(ctx context.Context, id uint) error
	Restore(ctx context.Context, id uint) error
	ToggleStatus(ctx context.Context, id uint) error
	BulkDelete(ctx context.Context, ids []uint) error
	BulkRestore(ctx context.Context, ids []uint) error
	GetDependencyInfo(ctx context.Context, id uint) (map[string]interface{}, error)
	CheckUnique(ctx context.Context, academicYearID uint, name string, excludeID uint) (bool, error)
}

type semesterService struct {
	db       *gorm.DB
	repo     repository.SemesterRepo
	auditSvc activitylogservice.ActivityLogService
}

func NewSemesterService(db *gorm.DB, repo repository.SemesterRepo, auditSvc activitylogservice.ActivityLogService) SemesterService {
	return &semesterService{db: db, repo: repo, auditSvc: auditSvc}
}

func (s *semesterService) log(ctx context.Context, input *activitylogdto.ActivityLogInput) {
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

func (s *semesterService) GetAll(ctx context.Context, req dto.SemesterQueryReq) ([]dto.SemesterRes, int, error) {
	models, total, err := s.repo.FindAll(ctx, req.Page, req.Limit, req.Search, req.Status, req.AcademicYearID, req.Sort)
	if err != nil {
		return nil, 0, err
	}
	entities := mapper.ModelListToEntity(models)
	s.fillAcademicYear(ctx, entities)
	s.fillCounts(ctx, entities)
	return mapper.EntitiesToRes(entities), total, nil
}

func (s *semesterService) fillAcademicYear(ctx context.Context, entities []entity.Semester) {
	ids := make([]uint, 0, len(entities))
	for i := range entities {
		if entities[i].AcademicYearID > 0 {
			ids = append(ids, entities[i].AcademicYearID)
		}
	}
	if len(ids) == 0 {
		return
	}
	type ay struct {
		ID   uint
		Name string
	}
	var rows []ay
	if err := s.db.WithContext(ctx).Model(&struct{}{}).Table("academic_years").
		Select("id, name").Where("id IN ?", ids).Scan(&rows).Error; err != nil {
		return
	}
	m := make(map[uint]string, len(rows))
	for _, r := range rows {
		m[r.ID] = r.Name
	}
	for i := range entities {
		if n, ok := m[entities[i].AcademicYearID]; ok {
			entities[i].AcademicYear = &entity.AcademicYear{ID: entities[i].AcademicYearID, Name: n}
		}
	}
}

func (s *semesterService) fillCounts(ctx context.Context, entities []entity.Semester) {
	for i := range entities {
		cm, _ := s.repo.CountClassMemberships(ctx, entities[i].ID)
		br, _ := s.repo.CountBillingRules(ctx, entities[i].ID)
		inv, _ := s.repo.CountInvoices(ctx, entities[i].ID)
		batch, _ := s.repo.CountBatches(ctx, entities[i].ID)
		entities[i].ClassMembershipCount = cm
		entities[i].BillingRuleCount = br
		entities[i].InvoiceCount = inv
		entities[i].BatchCount = batch
	}
}

func (s *semesterService) Create(ctx context.Context, req dto.SemesterCreateReq) (*dto.SemesterRes, error) {
	v := helper.NewValidationError()
	if req.Name == "" {
		v.Add("name", "Nama semester wajib diisi")
	}
	if req.AcademicYearID == 0 {
		v.Add("academic_year_id", "Tahun ajaran wajib dipilih")
	}
	if req.EndDate.Before(req.StartDate) {
		v.Add("end_date", "Tanggal akhir harus sama atau setelah tanggal mulai")
	}
	if len(v.Errors) > 0 {
		return nil, v
	}

	exists, err := s.repo.Exists(ctx, req.AcademicYearID, req.Name, 0)
	if err != nil {
		return nil, err
	}
	if exists {
		v.Add("name", "Semester untuk tahun ajaran tersebut sudah ada")
	}
	if len(v.Errors) > 0 {
		return nil, v
	}

	e := mapper.CreateReqToEntity(&req)
	if err := s.repo.Create(ctx, mapper.EntityToModel(e)); err != nil {
		return nil, err
	}

	s.log(ctx, &activitylogdto.ActivityLogInput{
		Action:      "semester.create",
		EntityType:  "semester",
		EntityID:    &e.ID,
		EntityLabel: e.Name,
		Description: fmt.Sprintf("Membuat semester: %s", e.Name),
	})

	created, _ := s.repo.FindByID(ctx, e.ID)
	res := mapper.EntityToRes(mapper.ModelToEntity(created))
	return &res, nil
}

func (s *semesterService) Update(ctx context.Context, id uint, req dto.SemesterUpdateReq) (*dto.SemesterRes, error) {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, helper.NewNotFoundError("Semester tidak ditemukan")
	}

	v := helper.NewValidationError()
	if req.Name == "" {
		v.Add("name", "Nama semester wajib diisi")
	}
	if req.AcademicYearID == 0 {
		v.Add("academic_year_id", "Tahun ajaran wajib dipilih")
	}
	if req.EndDate.Before(req.StartDate) {
		v.Add("end_date", "Tanggal akhir harus sama atau setelah tanggal mulai")
	}
	if len(v.Errors) > 0 {
		return nil, v
	}

	exists, err := s.repo.Exists(ctx, req.AcademicYearID, req.Name, id)
	if err != nil {
		return nil, err
	}
	if exists {
		v.Add("name", "Semester untuk tahun ajaran tersebut sudah ada")
	}
	if len(v.Errors) > 0 {
		return nil, v
	}

	e := mapper.ModelToEntity(existing)
	mapper.UpdateReqToEntity(&req, e)
	if err := s.repo.Update(ctx, mapper.EntityToModel(e)); err != nil {
		return nil, err
	}

	s.log(ctx, &activitylogdto.ActivityLogInput{
		Action:      "semester.update",
		EntityType:  "semester",
		EntityID:    &existing.ID,
		EntityLabel: e.Name,
		Description: fmt.Sprintf("Memperbarui semester: %s", e.Name),
	})

	updated, _ := s.repo.FindByID(ctx, id)
	res := mapper.EntityToRes(mapper.ModelToEntity(updated))
	return &res, nil
}

func (s *semesterService) Delete(ctx context.Context, id uint) error {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return helper.NewNotFoundError("Semester tidak ditemukan")
	}
	// FK semester -> child tables pakai SET NULL, jadi aman dihapus.
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	s.log(ctx, &activitylogdto.ActivityLogInput{
		Action:      "semester.delete",
		EntityType:  "semester",
		EntityID:    &existing.ID,
		EntityLabel: existing.Name,
		Description: fmt.Sprintf("Menghapus semester: %s", existing.Name),
	})
	return nil
}

func (s *semesterService) Restore(ctx context.Context, id uint) error {
	existing, err := s.repo.FindByIDUnscoped(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return helper.NewNotFoundError("Semester tidak ditemukan")
	}
	if err := s.repo.Restore(ctx, id); err != nil {
		return err
	}
	s.log(ctx, &activitylogdto.ActivityLogInput{
		Action:      "semester.restore",
		EntityType:  "semester",
		EntityID:    &existing.ID,
		EntityLabel: existing.Name,
		Description: fmt.Sprintf("Memulihkan semester: %s", existing.Name),
	})
	return nil
}

func (s *semesterService) ToggleStatus(ctx context.Context, id uint) error {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return helper.NewNotFoundError("Semester tidak ditemukan")
	}
	if err := s.repo.ToggleStatus(ctx, id); err != nil {
		return err
	}
	s.log(ctx, &activitylogdto.ActivityLogInput{
		Action:      "semester.update",
		EntityType:  "semester",
		EntityID:    &existing.ID,
		EntityLabel: existing.Name,
		Description: fmt.Sprintf("Mengubah status aktif/nonaktif semester: %s", existing.Name),
	})
	return nil
}

func (s *semesterService) BulkDelete(ctx context.Context, ids []uint) error {
	if err := s.repo.BulkDelete(ctx, ids); err != nil {
		return err
	}
	for _, id := range ids {
		s.log(ctx, &activitylogdto.ActivityLogInput{
			Action:      "semester.delete",
			EntityType:  "semester",
			EntityLabel: fmt.Sprintf("%d", id),
			Description: "Menghapus massal template kelas: " + fmt.Sprintf("%d", id),
		})
	}
	return nil
}

func (s *semesterService) BulkRestore(ctx context.Context, ids []uint) error {
	if err := s.repo.BulkRestore(ctx, ids); err != nil {
		return err
	}
	for _, id := range ids {
		s.log(ctx, &activitylogdto.ActivityLogInput{
			Action:      "semester.restore",
			EntityType:  "semester",
			EntityLabel: fmt.Sprintf("%d", id),
			Description: "Memulihkan massal template kelas: " + fmt.Sprintf("%d", id),
		})
	}
	return nil
}

func (s *semesterService) GetDependencyInfo(ctx context.Context, id uint) (map[string]interface{}, error) {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, helper.NewNotFoundError("Semester tidak ditemukan")
	}
	return map[string]interface{}{}, nil
}

func (s *semesterService) CheckUnique(ctx context.Context, academicYearID uint, name string, excludeID uint) (bool, error) {
	exists, err := s.repo.Exists(ctx, academicYearID, name, excludeID)
	if err != nil {
		return false, err
	}
	return !exists, nil
}
