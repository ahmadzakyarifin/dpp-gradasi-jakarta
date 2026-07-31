package dto

// RegisterRequest untuk pendaftaran user baru
type RegisterRequest struct {
	Name            string `json:"name" binding:"required,min=3"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
}

// LoginRequest untuk login
type LoginRequest struct {
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required"`
	RememberMe *bool  `json:"remember_me,omitempty"`
}

// ForgotPasswordRequest untuk lupa password
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ResetPasswordRequest untuk reset password
type ResetPasswordRequest struct {
	Token           string `json:"token" binding:"required"`
	Password        string `json:"password" binding:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
}

// ActivateAccountRequest untuk aktivasi akun pertama kali oleh Admin baru
type ActivateAccountRequest struct {
	Token           string `json:"token" binding:"required"`
	Password        string `json:"password" binding:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
}

// ChangePasswordRequest untuk ganti password saat login
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	Password        string `json:"password" binding:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
}

// RefreshTokenResponse DTO
type RefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
}

// AuthUserResponse menampilkan data user
type AuthUserResponse struct {
	ID        uint     `json:"id"`
	Name      string   `json:"name"`
	Email     string   `json:"email"`
	PhotoPath string   `json:"photo_path,omitempty"`
	Status    string   `json:"status"`
	Role      RoleInfo `json:"role"`
	CreatedAt string   `json:"created_at"`
}

// RoleInfo menampilkan role
type RoleInfo struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
}

// AuthResponse response login/register
type AuthResponse struct {
	AccessToken string           `json:"access_token"`
	User        AuthUserResponse `json:"user"`
}

// ValidateTokenQuery parameter query validate-reset-token
type ValidateTokenQuery struct {
	Token string `form:"token" binding:"required"`
}
