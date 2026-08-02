package dto

import "mime/multipart"

type ProfileUpdateRequest struct {
	Name  string                `form:"name" binding:"required,min=3,max=100"`
	Email string                `form:"email" binding:"required,email,max=150"`
	Photo *multipart.FileHeader `form:"photo"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

type AdminCreateRequest struct {
	Name   string `json:"name" binding:"required,min=3,max=100"`
	Email  string `json:"email" binding:"required,email,max=150"`
	RoleID uint   `json:"role_id" binding:"required,oneof=2 5 6"`
}

type AdminStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=active inactive"`
}

type UserResponse struct {
	ID                 uint   `json:"id"`
	RoleID             uint   `json:"role_id"`
	RoleName           string `json:"role,omitempty"`
	Name               string `json:"name"`
	Email              string `json:"email"`
	PhotoPath          string `json:"photo_path,omitempty"`
	Status             string `json:"status"`
	MustChangePassword bool   `json:"must_change_password"`
}

// ListUsersQuery adalah query params untuk GET /admin/users
type ListUsersQuery struct {
	Tab    string `form:"tab" binding:"omitempty,oneof=active pending trash"`
	Search string `form:"search"`
	Page   int    `form:"page" binding:"omitempty,min=1"`
	Limit  int    `form:"limit" binding:"omitempty,min=1,max=100"`
}

// UserListResponse adalah response daftar admin dengan pagination
type UserListResponse struct {
	Items      []UserResponse `json:"items"`
	Pagination Pagination     `json:"pagination"`
}

type Pagination struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

type VerifyEmailRequest struct {
	Token string `json:"token" binding:"required"`
}

type BulkRequest struct {
	IDs []uint `json:"ids" binding:"required,min=1"`
}
