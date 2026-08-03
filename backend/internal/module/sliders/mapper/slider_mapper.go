package mapper

import (
	"time"

	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/sliders/dto"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/sliders/entity"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/sliders/model"
	"gorm.io/gorm"
)

func toGormDeletedAt(t *time.Time) gorm.DeletedAt {
	if t == nil {
		return gorm.DeletedAt{}
	}
	return gorm.DeletedAt{Time: *t, Valid: true}
}

func fromGormDeletedAt(d gorm.DeletedAt) *time.Time {
	if !d.Valid {
		return nil
	}
	t := d.Time
	return &t
}

// CreateReqToEntity mengonversi request create menjadi entity.
func CreateReqToEntity(req *dto.SliderRequest) *entity.Slider {
	if req == nil {
		return nil
	}
	return &entity.Slider{
		Title:     req.Title,
		Subtitle:  req.Subtitle,
		Tag:       req.Tag,
		IsNew:     req.IsNew,
		EventDate: req.EventDate,
		Location:  req.Location,
		ImagePath: req.ImagePath,
		LinkURL:   req.LinkURL,
		SortOrder: req.SortOrder,
		IsActive:  req.IsActive,
	}
}

// UpdateReqToEntity menerapkan field request ke entity existing.
func UpdateReqToEntity(req *dto.SliderRequest, sl *entity.Slider) {
	if req == nil || sl == nil {
		return
	}
	sl.Title = req.Title
	sl.Subtitle = req.Subtitle
	sl.Tag = req.Tag
	sl.IsNew = req.IsNew
	sl.EventDate = req.EventDate
	sl.Location = req.Location
	sl.ImagePath = req.ImagePath
	sl.LinkURL = req.LinkURL
	sl.SortOrder = req.SortOrder
	sl.IsActive = req.IsActive
}

// EntityToModel mengonversi entity menjadi model GORM.
func EntityToModel(e *entity.Slider) *model.Slider {
	if e == nil {
		return nil
	}
	return &model.Slider{
		ID:        e.ID,
		Title:     e.Title,
		Subtitle:  e.Subtitle,
		Tag:       e.Tag,
		IsNew:     e.IsNew,
		EventDate: e.EventDate,
		Location:  e.Location,
		ImagePath: e.ImagePath,
		LinkURL:   e.LinkURL,
		SortOrder: e.SortOrder,
		IsActive:  e.IsActive,
		CreatedBy: e.CreatedBy,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
		DeletedAt: toGormDeletedAt(e.DeletedAt),
	}
}

// ModelToEntity mengonversi model GORM menjadi entity.
func ModelToEntity(m *model.Slider) *entity.Slider {
	if m == nil {
		return nil
	}
	return &entity.Slider{
		ID:        m.ID,
		Title:     m.Title,
		Subtitle:  m.Subtitle,
		Tag:       m.Tag,
		IsNew:     m.IsNew,
		EventDate: m.EventDate,
		Location:  m.Location,
		ImagePath: m.ImagePath,
		LinkURL:   m.LinkURL,
		SortOrder: m.SortOrder,
		IsActive:  m.IsActive,
		CreatedBy: m.CreatedBy,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		DeletedAt: fromGormDeletedAt(m.DeletedAt),
	}
}

// EntityToResponse mengonversi entity menjadi response DTO.
func EntityToResponse(e *entity.Slider) dto.SliderResponse {
	if e == nil {
		return dto.SliderResponse{}
	}
	r := dto.SliderResponse{
		ID:        e.ID,
		Title:     e.Title,
		IsNew:     e.IsNew,
		ImagePath: e.ImagePath,
		SortOrder: e.SortOrder,
		IsActive:  e.IsActive,
	}
	if e.Subtitle != nil {
		r.Subtitle = *e.Subtitle
	}
	if e.Tag != nil {
		r.Tag = *e.Tag
	}
	if e.EventDate != nil {
		r.EventDate = *e.EventDate
	}
	if e.Location != nil {
		r.Location = *e.Location
	}
	if e.LinkURL != nil {
		r.LinkURL = *e.LinkURL
	}
	return r
}

// EntityListToResponse mengonversi daftar entity menjadi daftar response DTO.
func EntityListToResponse(entities []entity.Slider) []dto.SliderResponse {
	list := make([]dto.SliderResponse, 0, len(entities))
	for i := range entities {
		list = append(list, EntityToResponse(&entities[i]))
	}
	return list
}
