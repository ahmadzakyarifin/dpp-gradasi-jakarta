package service

import (
	"context"
	"fmt"

	"github.com/ahmadzakyarifin/schoolpay/backend/internal/helper"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/academicyear/dto"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/academicyear/entity"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/academicyear/mapper"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/academicyear/repository"
	activitylogdto "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/activitylog/dto"
	activitylogservice "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/activitylog/service"
	"gorm.io/gorm"
)

type AcademicYearService interface {
	GetAll(ctx context.Context, req dto.AcademicYearQueryReq) ([]dto.AcademicYearRes, int, error)
	Create(ctx context.Context, req dto.AcademicYearCreateReq) (*dto.AcademicYearRes, error)
	Update(ctx context.Context, id uint, req dto.AcademicYearUpdateReq) (*dto.AcademicYearRes, error)
	Delete(ctx context.Context, id uint) error
	Restore(ctx context.Context, id uint) error
	ToggleStatus(ctx context.Context, id uint) error
	BulkDelete(ctx context.Context, ids []uint) error
	BulkRestore(ctx context.Context, ids []uint) error
	GetDependencyInfo(ctx context.Context, id uint) (map[string]interface{}, error)
	CheckUnique(ctx context.Context, name string, excludeID uint) (bool, error)
}

type academicYearService struct {
	db       *gorm.DB
	repo     repository.AcademicYearRepo
	auditSvc activitylogservice.ActivityLogService
}

func NewAcademicYearService(db *gorm.DB, repo repository.AcademicYearRepo, auditSvc activitylogservice.ActivityLogService) AcademicYearService {
	return &academicYearService{db: db, repo: repo, auditSvc: auditSvc}
}

func (s *academicYearService) log(ctx context.Context, input *activitylogdto.ActivityLogInput) {
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

func (s *academicYearService) GetAll(ctx context.Context, req dto.AcademicYearQueryReq) ([]dto.AcademicYearRes, int, error) {
	models, total, err := s.repo.FindAll(ctx, req.Page, req.Limit, req.Search, req.Status, req.Sort)
	if err != nil {
		return nil, 0, err
	}
	entities := mapper.ModelListToEntity(models)
	for i := range entities {
		s.fillCounts(ctx, &entities[i])
	}
	return mapper.EntitiesToRes(entities), total, nil
}

func (s *academicYearService) fillCounts(ctx context.Context, e *entity.AcademicYear) {
	sem, _ := s.repo.CountSemesters(ctx, e.ID)
	classes, _ := s.repo.CountActiveClasses(ctx, e.ID)
	students, _ := s.repo.CountStudents(ctx, e.ID)
	billings, _ := s.repo.CountBillingRules(ctx, e.ID)
	invoices, _ := s.repo.CountInvoices(ctx, e.ID)
	e.SemesterCount = sem
	e.ActiveClassCount = classes
	e.StudentCount = students
	e.BillingRuleCount = billings
	e.InvoiceCount = invoices
}

func (s *academicYearService) Create(ctx context.Context, req dto.AcademicYearCreateReq) (*dto.AcademicYearRes, error) {
	v := helper.NewValidationError()
	if req.Name == "" {
		v.Add("name", "Nama tahun ajaran wajib diisi")
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
		v.Add("name", "Nama tahun ajaran sudah digunakan")
	}
	if len(v.Errors) > 0 {
		return nil, v
	}

	e := mapper.CreateReqToEntity(&req)
	if err := s.repo.Create(ctx, mapper.EntityToModel(e)); err != nil {
		return nil, err
	}

	s.log(ctx, &activitylogdto.ActivityLogInput{
		Action:      "academic_year.create",
		EntityType:  "academic_year",
		EntityID:    &e.ID,
		EntityLabel: e.Name,
		Description: fmt.Sprintf("Membuat tahun ajaran: %s", e.Name),
	})

	created, _ := s.repo.FindByID(ctx, e.ID)
	res := mapper.EntityToRes(mapper.ModelToEntity(created))
	return &res, nil
}

func (s *academicYearService) Update(ctx context.Context, id uint, req dto.AcademicYearUpdateReq) (*dto.AcademicYearRes, error) {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, helper.NewNotFoundError("Tahun ajaran tidak ditemukan")
	}

	v := helper.NewValidationError()
	if req.Name == "" {
		v.Add("name", "Nama tahun ajaran wajib diisi")
	}
	if req.EndDate.Before(req.StartDate) {
		v.Add("end_date", "Tanggal akhir harus sama atau setelah tanggal mulai")
	}
	if len(v.Errors) > 0 {
		return nil, v
	}

	exists, err := s.repo.Exists(ctx, req.Name, id)
	if err != nil {
		return nil, err
	}
	if exists {
		v.Add("name", "Nama tahun ajaran sudah digunakan")
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
		Action:      "academic_year.update",
		EntityType:  "academic_year",
		EntityID:    &id,
		EntityLabel: e.Name,
		Description: fmt.Sprintf("Memperbarui tahun ajaran: %s", e.Name),
	})

	updated, _ := s.repo.FindByID(ctx, id)
	res := mapper.EntityToRes(mapper.ModelToEntity(updated))
	return &res, nil
}

