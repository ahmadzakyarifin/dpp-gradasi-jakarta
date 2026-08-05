package mapper

import (
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/auth/dto"
	userentity "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/user/entity"
)

func UserEntityToAuth(user userentity.User) dto.AuthUser {
	var roleID uint = 10
	roleDisplayName := "Admin"
	if user.Role == "superadmin" || user.Role == "super_admin" {
		roleID = 9
		roleDisplayName = "Super Administrator"
	}

	return dto.AuthUser{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		PhotoPath: user.PhotoPath,
		Role: dto.AuthRole{
			ID:          roleID,
			Name:        user.Role,
			DisplayName: roleDisplayName,
		},
		Status:             user.Status,
		MustChangePassword: user.MustChangePassword,
		CreatedAt:          user.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
