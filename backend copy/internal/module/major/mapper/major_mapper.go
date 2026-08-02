package mapper

import (
	"time"

	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/major/dto"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/major/entity"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/major/model"
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

func CreateReqToEntity(req *dto.MajorCreateReq) *entity.Major {
	if req == nil {
		return nil
	}

	return &entity.Major{
		Code:     req.Code,
		Name:     req.Name,
		IsActive: true,
	}
}

func UpdateReqToEntity(req *dto.MajorUpdateReq, major *entity.Major) {
	if req == nil || major == nil {
		return
	}

	major.Code = req.Code
	major.Name = req.Name
}

//
// =========================================
// Entity -> Model
// =========================================
//

func EntityToModel(e *entity.Major) *model.Major {
	if e == nil {
		return nil
	}

	var codePtr *string
	if e.Code != "" {
		c := e.Code
		codePtr = &c
	}

	return &model.Major{
		ID:        e.ID,
		Code:      codePtr,
		Name:      e.Name,
		IsActive:  e.IsActive,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
		DeletedAt: toGormDeletedAt(e.DeletedAt),
	}
}

//
// =========================================
// Model -> Entity
// =========================================
//

func ModelToEntity(m *model.Major) *entity.Major {
	if m == nil {
		return nil
	}

	code := ""
	if m.Code != nil {
		code = *m.Code
	}

	return &entity.Major{
		ID:                m.ID,
		Code:              code,
		Name:              m.Name,
		IsActive:          m.IsActive,
		CreatedAt:         m.CreatedAt,
		UpdatedAt:         m.UpdatedAt,
		DeletedAt:         toDeletedAtPtr(m.DeletedAt),
		AcademicYearCount: m.AcademicYearCount,
		ClassCount:        m.ClassCount,
		StudentCount:      m.StudentCount,
	}
}

func ModelListToEntity(models []model.Major) []entity.Major {
	list := make([]entity.Major, 0, len(models))

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

func EntityToResponse(e *entity.Major) *dto.MajorRes {
	if e == nil {
		return nil
	}

	return &dto.MajorRes{
		ID:                e.ID,
		Code:              e.Code,
		Name:              e.Name,
		IsActive:          e.IsActive,
		AcademicYearCount: e.AcademicYearCount,
		ClassCount:        e.ClassCount,
		StudentCount:      e.StudentCount,
		CreatedAt:         e.CreatedAt,
		UpdatedAt:         e.UpdatedAt,
		DeletedAt:         e.DeletedAt,
	}
}

func EntityListToResponse(entities []entity.Major) []dto.MajorRes {
	list := make([]dto.MajorRes, 0, len(entities))

	for i := range entities {
		list = append(list, *EntityToResponse(&entities[i]))
	}

	return list
}

//
// =========================================
// Request -> Status Entity
// =========================================
//

func StatusReqToEntity(req *dto.MajorStatusReq, major *entity.Major) {
	if req == nil || major == nil {
		return
	}

	major.IsActive = req.IsActive
}
