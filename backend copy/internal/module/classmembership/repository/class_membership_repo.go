package repository

import (
	"context"
	"time"

	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/classmembership/entity"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/classmembership/mapper"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/classmembership/model"
	"gorm.io/gorm"
)

type ClassMembershipRepo interface {
	FindAll(ctx context.Context, page, limit int, activeClassID, studentID uint, status string) ([]model.ClassMembership, int, error)
	FindByID(ctx context.Context, id uint) (*model.ClassMembership, error)
	FindActiveByStudent(ctx context.Context, studentID uint) (*model.ClassMembership, error)
	Create(ctx context.Context, m *model.ClassMembership) error
	UpdateStatus(ctx context.Context, id uint, status string, endDate *time.Time, note *string) error
	UpdateActiveClass(ctx context.Context, id uint, activeClassID uint, semesterID *uint, attendanceNumber *int, endDate *time.Time, note *string) error
}

type classMembershipRepoImpl struct {
	db *gorm.DB
}

func NewClassMembershipRepo(db *gorm.DB) ClassMembershipRepo {
	return &classMembershipRepoImpl{db: db}
}

func (r *classMembershipRepoImpl) FindAll(ctx context.Context, page, limit int, activeClassID, studentID uint, status string) ([]model.ClassMembership, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 50
	}
	q := r.db.WithContext(ctx).Model(&model.ClassMembership{})
	if activeClassID > 0 {
		q = q.Where("active_class_id = ?", activeClassID)
	}
	if studentID > 0 {
		q = q.Where("student_id = ?", studentID)
	}
	if status != "" {
		q = q.Where("status = ?", status)
	} else {
		q = q.Where("deleted_at IS NULL")
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var models []model.ClassMembership
	if err := q.Order("created_at DESC").Offset((page - 1) * limit).Limit(limit).Find(&models).Error; err != nil {
		return nil, 0, err
	}
	return models, int(total), nil
}

func (r *classMembershipRepoImpl) FindByID(ctx context.Context, id uint) (*model.ClassMembership, error) {
	var m model.ClassMembership
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (r *classMembershipRepoImpl) FindActiveByStudent(ctx context.Context, studentID uint) (*model.ClassMembership, error) {
	var m model.ClassMembership
	if err := r.db.WithContext(ctx).Where("student_id = ? AND status = ? AND deleted_at IS NULL", studentID, "active").First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (r *classMembershipRepoImpl) Create(ctx context.Context, m *model.ClassMembership) error {
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	return nil
}

func (r *classMembershipRepoImpl) UpdateStatus(ctx context.Context, id uint, status string, endDate *time.Time, note *string) error {
	updates := map[string]interface{}{"status": status}
	if endDate != nil {
		updates["end_date"] = endDate
	}
	if note != nil {
		updates["note"] = note
	}
	return r.db.WithContext(ctx).Model(&model.ClassMembership{}).Where("id = ?", id).Updates(updates).Error
}

func (r *classMembershipRepoImpl) UpdateActiveClass(ctx context.Context, id uint, activeClassID uint, semesterID *uint, attendanceNumber *int, endDate *time.Time, note *string) error {
	updates := map[string]interface{}{
		"active_class_id":   activeClassID,
		"semester_id":       semesterID,
		"attendance_number": attendanceNumber,
		"end_date":          endDate,
		"note":              note,
	}
	return r.db.WithContext(ctx).Model(&model.ClassMembership{}).Where("id = ?", id).Updates(updates).Error
}

var _ = mapper.ModelToEntity
var _ = entity.ClassMembership{}
