package dto

import (
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/user/entity"
)

type UserCreateReq struct {
	Name   string `json:"name" binding:"required,min=2"`
	Email  string `json:"email" binding:"required,email"`
	Role   string `json:"role" binding:"required,oneof=super_admin admin"`
	Status string `json:"status" binding:"omitempty,oneof=active inactive"`
}

type UserUpdateReq struct {
	Name   string `json:"name" binding:"required,min=2"`
	Email  string `json:"email" binding:"required,email"`
	Role   string `json:"role" binding:"required,oneof=super_admin admin"`
	Status string `json:"status" binding:"omitempty,oneof=active inactive"`
}

func (req *UserCreateReq) ToEntity() *entity.User {
	status := req.Status
	if status == "" {
		status = "inactive"
	}
	return &entity.User{
		Name:   req.Name,
		Email:  req.Email,
		Role:   req.Role,
		Status: status,
	}
}

func (req *UserUpdateReq) Apply(e *entity.User) {
	if req == nil || e == nil {
		return
	}
	e.Name = req.Name
	e.Email = req.Email
	e.Role = req.Role
	if req.Status != "" {
		e.Status = req.Status
	}
}

type UserQueryReq struct {
	Page    int    `form:"page"`
	Limit   int    `form:"limit"`
	Search  string `form:"search"`
	Role    string `form:"role"`
	Status  string `form:"status"`
	Sort    string `form:"sort"`
	Trashed bool   `form:"trashed"`
}

func (q *UserQueryReq) Normalize() {
	if q.Page <= 0 {
		q.Page = 1
	}
	if q.Limit <= 0 {
		q.Limit = 10
	}
	if q.Limit > 100 {
		q.Limit = 100
	}
}

type UserResponse struct {
	ID              uint    `json:"id"`
	Role            string  `json:"role"`
	IsSystem        bool    `json:"is_system"`
	RoleName        string  `json:"role_name"`
	RoleDisplayName string  `json:"role_display_name"`
	Name            string  `json:"name"`
	Email           string  `json:"email"`
	EmailPending    *string `json:"email_pending,omitempty"`
	PhotoPath       *string `json:"photo_path"`
	Status          string  `json:"status"`

	HasPassword        bool `json:"has_password"`
	MustChangePassword bool `json:"must_change_password"`

	EmailVerifiedAt *string `json:"email_verified_at,omitempty"`
	LastLoginAt     *string `json:"last_login_at,omitempty"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
}
