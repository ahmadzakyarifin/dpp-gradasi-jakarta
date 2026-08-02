package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/student/entity"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/student/model"
	"gorm.io/gorm"
)

// StudentRepo interface defines repository operations for students using uint ID.
type StudentRepo interface {
	Create(ctx context.Context, tx *gorm.DB, s *model.Student) error
	FindAllPaginated(ctx context.Context, page, limit int, search, filter, status string, entryYear int, classID, majorID uint, sort string) ([]model.Student, int64, error)
	FindByID(ctx context.Context, id uint) (*model.Student, error)
	FindByIDUnscoped(ctx context.Context, id uint) (*model.Student, error)
	FindByNIK(ctx context.Context, nik string) (*model.Student, error)
	FindByNISN(ctx context.Context, nisn string) (*model.Student, error)
	FindByNIS(ctx context.Context, nis string) (*model.Student, error)
	FindByEmail(ctx context.Context, email string) (*model.Student, error)
	FindByPhone(ctx context.Context, phone string) (*model.Student, error)
	Update(ctx context.Context, tx *gorm.DB, s *model.Student) error
	Delete(ctx context.Context, tx *gorm.DB, id uint) error
	ToggleStatus(ctx context.Context, tx *gorm.DB, id uint) error
	Restore(ctx context.Context, tx *gorm.DB, id uint) error
	BulkRestore(ctx context.Context, tx *gorm.DB, ids []uint) error
	GetDistinctEntryYears(ctx context.Context) ([]int, error)
	GetClassHistory(ctx context.Context, studentID uint) ([]entity.ClassHistory, error)
	AddClassHistory(ctx context.Context, tx *gorm.DB, studentID, classID uint) error
	UpdateActiveHistory(ctx context.Context, tx *gorm.DB, studentID, classID uint) error
	BulkGraduate(ctx context.Context, tx *gorm.DB, ids []uint) error
	BulkPromote(ctx context.Context, tx *gorm.DB, targetClassID uint, ids []uint) error
	DeactivateHistory(ctx context.Context, tx *gorm.DB, ids []uint) error
	FindActiveByIDs(ctx context.Context, ids []uint) ([]model.Student, error)
	CountActive(ctx context.Context) (int64, error)
	FindUnpaginated(ctx context.Context, search, filter, status string, entryYear int, classID, majorID uint) ([]model.Student, error)
	ExistsByEmail(ctx context.Context, email string, excludeID uint) error
	ExistsByPhone(ctx context.Context, phone string, excludeID uint) error
	ExistsByNIK(ctx context.Context, nik string, excludeID uint) error
	ExistsByNIS(ctx context.Context, nis string, excludeID uint) error
}

type studentRepo struct {
	db *gorm.DB
}

func NewStudentRepo(db *gorm.DB) StudentRepo {
	return &studentRepo{db: db}
}

func (r *studentRepo) Create(ctx context.Context, tx *gorm.DB, s *model.Student) error {
	return tx.WithContext(ctx).Create(s).Error
}

func (r *studentRepo) findAllQuery(ctx context.Context, search, filter, status string, entryYear int, classID, majorID uint) *gorm.DB {
	q := r.db.WithContext(ctx).Model(&model.Student{}).
		Select("students.*, COALESCE(u.name,'') as parent_name, c.name as class_name, m.name as major_name").
		Joins("LEFT JOIN users u ON students.user_id = u.id").
		Joins("LEFT JOIN active_classes c ON students.current_active_class_id = c.id").
		Joins("LEFT JOIN class_templates ct ON c.class_template_id = ct.id").
		Joins("LEFT JOIN majors m ON COALESCE(ct.major_id, students.current_major_id) = m.id")

	if search != "" {
		s := "%" + search + "%"
		q = q.Where("(students.name LIKE ? OR students.nisn LIKE ? OR students.nis LIKE ? OR students.email LIKE ?)", s, s, s, s)
	}
	if filter == "no_parent" {
		q = q.Where("students.user_id IS NULL")
	}
	if status == "trash" {
		q = q.Unscoped().Where("students.deleted_at IS NOT NULL")
	} else if status != "" {
		q = q.Where("students.status = ?", status)
	}
	if entryYear > 0 {
		q = q.Where("students.entry_year = ?", entryYear)
	}
	if classID > 0 {
		q = q.Where("students.current_active_class_id = ?", classID)
	}
	if majorID > 0 {
		q = q.Where("COALESCE(ct.major_id, students.current_major_id) = ?", majorID)
	}
	return q
}

