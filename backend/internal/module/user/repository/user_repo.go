package repository

import (
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/user/model"
	"gorm.io/gorm"
)

type UserRepo interface {
	FindByID(id uint) (*model.User, error)
	FindByEmail(email string) (*model.User, error)
	FindAllAdmins() ([]model.User, error)
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
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) FindByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) FindAllAdmins() ([]model.User, error) {
	var users []model.User
	// Ambil semua user dengan role admin pengelola (1=super_admin, 2=admin, 3=admin_berita, 4=admin_kegiatan)
	// Gunakan NOT NULL pada role_id untuk fleksibilitas penambahan role di masa depan.
	err := r.db.Where("role_id IN (1, 2, 3, 4)").Find(&users).Error
	return users, err
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
