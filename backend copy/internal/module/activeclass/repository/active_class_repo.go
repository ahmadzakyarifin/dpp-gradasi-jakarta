package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/activeclass/mapper"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/activeclass/model"
	"gorm.io/gorm"
)

type ActiveClassRepo interface {
	FindAll(ctx context.Context, page, limit int, search, status string, academicYearID uint, sort string) ([]model.ActiveClass, int, error)
	FindByID(ctx context.Context, id uint) (*model.ActiveClass, error)
	FindByAcademicYear(ctx context.Context, academicYearID uint) ([]model.ActiveClass, error)
	FindByNameInYear(ctx context.Context, academicYearID uint, name string, excludeID uint) (*model.ActiveClass, error)
	ExistsNameInYear(ctx context.Context, academicYearID uint, name string, excludeID uint) (bool, error)
	ExistsTemplateHomeroom(ctx context.Context, academicYearID, classTemplateID uint, homeroom *string, excludeID uint) (bool, error)
	Create(ctx context.Context, m *model.ActiveClass) error
	Update(ctx context.Context, m *model.ActiveClass) error
	Delete(ctx context.Context, id uint) error
	ToggleStatus(ctx context.Context, id uint) error
	BulkUpsert(ctx context.Context, academicYearID uint, models []model.ActiveClass) error
}

type activeClassRepoImpl struct {
	db *gorm.DB
}

func NewActiveClassRepo(db *gorm.DB) ActiveClassRepo {
	return &activeClassRepoImpl{db: db}
}

func (r *activeClassRepoImpl) FindAll(ctx context.Context, page, limit int, search, status string, academicYearID uint, sort string) ([]model.ActiveClass, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	q := r.db.WithContext(ctx).Model(&model.ActiveClass{})
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
	case "deleted":
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
	var models []model.ActiveClass
	if err := q.Order(order).Offset((page - 1) * limit).Limit(limit).Find(&models).Error; err != nil {
		return nil, 0, err
	}
	return models, int(total), nil
}

func (r *activeClassRepoImpl) FindByID(ctx context.Context, id uint) (*model.ActiveClass, error) {
	var m model.ActiveClass
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (r *activeClassRepoImpl) FindByAcademicYear(ctx context.Context, academicYearID uint) ([]model.ActiveClass, error) {
	var models []model.ActiveClass
	if err := r.db.WithContext(ctx).Where("academic_year_id = ? AND deleted_at IS NULL", academicYearID).Find(&models).Error; err != nil {
		return nil, err
	}
	return models, nil
}

func (r *activeClassRepoImpl) FindByNameInYear(ctx context.Context, academicYearID uint, name string, excludeID uint) (*model.ActiveClass, error) {
	var m model.ActiveClass
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

func (r *activeClassRepoImpl) ExistsNameInYear(ctx context.Context, academicYearID uint, name string, excludeID uint) (bool, error) {
	var c int64
	q := r.db.WithContext(ctx).Model(&model.ActiveClass{}).Where("academic_year_id = ? AND name = ?", academicYearID, name)
	if excludeID > 0 {
		q = q.Where("id != ?", excludeID)
	}
	if err := q.Count(&c).Error; err != nil {
		return false, err
	}
	return c > 0, nil
}

func (r *activeClassRepoImpl) ExistsTemplateHomeroom(ctx context.Context, academicYearID, classTemplateID uint, homeroom *string, excludeID uint) (bool, error) {
	var c int64
	q := r.db.WithContext(ctx).Model(&model.ActiveClass{}).
		Where("academic_year_id = ? AND class_template_id = ?", academicYearID, classTemplateID)
	if homeroom != nil && *homeroom != "" {
		q = q.Where("homeroom_number = ?", *homeroom)
	} else {
		q = q.Where("homeroom_number IS NULL")
	}
	if excludeID > 0 {
		q = q.Where("id != ?", excludeID)
	}
	if err := q.Count(&c).Error; err != nil {
		return false, err
	}
	return c > 0, nil
}

func (r *activeClassRepoImpl) Create(ctx context.Context, m *model.ActiveClass) error {
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return fmt.Errorf("nama kelas sudah digunakan di tahun ajaran tersebut")
		}
		return err
	}
	return nil
}

func (r *activeClassRepoImpl) Update(ctx context.Context, m *model.ActiveClass) error {
	if err := r.db.WithContext(ctx).Model(&model.ActiveClass{}).Where("id = ?", m.ID).Updates(map[string]interface{}{
		"class_template_id":     m.ClassTemplateID,
		"name":                  m.Name,
		"homeroom_number":       m.HomeroomNumber,
		"homeroom_teacher_name": m.HomeroomTeacherName,
		"room":                  m.Room,
		"capacity":              m.Capacity,
		"is_active":             m.IsActive,
	}).Error; err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			return fmt.Errorf("nama kelas atau kombinasi rombel sudah digunakan di tahun ajaran tersebut")
		}
		return err
	}
	return nil
}

func (r *activeClassRepoImpl) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.ActiveClass{}, id).Error
}

func (r *activeClassRepoImpl) ToggleStatus(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&model.ActiveClass{}).
		Where("id = ?", id).Update("is_active", gorm.Expr("NOT is_active")).Error
}

// BulkUpsert: replace semua active_class milik tahun ajaran dengan daftar baru.
// Items dengan ID>0 di-update, ID=0 di-create.
func (r *activeClassRepoImpl) BulkUpsert(ctx context.Context, academicYearID uint, models []model.ActiveClass) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Soft-delete yang tidak ada di daftar baru (ID lama tidak disertakan).
		var incomingIDs []uint
		for i := range models {
			if models[i].ID > 0 {
				incomingIDs = append(incomingIDs, models[i].ID)
			}
		}
		delQ := tx.Where("academic_year_id = ?", academicYearID)
		if len(incomingIDs) > 0 {
			delQ = delQ.Where("id NOT IN ?", incomingIDs)
		}
		if err := delQ.Delete(&model.ActiveClass{}).Error; err != nil {
			return err
		}
		for i := range models {
			m := models[i]
			m.AcademicYearID = academicYearID
			if m.ID > 0 {
				if err := tx.Model(&model.ActiveClass{}).Where("id = ?", m.ID).Updates(map[string]interface{}{
					"class_template_id":     m.ClassTemplateID,
					"name":                  m.Name,
					"homeroom_number":       m.HomeroomNumber,
					"homeroom_teacher_name": m.HomeroomTeacherName,
					"room":                  m.Room,
					"capacity":              m.Capacity,
					"is_active":             m.IsActive,
				}).Error; err != nil {
					return err
				}
			} else {
				if err := tx.Create(&m).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

var _ = mapper.ModelToEntity
