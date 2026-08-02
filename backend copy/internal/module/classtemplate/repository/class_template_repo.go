package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/classtemplate/entity"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/classtemplate/mapper"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/classtemplate/model"
	"gorm.io/gorm"
)

type ClassTemplateRepo interface {
	FindAll(ctx context.Context, page, limit int, search, status, majorID, gradeLevel, sort string) ([]entity.ClassTemplate, int, error)
	FindByID(ctx context.Context, id uint) (*entity.ClassTemplate, error)
	FindByName(ctx context.Context, name string, majorID *uint) (*entity.ClassTemplate, error)
	Create(ctx context.Context, e *entity.ClassTemplate) error
	Update(ctx context.Context, e *entity.ClassTemplate) error
	Delete(ctx context.Context, id uint) error
	Restore(ctx context.Context, id uint) error
	ToggleStatus(ctx context.Context, id uint) error
	BulkDelete(ctx context.Context, ids []uint) error
	BulkRestore(ctx context.Context, ids []uint) error
	CountActiveClasses(ctx context.Context, classTemplateID uint) (int, error)
	Exists(ctx context.Context, name string, majorID *uint, excludeID uint) (bool, error)
	FindNamesByIDs(ctx context.Context, ids []uint) (map[uint]string, error)
}

type classTemplateRepoImpl struct {
	db *gorm.DB
}

func NewClassTemplateRepo(db *gorm.DB) ClassTemplateRepo {
	return &classTemplateRepoImpl{db: db}
}

func (r *classTemplateRepoImpl) FindAll(ctx context.Context, page, limit int, search, status, majorID, gradeLevel, sort string) ([]entity.ClassTemplate, int, error) {
	var models []model.ClassTemplateModel
	query := r.db.WithContext(ctx).
		Model(&model.ClassTemplateModel{}).
		Select("class_templates.*, m.name AS major_name").
		Joins("LEFT JOIN majors AS m ON m.id = class_templates.major_id")

	if search != "" {
		s := "%" + search + "%"
		query = query.Where("class_templates.name LIKE ? OR class_templates.description LIKE ?", s, s)
	}

	switch status {
	case "active":
		query = query.Where("class_templates.is_active = ?", true)
	case "inactive":
		query = query.Where("class_templates.is_active = ?", false)
	case "trash", "deleted":
		query = query.Unscoped().Where("class_templates.deleted_at IS NOT NULL")
	}

	if majorID != "" {
		query = query.Where("class_templates.major_id = ?", majorID)
	}
	if gradeLevel != "" {
		query = query.Where("class_templates.grade_level = ?", gradeLevel)
	}

	switch sort {
	case "name_desc":
		query = query.Order("class_templates.name DESC")
	case "name_asc":
		query = query.Order("class_templates.name ASC")
	case "created_at_asc":
		query = query.Order("class_templates.created_at ASC")
	default:
		query = query.Order("class_templates.created_at DESC")
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

	entities := make([]entity.ClassTemplate, 0, len(models))
	for i := range models {
		entities = append(entities, *mapper.ModelToEntity(&models[i]))
	}
	return entities, int(total), nil
}

func (r *classTemplateRepoImpl) FindByID(ctx context.Context, id uint) (*entity.ClassTemplate, error) {
	var m model.ClassTemplateModel
	err := r.db.WithContext(ctx).
		Model(&model.ClassTemplateModel{}).
		Select("class_templates.*, m.name AS major_name").
		Joins("LEFT JOIN majors AS m ON m.id = class_templates.major_id").
		Where("class_templates.id = ?", id).
		First(&m).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return mapper.ModelToEntity(&m), nil
}

func (r *classTemplateRepoImpl) FindByName(ctx context.Context, name string, majorID *uint) (*entity.ClassTemplate, error) {
	var m model.ClassTemplateModel
	q := r.db.WithContext(ctx).Where("name = ?", name)
	if majorID != nil {
		q = q.Where("major_id = ?", *majorID)
	} else {
		q = q.Where("major_id IS NULL")
	}
	if err := q.First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return mapper.ModelToEntity(&m), nil
}

func (r *classTemplateRepoImpl) Create(ctx context.Context, e *entity.ClassTemplate) error {
	m := mapper.EntityToModel(e)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return fmt.Errorf("nama kelas tersebut sudah terdaftar pada jurusan dan tingkat yang sama")
		}
		return err
	}
	e.ID = m.ID
	e.CreatedAt = m.CreatedAt
	e.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *classTemplateRepoImpl) Update(ctx context.Context, e *entity.ClassTemplate) error {
	m := mapper.EntityToModel(e)
	err := r.db.WithContext(ctx).
		Model(&model.ClassTemplateModel{}).
		Where("id = ?", e.ID).
		Updates(map[string]interface{}{
			"name":        m.Name,
			"major_id":    m.MajorID,
			"grade_level": m.GradeLevel,
			"description": m.Description,
			"updated_at":  m.UpdatedAt,
		}).Error
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return fmt.Errorf("nama kelas tersebut sudah terdaftar pada jurusan dan tingkat yang sama")
		}
		return err
	}
	return nil
}

func (r *classTemplateRepoImpl) Delete(ctx context.Context, id uint) error {
	res := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		Delete(&model.ClassTemplateModel{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("template kelas tidak ditemukan atau sudah terhapus")
	}
	return nil
}

func (r *classTemplateRepoImpl) Restore(ctx context.Context, id uint) error {
	res := r.db.WithContext(ctx).
		Unscoped().
		Model(&model.ClassTemplateModel{}).
		Where("id = ? AND deleted_at IS NOT NULL", id).
		Update("deleted_at", nil)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("template kelas tidak ditemukan di riwayat penghapusan")
	}
	return nil
}

func (r *classTemplateRepoImpl) ToggleStatus(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).
		Model(&model.ClassTemplateModel{}).
		Where("id = ?", id).
		Update("is_active", gorm.Expr("NOT is_active")).Error
}

func (r *classTemplateRepoImpl) BulkDelete(ctx context.Context, ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	res := r.db.WithContext(ctx).
		Where("id IN ? AND deleted_at IS NULL", ids).
		Delete(&model.ClassTemplateModel{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("tidak ada data yang dapat dihapus")
	}
	return nil
}

func (r *classTemplateRepoImpl) BulkRestore(ctx context.Context, ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	res := r.db.WithContext(ctx).
		Unscoped().
		Model(&model.ClassTemplateModel{}).
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

func (r *classTemplateRepoImpl) CountActiveClasses(ctx context.Context, classTemplateID uint) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Table("active_classes").
		Where("class_template_id = ? AND deleted_at IS NULL", classTemplateID).
		Count(&count).Error
	return int(count), err
}

func (r *classTemplateRepoImpl) Exists(ctx context.Context, name string, majorID *uint, excludeID uint) (bool, error) {
	q := r.db.WithContext(ctx).
		Unscoped().
		Model(&model.ClassTemplateModel{}).
		Where("name = ?", name)
	if majorID != nil && *majorID != 0 {
		q = q.Where("major_id = ?", *majorID)
	} else {
		q = q.Where("major_id IS NULL")
	}
	if excludeID > 0 {
		q = q.Where("id != ?", excludeID)
	}
	var count int64
	if err := q.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *classTemplateRepoImpl) FindNamesByIDs(ctx context.Context, ids []uint) (map[uint]string, error) {
	type row struct {
		ID   uint   `gorm:"column:id"`
		Name string `gorm:"column:name"`
	}
	var rows []row
	if err := r.db.WithContext(ctx).
		Unscoped().
		Table("class_templates").
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
