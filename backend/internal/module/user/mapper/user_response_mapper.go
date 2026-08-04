package mapper

import (
	"time"

	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/user/dto"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/user/entity"
)

const dateTimeLayout = "2006-01-02T15:04:05Z"

// CreateReqToEntity memetakan request pembuatan user ke domain entity.
func CreateReqToEntity(req *dto.UserCreateReq) *entity.User {
	if req == nil {
		return nil
	}
	status := req.Status
	if status == "" {
		status = "inactive"
	}
	return &entity.User{
		Name:   req.Name,
		Email:  req.Email,
		RoleID: req.RoleID,
		Status: status,
	}
}

// UpdateReqToEntity menerapkan perubahan dari request ke entity yang sudah ada.
func UpdateReqToEntity(req *dto.UserUpdateReq, user *entity.User) {
	if req == nil || user == nil {
		return
	}
	user.Name = req.Name
	user.Email = req.Email
	user.RoleID = req.RoleID
	if req.Status != "" {
		user.Status = req.Status
	}
}

func UserEntityToResponse(user entity.User) dto.UserResponse {
	roleName := user.Role
	roleDisplayName := "Admin"
	if user.Role == "super_admin" {
		roleDisplayName = "Super Administrator"
	}
	return dto.UserResponse{
		ID:              user.ID,
		Role:            user.Role,
		IsSystem:        user.IsSystem,
		RoleName:        roleName,
		RoleDisplayName: roleDisplayName,
		Name:            user.Name,
		Email:           user.Email,
		PhotoPath:       user.PhotoPath,
		Status:          user.Status,

		HasPassword:        user.Password != "",
		MustChangePassword: user.MustChangePassword,

		EmailVerifiedAt: formatTimePtr(user.EmailVerifiedAt),
		LastLoginAt:     formatTimePtr(user.LastLoginAt),

		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

func UsersEntityToResponse(users []entity.User) []dto.UserResponse {
	res := make([]dto.UserResponse, len(users))

	for i := range users {
		res[i] = UserEntityToResponse(users[i])
	}

	return res
}

func formatTimePtr(t *time.Time) *string {
	if t == nil {
		return nil
	}

	s := t.Format(dateTimeLayout)
	return &s
}
