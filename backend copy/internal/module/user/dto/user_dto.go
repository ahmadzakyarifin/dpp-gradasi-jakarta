package dto

import (
	"time"

	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/user/entity"
)

type UserCreateReq struct {
	Name                string `json:"name" binding:"required,min=2"`
	Email               string `json:"email" binding:"required,email"`
	Phone               string `json:"phone" binding:"required,min=9,max=15"`
	RoleID              uint   `json:"role_id" binding:"required"`
	Status              string `json:"status" binding:"omitempty,oneof=active inactive"`
	DateOfBirth         string `json:"date_of_birth" binding:"omitempty,datetime=2006-01-02"`
	CountryCode         string `json:"country_code" binding:"omitempty,min=2,max=2"`
	NotificationChannel string `json:"notification_channel" binding:"omitempty,oneof=email whatsapp all"`
}

type UserUpdateReq struct {
	Name                string `json:"name" binding:"required,min=2"`
	Email               string `json:"email" binding:"required,email"`
	Phone               string `json:"phone" binding:"required,min=9,max=15"`
	RoleID              uint   `json:"role_id" binding:"required"`
	Status              string `json:"status" binding:"omitempty,oneof=active inactive"`
	DateOfBirth         string `json:"date_of_birth" binding:"omitempty,datetime=2006-01-02"`
	CountryCode         string `json:"country_code" binding:"omitempty,min=2,max=2"`
	NotificationChannel string `json:"notification_channel" binding:"omitempty,oneof=email whatsapp all"`
}

func (req *UserCreateReq) ToEntity() *entity.User {
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

func (req *UserUpdateReq) Apply(e *entity.User) {
	if req == nil || e == nil {
		return
	}
	e.Name = req.Name
	e.Email = req.Email
	e.Phone = req.Phone
	e.RoleID = req.RoleID
	if req.Status != "" {
		e.Status = req.Status
	}
	if req.DateOfBirth != "" {
		if t, err := time.Parse("2006-01-02", req.DateOfBirth); err == nil {
			e.DateOfBirth = &t
		}
	} else {
		e.DateOfBirth = nil
	}
	if req.CountryCode != "" {
		cc := req.CountryCode
		e.CountryCode = &cc
	} else {
		e.CountryCode = nil
	}
}

type UserQueryReq struct {
	Page     int    `form:"page"`
	Limit    int    `form:"limit"`
	Search   string `form:"search"`
	Role     string `form:"role"`
	Status   string `form:"status"`
	Sort     string `form:"sort"`
	Relation string `form:"relation"`
	Trashed  bool   `form:"trashed"`
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
	RoleID          uint    `json:"role_id"`
	IsSystem        bool    `json:"is_system"`
	RoleName        string  `json:"role_name"`
	RoleDisplayName string  `json:"role_display_name"`
	Name            string  `json:"name"`
	Email           string  `json:"email"`
	Phone           string  `json:"phone"`
	DateOfBirth     *string `json:"date_of_birth,omitempty"`
	CountryCode     *string `json:"country_code,omitempty"`
	PhotoPath       *string `json:"photo_path"`
	Status          string  `json:"status"`

	HasPassword bool `json:"has_password"`

	EmailVerifiedAt *string `json:"email_verified_at,omitempty"`
	PhoneVerifiedAt *string `json:"phone_verified_at,omitempty"`
	LastLoginAt     *string `json:"last_login_at,omitempty"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`

	StudentCount int    `json:"student_count"`
	StudentNames string `json:"student_names"`
}
