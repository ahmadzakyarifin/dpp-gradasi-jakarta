package repository

import (
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/sliders/model"
	"gorm.io/gorm"
)

type SlidersRepo interface {
	FindAll(publishedOnly bool) ([]model.Slider, error)
	FindByID(id uint) (*model.Slider, error)
	Create(slider *model.Slider) error
	Update(slider *model.Slider) error
	Delete(id uint) error
	Restore(id uint) error
	BulkSoftDelete(ids []uint) error
	BulkRestore(ids []uint) error
	Reorder(ids []uint) error
}

type slidersRepo struct {
	db *gorm.DB
}

func NewSlidersRepo(db *gorm.DB) SlidersRepo {
	return &slidersRepo{db: db}
}

func (r *slidersRepo) FindAll(publishedOnly bool) ([]model.Slider, error) {
	var sliders []model.Slider
	q := r.db.Order("sort_order ASC, created_at DESC")
	if publishedOnly {
		q = q.Where("is_published = ?", true)
	}
	err := q.Find(&sliders).Error
	return sliders, err
}

func (r *slidersRepo) FindByID(id uint) (*model.Slider, error) {
	var s model.Slider
	err := r.db.First(&s, id).Error
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *slidersRepo) Create(slider *model.Slider) error {
	return r.db.Create(slider).Error
}

func (r *slidersRepo) Update(slider *model.Slider) error {
	return r.db.Save(slider).Error
}

func (r *slidersRepo) Delete(id uint) error {
	return r.db.Delete(&model.Slider{}, id).Error
}

func (r *slidersRepo) Reorder(ids []uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for i, id := range ids {
			if err := tx.Model(&model.Slider{}).Where("id = ?", id).Update("sort_order", i).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *slidersRepo) Restore(id uint) error {
	return r.db.Unscoped().Model(&model.Slider{}).Where("id = ?", id).Update("deleted_at", nil).Error
}

func (r *slidersRepo) BulkSoftDelete(ids []uint) error {
	return r.db.Where("id IN ?", ids).Delete(&model.Slider{}).Error
}

func (r *slidersRepo) BulkRestore(ids []uint) error {
	return r.db.Unscoped().Model(&model.Slider{}).Where("id IN ?", ids).Update("deleted_at", nil).Error
}
