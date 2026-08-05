package repository

import (
	"context"
	"fmt"

	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/user/entity"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/user/mapper"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/user/model"
	"gorm.io/gorm"
)

type UserRepo interface {
	Create(ctx context.Context, db *gorm.DB, user *entity.User) error
	CreateTx(ctx context.Context, user *entity.User) error
	FindAll(ctx context.Context, role string) ([]entity.User, error)
	FindByID(ctx context.Context, id uint) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	Update(ctx context.Context, db *gorm.DB, user *entity.User) error
	UpdateTx(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id uint) error
	UpdatePassword(ctx context.Context, id uint, passwordHash string) error
	Activate(ctx context.Context, id uint) error
	ToggleStatus(ctx context.Context, id uint) error
	FindPaginated(ctx context.Context, page, limit int, search, role, status, sort string, trashed bool) ([]entity.User, int, error)
	CountActive(ctx context.Context) (int, error)
	BulkDelete(ctx context.Context, ids []uint) error
	Restore(ctx context.Context, id uint) error
	BulkRestore(ctx context.Context, ids []uint) error
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) UserRepo {
	return &userRepo{db: db}
}

func (r *userRepo) Create(ctx context.Context, db *gorm.DB, u *entity.User) error {
	if db == nil {
		db = r.db
	}
	m := mapper.EntityToUserModel(u)
	if err := db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	u.ID = m.ID
	u.CreatedAt = m.CreatedAt
	u.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *userRepo) FindAll(ctx context.Context, role string) ([]entity.User, error) {
	q := r.db.WithContext(ctx)
	if role != "" {
		q = q.Where("role = ?", role)
	}
	var models []model.UserModel
	if err := q.Find(&models).Error; err != nil {
		return nil, err
	}
	users := make([]entity.User, len(models))
	for i := range models {
		users[i] = *mapper.ModelToUserEntity(&models[i])
	}
	return users, nil
}

func (r *userRepo) FindByID(ctx context.Context, id uint) (*entity.User, error) {
	var m model.UserModel
	err := r.db.WithContext(ctx).Unscoped().Where("users.id = ?", id).First(&m).Error
	if err != nil {
		return nil, err
	}
	return mapper.ModelToUserEntity(&m), nil
}

func (r *userRepo) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var m model.UserModel
	err := r.db.WithContext(ctx).Unscoped().Where("email = ?", email).First(&m).Error
	if err != nil {
		return nil, err
	}
	return mapper.ModelToUserEntity(&m), nil
}

func (r *userRepo) Update(ctx context.Context, db *gorm.DB, u *entity.User) error {
	if db == nil {
		db = r.db
	}
	m := mapper.EntityToUserModel(u)
	updates := map[string]any{
		"role":          m.Role,
		"name":          m.Name,
		"email":         m.Email,
		"email_pending": m.EmailPending,
		"photo_path":    m.PhotoPath,
		"status":        m.Status,
	}
	return db.WithContext(ctx).Model(&model.UserModel{}).
		Where("id = ?", m.ID).
		Updates(updates).Error
}

func (r *userRepo) CreateTx(ctx context.Context, u *entity.User) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return r.Create(ctx, tx, u)
	})
}

func (r *userRepo) UpdateTx(ctx context.Context, u *entity.User) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return r.Update(ctx, tx, u)
	})
}

func (r *userRepo) Delete(ctx context.Context, id uint) error {
	res := r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.UserModel{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("pengguna sudah terhapus atau tidak ditemukan")
	}
	return nil
}

func (r *userRepo) BulkDelete(ctx context.Context, ids []uint) error {
	res := r.db.WithContext(ctx).Where("id IN ?", ids).Delete(&model.UserModel{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("tidak ada data yang perlu dihapus")
	}
	return nil
}

func (r *userRepo) Restore(ctx context.Context, id uint) error {
	res := r.db.WithContext(ctx).Unscoped().Model(&model.UserModel{}).
		Where("id = ? AND deleted_at IS NOT NULL", id).
		Update("deleted_at", nil)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("data tidak ditemukan di tempat sampah atau sudah aktif")
	}
	return nil
}

func (r *userRepo) BulkRestore(ctx context.Context, ids []uint) error {
	res := r.db.WithContext(ctx).Unscoped().Model(&model.UserModel{}).
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

func (r *userRepo) UpdatePassword(ctx context.Context, id uint, hash string) error {
	res := r.db.WithContext(ctx).Model(&model.UserModel{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"password": hash,
			"status":   "active",
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("pengguna tidak ditemukan")
	}
	return nil
}

func (r *userRepo) Activate(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&model.UserModel{}).
		Where("id = ?", id).
		Update("status", "active").Error
}

func (r *userRepo) ToggleStatus(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&model.UserModel{}).
		Where("id = ?", id).
		Update("status", gorm.Expr("CASE WHEN status = 'active' THEN 'inactive' ELSE 'active' END")).Error
}

func (r *userRepo) FindPaginated(ctx context.Context, page, limit int, search, role, status, sort string, trashed bool) ([]entity.User, int, error) {
	q := r.db.WithContext(ctx).Model(&model.UserModel{})

	if trashed {
		q = q.Unscoped().Where("users.deleted_at IS NOT NULL")
	} else {
		switch status {
		case "active":
			q = q.Where("users.status = ? OR (users.status = ? AND users.password != ?)", "active", "inactive", "")
		case "inactive":
			q = q.Where("users.status = ? AND (users.password = ? OR users.password IS NULL)", "inactive", "")
		}
	}

	if search != "" {
		s := "%" + search + "%"
		q = q.Where("users.name LIKE ? OR users.email LIKE ?", s, s)
	}

	if role != "" {
		q = q.Where("users.role = ?", role)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	switch sort {
	case "name_asc":
		q = q.Order("users.name ASC")
	case "name_desc":
		q = q.Order("users.name DESC")
	case "created_asc":
		q = q.Order("users.created_at ASC")
	case "created_desc":
		q = q.Order("users.created_at DESC")
	default:
		q = q.Order("users.created_at DESC")
	}

	var models []model.UserModel
	err := q.
		Limit(limit).
		Offset((page - 1) * limit).
		Find(&models).Error
	if err != nil {
		return nil, 0, err
	}

	users := make([]entity.User, len(models))
	for i := range models {
		users[i] = *mapper.ModelToUserEntity(&models[i])
	}

	return users, int(total), nil
}

func (r *userRepo) CountActive(ctx context.Context) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.UserModel{}).Where("status = ?", "active").Count(&count).Error
	return int(count), err
}
