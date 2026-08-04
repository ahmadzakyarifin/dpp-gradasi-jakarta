package mapper

import (
	"time"

	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/user/entity"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/user/model"
)

// ModelToRoleEntity memetakan RoleModel -> entity.Role.
// Disalin dari modul role (dipakai saat join user->role).
func ModelToRoleEntity(m *model.RoleModel) *entity.Role {
	if m == nil {
		return nil
	}

	var deletedAt *time.Time
	if m.DeletedAt.Valid {
		deletedAt = &m.DeletedAt.Time
	}

	return &entity.Role{
		ID:          m.ID,
		Name:        m.Name,
		DisplayName: m.DisplayName,
		IsActive:    m.IsActive,
		IsSystem:    m.IsSystem,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		DeletedAt:   deletedAt,
	}
}
