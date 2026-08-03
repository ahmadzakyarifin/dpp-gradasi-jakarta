package mapper

import (
	"time"

	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/pengurus/dto"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/pengurus/entity"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/pengurus/model"
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

// StrPtr mengonversi string kosong menjadi nil pointer.
func StrPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// CreateReqToEntity mengonversi request create menjadi entity.
func CreateReqToEntity(req *dto.PengurusRequest, imageURL string, isActive bool) *entity.Pengurus {
	if req == nil {
		return nil
	}
	return &entity.Pengurus{
		Name:         req.Name,
		Role:         req.Role,
		Department:   StrPtr(req.Department),
		Level:        req.Level,
		Provinsi:     StrPtr(req.Provinsi),
		Kabupaten:    StrPtr(req.Kabupaten),
		ImagePath:    imageURL,
		FacebookURL:  StrPtr(req.FacebookURL),
		InstagramURL: StrPtr(req.InstagramURL),
		LinkedinURL:  StrPtr(req.LinkedinURL),
		Whatsapp:     StrPtr(req.Whatsapp),
		Periode:      req.Periode,
		SortOrder:    req.SortOrder,
		IsActive:     isActive,
	}
}

// ApplyUpdate menerapkan request update ke entity yang sudah ada.
func ApplyUpdate(e *entity.Pengurus, req *dto.PengurusRequest) {
	if e == nil || req == nil {
		return
	}
	e.Name = req.Name
	e.Role = req.Role
	e.Department = StrPtr(req.Department)
	e.Level = req.Level
	e.Provinsi = StrPtr(req.Provinsi)
	e.Kabupaten = StrPtr(req.Kabupaten)
	e.FacebookURL = StrPtr(req.FacebookURL)
	e.InstagramURL = StrPtr(req.InstagramURL)
	e.LinkedinURL = StrPtr(req.LinkedinURL)
	e.Whatsapp = StrPtr(req.Whatsapp)
	e.Periode = req.Periode
	e.SortOrder = req.SortOrder
	if req.IsActive != nil {
		e.IsActive = *req.IsActive
	}
}

// EntityToModel mengonversi entity menjadi model GORM.
func EntityToModel(e *entity.Pengurus) *model.Pengurus {
	if e == nil {
		return nil
	}
	return &model.Pengurus{
		ID:           e.ID,
		Name:         e.Name,
		Role:         e.Role,
		Department:   e.Department,
		Level:        e.Level,
		Provinsi:     e.Provinsi,
		Kabupaten:    e.Kabupaten,
		ImagePath:    e.ImagePath,
		FacebookURL:  e.FacebookURL,
		InstagramURL: e.InstagramURL,
		LinkedinURL:  e.LinkedinURL,
		Whatsapp:     e.Whatsapp,
		Periode:      e.Periode,
		SortOrder:    e.SortOrder,
		IsActive:     e.IsActive,
		CreatedAt:    e.CreatedAt,
		UpdatedAt:    e.UpdatedAt,
		DeletedAt:    toGormDeletedAt(e.DeletedAt),
	}
}

// ModelToEntity mengonversi model GORM menjadi entity.
func ModelToEntity(m *model.Pengurus) *entity.Pengurus {
	if m == nil {
		return nil
	}
	return &entity.Pengurus{
		ID:           m.ID,
		Name:         m.Name,
		Role:         m.Role,
		Department:   m.Department,
		Level:        m.Level,
		Provinsi:     m.Provinsi,
		Kabupaten:    m.Kabupaten,
		ImagePath:    m.ImagePath,
		FacebookURL:  m.FacebookURL,
		InstagramURL: m.InstagramURL,
		LinkedinURL:  m.LinkedinURL,
		Whatsapp:     m.Whatsapp,
		Periode:      m.Periode,
		SortOrder:    m.SortOrder,
		IsActive:     m.IsActive,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
		DeletedAt:    fromGormDeletedAt(m.DeletedAt),
	}
}

// ModelListToEntity mengonversi daftar model menjadi daftar entity.
func ModelListToEntity(models []model.Pengurus) []entity.Pengurus {
	list := make([]entity.Pengurus, 0, len(models))
	for i := range models {
		list = append(list, *ModelToEntity(&models[i]))
	}
	return list
}

// EntityToResponse mengonversi entity menjadi response.
func EntityToResponse(e *entity.Pengurus) dto.PengurusResponse {
	if e == nil {
		return dto.PengurusResponse{}
	}
	resp := dto.PengurusResponse{
		ID:        e.ID,
		Name:      e.Name,
		Role:      e.Role,
		Level:     e.Level,
		ImagePath: e.ImagePath,
		Periode:   e.Periode,
		SortOrder: e.SortOrder,
		IsActive:  e.IsActive,
		CreatedAt: e.CreatedAt.Format(time.RFC3339),
		UpdatedAt: e.UpdatedAt.Format(time.RFC3339),
	}
	if e.Department != nil {
		resp.Department = *e.Department
	}
	if e.Provinsi != nil {
		resp.Provinsi = *e.Provinsi
	}
	if e.Kabupaten != nil {
		resp.Kabupaten = *e.Kabupaten
	}
	if e.FacebookURL != nil {
		resp.FacebookURL = *e.FacebookURL
	}
	if e.InstagramURL != nil {
		resp.InstagramURL = *e.InstagramURL
	}
	if e.LinkedinURL != nil {
		resp.LinkedinURL = *e.LinkedinURL
	}
	if e.Whatsapp != nil {
		resp.Whatsapp = *e.Whatsapp
	}
	return resp
}

// EntityListToResponse mengonversi daftar entity menjadi daftar response.
func EntityListToResponse(entities []entity.Pengurus) []dto.PengurusResponse {
	list := make([]dto.PengurusResponse, 0, len(entities))
	for i := range entities {
		list = append(list, EntityToResponse(&entities[i]))
	}
	return list
}
