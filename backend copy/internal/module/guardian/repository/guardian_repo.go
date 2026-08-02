package repository

import (
	"context"

	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/guardian/entity"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/guardian/mapper"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/guardian/model"
	"gorm.io/gorm"
)

type GuardianRepo interface {
	FindAll(ctx context.Context, page, limit int, search string) ([]entity.Guardian, int64, error)
	FindByID(ctx context.Context, id uint) (*entity.Guardian, error)
	Create(ctx context.Context, e *entity.Guardian) error
	Update(ctx context.Context, e *entity.Guardian) error
	Delete(ctx context.Context, id uint) error
}

type guardianRepoImpl struct {
	db *gorm.DB
}

func NewGuardianRepo(db *gorm.DB) GuardianRepo {
	return &guardianRepoImpl{db: db}
}

func (r *guardianRepoImpl) FindAll(ctx context.Context, page, limit int, search string) ([]entity.Guardian, int64, error) {
	q := r.db.WithContext(ctx).Model(&model.GuardianModel{})
	if search != "" {
		s := "%" + search + "%"
		q = q.Where("name LIKE ? OR email LIKE ? OR phone LIKE ?", s, s, s)
	}
	var total int64
	q.Count(&total)

	var models []model.GuardianModel
	if err := q.Offset((page - 1) * limit).Limit(limit).Order("name ASC").Find(&models).Error; err != nil {
		return nil, 0, err
	}

	entities := make([]entity.Guardian, 0, len(models))
	for i := range models {
		entities = append(entities, *mapper.ModelToEntity(&models[i]))
	}
	return entities, total, nil
}

func (r *guardianRepoImpl) FindByID(ctx context.Context, id uint) (*entity.Guardian, error) {
	var m model.GuardianModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, err
	}
	return mapper.ModelToEntity(&m), nil
}

func (r *guardianRepoImpl) Create(ctx context.Context, e *entity.Guardian) error {
	m := mapper.EntityToModel(e)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	e.ID = m.ID
	e.CreatedAt = m.CreatedAt
	e.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *guardianRepoImpl) Update(ctx context.Context, e *entity.Guardian) error {
	updates := map[string]any{
		"name":         e.Name,
		"phone":        e.Phone,
		"email":        e.Email,
		"nik":          e.NIK,
		"education":    e.Education,
		"occupation":   e.Occupation,
		"income_range": e.IncomeRange,
	}
	return r.db.WithContext(ctx).Model(&model.GuardianModel{}).Where("id = ?", e.ID).Updates(updates).Error
}

func (r *guardianRepoImpl) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.GuardianModel{}).Error
}
