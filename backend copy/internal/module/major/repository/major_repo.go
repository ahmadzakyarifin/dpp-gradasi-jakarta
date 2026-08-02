package repository

import (
	"context"
	"fmt"

	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/major/entity"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/major/mapper"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/major/model"
	"gorm.io/gorm"
)

type MajorRepo interface {
	Create(ctx context.Context, j *entity.Major) error
	FindAll(ctx context.Context, page, limit int, search, status, sort string) ([]entity.Major, int, error)
	FindByID(ctx context.Context, id uint) (*entity.Major, error)
	FindByIDUnscoped(ctx context.Context, id uint) (*entity.Major, error)
	FindByName(ctx context.Context, name string) (*entity.Major, error)
	Update(ctx context.Context, j *entity.Major) error
	Delete(ctx context.Context, id uint) error
	Restore(ctx context.Context, id uint) error
	BulkDelete(ctx context.Context, ids []uint) error
	BulkRestore(ctx context.Context, ids []uint) error
	CountClasses(ctx context.Context, majorID uint) (int, error)
	CountStudents(ctx context.Context, majorID uint) (int, error)
	CountAcademicYears(ctx context.Context, majorID uint) (int, error)
	// FindNamesByIDs mengambil id+name untuk kebutuhan log bulk di service layer.
	FindNamesByIDs(ctx context.Context, ids []uint) (map[uint]string, error)
}

type majorRepo struct {
	db *gorm.DB
}

func NewMajorRepo(db *gorm.DB) MajorRepo {
	return &majorRepo{db: db}
}

