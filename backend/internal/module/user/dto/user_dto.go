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
	RoleID uint   `json:"role_id" binding:"required,oneof=2 3 4"`
}

type AdminStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=active inactive"`
}

type UserResponse struct {
	ID        uint   `json:"id"`
	RoleID    uint   `json:"role_id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	PhotoPath string `json:"photo_path,omitempty"`
	Status    string `json:"status"`
}

type VerifyEmailRequest struct {
	Token string `json:"token" binding:"required"`
}

type BulkRequest struct {
	IDs []uint `json:"ids" binding:"required,min=1"`
}
