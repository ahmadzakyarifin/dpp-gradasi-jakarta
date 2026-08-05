package dto

import "time"

type LoginRequest struct {
	Email          string `json:"email" binding:"required,email"`
	Password       string `json:"password" binding:"required,min=6"`
	TurnstileToken string `json:"turnstile_token"`
	RememberMe     bool   `json:"remember_me"`
}

type LoginResponse struct {
	AccessToken        string    `json:"access_token"`
	RefreshToken       string    `json:"-"`
	RefreshTokenExpiry time.Time `json:"-"`
	User               AuthUser  `json:"user"`
}

type AuthRole struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
}

type AuthUser struct {
	ID                 uint     `json:"id"`
	Name               string   `json:"name"`
	Email              string   `json:"email"`
	EmailPending       *string  `json:"email_pending,omitempty"`
	EmailVerifiedAt    *string  `json:"email_verified_at,omitempty"`
	PhotoPath          *string  `json:"photo_path"`
	Role               AuthRole `json:"role"`
	Status             string   `json:"status"`
	MustChangePassword bool     `json:"must_change_password"`
	CreatedAt          string   `json:"created_at"`
}