func (r *majorRepo) Create(ctx context.Context, j *entity.Major) error {
	m := mapper.EntityToModel(j)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	j.ID = m.ID
	j.CreatedAt = m.CreatedAt
	j.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *majorRepo) FindAll(ctx context.Context, page, limit int, search, status, sort string) ([]entity.Major, int, error) {
	var list []model.Major
	q := r.db.WithContext(ctx).Model(&model.Major{})

	if status == "trash" {
		q = q.Unscoped().Where("deleted_at IS NOT NULL")
	} else {
		q = q.Where("deleted_at IS NULL")
		switch status {
		case "active":
			q = q.Where("is_active = ?", true)
		case "inactive":
			q = q.Where("is_active = ?", false)
		case "has_class":
			q = q.Where("EXISTS (SELECT 1 FROM class_templates WHERE major_id = majors.id AND deleted_at IS NULL)")
		case "no_class":
			q = q.Where("NOT EXISTS (SELECT 1 FROM class_templates WHERE major_id = majors.id AND deleted_at IS NULL)")
		}
	}

	if search != "" {
		s := "%" + search + "%"
		q = q.Where("name LIKE ? OR code LIKE ?", s, s)
	}

	switch sort {
	case "created_asc":
		q = q.Order("created_at ASC, id ASC")
	case "created_desc":
		q = q.Order("created_at DESC, id DESC")
	case "name_asc":
		q = q.Order("name ASC, id ASC")
	case "name_desc":
		q = q.Order("name DESC, id DESC")
	case "code_asc":
		q = q.Order("code ASC, name ASC")
	case "code_desc":
		q = q.Order("code DESC, name ASC")
	case "oldest":
		q = q.Order("created_at ASC, id ASC")
	case "newest":
		q = q.Order("created_at DESC, id DESC")
	default:
		q = q.Order("created_at DESC, id DESC")
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := q.
		Limit(limit).
		Offset((page - 1) * limit).
		Find(&list).Error; err != nil {
		return nil, int(total), err
	}

	if len(list) > 0 {
		var majorIDs []uint
		majorMap := make(map[uint]int)
		for i, m := range list {
			majorIDs = append(majorIDs, m.ID)
			majorMap[m.ID] = i
		}

		// Batch fetch count of active academic years
		type ayCount struct {
			MajorID uint `gorm:"column:major_id"`
			Count   int  `gorm:"column:count"`
		}
		var ayCounts []ayCount
		r.db.WithContext(ctx).
			Table("active_classes AS ac").
			Select("ct.major_id AS major_id, COUNT(DISTINCT ac.academic_year_id) AS count").
			Joins("JOIN class_templates AS ct ON ct.id = ac.class_template_id").
			Where("ct.major_id IN ? AND ac.deleted_at IS NULL AND ct.deleted_at IS NULL", majorIDs).
			Group("ct.major_id").
			Scan(&ayCounts)
		for _, ac := range ayCounts {
			if idx, ok := majorMap[ac.MajorID]; ok {
				list[idx].AcademicYearCount = ac.Count
			}
		}

		// Batch fetch class templates
		type classCount struct {
			MajorID uint `gorm:"column:major_id"`
			Count   int  `gorm:"column:count"`
		}
		var classCounts []classCount
		r.db.WithContext(ctx).
			Table("class_templates").
			Select("major_id, COUNT(*) AS count").
			Where("major_id IN ? AND deleted_at IS NULL", majorIDs).
			Group("major_id").
			Scan(&classCounts)
		for _, cc := range classCounts {
			if idx, ok := majorMap[cc.MajorID]; ok {
				list[idx].ClassCount = cc.Count
			}
		}

		// Batch fetch students
		type studentCount struct {
			MajorID uint `gorm:"column:major_id"`
			Count   int  `gorm:"column:count"`
		}
		var studentCounts []studentCount
		r.db.WithContext(ctx).
			Table("students").
			Select("current_major_id AS major_id, COUNT(*) AS count").
			Where("current_major_id IN ? AND deleted_at IS NULL AND status = 'active'", majorIDs).
			Group("current_major_id").
			Scan(&studentCounts)
		for _, sc := range studentCounts {
			if idx, ok := majorMap[sc.MajorID]; ok {
				list[idx].StudentCount = sc.Count
			}
		}
	}

	return mapper.ModelListToEntity(list), int(total), nil
}

func (r *majorRepo) FindByID(ctx context.Context, id uint) (*entity.Major, error) {
	var m model.Major
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return mapper.ModelToEntity(&m), nil
}

func (r *majorRepo) FindByIDUnscoped(ctx context.Context, id uint) (*entity.Major, error) {
	var m model.Major
	if err := r.db.WithContext(ctx).Unscoped().Where("id = ?", id).First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return mapper.ModelToEntity(&m), nil
}

func (r *majorRepo) FindByName(ctx context.Context, name string) (*entity.Major, error) {
	var m model.Major
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return mapper.ModelToEntity(&m), nil
}

func (r *majorRepo) Update(ctx context.Context, j *entity.Major) error {
	m := mapper.EntityToModel(j)
	return r.db.WithContext(ctx).
		Model(&model.Major{}).
		Where("id = ?", j.ID).
		Updates(map[string]interface{}{
			"name":       m.Name,
			"code":       m.Code,
			"updated_at": m.UpdatedAt,
		}).Error
}

func (r *majorRepo) Delete(ctx context.Context, id uint) error {
	res := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		Delete(&model.Major{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("jurusan tidak ditemukan atau sudah terhapus")
	}
	return nil
}

func (r *majorRepo) Restore(ctx context.Context, id uint) error {
	res := r.db.WithContext(ctx).
		Unscoped().
		Model(&model.Major{}).
		Where("id = ? AND deleted_at IS NOT NULL", id).
		Update("deleted_at", nil)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("jurusan tidak ditemukan di riwayat penghapusan")
	}
	return nil
}

func (r *majorRepo) BulkDelete(ctx context.Context, ids []uint) error {
	res := r.db.WithContext(ctx).
		Where("id IN ? AND deleted_at IS NULL", ids).
		Delete(&model.Major{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("tidak ada data yang dapat dihapus")
	}
	return nil
}

func (r *majorRepo) BulkRestore(ctx context.Context, ids []uint) error {
	res := r.db.WithContext(ctx).
		Unscoped().
		Model(&model.Major{}).
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

func (r *majorRepo) CountClasses(ctx context.Context, majorID uint) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Table("class_templates").
		Where("major_id = ? AND deleted_at IS NULL", majorID).
		Count(&count).Error
	return int(count), err
}

func (r *majorRepo) CountStudents(ctx context.Context, majorID uint) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Table("students").
		Where("current_major_id = ? AND status = 'active' AND deleted_at IS NULL", majorID).
		Count(&count).Error
	return int(count), err
}

func (r *majorRepo) CountAcademicYears(ctx context.Context, majorID uint) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Table("active_classes AS ac").
		Select("ac.academic_year_id").
		Joins("JOIN class_templates ct ON ct.id = ac.class_template_id").
		Where("ct.major_id = ? AND ac.deleted_at IS NULL AND ct.deleted_at IS NULL", majorID).
		Group("ac.academic_year_id").
		Count(&count).Error
	return int(count), err
}

func (r *majorRepo) FindNamesByIDs(ctx context.Context, ids []uint) (map[uint]string, error) {
	type row struct {
		ID   uint   `gorm:"column:id"`
		Name string `gorm:"column:name"`
	}
	var rows []row
	if err := r.db.WithContext(ctx).
		Unscoped().
		Table("majors").
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
