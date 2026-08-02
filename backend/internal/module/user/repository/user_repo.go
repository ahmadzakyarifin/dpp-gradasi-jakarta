package repository

import (
	"strings"

	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/user/dto"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/user/model"
	"gorm.io/gorm"
)

type UserRepo interface {
	FindByID(id uint) (*model.User, error)
	FindByEmail(email string) (*model.User, error)
	FindAllAdmins(q dto.ListUsersQuery) ([]model.User, int64, error)
	Update(user *model.User) error
	Create(user *model.User) error
	Delete(id uint) error
	Restore(id uint) error
	BulkSoftDelete(ids []uint) error
	BulkRestore(ids []uint) error
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) UserRepo {
	return &userRepo{db: db}
}

func (r *userRepo) FindByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.Preload("Role").First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Preload("Role").Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) FindAllAdmins(q dto.ListUsersQuery) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	base := r.db.Model(&model.User{}).Where("role_id IN (1, 2, 5, 6)")

	switch q.Tab {
	case "pending":
		base = base.Where("status = ?", "pending_activation")
	case "trash":
		base = base.Unscoped().Where("deleted_at IS NOT NULL")
	default: // active
		base = base.Where("status != ?", "pending_activation")
	}

	if s := strings.TrimSpace(q.Search); s != "" {
		like := "%" + s + "%"
		base = base.Where("name LIKE ? OR email LIKE ?", like, like)
	}

	if err := base.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	page, limit := q.Page, q.Limit
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	err := base.
		Preload("Role").
		Order("id ASC").
		Offset((page - 1) * limit).
		Limit(limit).
		Find(&users).Error
	if err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (r *userRepo) Update(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *userRepo) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *userRepo) Delete(id uint) error {
	return r.db.Delete(&model.User{}, id).Error
}

func (r *userRepo) Restore(id uint) error {
	return r.db.Unscoped().Model(&model.User{}).Where("id = ?", id).Update("deleted_at", nil).Error
}

func (r *userRepo) BulkSoftDelete(ids []uint) error {
	return r.db.Where("id IN ?", ids).Delete(&model.User{}).Error
}

func (r *userRepo) BulkRestore(ids []uint) error {
	return r.db.Unscoped().Model(&model.User{}).Where("id IN ?", ids).Update("deleted_at", nil).Error
}
