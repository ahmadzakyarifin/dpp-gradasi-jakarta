package service

import (
	"context"
	"fmt"

	"github.com/ahmadzakyarifin/schoolpay/backend/internal/helper"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/activeclass/dto"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/activeclass/entity"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/activeclass/mapper"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/activeclass/model"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/activeclass/repository"
	activitylogdto "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/activitylog/dto"
	activitylogservice "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/activitylog/service"
	"gorm.io/gorm"
)

type ActiveClassService interface {
	GetAll(ctx context.Context, req dto.ActiveClassQueryReq) ([]dto.ActiveClassRes, int, error)
	GetByID(ctx context.Context, id uint) (*dto.ActiveClassRes, error)
	Create(ctx context.Context, req dto.ActiveClassCreateReq) (*dto.ActiveClassRes, error)
	Update(ctx context.Context, id uint, req dto.ActiveClassUpdateReq) (*dto.ActiveClassRes, error)
	Delete(ctx context.Context, id uint) error
	ToggleStatus(ctx context.Context, id uint) error
	BulkUpsertByYear(ctx context.Context, academicYearID uint, items []dto.BulkUpsertItem) ([]dto.ActiveClassRes, error)
}

type activeClassService struct {
	db       *gorm.DB
	repo     repository.ActiveClassRepo
	auditSvc activitylogservice.ActivityLogService
}

func NewActiveClassService(db *gorm.DB, repo repository.ActiveClassRepo, auditSvc activitylogservice.ActivityLogService) ActiveClassService {
	return &activeClassService{db: db, repo: repo, auditSvc: auditSvc}
}

func (s *activeClassService) log(ctx context.Context, input *activitylogdto.ActivityLogInput) {
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

func (s *activeClassService) GetAll(ctx context.Context, req dto.ActiveClassQueryReq) ([]dto.ActiveClassRes, int, error) {
	models, total, err := s.repo.FindAll(ctx, req.Page, req.Limit, req.Search, req.Status, req.AcademicYearID, req.Sort)
	if err != nil {
		return nil, 0, err
	}
	entities := mapper.ModelListToEntity(models)
	s.fillNames(ctx, entities)
	return mapper.EntitiesToRes(entities), total, nil
}

func (s *activeClassService) fillNames(ctx context.Context, entities []entity.ActiveClass) {
	type row struct {
		ID   uint
		Name string
	}
	var ayIDs, ctIDs []uint
	for i := range entities {
		if entities[i].AcademicYearID > 0 {
			ayIDs = append(ayIDs, entities[i].AcademicYearID)
		}
		if entities[i].ClassTemplateID > 0 {
			ctIDs = append(ctIDs, entities[i].ClassTemplateID)
		}
	}
	ayMap := s.mapNames(ctx, "academic_years", ayIDs)
	ctMap := s.mapNames(ctx, "class_templates", ctIDs)
	for i := range entities {
		if n, ok := ayMap[entities[i].AcademicYearID]; ok {
			entities[i].AcademicYear = &entity.AcademicYear{ID: entities[i].AcademicYearID, Name: n}
		}
		if n, ok := ctMap[entities[i].ClassTemplateID]; ok {
			entities[i].ClassTemplate = &entity.ClassTemplate{ID: entities[i].ClassTemplateID, Name: n}
		}
	}
}

func (s *activeClassService) mapNames(ctx context.Context, table string, ids []uint) map[uint]string {
	out := make(map[uint]string)
	if len(ids) == 0 {
		return out
	}
	var rows []struct {
		ID   uint
		Name string
	}
	if err := s.db.WithContext(ctx).Table(table).Select("id, name").Where("id IN ?", ids).Scan(&rows).Error; err != nil {
		return out
	}
	for _, r := range rows {
		out[r.ID] = r.Name
	}
	return out
}

func (s *activeClassService) GetByID(ctx context.Context, id uint) (*dto.ActiveClassRes, error) {
	m, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, helper.NewNotFoundError("Kelas aktif tidak ditemukan")
	}
	entities := []entity.ActiveClass{*mapper.ModelToEntity(m)}
	s.fillNames(ctx, entities)
	res := mapper.EntityToRes(&entities[0])
	return &res, nil
}

