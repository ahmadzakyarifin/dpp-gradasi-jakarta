package mapper

import (
	"time"

	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/user/dto"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/user/entity"
)

const dateTimeLayout = "2006-01-02 15:04:05"

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
	return dto.UserResponse{
		ID:              user.ID,
		RoleID:          user.RoleID,
		IsSystem:        user.IsSystem,
		RoleName:        user.RoleName,
		RoleDisplayName: user.RoleDisplayName,
		Name:            user.Name,
		Email:           user.Email,
		PhotoPath:       user.PhotoPath,
		Status:          user.Status,

		HasPassword: user.PasswordHash != "",

		EmailVerifiedAt: formatTimePtr(user.EmailVerifiedAt),
		LastLoginAt:     formatTimePtr(user.LastLoginAt),

		CreatedAt: user.CreatedAt.Format("02/01/2006 15:04"),
		UpdatedAt: user.UpdatedAt.Format("02/01/2006 15:04"),
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
