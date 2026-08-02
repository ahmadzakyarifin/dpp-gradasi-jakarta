package mapper

import (
	"time"

	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/user/dto"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/user/entity"
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
	e := &entity.User{
		Name:   req.Name,
		Email:  req.Email,
		Phone:  req.Phone,
		RoleID: req.RoleID,
		Status: status,
	}
	if req.DateOfBirth != "" {
		if t, err := time.Parse("2006-01-02", req.DateOfBirth); err == nil {
			e.DateOfBirth = &t
		}
	}
	if req.CountryCode != "" {
		cc := req.CountryCode
		e.CountryCode = &cc
	}
	return e
}

// UpdateReqToEntity menerapkan perubahan dari request ke entity yang sudah ada.
func UpdateReqToEntity(req *dto.UserUpdateReq, user *entity.User) {
	if req == nil || user == nil {
		return
	}
	user.Name = req.Name
	user.Email = req.Email
	user.Phone = req.Phone
	user.RoleID = req.RoleID
	if req.Status != "" {
		user.Status = req.Status
	}
	if req.DateOfBirth != "" {
		if t, err := time.Parse("2006-01-02", req.DateOfBirth); err == nil {
			user.DateOfBirth = &t
		}
	} else {
		user.DateOfBirth = nil
	}
	if req.CountryCode != "" {
		cc := req.CountryCode
		user.CountryCode = &cc
	} else {
		user.CountryCode = nil
	}
}

func UserEntityToResponse(user entity.User) dto.UserResponse {
	dateOfBirthStr := formatDatePtr(user.DateOfBirth)
	return dto.UserResponse{
		ID:              user.ID,
		RoleID:          user.RoleID,
		IsSystem:        user.IsSystem,
		RoleName:        user.RoleName,
		RoleDisplayName: user.RoleDisplayName,
		Name:            user.Name,
		Email:           user.Email,
		Phone:           user.Phone,
		DateOfBirth:     dateOfBirthStr,
		CountryCode:     user.CountryCode,
		PhotoPath:       user.PhotoPath,
		Status:          user.Status,

		HasPassword: user.PasswordHash != "",

		EmailVerifiedAt: formatTimePtr(user.EmailVerifiedAt),
		PhoneVerifiedAt: formatTimePtr(user.PhoneVerifiedAt),
		LastLoginAt:     formatTimePtr(user.LastLoginAt),

		CreatedAt: user.CreatedAt.Format("02/01/2006 15:04"),
		UpdatedAt: user.UpdatedAt.Format("02/01/2006 15:04"),

		StudentCount: user.StudentCount,
		StudentNames: user.StudentNames,
	}
}

func UsersEntityToResponse(users []entity.User) []dto.UserResponse {
	res := make([]dto.UserResponse, len(users))

	for i := range users {
		res[i] = UserEntityToResponse(users[i])
	}

	return res
}

func formatDatePtr(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := t.Format("2006-01-02")
	return &s
}

func formatTimePtr(t *time.Time) *string {
	if t == nil {
		return nil
	}

	s := t.Format(dateTimeLayout)
	return &s
}
