package mapper

import (
	"time"

	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/kontak/dto"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/kontak/entity"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/kontak/model"
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

// CreateReqToEntity mengonversi request submit menjadi entity.
func CreateReqToEntity(req *dto.KontakRequest) *entity.PesanKontak {
	if req == nil {
		return nil
	}
	return &entity.PesanKontak{
		Nama:   req.Nama,
		Email:  req.Email,
		Subjek: req.Subjek,
		Pesan:  req.Pesan,
	}
}

// EntityToModel mengonversi entity menjadi model GORM.
func EntityToModel(e *entity.PesanKontak) *model.PesanKontak {
	if e == nil {
		return nil
	}
	return &model.PesanKontak{
		ID:           e.ID,
		Nama:         e.Nama,
		Email:        e.Email,
		Subjek:       e.Subjek,
		Pesan:        e.Pesan,
		IsRead:       e.IsRead,
		ResponseNote: e.ResponseNote,
		CreatedAt:    e.CreatedAt,
		UpdatedAt:    e.UpdatedAt,
		DeletedAt:    toGormDeletedAt(e.DeletedAt),
	}
}

// ModelToEntity mengonversi model GORM menjadi entity.
func ModelToEntity(m *model.PesanKontak) *entity.PesanKontak {
	if m == nil {
		return nil
	}
	return &entity.PesanKontak{
		ID:           m.ID,
		Nama:         m.Nama,
		Email:        m.Email,
		Subjek:       m.Subjek,
		Pesan:        m.Pesan,
		IsRead:       m.IsRead,
		ResponseNote: m.ResponseNote,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
		DeletedAt:    fromGormDeletedAt(m.DeletedAt),
	}
}

// ModelListToEntity mengonversi daftar model menjadi daftar entity.
func ModelListToEntity(models []model.PesanKontak) []entity.PesanKontak {
	list := make([]entity.PesanKontak, 0, len(models))
	for i := range models {
		list = append(list, *ModelToEntity(&models[i]))
	}
	return list
}

// EntityToListItem mengonversi entity menjadi list item response.
func EntityToListItem(e *entity.PesanKontak) dto.KontakListItem {
	if e == nil {
		return dto.KontakListItem{}
	}
	return dto.KontakListItem{
		ID:        e.ID,
		Nama:      e.Nama,
		Email:     e.Email,
		Subjek:    e.Subjek,
		Pesan:     e.Pesan,
		IsRead:    e.IsRead,
		CreatedAt: e.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

// EntityListToItem mengonversi daftar entity menjadi daftar list item response.
func EntityListToItem(entities []entity.PesanKontak) []dto.KontakListItem {
	list := make([]dto.KontakListItem, 0, len(entities))
	for i := range entities {
		list = append(list, EntityToListItem(&entities[i]))
	}
	return list
}

// EntityToDetail mengonversi entity menjadi detail response.
func EntityToDetail(e *entity.PesanKontak) dto.KontakDetailResponse {
	if e == nil {
		return dto.KontakDetailResponse{}
	}
	resp := dto.KontakDetailResponse{
		ID:        e.ID,
		Nama:      e.Nama,
		Email:     e.Email,
		Subjek:    e.Subjek,
		Pesan:     e.Pesan,
		IsRead:    e.IsRead,
		CreatedAt: e.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: e.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
	if e.ResponseNote != nil {
		resp.ResponseNote = *e.ResponseNote
	}
	return resp
}
