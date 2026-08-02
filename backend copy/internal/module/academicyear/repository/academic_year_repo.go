package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/academicyear/mapper"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/academicyear/model"
	"gorm.io/gorm"
)

type AcademicYearRepo interface {
	FindAll(ctx context.Context, page, limit int, search, status, sort string) ([]model.AcademicYear, int, error)
	FindByID(ctx context.Context, id uint) (*model.AcademicYear, error)
	FindByIDUnscoped(ctx context.Context, id uint) (*model.AcademicYear, error)
	FindByName(ctx context.Context, name string, excludeID uint) (*model.AcademicYear, error)
	Exists(ctx context.Context, name string, excludeID uint) (bool, error)
	Create(ctx context.Context, m *model.AcademicYear) error
	Update(ctx context.Context, m *model.AcademicYear) error
	Delete(ctx context.Context, id uint) error
	Restore(ctx context.Context, id uint) error
	ToggleStatus(ctx context.Context, id uint) error
	BulkDelete(ctx context.Context, ids []uint) error
	BulkRestore(ctx context.Context, ids []uint) error
	FindNamesByIDs(ctx context.Context, ids []uint) (map[uint]string, error)
	CountSemesters(ctx context.Context, id uint) (int, error)
	CountActiveClasses(ctx context.Context, id uint) (int, error)
	CountStudents(ctx context.Context, id uint) (int, error)
	CountBillingRules(ctx context.Context, id uint) (int, error)
	CountInvoices(ctx context.Context, id uint) (int, error)
}

type academicYearRepoImpl struct {
	db *gorm.DB
}

func NewAcademicYearRepo(db *gorm.DB) AcademicYearRepo {
	return &academicYearRepoImpl{db: db}
}

func (r *academicYearRepoImpl) buildQuery(q *gorm.DB) *gorm.DB {
	return q.Model(&model.AcademicYear{})
}

func (r *academicYearRepoImpl) FindAll(ctx context.Context, page, limit int, search, status, sort string) ([]model.AcademicYear, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	q := r.buildQuery(r.db.WithContext(ctx))
	if search != "" {
		q = q.Where("name LIKE ?", "%"+search+"%")
	}
	switch status {
	case "active":
		q = q.Where("is_active = ?", true)
	case "inactive":
		q = q.Where("is_active = ?", false)
	case "trash":
		q = q.Unscoped().Where("deleted_at IS NOT NULL")
	default:
		q = q.Where("deleted_at IS NULL")
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	order := "created_at DESC"
	if sort == "name_asc" {
		order = "name ASC"
	}
	var models []model.AcademicYear
	if err := q.Order(order).Offset((page - 1) * limit).Limit(limit).Find(&models).Error; err != nil {
		return nil, 0, err
	}
	return models, int(total), nil
}

func (r *academicYearRepoImpl) FindByID(ctx context.Context, id uint) (*model.AcademicYear, error) {
	var m model.AcademicYear
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (r *academicYearRepoImpl) FindByIDUnscoped(ctx context.Context, id uint) (*model.AcademicYear, error) {
	var m model.AcademicYear
	if err := r.db.WithContext(ctx).Unscoped().Where("id = ?", id).First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (r *academicYearRepoImpl) FindByName(ctx context.Context, name string, excludeID uint) (*model.AcademicYear, error) {
	var m model.AcademicYear
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
	return &m, nil
}

func (r *academicYearRepoImpl) Exists(ctx context.Context, name string, excludeID uint) (bool, error) {
	var count int64
	q := r.db.WithContext(ctx).Model(&model.AcademicYear{}).Where("name = ?", name)
	if excludeID > 0 {
		q = q.Where("id != ?", excludeID)
	}
	if err := q.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *academicYearRepoImpl) Create(ctx context.Context, m *model.AcademicYear) error {
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return fmt.Errorf("nama tahun ajaran sudah digunakan")
		}
		return err
	}
	return nil
}

func (r *academicYearRepoImpl) Update(ctx context.Context, m *model.AcademicYear) error {
	if err := r.db.WithContext(ctx).Model(&model.AcademicYear{}).Where("id = ?", m.ID).Updates(map[string]interface{}{
		"name":       m.Name,
		"start_date": m.StartDate,
		"end_date":   m.EndDate,
	}).Error; err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return fmt.Errorf("nama tahun ajaran sudah digunakan")
		}
		return err
	}
	return nil
}

func (r *academicYearRepoImpl) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.AcademicYear{}, id).Error
}

func (r *academicYearRepoImpl) Restore(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Unscoped().Model(&model.AcademicYear{}).Where("id = ?", id).Update("deleted_at", nil).Error
}

func (r *academicYearRepoImpl) ToggleStatus(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&model.AcademicYear{}).Where("id = ?", id).
		Update("is_active", gorm.Expr("NOT is_active")).Error
}

func (r *academicYearRepoImpl) BulkDelete(ctx context.Context, ids []uint) error {
	return r.db.WithContext(ctx).Delete(&model.AcademicYear{}, ids).Error
}

func (r *academicYearRepoImpl) BulkRestore(ctx context.Context, ids []uint) error {
	return r.db.WithContext(ctx).Unscoped().Model(&model.AcademicYear{}).Where("id IN ?", ids).Update("deleted_at", nil).Error
}

func (r *academicYearRepoImpl) FindNamesByIDs(ctx context.Context, ids []uint) (map[uint]string, error) {
	var models []model.AcademicYear
	if err := r.db.WithContext(ctx).Unscoped().Where("id IN ?", ids).Find(&models).Error; err != nil {
		return nil, err
	}
	out := make(map[uint]string, len(models))
	for i := range models {
		out[models[i].ID] = models[i].Name
	}
	return out, nil
}

func (r *academicYearRepoImpl) CountSemesters(ctx context.Context, id uint) (int, error) {
	var c int64
	if err := r.db.WithContext(ctx).Model(&struct{}{}).Table("semesters").Where("academic_year_id = ? AND deleted_at IS NULL", id).Count(&c).Error; err != nil {
		return 0, err
	}
	return int(c), nil
}

func (r *academicYearRepoImpl) CountActiveClasses(ctx context.Context, id uint) (int, error) {
	var c int64
	if err := r.db.WithContext(ctx).Model(&struct{}{}).Table("active_classes").Where("academic_year_id = ? AND deleted_at IS NULL", id).Count(&c).Error; err != nil {
		return 0, err
	}
	return int(c), nil
}

func (r *academicYearRepoImpl) CountStudents(ctx context.Context, id uint) (int, error) {
	var c int64
	if err := r.db.WithContext(ctx).Model(&struct{}{}).Table("class_memberships").Where("academic_year_id = ? AND deleted_at IS NULL", id).Count(&c).Error; err != nil {
		return 0, err
	}
	return int(c), nil
}

func (r *academicYearRepoImpl) CountBillingRules(ctx context.Context, id uint) (int, error) {
	var c int64
	if err := r.db.WithContext(ctx).Model(&struct{}{}).Table("billing_rules").Where("academic_year_id = ? AND deleted_at IS NULL", id).Count(&c).Error; err != nil {
		return 0, err
	}
	return int(c), nil
}

func (r *academicYearRepoImpl) CountInvoices(ctx context.Context, id uint) (int, error) {
	var c int64
	if err := r.db.WithContext(ctx).Model(&struct{}{}).Table("invoices").Where("academic_year_id = ? AND deleted_at IS NULL", id).Count(&c).Error; err != nil {
		return 0, err
	}
	return int(c), nil
}

var _ = mapper.ModelToEntity
