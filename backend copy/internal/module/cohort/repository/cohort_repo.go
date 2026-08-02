package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/cohort/entity"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/cohort/mapper"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/cohort/model"
	"gorm.io/gorm"
)

type CohortRepo interface {
	FindAll(ctx context.Context, page, limit int, search, status, sort string) ([]entity.Cohort, int, error)
	FindByID(ctx context.Context, id uint) (*entity.Cohort, error)
	FindByIDUnscoped(ctx context.Context, id uint) (*entity.Cohort, error)
	FindByName(ctx context.Context, name string, excludeID uint) (*entity.Cohort, error)
	Create(ctx context.Context, e *entity.Cohort) error
	Update(ctx context.Context, e *entity.Cohort) error
	Delete(ctx context.Context, id uint) error
	Restore(ctx context.Context, id uint) error
	ToggleStatus(ctx context.Context, id uint) error
	BulkDelete(ctx context.Context, ids []uint) error
	BulkRestore(ctx context.Context, ids []uint) error
	Exists(ctx context.Context, name string, excludeID uint) (bool, error)
	CountStudents(ctx context.Context, cohortID uint) (int, error)
	CountBillingRules(ctx context.Context, cohortID uint) (int, error)
	FindNamesByIDs(ctx context.Context, ids []uint) (map[uint]string, error)
}

type cohortRepoImpl struct {
	db *gorm.DB
}

func NewCohortRepo(db *gorm.DB) CohortRepo {
	return &cohortRepoImpl{db: db}
}

func (r *cohortRepoImpl) FindAll(ctx context.Context, page, limit int, search, status, sort string) ([]entity.Cohort, int, error) {
	var models []model.Cohort
	query := r.db.WithContext(ctx).Model(&model.Cohort{})

	if search != "" {
		s := "%" + search + "%"
		query = query.Where("name LIKE ? OR description LIKE ?", s, s)
	}

	switch status {
	case "trash":
		query = query.Unscoped().Where("deleted_at IS NOT NULL")
	case "inactive":
		query = query.Where("is_active = ?", false)
	case "active":
		query = query.Where("is_active = ?", true)
	}

	switch sort {
	case "name_desc", "year_desc":
		query = query.Order("name DESC")
	case "name_asc", "year_asc":
		query = query.Order("name ASC")
	case "created_desc":
		query = query.Order("created_at DESC")
	case "created_asc", "created_at_asc":
		query = query.Order("created_at ASC")
	default:
		query = query.Order("created_at DESC")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.
		Limit(limit).
		Offset((page - 1) * limit).
		Find(&models).Error; err != nil {
		return nil, int(total), err
	}

	entities := make([]entity.Cohort, 0, len(models))
	for i := range models {
		entities = append(entities, *mapper.ModelToEntity(&models[i]))
	}
	return entities, int(total), nil
}

func (r *cohortRepoImpl) FindByID(ctx context.Context, id uint) (*entity.Cohort, error) {
	var m model.Cohort
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&m).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return mapper.ModelToEntity(&m), nil
}

func (r *cohortRepoImpl) FindByIDUnscoped(ctx context.Context, id uint) (*entity.Cohort, error) {
	var m model.Cohort
	err := r.db.WithContext(ctx).Unscoped().
		Where("id = ?", id).
		First(&m).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return mapper.ModelToEntity(&m), nil
}

func (r *cohortRepoImpl) FindByName(ctx context.Context, name string, excludeID uint) (*entity.Cohort, error) {
	var m model.Cohort
	q := r.db.WithContext(ctx).Where("name = ?", name)
	if excludeID > 0 {
		q = q.Where("id != ?", excludeID)
	}
	if err := q.First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return mapper.ModelToEntity(&m), nil
}

func (r *cohortRepoImpl) Create(ctx context.Context, e *entity.Cohort) error {
	m := mapper.EntityToModel(e)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return fmt.Errorf("nama angkatan sudah digunakan")
		}
		return err
	}
	e.ID = m.ID
	e.CreatedAt = m.CreatedAt
	e.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *cohortRepoImpl) Update(ctx context.Context, e *entity.Cohort) error {
	m := mapper.EntityToModel(e)
	err := r.db.WithContext(ctx).
		Model(&model.Cohort{}).
		Where("id = ?", e.ID).
		Updates(map[string]interface{}{
			"name":        m.Name,
			"start_date":  m.StartDate,
			"end_date":    m.EndDate,
			"description": m.Description,
			"updated_at":  m.UpdatedAt,
		}).Error
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return fmt.Errorf("nama angkatan sudah digunakan")
		}
		return err
	}
	return nil
}

func (r *cohortRepoImpl) Delete(ctx context.Context, id uint) error {
	res := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		Delete(&model.Cohort{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("angkatan tidak ditemukan atau sudah terhapus")
	}
	return nil
}

func (r *cohortRepoImpl) Restore(ctx context.Context, id uint) error {
	res := r.db.WithContext(ctx).
		Unscoped().
		Model(&model.Cohort{}).
		Where("id = ? AND deleted_at IS NOT NULL", id).
		Update("deleted_at", nil)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("angkatan tidak ditemukan di riwayat penghapusan")
	}
	return nil
}

func (r *cohortRepoImpl) ToggleStatus(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).
		Model(&model.Cohort{}).
		Where("id = ?", id).
		Update("is_active", gorm.Expr("NOT is_active")).Error
}

func (r *cohortRepoImpl) BulkDelete(ctx context.Context, ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	res := r.db.WithContext(ctx).
		Where("id IN ? AND deleted_at IS NULL", ids).
		Delete(&model.Cohort{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("tidak ada data yang dapat dihapus")
	}
	return nil
}

func (r *cohortRepoImpl) BulkRestore(ctx context.Context, ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	res := r.db.WithContext(ctx).
		Unscoped().
		Model(&model.Cohort{}).
		Where("id IN ? AND deleted_at IS NOT NULL", ids).
		Update("deleted_at", nil)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("tidak ada data yang perlu dipulihkan")
	}
	return nil
}

func (r *cohortRepoImpl) Exists(ctx context.Context, name string, excludeID uint) (bool, error) {
	q := r.db.WithContext(ctx).
		Unscoped().
		Model(&model.Cohort{}).
		Where("name = ?", name)
	if excludeID > 0 {
		q = q.Where("id != ?", excludeID)
	}
	var count int64
	if err := q.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *cohortRepoImpl) CountStudents(ctx context.Context, cohortID uint) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Table("students").
		Where("cohort_id = ? AND deleted_at IS NULL", cohortID).
		Count(&count).Error
	return int(count), err
}

func (r *cohortRepoImpl) CountBillingRules(ctx context.Context, cohortID uint) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Table("billing_rules").
		Where("cohort_id = ? AND deleted_at IS NULL", cohortID).
		Count(&count).Error
	return int(count), err
}

func (r *cohortRepoImpl) FindNamesByIDs(ctx context.Context, ids []uint) (map[uint]string, error) {
	type row struct {
		ID   uint   `gorm:"column:id"`
		Name string `gorm:"column:name"`
	}
	var rows []row
	if err := r.db.WithContext(ctx).
		Unscoped().
		Table("cohorts").
		Select("id, name").
		Where("id IN ?", ids).
		Scan(&rows).Error; err != nil {
		return nil, err
	}
	m := make(map[uint]string, len(rows))
	for _, r := range rows {
		m[r.ID] = r.Name
	}
	return m, nil
}
