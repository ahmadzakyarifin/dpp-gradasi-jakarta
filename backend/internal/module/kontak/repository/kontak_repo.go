package repository

import (
	"strings"

	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/kontak/dto"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/kontak/entity"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/kontak/mapper"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/kontak/model"
	"gorm.io/gorm"
)

type KontakRepo interface {
	FindAll(q dto.KontakQuery) ([]entity.PesanKontak, int64, error)
	FindByID(id uint) (*entity.PesanKontak, error)
	Create(p *entity.PesanKontak) error
	MarkAsRead(id uint) error
	Delete(id uint) error
	Restore(id uint) error
	BulkSoftDelete(ids []uint) error
	BulkRestore(ids []uint) error
}

type kontakRepo struct{ db *gorm.DB }

func NewKontakRepo(db *gorm.DB) KontakRepo {
	return &kontakRepo{db: db}
}

func (r *kontakRepo) FindAll(q dto.KontakQuery) ([]entity.PesanKontak, int64, error) {
	var items []model.PesanKontak
	var total int64

	db := r.db.Model(&model.PesanKontak{})

	if q.Status == "unread" {
		db = db.Where("is_read = ?", false)
	} else if q.Status == "read" {
		db = db.Where("is_read = ?", true)
	}

	if q.Search != "" {
		s := "%" + strings.TrimSpace(q.Search) + "%"
		db = db.Where("nama LIKE ? OR email LIKE ? OR subjek LIKE ? OR pesan LIKE ?", s, s, s, s)
	}

	if q.Status == "trashed" {
		db = r.db.Unscoped().Model(&model.PesanKontak{}).Where("deleted_at IS NOT NULL")
		db.Count(&total)
		err := db.Order("deleted_at DESC").Find(&items).Error
		return mapper.ModelListToEntity(items), total, err
	}

	db = db.Where("deleted_at IS NULL")
	db.Count(&total)

	switch q.Sort {
	case "oldest":
		db = db.Order("created_at ASC")
	default:
		db = db.Order("created_at DESC")
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

	err := db.Limit(limit).Offset(offset).Find(&items).Error
	return mapper.ModelListToEntity(items), total, err
}

func (r *kontakRepo) FindByID(id uint) (*entity.PesanKontak, error) {
	var p model.PesanKontak
	err := r.db.First(&p, id).Error
	if err != nil {
		return nil, err
	}
	return mapper.ModelToEntity(&p), nil
}

func (r *kontakRepo) Create(p *entity.PesanKontak) error {
	return r.db.Create(mapper.EntityToModel(p)).Error
}

func (r *kontakRepo) MarkAsRead(id uint) error {
	return r.db.Model(&model.PesanKontak{}).Where("id = ?", id).
		Update("is_read", true).Error
}

func (r *kontakRepo) Delete(id uint) error {
	return r.db.Delete(&model.PesanKontak{}, id).Error
}

func (r *kontakRepo) Restore(id uint) error {
	return r.db.Unscoped().Model(&model.PesanKontak{}).Where("id = ?", id).Update("deleted_at", nil).Error
}

func (r *kontakRepo) BulkSoftDelete(ids []uint) error {
	return r.db.Where("id IN ?", ids).Delete(&model.PesanKontak{}).Error
}

func (r *kontakRepo) BulkRestore(ids []uint) error {
	return r.db.Unscoped().Model(&model.PesanKontak{}).Where("id IN ?", ids).Update("deleted_at", nil).Error
}
