package mapper

import (
	"time"

	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/activeclass/dto"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/activeclass/entity"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/activeclass/model"
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

func CreateReqToModel(req *dto.ActiveClassCreateReq) *model.ActiveClass {
	if req == nil {
		return nil
	}
	return &model.ActiveClass{
		AcademicYearID:      req.AcademicYearID,
		ClassTemplateID:     req.ClassTemplateID,
		Name:                req.Name,
		HomeroomNumber:      req.HomeroomNumber,
		HomeroomTeacherName: req.HomeroomTeacherName,
		Room:                req.Room,
		Capacity:            req.Capacity,
		IsActive:            req.IsActive,
	}
}

func UpdateReqToModel(req *dto.ActiveClassUpdateReq, m *model.ActiveClass) {
	if req == nil || m == nil {
		return
	}
	m.ClassTemplateID = req.ClassTemplateID
	m.Name = req.Name
	m.HomeroomNumber = req.HomeroomNumber
	m.HomeroomTeacherName = req.HomeroomTeacherName
	m.Room = req.Room
	m.Capacity = req.Capacity
	m.IsActive = req.IsActive
}

func BulkItemToModel(academicYearID uint, item dto.BulkUpsertItem) *model.ActiveClass {
	return &model.ActiveClass{
		ID:                  item.ID,
		AcademicYearID:      academicYearID,
		ClassTemplateID:     item.ClassTemplateID,
		Name:                item.Name,
		HomeroomNumber:      item.HomeroomNumber,
		HomeroomTeacherName: item.HomeroomTeacherName,
		Room:                item.Room,
		Capacity:            item.Capacity,
		IsActive:            item.IsActive,
	}
}

func ModelToEntity(m *model.ActiveClass) *entity.ActiveClass {
	if m == nil {
		return nil
	}
	return &entity.ActiveClass{
		ID:                  m.ID,
		AcademicYearID:      m.AcademicYearID,
		ClassTemplateID:     m.ClassTemplateID,
		Name:                m.Name,
		HomeroomNumber:      m.HomeroomNumber,
		HomeroomTeacherName: m.HomeroomTeacherName,
		Room:                m.Room,
		Capacity:            m.Capacity,
		StudentCount:        m.StudentCount,
		IsActive:            m.IsActive,
		CreatedAt:           m.CreatedAt,
		UpdatedAt:           m.UpdatedAt,
		DeletedAt:           toDeletedAtPtr(m.DeletedAt),
	}
}

func ModelListToEntity(models []model.ActiveClass) []entity.ActiveClass {
	list := make([]entity.ActiveClass, 0, len(models))
	for i := range models {
		list = append(list, *ModelToEntity(&models[i]))
	}
	return list
}

func EntityToRes(e *entity.ActiveClass) dto.ActiveClassRes {
	if e == nil {
		return dto.ActiveClassRes{}
	}
	ayName, ctName := "", ""
	if e.AcademicYear != nil {
		ayName = e.AcademicYear.Name
	}
	if e.ClassTemplate != nil {
		ctName = e.ClassTemplate.Name
	}
	return dto.ActiveClassRes{
		ID:                  e.ID,
		AcademicYearID:      e.AcademicYearID,
		AcademicYearName:    ayName,
		ClassTemplateID:     e.ClassTemplateID,
		ClassTemplateName:   ctName,
		Name:                e.Name,
		HomeroomNumber:      e.HomeroomNumber,
		HomeroomTeacherName: e.HomeroomTeacherName,
		Room:                e.Room,
		Capacity:            e.Capacity,
		StudentCount:        e.StudentCount,
		IsActive:            e.IsActive,
		CreatedAt:           e.CreatedAt,
		UpdatedAt:           e.UpdatedAt,
		DeletedAt:           e.DeletedAt,
	}
}

func EntitiesToRes(entities []entity.ActiveClass) []dto.ActiveClassRes {
	list := make([]dto.ActiveClassRes, 0, len(entities))
	for i := range entities {
		list = append(list, EntityToRes(&entities[i]))
	}
	return list
}
