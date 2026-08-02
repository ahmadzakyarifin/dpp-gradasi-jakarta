package repository

import (
	"context"
	"fmt"

	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/helper"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/role/entity"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/role/mapper"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/role/model"
	usermodel "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/user/model"
	"gorm.io/gorm"
)

type RoleRepo interface {
	GetAllRoles(ctx context.Context, page, limit int, search, status, sort string) ([]entity.Role, int, error)
	GetRoleByID(ctx context.Context, id uint) (*entity.Role, error)
	GetRoleByIDUnscoped(ctx context.Context, id uint) (*entity.Role, error)
	GetRoleByNameWithDeleted(ctx context.Context, name string) (*model.RoleModel, error)
	GetRoleByDisplayNameWithDeleted(ctx context.Context, displayName string) (*model.RoleModel, error)
	GetPermissions(ctx context.Context) ([]entity.Permission, error)
	CreateRole(ctx context.Context, role *entity.Role, permissionIDs []uint) error
	UpdateRole(ctx context.Context, role *entity.Role, permissionIDs []uint) error
	DeleteRole(ctx context.Context, id uint) error
	Restore(ctx context.Context, id uint) error
	ToggleStatus(ctx context.Context, id uint) error
	BulkDelete(ctx context.Context, ids []uint) error
	BulkRestore(ctx context.Context, ids []uint) error
	HasUsers(ctx context.Context, roleID uint) (bool, error)
	CountUsers(ctx context.Context, roleID uint) (int, error)
	GetPermissionsByRoleID(ctx context.Context, roleID uint) ([]string, error)
	ExistsByName(ctx context.Context, name string, excludeID uint) (bool, error)
	ExistsByDisplayName(ctx context.Context, displayName string, excludeID uint) (bool, error)
}

type roleRepo struct {
	db *gorm.DB
}

func NewRoleRepo(db *gorm.DB) RoleRepo {
	return &roleRepo{db: db}
}

// Cache Helpers

func fetchRolePermissionsFromDB(ctx context.Context, db *gorm.DB, roleID uint) ([]string, error) {
	var names []string
	err := db.WithContext(ctx).
		Table("role_permissions AS rp").
		Select("p.name").
		Joins("JOIN permissions AS p ON rp.permission_id = p.id").
		Where("rp.role_id = ?", roleID).
		Scan(&names).Error
	return names, err
}

