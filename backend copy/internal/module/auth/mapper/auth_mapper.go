package mapper

import (
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/auth/dto"
	roleentity "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/role/entity"
	userentity "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/user/entity"
)

func UserEntityToAuth(user userentity.User, permissions []string) dto.AuthUser {
	return dto.AuthUser{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role: dto.AuthRole{
			ID:          user.RoleID,
			Name:        user.RoleName,
			DisplayName: user.RoleDisplayName,
		},
		Permissions: permissions,
		Status:      user.Status,
	}
}

func PermissionsToNames(perms []roleentity.Permission) []string {
	names := make([]string, len(perms))
	for i := range perms {
		names[i] = perms[i].Name
	}
	return names
}
