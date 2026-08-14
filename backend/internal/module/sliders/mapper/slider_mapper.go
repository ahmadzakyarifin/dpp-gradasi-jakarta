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
		SortOrder: req.SortOrder,
		IsPublished:  req.IsPublished,
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
	sl.SortOrder = req.SortOrder
	sl.IsPublished = req.IsPublished
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
		SortOrder: e.SortOrder,
		IsPublished:  e.IsPublished,
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
		SortOrder: m.SortOrder,
		IsPublished:  m.IsPublished,
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
		IsPublished:  e.IsPublished,
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
	// Always map updated_at (it's not a pointer, but we format it)
	r.UpdatedAt = e.UpdatedAt.Format(time.RFC3339)

	if e.DeletedAt != nil {
		r.DeletedAt = e.DeletedAt.Format(time.RFC3339)
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
