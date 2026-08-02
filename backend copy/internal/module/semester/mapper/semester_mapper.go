package mapper

import (
	"time"

	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/semester/dto"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/semester/entity"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/semester/model"
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

func CreateReqToEntity(req *dto.SemesterCreateReq) *entity.Semester {
	if req == nil {
		return nil
	}
	return &entity.Semester{
		AcademicYearID: req.AcademicYearID,
		Name:           req.Name,
		StartDate:      req.StartDate,
		EndDate:        req.EndDate,
		IsActive:       true,
	}
}

func UpdateReqToEntity(req *dto.SemesterUpdateReq, c *entity.Semester) {
	if req == nil || c == nil {
		return
	}
	c.AcademicYearID = req.AcademicYearID
	c.Name = req.Name
	c.StartDate = req.StartDate
	c.EndDate = req.EndDate
}

func EntityToModel(e *entity.Semester) *model.Semester {
	if e == nil {
		return nil
	}
	m := &model.Semester{
		ID:             e.ID,
		AcademicYearID: e.AcademicYearID,
		Name:           e.Name,
		StartDate:      e.StartDate,
		EndDate:        e.EndDate,
		IsActive:       e.IsActive,
		CreatedAt:      e.CreatedAt,
		UpdatedAt:      e.UpdatedAt,
		DeletedAt:      toGormDeletedAt(e.DeletedAt),
	}
	return m
}

func ModelToEntity(m *model.Semester) *entity.Semester {
	if m == nil {
		return nil
	}
	return &entity.Semester{
		ID:             m.ID,
		AcademicYearID: m.AcademicYearID,
		Name:           m.Name,
		StartDate:      m.StartDate,
		EndDate:        m.EndDate,
		IsActive:       m.IsActive,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
		DeletedAt:      toDeletedAtPtr(m.DeletedAt),

		ClassMembershipCount: m.ClassMembershipCount,
		BillingRuleCount:     m.BillingRuleCount,
		InvoiceCount:         m.InvoiceCount,
		BatchCount:           m.BatchCount,
	}
}

func ModelListToEntity(models []model.Semester) []entity.Semester {
	list := make([]entity.Semester, 0, len(models))
	for i := range models {
		list = append(list, *ModelToEntity(&models[i]))
	}
	return list
}

func EntityToRes(e *entity.Semester) dto.SemesterRes {
	if e == nil {
		return dto.SemesterRes{}
	}
	ayName := ""
	if e.AcademicYear != nil {
		ayName = e.AcademicYear.Name
	}
	return dto.SemesterRes{
		ID:               e.ID,
		AcademicYearID:   e.AcademicYearID,
		AcademicYearName: ayName,
		Name:             e.Name,
		StartDate:        e.StartDate,
		EndDate:          e.EndDate,
		IsActive:         e.IsActive,
		CreatedAt:        e.CreatedAt,
		UpdatedAt:        e.UpdatedAt,
		DeletedAt:        e.DeletedAt,

		ClassMembershipCount: e.ClassMembershipCount,
		BillingRuleCount:     e.BillingRuleCount,
		InvoiceCount:         e.InvoiceCount,
		BatchCount:           e.BatchCount,
	}
}

func EntitiesToRes(entities []entity.Semester) []dto.SemesterRes {
	list := make([]dto.SemesterRes, 0, len(entities))
	for i := range entities {
		list = append(list, EntityToRes(&entities[i]))
	}
	return list
}