func (r *studentRepo) FindAllPaginated(ctx context.Context, page, limit int, search, filter, status string, entryYear int, classID, majorID uint, sort string) ([]model.Student, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	var list []model.Student
	var total int64

	q := r.findAllQuery(ctx, search, filter, status, entryYear, classID, majorID)
	q.Count(&total)

	switch sort {
	case "created_desc":
		q = q.Order("students.created_at DESC, students.id DESC")
	case "created_asc":
		q = q.Order("students.created_at ASC, students.id ASC")
	case "name_desc":
		q = q.Order("students.name DESC, students.id DESC")
	case "entry_year_desc":
		q = q.Order("students.entry_year DESC, students.name ASC")
	case "entry_year_asc":
		q = q.Order("students.entry_year ASC, students.name ASC")
	default:
		q = q.Order("students.name ASC, students.id ASC")
	}

	err := q.Offset((page - 1) * limit).Limit(limit).Find(&list).Error
	return list, total, err
}

func (r *studentRepo) FindByID(ctx context.Context, id uint) (*model.Student, error) {
	var s model.Student
	err := r.db.WithContext(ctx).
		Select("students.*, COALESCE(u.name,'') as parent_name, c.name as class_name, m.name as major_name").
		Joins("LEFT JOIN users u ON students.user_id = u.id").
		Joins("LEFT JOIN active_classes c ON students.current_active_class_id = c.id").
		Joins("LEFT JOIN class_templates ct ON c.class_template_id = ct.id").
		Joins("LEFT JOIN majors m ON COALESCE(ct.major_id, students.current_major_id) = m.id").
		First(&s, id).Error
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *studentRepo) FindByIDUnscoped(ctx context.Context, id uint) (*model.Student, error) {
	var s model.Student
	err := r.db.WithContext(ctx).Unscoped().
		Select("students.*, COALESCE(u.name,'') as parent_name, c.name as class_name, m.name as major_name").
		Joins("LEFT JOIN users u ON students.user_id = u.id").
		Joins("LEFT JOIN active_classes c ON students.current_active_class_id = c.id").
		Joins("LEFT JOIN class_templates ct ON c.class_template_id = ct.id").
		Joins("LEFT JOIN majors m ON COALESCE(ct.major_id, students.current_major_id) = m.id").
		First(&s, id).Error
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *studentRepo) FindByNIK(ctx context.Context, nik string) (*model.Student, error) {
	var s model.Student
	err := r.db.WithContext(ctx).Where("nik = ?", nik).First(&s).Error
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *studentRepo) FindByNISN(ctx context.Context, nisn string) (*model.Student, error) {
	var s model.Student
	err := r.db.WithContext(ctx).Where("nisn = ?", nisn).First(&s).Error
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *studentRepo) FindByNIS(ctx context.Context, nis string) (*model.Student, error) {
	var s model.Student
	err := r.db.WithContext(ctx).Where("nis = ?", nis).First(&s).Error
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *studentRepo) FindByEmail(ctx context.Context, email string) (*model.Student, error) {
	var s model.Student
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&s).Error
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *studentRepo) FindByPhone(ctx context.Context, phone string) (*model.Student, error) {
	var s model.Student
	err := r.db.WithContext(ctx).Where("phone = ?", phone).First(&s).Error
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *studentRepo) Update(ctx context.Context, tx *gorm.DB, s *model.Student) error {
	return tx.WithContext(ctx).Model(&model.Student{}).Where("id = ?", s.ID).Updates(s).Error
}

func (r *studentRepo) Delete(ctx context.Context, tx *gorm.DB, id uint) error {
	res := tx.WithContext(ctx).Where("id = ?", id).Delete(&model.Student{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("siswa sudah terhapus atau tidak ditemukan")
	}
	return nil
}

func (r *studentRepo) ToggleStatus(ctx context.Context, tx *gorm.DB, id uint) error {
	var s model.Student
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&s).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("siswa tidak ditemukan")
		}
		return err
	}

	nextStatus := "active"
	if s.Status == "active" {
		nextStatus = "inactive"
	} else if s.Status == "inactive" {
		nextStatus = "active"
	} else {
		return fmt.Errorf("status siswa tidak dapat diubah: %s", s.Status)
	}
	return tx.WithContext(ctx).Model(&model.Student{}).Where("id = ?", id).UpdateColumn("status", nextStatus).Error
}

func (r *studentRepo) Restore(ctx context.Context, tx *gorm.DB, id uint) error {
	return tx.WithContext(ctx).Unscoped().Where("id = ?", id).UpdateColumns(map[string]interface{}{
		"deleted_at": nil,
		"status":     "active",
	}).Error
}

