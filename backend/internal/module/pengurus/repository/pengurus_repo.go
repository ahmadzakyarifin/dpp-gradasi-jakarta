package repository

import (
	"strings"

	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/pengurus/dto"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/pengurus/entity"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/pengurus/mapper"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/pengurus/model"
	"gorm.io/gorm"
)

type PengurusRepo interface {
	FindAll(query dto.PengurusQuery, adminMode bool) ([]entity.Pengurus, int64, error)
	GetRegions() ([]entity.Pengurus, error)
	FindByID(id uint) (*entity.Pengurus, error)
	Create(pengurus *entity.Pengurus) error
	Update(pengurus *entity.Pengurus) error
	SoftDelete(id uint) error
	Restore(id uint) error
	BulkSoftDelete(ids []uint) error
	BulkRestore(ids []uint) error
	Reorder(ids []uint) error
}

type pengurusRepo struct {
	db *gorm.DB
}

func NewPengurusRepo(db *gorm.DB) PengurusRepo {
	return &pengurusRepo{db: db}
}

func (r *pengurusRepo) FindAll(q dto.PengurusQuery, adminMode bool) ([]entity.Pengurus, int64, error) {
	var results []model.Pengurus
	var total int64

	db := r.db.Model(&model.Pengurus{})

	if !adminMode {
		db = db.Where("is_active = ?", true)
		db = db.Where("deleted_at IS NULL")
	} else {
		if q.Trashed {
			db = db.Unscoped().Where("deleted_at IS NOT NULL")
		} else {
			db = db.Where("deleted_at IS NULL")
		}

		switch q.Status {
		case "active":
			db = db.Where("is_active = ?", true)
		case "inactive":
			db = db.Where("is_active = ?", false)
		}
	}

	if q.Level != "" {
		db = db.Where("level = ?", q.Level)
	}
	if q.Provinsi != "" {
		db = db.Where("provinsi = ?", q.Provinsi)
	}
	if q.Kabupaten != "" {
		db = db.Where("kabupaten = ?", q.Kabupaten)
	}

	if q.Search != "" {
		search := "%" + strings.TrimSpace(q.Search) + "%"
		db = db.Where("name LIKE ? OR role LIKE ?", search, search)
	}

	db.Count(&total)

	switch q.Sort {
	case "name_asc":
		db = db.Order("name ASC")
	case "name_desc":
		db = db.Order("name DESC")
	default: // sort_order
		db = db.Order("sort_order ASC, name ASC")
	}

	page := q.Page
	if page <= 0 {
		page = 1
	}
	limit := q.Limit
	if limit <= 0 {
		limit = 10
	}
	offset := (page - 1) * limit

	err := db.Limit(limit).Offset(offset).Find(&results).Error
	return mapper.ModelListToEntity(results), total, err
}

func (r *pengurusRepo) GetRegions() ([]entity.Pengurus, error) {
	var results []model.Pengurus
	// Hanya ambil data unik provinsi & kabupaten dari pengurus aktif
	err := r.db.Select("provinsi, kabupaten").
		Where("is_active = ? AND deleted_at IS NULL", true).
		Group("provinsi, kabupaten").
		Find(&results).Error
	return mapper.ModelListToEntity(results), err
}

func (r *pengurusRepo) FindByID(id uint) (*entity.Pengurus, error) {
	var p model.Pengurus
	err := r.db.First(&p, id).Error
	if err != nil {
		return nil, err
	}
	return mapper.ModelToEntity(&p), nil
}

func (r *pengurusRepo) Create(pengurus *entity.Pengurus) error {
	m := mapper.EntityToModel(pengurus)
	if err := r.db.Create(m).Error; err != nil {
		return err
	}
	pengurus.ID = m.ID
	return nil
}

func (r *pengurusRepo) Update(pengurus *entity.Pengurus) error {
	return r.db.Save(mapper.EntityToModel(pengurus)).Error
}

func (r *pengurusRepo) SoftDelete(id uint) error {
	return r.db.Delete(&model.Pengurus{}, id).Error
}

func (r *pengurusRepo) Restore(id uint) error {
	return r.db.Unscoped().Model(&model.Pengurus{}).Where("id = ?", id).Update("deleted_at", nil).Error
}

func (r *pengurusRepo) BulkSoftDelete(ids []uint) error {
	return r.db.Where("id IN ?", ids).Delete(&model.Pengurus{}).Error
}

func (r *pengurusRepo) BulkRestore(ids []uint) error {
	return r.db.Unscoped().Model(&model.Pengurus{}).Where("id IN ?", ids).Update("deleted_at", nil).Error
}

func (r *pengurusRepo) Reorder(ids []uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for i, id := range ids {
			if err := tx.Model(&model.Pengurus{}).Where("id = ?", id).Update("sort_order", i).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
