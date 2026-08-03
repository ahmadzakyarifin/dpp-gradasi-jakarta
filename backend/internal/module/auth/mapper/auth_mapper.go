package mapper

import (
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/auth/dto"
	userentity "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/user/entity"
)

func UserEntityToAuth(user userentity.User) dto.AuthUser {
	return dto.AuthUser{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		PhotoPath: user.PhotoPath,
		Role: dto.AuthRole{
			ID:          user.RoleID,
			Name:        user.RoleName,
			DisplayName: user.RoleDisplayName,
		},
		Status:             user.Status,
		MustChangePassword: user.MustChangePassword,
		CreatedAt:          user.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
