package mapper

import (
	"time"

	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/academicyear/dto"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/academicyear/entity"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/academicyear/model"
	"gorm.io/gorm"
)

func toGormDeletedAt(t *time.Time) gorm.DeletedAt {
	if t == nil {
		return gorm.DeletedAt{}
	}
	return gorm.DeletedAt{Time: *t, Valid: true}
}

func toDeletedAtPtr(d gorm.DeletedAt) *time.Time {
	if !d.Valid {
		return nil
	}
	t := d.Time
	return &t
}

func CreateReqToEntity(req *dto.AcademicYearCreateReq) *entity.AcademicYear {
	if req == nil {
		return nil
	}
	return &entity.AcademicYear{
		Name:      req.Name,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		IsActive:  true,
	}
}

func UpdateReqToEntity(req *dto.AcademicYearUpdateReq, c *entity.AcademicYear) {
	if req == nil || c == nil {
		return
	}
	c.Name = req.Name
	c.StartDate = req.StartDate
	c.EndDate = req.EndDate
}

func EntityToModel(e *entity.AcademicYear) *model.AcademicYear {
	if e == nil {
		return nil
	}
	return &model.AcademicYear{
		ID:        e.ID,
		Name:      e.Name,
		StartDate: e.StartDate,
		EndDate:   e.EndDate,
		IsActive:  e.IsActive,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
		DeletedAt: toGormDeletedAt(e.DeletedAt),
	}
}

func ModelToEntity(m *model.AcademicYear) *entity.AcademicYear {
	if m == nil {
		return nil
	}
	return &entity.AcademicYear{
		ID:        m.ID,
		Name:      m.Name,
		StartDate: m.StartDate,
		EndDate:   m.EndDate,
		IsActive:  m.IsActive,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		DeletedAt: toDeletedAtPtr(m.DeletedAt),

		SemesterCount:    m.SemesterCount,
		ActiveClassCount: m.ActiveClassCount,
		StudentCount:     m.StudentCount,
		BillingRuleCount: m.BillingRuleCount,
		InvoiceCount:     m.InvoiceCount,
	}
}

func ModelListToEntity(models []model.AcademicYear) []entity.AcademicYear {
	list := make([]entity.AcademicYear, 0, len(models))
	for i := range models {
		list = append(list, *ModelToEntity(&models[i]))
	}
	return list
}

func EntityToRes(e *entity.AcademicYear) dto.AcademicYearRes {
	if e == nil {
		return dto.AcademicYearRes{}
	}
	return dto.AcademicYearRes{
		ID:        e.ID,
		Name:      e.Name,
		StartDate: e.StartDate,
		EndDate:   e.EndDate,
		IsActive:  e.IsActive,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
		DeletedAt: e.DeletedAt,

		SemesterCount:    e.SemesterCount,
		ActiveClassCount: e.ActiveClassCount,
		StudentCount:     e.StudentCount,
		BillingRuleCount: e.BillingRuleCount,
		InvoiceCount:     e.InvoiceCount,
	}
}

func EntitiesToRes(entities []entity.AcademicYear) []dto.AcademicYearRes {
	list := make([]dto.AcademicYearRes, 0, len(entities))
	for i := range entities {
		list = append(list, EntityToRes(&entities[i]))
	}
	return list
}
