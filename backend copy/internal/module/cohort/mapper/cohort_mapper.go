package mapper

import (
	"time"

	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/cohort/dto"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/cohort/entity"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/cohort/model"
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

func CreateReqToEntity(req *dto.CohortCreateReq) *entity.Cohort {
	if req == nil {
		return nil
	}

	return &entity.Cohort{
		Name:        req.Name,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		Description: req.Description,
		IsActive:    true,
	}
}

func UpdateReqToEntity(req *dto.CohortUpdateReq, c *entity.Cohort) {
	if req == nil || c == nil {
		return
	}

	c.Name = req.Name
	c.StartDate = req.StartDate
	c.EndDate = req.EndDate
	c.Description = req.Description
}

//
// =========================================
// Entity -> Model
// =========================================
//

func EntityToModel(e *entity.Cohort) *model.Cohort {
	if e == nil {
		return nil
	}

	return &model.Cohort{
		ID:          e.ID,
		Name:        e.Name,
		StartDate:   e.StartDate,
		EndDate:     e.EndDate,
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

func ModelToEntity(m *model.Cohort) *entity.Cohort {
	if m == nil {
		return nil
	}

	return &entity.Cohort{
		ID:           m.ID,
		Name:         m.Name,
		StartDate:    m.StartDate,
		EndDate:      m.EndDate,
		Description:  m.Description,
		IsActive:     m.IsActive,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
		DeletedAt:    toDeletedAtPtr(m.DeletedAt),
		StudentCount: m.StudentCount,
	}
}

func ModelListToEntity(models []model.Cohort) []entity.Cohort {
	list := make([]entity.Cohort, 0, len(models))
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

func EntityToRes(e *entity.Cohort) dto.CohortRes {
	if e == nil {
		return dto.CohortRes{}
	}

	return dto.CohortRes{
		ID:           e.ID,
		Name:         e.Name,
		StartDate:    e.StartDate,
		EndDate:      e.EndDate,
		Description:  e.Description,
		IsActive:     e.IsActive,
		CreatedAt:    e.CreatedAt,
		UpdatedAt:    e.UpdatedAt,
		DeletedAt:    e.DeletedAt,
		StudentCount: e.StudentCount,
	}
}

func EntitiesToRes(entities []entity.Cohort) []dto.CohortRes {
	list := make([]dto.CohortRes, 0, len(entities))
	for i := range entities {
		list = append(list, EntityToRes(&entities[i]))
	}
	return list
}
