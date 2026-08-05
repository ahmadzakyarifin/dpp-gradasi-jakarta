package mapper

import (
	"time"

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

	var deletedAt *time.Time
	if m.DeletedAt.Valid {
		t := m.DeletedAt.Time
		deletedAt = &t
	}

	return &entity.User{
		ID:                 m.ID,
		Role:               m.Role,
		IsSystem:           m.Role == "super_admin",
		Name:               m.Name,
		Email:              m.Email,
		EmailPending:       m.EmailPending,
		PhotoPath:          m.PhotoPath,
		Password:           m.Password,
		Status:             m.Status,
		MustChangePassword: m.MustChangePassword,
		EmailVerifiedAt:    m.EmailVerifiedAt,
		LastLoginAt:        m.LastLoginAt,
		CreatedAt:          m.CreatedAt,
		UpdatedAt:          m.UpdatedAt,
		DeletedAt:          deletedAt,
	}
}

func EntityToUserModel(e *entity.User) *model.UserModel {
	if e == nil {
		return nil
	}

	m := &model.UserModel{
		ID:                 e.ID,
		Role:               e.Role,
		Name:               e.Name,
		Email:              e.Email,
		EmailPending:       e.EmailPending,
		PhotoPath:          e.PhotoPath,
		Password:           e.Password,
		Status:             e.Status,
		MustChangePassword: e.MustChangePassword,
		EmailVerifiedAt:    e.EmailVerifiedAt,
		LastLoginAt:        e.LastLoginAt,
		CreatedAt:          e.CreatedAt,
		UpdatedAt:          e.UpdatedAt,
	}
	if e.DeletedAt != nil {
		m.DeletedAt.Time = *e.DeletedAt
		m.DeletedAt.Valid = true
	}
	return m
}
