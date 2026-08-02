package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/semester/entity"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/semester/mapper"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/semester/model"
	"gorm.io/gorm"
)

type SemesterRepo interface {
	FindAll(ctx context.Context, page, limit int, search, status string, academicYearID uint, sort string) ([]model.Semester, int, error)
	FindByID(ctx context.Context, id uint) (*model.Semester, error)
	FindByIDUnscoped(ctx context.Context, id uint) (*model.Semester, error)
	FindByName(ctx context.Context, academicYearID uint, name string, excludeID uint) (*model.Semester, error)
	Exists(ctx context.Context, academicYearID uint, name string, excludeID uint) (bool, error)
	Create(ctx context.Context, m *model.Semester) error
	Update(ctx context.Context, m *model.Semester) error
	Delete(ctx context.Context, id uint) error
	Restore(ctx context.Context, id uint) error
	ToggleStatus(ctx context.Context, id uint) error
	BulkDelete(ctx context.Context, ids []uint) error
	BulkRestore(ctx context.Context, ids []uint) error
	FindNamesByIDs(ctx context.Context, ids []uint) (map[uint]string, error)
	CountClassMemberships(ctx context.Context, id uint) (int, error)
	CountBillingRules(ctx context.Context, id uint) (int, error)
	CountInvoices(ctx context.Context, id uint) (int, error)
	CountBatches(ctx context.Context, id uint) (int, error)
}

type semesterRepoImpl struct {
	db *gorm.DB
}

func NewSemesterRepo(db *gorm.DB) SemesterRepo {
	return &semesterRepoImpl{db: db}
}

func (r *semesterRepoImpl) FindAll(ctx context.Context, page, limit int, search, status string, academicYearID uint, sort string) ([]model.Semester, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	q := r.db.WithContext(ctx).Model(&model.Semester{})
	if search != "" {
		q = q.Where("name LIKE ?", "%"+search+"%")
	}
	if academicYearID > 0 {
		q = q.Where("academic_year_id = ?", academicYearID)
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
	var models []model.Semester
	if err := q.Order(order).Offset((page - 1) * limit).Limit(limit).Find(&models).Error; err != nil {
		return nil, 0, err
	}
	return models, int(total), nil
}

func (r *semesterRepoImpl) FindByID(ctx context.Context, id uint) (*model.Semester, error) {
	var m model.Semester
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (r *semesterRepoImpl) FindByIDUnscoped(ctx context.Context, id uint) (*model.Semester, error) {
	var m model.Semester
	if err := r.db.WithContext(ctx).Unscoped().Where("id = ?", id).First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (r *semesterRepoImpl) FindByName(ctx context.Context, academicYearID uint, name string, excludeID uint) (*model.Semester, error) {
	var m model.Semester
	q := r.db.WithContext(ctx).Where("academic_year_id = ? AND name = ?", academicYearID, name)
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

func (r *semesterRepoImpl) Exists(ctx context.Context, academicYearID uint, name string, excludeID uint) (bool, error) {
	var count int64
	q := r.db.WithContext(ctx).Model(&model.Semester{}).Where("academic_year_id = ? AND name = ?", academicYearID, name)
	if excludeID > 0 {
		q = q.Where("id != ?", excludeID)
	}
	if err := q.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *semesterRepoImpl) Create(ctx context.Context, m *model.Semester) error {
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return fmt.Errorf("semester untuk tahun ajaran tersebut sudah ada")
		}
		return err
	}
	return nil
}

func (r *semesterRepoImpl) Update(ctx context.Context, m *model.Semester) error {
	if err := r.db.WithContext(ctx).Model(&model.Semester{}).Where("id = ?", m.ID).Updates(map[string]interface{}{
		"academic_year_id": m.AcademicYearID,
		"name":             m.Name,
		"start_date":       m.StartDate,
		"end_date":         m.EndDate,
	}).Error; err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return fmt.Errorf("semester untuk tahun ajaran tersebut sudah ada")
		}
		return err
	}
	return nil
}

func (r *semesterRepoImpl) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Semester{}, id).Error
}

func (r *semesterRepoImpl) Restore(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Unscoped().Model(&model.Semester{}).Where("id = ?", id).Update("deleted_at", nil).Error
}

func (r *semesterRepoImpl) ToggleStatus(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&model.Semester{}).Where("id = ?", id).
		Update("is_active", gorm.Expr("NOT is_active")).Error
}

func (r *semesterRepoImpl) BulkDelete(ctx context.Context, ids []uint) error {
	return r.db.WithContext(ctx).Delete(&model.Semester{}, ids).Error
}

func (r *semesterRepoImpl) BulkRestore(ctx context.Context, ids []uint) error {
	return r.db.WithContext(ctx).Unscoped().Model(&model.Semester{}).Where("id IN ?", ids).Update("deleted_at", nil).Error
}

func (r *semesterRepoImpl) FindNamesByIDs(ctx context.Context, ids []uint) (map[uint]string, error) {
	var models []model.Semester
	if err := r.db.WithContext(ctx).Unscoped().Where("id IN ?", ids).Find(&models).Error; err != nil {
		return nil, err
	}
	out := make(map[uint]string, len(models))
	for i := range models {
		out[models[i].ID] = models[i].Name
	}
	return out, nil
}

func (r *semesterRepoImpl) CountClassMemberships(ctx context.Context, id uint) (int, error) {
	var c int64
	if err := r.db.WithContext(ctx).Model(&struct{}{}).Table("class_memberships").Where("semester_id = ? AND deleted_at IS NULL", id).Count(&c).Error; err != nil {
		return 0, err
	}
	return int(c), nil
}

func (r *semesterRepoImpl) CountBillingRules(ctx context.Context, id uint) (int, error) {
	var c int64
	if err := r.db.WithContext(ctx).Model(&struct{}{}).Table("billing_rules").Where("semester_id = ? AND deleted_at IS NULL", id).Count(&c).Error; err != nil {
		return 0, err
	}
	return int(c), nil
}

func (r *semesterRepoImpl) CountInvoices(ctx context.Context, id uint) (int, error) {
	var c int64
	if err := r.db.WithContext(ctx).Model(&struct{}{}).Table("invoices").Where("semester_id = ? AND deleted_at IS NULL", id).Count(&c).Error; err != nil {
		return 0, err
	}
	return int(c), nil
}

func (r *semesterRepoImpl) CountBatches(ctx context.Context, id uint) (int, error) {
	var c int64
	if err := r.db.WithContext(ctx).Model(&struct{}{}).Table("invoice_generation_batches").Where("semester_id = ? AND deleted_at IS NULL", id).Count(&c).Error; err != nil {
		return 0, err
	}
	return int(c), nil
}

var _ = mapper.ModelToEntity
var _ = entity.Semester{}