func (s *activeClassService) Create(ctx context.Context, req dto.ActiveClassCreateReq) (*dto.ActiveClassRes, error) {
	v := helper.NewValidationError()
	if req.Name == "" {
		v.Add("name", "Nama kelas wajib diisi")
	}
	if req.AcademicYearID == 0 {
		v.Add("academic_year_id", "Tahun ajaran wajib dipilih")
	}
	if req.ClassTemplateID == 0 {
		v.Add("class_template_id", "Template kelas wajib dipilih")
	}
	if len(v.Errors) > 0 {
		return nil, v
	}

	if exists, _ := s.repo.ExistsNameInYear(ctx, req.AcademicYearID, req.Name, 0); exists {
		v.Add("name", "Nama kelas sudah digunakan di tahun ajaran tersebut")
	}
	if len(v.Errors) > 0 {
		return nil, v
	}

	m := mapper.CreateReqToModel(&req)
	if err := s.repo.Create(ctx, m); err != nil {
		return nil, err
	}
	s.log(ctx, &activitylogdto.ActivityLogInput{
		Action:      "active_class.create",
		EntityType:  "active_class",
		EntityID:    &m.ID,
		EntityLabel: m.Name,
		Description: fmt.Sprintf("Membuat kelas aktif: %s", m.Name),
	})
	return s.GetByID(ctx, m.ID)
}

func (s *activeClassService) Update(ctx context.Context, id uint, req dto.ActiveClassUpdateReq) (*dto.ActiveClassRes, error) {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, helper.NewNotFoundError("Kelas aktif tidak ditemukan")
	}
	v := helper.NewValidationError()
	if req.Name == "" {
		v.Add("name", "Nama kelas wajib diisi")
	}
	if req.ClassTemplateID == 0 {
		v.Add("class_template_id", "Template kelas wajib dipilih")
	}
	if len(v.Errors) > 0 {
		return nil, v
	}
	if exists, _ := s.repo.ExistsNameInYear(ctx, existing.AcademicYearID, req.Name, existing.ID); exists {
		v.Add("name", "Nama kelas sudah digunakan di tahun ajaran tersebut")
	}
	if len(v.Errors) > 0 {
		return nil, v
	}
	mapper.UpdateReqToModel(&req, existing)
	if err := s.repo.Update(ctx, existing); err != nil {
		return nil, err
	}
	s.log(ctx, &activitylogdto.ActivityLogInput{
		Action:      "active_class.update",
		EntityType:  "active_class",
		EntityID:    &existing.ID,
		EntityLabel: existing.Name,
		Description: fmt.Sprintf("Memperbarui kelas aktif: %s", existing.Name),
	})
	res := mapper.EntityToRes(mapper.ModelToEntity(existing))
	return &res, nil
}

func (s *activeClassService) Delete(ctx context.Context, id uint) error {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return helper.NewNotFoundError("Kelas aktif tidak ditemukan")
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	s.log(ctx, &activitylogdto.ActivityLogInput{
		Action:      "active_class.delete",
		EntityType:  "active_class",
		EntityID:    &existing.ID,
		EntityLabel: existing.Name,
		Description: fmt.Sprintf("Menghapus kelas aktif: %s", existing.Name),
	})
	return nil
}

func (s *activeClassService) ToggleStatus(ctx context.Context, id uint) error {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return helper.NewNotFoundError("Kelas aktif tidak ditemukan")
	}
	if err := s.repo.ToggleStatus(ctx, id); err != nil {
		return err
	}
	s.log(ctx, &activitylogdto.ActivityLogInput{
		Action:      "active_class.update",
		EntityType:  "active_class",
		EntityID:    &existing.ID,
		EntityLabel: existing.Name,
		Description: fmt.Sprintf("Mengubah status kelas aktif: %s", existing.Name),
	})
	return nil
}

func (s *activeClassService) BulkUpsertByYear(ctx context.Context, academicYearID uint, items []dto.BulkUpsertItem) ([]dto.ActiveClassRes, error) {
	v := helper.NewValidationError()
	if academicYearID == 0 {
		v.Add("academic_year_id", "Tahun ajaran wajib dipilih")
	}
	models := make([]model.ActiveClass, 0, len(items))
	for i, it := range items {
		if it.Name == "" {
			v.Add(fmt.Sprintf("items[%d].name", i), "Nama kelas wajib diisi")
		}
		if it.ClassTemplateID == 0 {
			v.Add(fmt.Sprintf("items[%d].class_template_id", i), "Template kelas wajib dipilih")
		}
		models = append(models, *mapper.BulkItemToModel(academicYearID, it))
	}
	if len(v.Errors) > 0 {
		return nil, v
	}
	if err := s.repo.BulkUpsert(ctx, academicYearID, models); err != nil {
		return nil, err
	}
	updated, err := s.repo.FindByAcademicYear(ctx, academicYearID)
	if err != nil {
		return nil, err
	}
	entities := mapper.ModelListToEntity(updated)
	s.fillNames(ctx, entities)
	return mapper.EntitiesToRes(entities), nil
}
