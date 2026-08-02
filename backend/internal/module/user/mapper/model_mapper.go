package mapper

import (
	"time"

	roleentity "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/role/entity"
	rolemapper "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/role/mapper"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/user/entity"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/user/model"
)

// ===============================
// User
// ===============================

func ModelToUserEntity(m *model.UserModel) *entity.User {
	if m == nil {
		return nil
	}

	var role *roleentity.Role
	roleName := ""
	roleDisplayName := ""

	if m.Role != nil {
		role = rolemapper.ModelToEntity(m.Role)
		if role != nil {
			roleName = role.Name
			roleDisplayName = role.DisplayName
		}
	}

	var deletedAt *time.Time
	if m.DeletedAt.Valid {
		t := m.DeletedAt.Time
		deletedAt = &t
	}

	return &entity.User{
		ID:              m.ID,
		RoleID:          m.RoleID,
		IsSystem:        role != nil && role.IsSystem,
		Role:            role,
		RoleName:        roleName,
		RoleDisplayName: roleDisplayName,
		Name:            m.Name,
		Email:           m.Email,
		PhotoPath:       m.PhotoPath,
		PasswordHash:    m.PasswordHash,
		Status:          m.Status,
		EmailVerifiedAt: m.EmailVerifiedAt,
		LastLoginAt:     m.LastLoginAt,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
		DeletedAt:       deletedAt,
	}
}

func EntityToUserModel(e *entity.User) *model.UserModel {
	if e == nil {
		return nil
	}

	m := &model.UserModel{
		ID:              e.ID,
		RoleID:          e.RoleID,
		Name:            e.Name,
		Email:           e.Email,
		PhotoPath:       e.PhotoPath,
		PasswordHash:    e.PasswordHash,
		Status:          e.Status,
		EmailVerifiedAt: e.EmailVerifiedAt,
		LastLoginAt:     e.LastLoginAt,
		CreatedAt:       e.CreatedAt,
		UpdatedAt:       e.UpdatedAt,
	}
	if e.DeletedAt != nil {
		m.DeletedAt.Time = *e.DeletedAt
		m.DeletedAt.Valid = true
	}
	return m
}
