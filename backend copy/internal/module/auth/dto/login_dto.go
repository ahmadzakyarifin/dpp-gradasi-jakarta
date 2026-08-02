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
	ID          uint     `json:"id"`
	Name        string   `json:"name"`
	Email       string   `json:"email"`
	Role        AuthRole `json:"role"`
	Permissions []string `json:"permissions"`
	Status      string   `json:"status"`
}