func (r *studentRepo) BulkRestore(ctx context.Context, tx *gorm.DB, ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	return tx.WithContext(ctx).Where("id IN ?", ids).Unscoped().UpdateColumn("deleted_at", nil).Error
}

func (r *studentRepo) GetDistinctEntryYears(ctx context.Context) ([]int, error) {
	var years []int
	err := r.db.WithContext(ctx).Model(&model.Student{}).
		Select("DISTINCT entry_year").Order("entry_year DESC").Find(&years).Error
	return years, err
}

func (r *studentRepo) GetClassHistory(ctx context.Context, studentID uint) ([]entity.ClassHistory, error) {
	var list []entity.ClassHistory
	err := r.db.WithContext(ctx).Raw(`
		SELECT sc.id, sc.student_id, sc.class_id, sc.is_active, sc.created_at,
		       c.name as class_name, c.grade,
		       CASE WHEN MONTH(sc.created_at) >= 7
		            THEN CONCAT(YEAR(sc.created_at), '/', YEAR(sc.created_at)+1)
		            ELSE CONCAT(YEAR(sc.created_at)-1, '/', YEAR(sc.created_at))
		       END as academic_year
		FROM student_classes sc
		JOIN active_classes c ON sc.class_id = c.id
		WHERE sc.student_id = ?
		ORDER BY sc.created_at DESC
	`, studentID).Find(&list).Error
	return list, err
}

func (r *studentRepo) AddClassHistory(ctx context.Context, tx *gorm.DB, studentID, classID uint) error {
	tx.WithContext(ctx).Exec("UPDATE student_classes SET is_active = 0 WHERE student_id = ? AND is_active = 1", studentID)
	return tx.WithContext(ctx).Exec("INSERT INTO student_classes (student_id, class_id, is_active, created_at) VALUES (?, ?, 1, NOW())", studentID, classID).Error
}

func (r *studentRepo) UpdateActiveHistory(ctx context.Context, tx *gorm.DB, studentID, classID uint) error {
	return tx.WithContext(ctx).Exec("UPDATE student_classes SET class_id = ? WHERE student_id = ? AND is_active = 1", classID, studentID).Error
}

func (r *studentRepo) BulkGraduate(ctx context.Context, tx *gorm.DB, ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	return tx.WithContext(ctx).Model(&model.Student{}).
		Where("id IN ? AND status = 'active'", ids).
		Update("status", "graduated").Error
}

func (r *studentRepo) BulkPromote(ctx context.Context, tx *gorm.DB, targetClassID uint, ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	return tx.WithContext(ctx).Model(&model.Student{}).
		Where("id IN ?", ids).
		Update("current_active_class_id", targetClassID).Error
}

func (r *studentRepo) DeactivateHistory(ctx context.Context, tx *gorm.DB, ids []uint) error {
	return tx.WithContext(ctx).Exec("UPDATE student_classes SET is_active = 0 WHERE student_id IN ? AND is_active = 1", ids).Error
}

func (r *studentRepo) FindActiveByIDs(ctx context.Context, ids []uint) ([]model.Student, error) {
	var list []model.Student
	err := r.db.WithContext(ctx).Where("status = 'active'").Where("id IN ?", ids).Find(&list).Error
	return list, err
}

func (r *studentRepo) CountActive(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Student{}).Where("status = 'active'").Count(&count).Error
	return count, err
}

func (r *studentRepo) FindUnpaginated(ctx context.Context, search, filter, status string, entryYear int, classID, majorID uint) ([]model.Student, error) {
	var list []model.Student
	q := r.findAllQuery(ctx, search, filter, status, entryYear, classID, majorID)
	return list, q.Find(&list).Error
}

func (r *studentRepo) ExistsByEmail(ctx context.Context, email string, excludeID uint) error {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Student{}).Where("email = ? AND id != ?", email, excludeID).Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("email already exists")
	}
	return nil
}

func (r *studentRepo) ExistsByPhone(ctx context.Context, phone string, excludeID uint) error {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Student{}).Where("phone = ? AND id != ?", phone, excludeID).Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("phone already exists")
	}
	return nil
}

func (r *studentRepo) ExistsByNIK(ctx context.Context, nik string, excludeID uint) error {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Student{}).Where("nik = ? AND id != ?", nik, excludeID).Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("nik already exists")
	}
	return nil
}

func (r *studentRepo) ExistsByNIS(ctx context.Context, nis string, excludeID uint) error {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Student{}).Where("nis = ? AND id != ?", nis, excludeID).Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("nis already exists")
	}
	return nil
}

// Ensure time import is used
var _ = time.Now