func (r *roleRepo) GetAllRoles(ctx context.Context, page, limit int, search, status, sort string) ([]entity.Role, int, error) {
	q := r.db.WithContext(ctx).Model(&model.RoleModel{})

	if status == "trash" {
		q = q.Unscoped().Where("roles.deleted_at IS NOT NULL")
	} else {
		switch status {
		case "active":
			q = q.Where("roles.is_active = ?", 1)
		case "inactive":
			q = q.Where("roles.is_active = ?", 0)
		}
	}

	if search != "" {
		s := "%" + search + "%"
		q = q.Where("roles.name LIKE ? OR roles.display_name LIKE ?", s, s)
	}

	switch sort {
	case "name_asc":
		q = q.Order("roles.display_name ASC").Order("roles.id ASC")
	case "name_desc":
		q = q.Order("roles.display_name DESC").Order("roles.id DESC")
	case "oldest":
		q = q.Order("roles.created_at ASC").Order("roles.id ASC")
	case "newest":
		q = q.Order("roles.created_at DESC").Order("roles.id DESC")
	default:
		q = q.Order("roles.created_at DESC").Order("roles.id DESC")
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if limit > 0 {
		q = q.Limit(limit).Offset((page - 1) * limit)
	}

	var models []model.RoleModel
	if err := q.Preload("Permissions").Find(&models).Error; err != nil {
		return nil, int(total), err
	}

	roles := make([]entity.Role, 0, len(models))
	for i := range models {
		domain := mapper.ModelToEntity(&models[i])
		if domain != nil {
			roles = append(roles, *domain)
		}
	}
	return roles, int(total), nil
}

func (r *roleRepo) GetRoleByID(ctx context.Context, id uint) (*entity.Role, error) {
	var m model.RoleModel
	err := r.db.WithContext(ctx).Preload("Permissions").Where("roles.id = ?", id).First(&m).Error
	if err != nil {
		return nil, err
	}
	domain := mapper.ModelToEntity(&m)
	return domain, nil
}

func (r *roleRepo) GetRoleByIDUnscoped(ctx context.Context, id uint) (*entity.Role, error) {
	var m model.RoleModel
	err := r.db.WithContext(ctx).Unscoped().Preload("Permissions").Where("roles.id = ?", id).First(&m).Error
	if err != nil {
		return nil, err
	}
	domain := mapper.ModelToEntity(&m)
	return domain, nil
}

func (r *roleRepo) GetRoleByNameWithDeleted(ctx context.Context, name string) (*model.RoleModel, error) {
	var m model.RoleModel
	err := r.db.WithContext(ctx).Unscoped().Where("name = ?", name).First(&m).Error
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *roleRepo) GetRoleByDisplayNameWithDeleted(ctx context.Context, displayName string) (*model.RoleModel, error) {
	var m model.RoleModel
	err := r.db.WithContext(ctx).Unscoped().Where("display_name = ?", displayName).First(&m).Error
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *roleRepo) GetPermissions(ctx context.Context) ([]entity.Permission, error) {
	var models []model.PermissionModel
	err := r.db.WithContext(ctx).Order("module ASC").Order("id ASC").Find(&models).Error
	if err != nil {
		return nil, err
	}

	perms := make([]entity.Permission, 0, len(models))
	for _, p := range models {
		perms = append(perms, entity.Permission{
			ID:          p.ID,
			Name:        p.Name,
			DisplayName: p.DisplayName,
			Module:      p.Module,
			Description: p.Description,
			CreatedAt:   p.CreatedAt,
			UpdatedAt:   p.UpdatedAt,
		})
	}
	return perms, nil
}

func (r *roleRepo) GetPermissionsByRoleID(ctx context.Context, roleID uint) ([]string, error) {
	return fetchRolePermissionsFromDB(ctx, r.db, roleID)
}

func (r *roleRepo) CreateRole(ctx context.Context, role *entity.Role, permissionIDs []uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		m := mapper.EntityToModel(role)
		if err := tx.Create(m).Error; err != nil {
			return err
		}
		role.ID = m.ID

		if len(permissionIDs) > 0 {
			rows := make([]model.RolePermissionModel, 0, len(permissionIDs))
			for _, pid := range permissionIDs {
				rows = append(rows, model.RolePermissionModel{RoleID: role.ID, PermissionID: pid})
			}
			if err := tx.Create(&rows).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *roleRepo) UpdateRole(ctx context.Context, role *entity.Role, permissionIDs []uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		m := mapper.EntityToModel(role)
		err := tx.Model(&model.RoleModel{}).Where("id = ?", m.ID).Updates(map[string]any{
			"name":         m.Name,
			"display_name": m.DisplayName,
			"is_active":    m.IsActive,
		}).Error
		if err != nil {
			return err
		}

		if permissionIDs != nil {
			if err := tx.Where("role_id = ?", role.ID).Delete(&model.RolePermissionModel{}).Error; err != nil {
				return err
			}
			if len(permissionIDs) > 0 {
				rows := make([]model.RolePermissionModel, 0, len(permissionIDs))
				for _, pid := range permissionIDs {
					rows = append(rows, model.RolePermissionModel{RoleID: role.ID, PermissionID: pid})
				}
				if err := tx.Create(&rows).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func (r *roleRepo) DeleteRole(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.RoleModel{}).Error
}

func (r *roleRepo) SaveRolePermissions(ctx context.Context, roleID uint, permissionIDs []uint) error {
	if err := r.db.WithContext(ctx).Where("role_id = ?", roleID).Delete(&model.RolePermissionModel{}).Error; err != nil {
		return err
	}

	if len(permissionIDs) == 0 {
		return nil
	}

	rows := make([]model.RolePermissionModel, 0, len(permissionIDs))
	for _, pid := range permissionIDs {
		rows = append(rows, model.RolePermissionModel{RoleID: roleID, PermissionID: pid})
	}

	if err := r.db.WithContext(ctx).Create(&rows).Error; err != nil {
		return err
	}
	return nil
}

func (r *roleRepo) HasUsers(ctx context.Context, roleID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&usermodel.UserModel{}).
		Where("role_id = ? AND deleted_at IS NULL", roleID).
		Count(&count).Error
	return count > 0, err
}

func (r *roleRepo) CountUsers(ctx context.Context, roleID uint) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).Table("users").
		Where("role_id = ? AND deleted_at IS NULL", roleID).
		Count(&count).Error
	return int(count), err
}

func (r *roleRepo) Restore(ctx context.Context, id uint) error {
	res := r.db.WithContext(ctx).Unscoped().Model(&model.RoleModel{}).
		Where("id = ? AND deleted_at IS NOT NULL", id).
		Update("deleted_at", nil)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return helper.NewNotFoundError("role tidak ditemukan di riwayat penghapusan")
	}
	return nil
}

func (r *roleRepo) ToggleStatus(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&model.RoleModel{}).Where("id = ?", id).Update("is_active", gorm.Expr("NOT is_active")).Error
}

func (r *roleRepo) BulkDelete(ctx context.Context, ids []uint) error {
	res := r.db.WithContext(ctx).
		Where("id IN ? AND deleted_at IS NULL AND is_system = 0", ids).
		Delete(&model.RoleModel{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("tidak ada data yang dapat dihapus")
	}
	return nil
}

func (r *roleRepo) BulkRestore(ctx context.Context, ids []uint) error {
	res := r.db.WithContext(ctx).Unscoped().Model(&model.RoleModel{}).
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

func (r *roleRepo) ExistsByName(ctx context.Context, name string, excludeID uint) (bool, error) {
	q := r.db.WithContext(ctx).
		Unscoped().
		Model(&model.RoleModel{}).
		Where("name = ?", name)
	if excludeID > 0 {
		q = q.Where("id != ?", excludeID)
	}
	var count int64
	if err := q.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *roleRepo) ExistsByDisplayName(ctx context.Context, displayName string, excludeID uint) (bool, error) {
	q := r.db.WithContext(ctx).
		Unscoped().
		Model(&model.RoleModel{}).
		Where("display_name = ?", displayName)
	if excludeID > 0 {
		q = q.Where("id != ?", excludeID)
	}
	var count int64
	if err := q.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
