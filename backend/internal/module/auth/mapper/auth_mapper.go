package mapper

import (
	authdto "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/auth/dto"
	authentity "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/auth/entity"
	authmodel "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/auth/model"
	rolemodel "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/role/model"
	usermodel "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/user/model"
)

// ---- Model -> Entity ----

func UserModelToEntity(m *usermodel.User) *authentity.UserEntity {
	if m == nil {
		return nil
	}
	return &authentity.UserEntity{
		ID:                 m.ID,
		RoleID:             m.RoleID,
		Name:               m.Name,
		Email:              m.Email,
		Password:           m.Password,
		PhotoPath:          m.PhotoPath,
		Status:             m.Status,
		MustChangePassword: m.MustChangePassword,
		EmailVerifiedAt:    m.EmailVerifiedAt,
		LastLoginAt:        m.LastLoginAt,
		CreatedAt:          m.CreatedAt,
		Role:               RoleModelToEntity(&m.Role),
	}
}

func RoleModelToEntity(m *rolemodel.Role) authentity.RoleEntity {
	if m == nil {
		return authentity.RoleEntity{}
	}
	return authentity.RoleEntity{
		ID:          m.ID,
		Name:        m.Name,
		DisplayName: m.DisplayName,
	}
}

func RefreshTokenModelToEntity(m *authmodel.RefreshToken) *authentity.RefreshTokenEntity {
	if m == nil {
		return nil
	}
	return &authentity.RefreshTokenEntity{
		ID:         m.ID,
		UserID:     m.UserID,
		TokenHash:  m.TokenHash,
		DeviceInfo: m.DeviceInfo,
		IPAddress:  m.IPAddress,
		ExpiresAt:  m.ExpiresAt,
		CreatedAt:  m.CreatedAt,
	}
}

func PasswordResetTokenModelToEntity(m *authmodel.PasswordResetToken) *authentity.PasswordResetTokenEntity {
	if m == nil {
		return nil
	}
	return &authentity.PasswordResetTokenEntity{
		ID:        m.ID,
		UserID:    m.UserID,
		TokenHash: m.TokenHash,
		ExpiresAt: m.ExpiresAt,
		UsedAt:    m.UsedAt,
		CreatedAt: m.CreatedAt,
	}
}

func ActivationTokenModelToEntity(m *authmodel.ActivationToken) *authentity.ActivationTokenEntity {
	if m == nil {
		return nil
	}
	return &authentity.ActivationTokenEntity{
		ID:        m.ID,
		UserID:    m.UserID,
		TokenHash: m.TokenHash,
		Channel:   m.Channel,
		ExpiresAt: m.ExpiresAt,
		UsedAt:    m.UsedAt,
		CreatedAt: m.CreatedAt,
	}
}

// ---- Entity -> Model ----

func UserEntityToModel(e *authentity.UserEntity) *usermodel.User {
	if e == nil {
		return nil
	}
	return &usermodel.User{
		ID:                 e.ID,
		RoleID:             e.RoleID,
		Name:               e.Name,
		Email:              e.Email,
		Password:           e.Password,
		PhotoPath:          e.PhotoPath,
		Status:             e.Status,
		MustChangePassword: e.MustChangePassword,
		EmailVerifiedAt:    e.EmailVerifiedAt,
		LastLoginAt:        e.LastLoginAt,
		CreatedAt:          e.CreatedAt,
		Role: rolemodel.Role{
			ID:          e.Role.ID,
			Name:        e.Role.Name,
			DisplayName: e.Role.DisplayName,
		},
	}
}

func RefreshTokenEntityToModel(e *authentity.RefreshTokenEntity) *authmodel.RefreshToken {
	if e == nil {
		return nil
	}
	return &authmodel.RefreshToken{
		ID:         e.ID,
		UserID:     e.UserID,
		TokenHash:  e.TokenHash,
		DeviceInfo: e.DeviceInfo,
		IPAddress:  e.IPAddress,
		ExpiresAt:  e.ExpiresAt,
		CreatedAt:  e.CreatedAt,
	}
}

func PasswordResetTokenEntityToModel(e *authentity.PasswordResetTokenEntity) *authmodel.PasswordResetToken {
	if e == nil {
		return nil
	}
	return &authmodel.PasswordResetToken{
		ID:        e.ID,
		UserID:    e.UserID,
		TokenHash: e.TokenHash,
		ExpiresAt: e.ExpiresAt,
		UsedAt:    e.UsedAt,
		CreatedAt: e.CreatedAt,
	}
}

func ActivationTokenEntityToModel(e *authentity.ActivationTokenEntity) *authmodel.ActivationToken {
	if e == nil {
		return nil
	}
	return &authmodel.ActivationToken{
		ID:        e.ID,
		UserID:    e.UserID,
		TokenHash: e.TokenHash,
		Channel:   e.Channel,
		ExpiresAt: e.ExpiresAt,
		UsedAt:    e.UsedAt,
		CreatedAt: e.CreatedAt,
	}
}

// ---- Entity -> DTO (response) ----

func UserEntityToAuthResponse(e *authentity.UserEntity) authdto.AuthUserResponse {
	if e == nil {
		return authdto.AuthUserResponse{}
	}
	return authdto.AuthUserResponse{
		ID:                 e.ID,
		Name:               e.Name,
		Email:              e.Email,
		PhotoPath:          e.PhotoPath,
		Status:             e.Status,
		MustChangePassword: e.MustChangePassword,
		Role: authdto.RoleInfo{
			ID:          e.Role.ID,
			Name:        e.Role.Name,
			DisplayName: e.Role.DisplayName,
		},
		CreatedAt: e.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