func (s *academicYearService) Delete(ctx context.Context, id uint) error {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return helper.NewNotFoundError("Tahun ajaran tidak ditemukan")
	}

	sem, _ := s.repo.CountSemesters(ctx, id)
	classes, _ := s.repo.CountActiveClasses(ctx, id)
	students, _ := s.repo.CountStudents(ctx, id)
	billings, _ := s.repo.CountBillingRules(ctx, id)
	invoices, _ := s.repo.CountInvoices(ctx, id)
	if sem > 0 || classes > 0 || students > 0 || billings > 0 || invoices > 0 {
		v := helper.NewValidationError()
		v.Add("general", fmt.Sprintf("tidak dapat menghapus tahun ajaran karena masih digunakan: %d semester, %d kelas, %d siswa, %d aturan tagihan, %d invoice", sem, classes, students, billings, invoices))
		return v
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	s.log(ctx, &activitylogdto.ActivityLogInput{
		Action:      "academic_year.delete",
		EntityType:  "academic_year",
		EntityID:    &id,
		EntityLabel: existing.Name,
		Description: fmt.Sprintf("Menghapus tahun ajaran: %s", existing.Name),
	})
	return nil
}

func (s *academicYearService) Restore(ctx context.Context, id uint) error {
	existing, err := s.repo.FindByIDUnscoped(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return helper.NewNotFoundError("Tahun ajaran tidak ditemukan")
	}
	if err := s.repo.Restore(ctx, id); err != nil {
		return err
	}
	s.log(ctx, &activitylogdto.ActivityLogInput{
		Action:      "academic_year.restore",
		EntityType:  "academic_year",
		EntityID:    &id,
		EntityLabel: existing.Name,
		Description: fmt.Sprintf("Memulihkan tahun ajaran: %s", existing.Name),
	})
	return nil
}

func (s *academicYearService) ToggleStatus(ctx context.Context, id uint) error {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return helper.NewNotFoundError("Tahun ajaran tidak ditemukan")
	}
	if err := s.repo.ToggleStatus(ctx, id); err != nil {
		return err
	}
	s.log(ctx, &activitylogdto.ActivityLogInput{
		Action:      "academic_year.update",
		EntityType:  "academic_year",
		EntityID:    &id,
		EntityLabel: existing.Name,
		Description: fmt.Sprintf("Mengubah status aktif/nonaktif tahun ajaran: %s", existing.Name),
	})
	return nil
}

func (s *academicYearService) BulkDelete(ctx context.Context, ids []uint) error {
	var safeIDs []uint
	for _, id := range ids {
		sem, _ := s.repo.CountSemesters(ctx, id)
		classes, _ := s.repo.CountActiveClasses(ctx, id)
		students, _ := s.repo.CountStudents(ctx, id)
		billings, _ := s.repo.CountBillingRules(ctx, id)
		invoices, _ := s.repo.CountInvoices(ctx, id)
		if sem == 0 && classes == 0 && students == 0 && billings == 0 && invoices == 0 {
			safeIDs = append(safeIDs, id)
		}
	}
	if len(safeIDs) == 0 {
		v := helper.NewValidationError()
		v.Add("general", "semua tahun ajaran terpilih masih memiliki relasi dan tidak dapat dihapus")
		return v
	}
	if err := s.repo.BulkDelete(ctx, safeIDs); err != nil {
		return err
	}
	nameMap, _ := s.repo.FindNamesByIDs(ctx, safeIDs)
	for _, id := range safeIDs {
		entityID := id
		s.log(ctx, &activitylogdto.ActivityLogInput{
			Action:      "academic_year.delete",
			EntityType:  "academic_year",
			EntityID:    &entityID,
			EntityLabel: nameMap[id],
			Description: fmt.Sprintf("Menghapus massal tahun ajaran: %s", nameMap[id]),
		})
	}
	return nil
}

func (s *academicYearService) BulkRestore(ctx context.Context, ids []uint) error {
	if err := s.repo.BulkRestore(ctx, ids); err != nil {
		return err
	}
	nameMap, _ := s.repo.FindNamesByIDs(ctx, ids)
	for _, id := range ids {
		entityID := id
		s.log(ctx, &activitylogdto.ActivityLogInput{
			Action:      "academic_year.restore",
			EntityType:  "academic_year",
			EntityID:    &entityID,
			EntityLabel: nameMap[id],
			Description: fmt.Sprintf("Memulihkan massal tahun ajaran: %s", nameMap[id]),
		})
	}
	return nil
}

func (s *academicYearService) GetDependencyInfo(ctx context.Context, id uint) (map[string]interface{}, error) {
	sem, err := s.repo.CountSemesters(ctx, id)
	if err != nil {
		return nil, err
	}
	classes, err := s.repo.CountActiveClasses(ctx, id)
	if err != nil {
		return nil, err
	}
	students, err := s.repo.CountStudents(ctx, id)
	if err != nil {
		return nil, err
	}
	billings, err := s.repo.CountBillingRules(ctx, id)
	if err != nil {
		return nil, err
	}
	invoices, err := s.repo.CountInvoices(ctx, id)
	if err != nil {
		return nil, err
	}
	hasDep := sem > 0 || classes > 0 || students > 0 || billings > 0 || invoices > 0
	return map[string]interface{}{
		"semester_count":     sem,
		"active_class_count": classes,
		"student_count":      students,
		"billing_rule_count": billings,
		"invoice_count":      invoices,
		"can_delete":         !hasDep,
		"has_dependencies":   hasDep,
		"message": fmt.Sprintf("%d semester, %d kelas, %d siswa, %d aturan tagihan, %d invoice menggunakan tahun ajaran ini",
			sem, classes, students, billings, invoices),
	}, nil
}

func (s *academicYearService) CheckUnique(ctx context.Context, name string, excludeID uint) (bool, error) {
	exists, err := s.repo.Exists(ctx, name, excludeID)
	return !exists, err
}
