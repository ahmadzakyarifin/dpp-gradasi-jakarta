package mapper

import (
	"time"

	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/classtemplate/dto"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/classtemplate/entity"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/classtemplate/model"
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

//
// =========================================
// Request -> Entity
// =========================================
//

func CreateReqToEntity(req *dto.ClassTemplateCreateReq) *entity.ClassTemplate {
	if req == nil {
		return nil
	}

	return &entity.ClassTemplate{
		MajorID:     req.MajorID,
		Name:        req.Name,
		GradeLevel:  req.GradeLevel,
		Description: req.Description,
		IsActive:    true,
	}
}

func UpdateReqToEntity(req *dto.ClassTemplateUpdateReq, tpl *entity.ClassTemplate) {
	if req == nil || tpl == nil {
		return
	}

	tpl.MajorID = req.MajorID
	tpl.Name = req.Name
	tpl.GradeLevel = req.GradeLevel
	tpl.Description = req.Description
}

//
// =========================================
// Entity -> Model
// =========================================
//

func EntityToModel(e *entity.ClassTemplate) *model.ClassTemplateModel {
	if e == nil {
		return nil
	}

	return &model.ClassTemplateModel{
		ID:          e.ID,
		MajorID:     e.MajorID,
		Name:        e.Name,
		GradeLevel:  e.GradeLevel,
		Description: e.Description,
		IsActive:    e.IsActive,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
		DeletedAt:   toGormDeletedAt(e.DeletedAt),
	}
}

//
// =========================================
// Model -> Entity
// =========================================
//

func ModelToEntity(m *model.ClassTemplateModel) *entity.ClassTemplate {
	if m == nil {
		return nil
	}

	return &entity.ClassTemplate{
		ID:          m.ID,
		MajorID:     m.MajorID,
		Name:        m.Name,
		GradeLevel:  m.GradeLevel,
		Description: m.Description,
		IsActive:    m.IsActive,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		DeletedAt:   toDeletedAtPtr(m.DeletedAt),

		MajorName: m.MajorName,
	}
}

func ModelListToEntity(models []model.ClassTemplateModel) []entity.ClassTemplate {
	list := make([]entity.ClassTemplate, 0, len(models))
	for i := range models {
		list = append(list, *ModelToEntity(&models[i]))
	}
	return list
}

//
// =========================================
// Entity -> Response
// =========================================
//

func EntityToResponse(e *entity.ClassTemplate) dto.ClassTemplateRes {
	if e == nil {
		return dto.ClassTemplateRes{}
	}

	return dto.ClassTemplateRes{
		ID:          e.ID,
		MajorID:     e.MajorID,
		Name:        e.Name,
		GradeLevel:  e.GradeLevel,
		Description: e.Description,
		IsActive:    e.IsActive,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
		DeletedAt:   e.DeletedAt,

		MajorName: e.MajorName,
	}
}

func EntityListToResponse(entities []entity.ClassTemplate) []dto.ClassTemplateRes {
	list := make([]dto.ClassTemplateRes, 0, len(entities))
	for i := range entities {
		list = append(list, EntityToResponse(&entities[i]))
	}
	return list
}
