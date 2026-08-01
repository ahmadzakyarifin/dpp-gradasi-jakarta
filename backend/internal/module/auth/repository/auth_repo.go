package repository

import (
	"time"

	authmodel "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/auth/model"
	rolemodel "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/role/model"
	usermodel "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/user/model"
	"gorm.io/gorm"
)

type AuthRepo interface {
	// User
	FindByEmail(email string) (*usermodel.User, error)
	FindByID(id uint) (*usermodel.User, error)
	CreateUser(name, email, hashedPassword string) (*usermodel.User, error)
	UpdateLastLogin(id uint) error

	// Refresh Token
	SaveRefreshToken(token *authmodel.RefreshToken) error
	FindRefreshTokenByHash(hash string) (*authmodel.RefreshToken, error)
	DeleteRefreshToken(id uint) error
	DeleteUserRefreshTokens(userID uint) error

	// Password Reset
	SavePasswordResetToken(token *authmodel.PasswordResetToken) error
	FindPasswordResetTokenByHash(hash string) (*authmodel.PasswordResetToken, error)
	MarkResetTokenUsed(id uint) error

	// Activation
	SaveActivationToken(token *authmodel.ActivationToken) error
	FindActivationTokenByHash(hash string) (*authmodel.ActivationToken, error)
	MarkActivationTokenUsed(id uint) error
	ActivateUser(userID uint) error
	SetUserPassword(userID uint, hashedPassword string) error

	// Role
	FindRoleByID(id uint) (*rolemodel.Role, error)
	FindRoleByName(name string) (*rolemodel.Role, error)
}

type authRepo struct {
	db *gorm.DB
}

func NewAuthRepo(db *gorm.DB) AuthRepo {
	return &authRepo{db: db}
}

func (r *authRepo) DB() *gorm.DB {
	return r.db
}

func (r *authRepo) FindByEmail(email string) (*usermodel.User, error) {
	var user usermodel.User
	err := r.db.Where("email = ?", email).Preload("Role").First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepo) FindByID(id uint) (*usermodel.User, error) {
	var user usermodel.User
	err := r.db.Preload("Role").First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepo) CreateUser(name, email, hashedPassword string) (*usermodel.User, error) {
	// Find default role (admin)
	role, err := r.FindRoleByName("admin")
	if err != nil {
		// Fallback to role_id=2
		role = &rolemodel.Role{ID: 2, Name: "admin", DisplayName: "Administrator"}
	}

	user := usermodel.User{
		RoleID:   role.ID,
		Name:     name,
		Email:    email,
		Password: hashedPassword,
		Status:   "active",
	}

	if err := r.db.Create(&user).Error; err != nil {
		return nil, err
	}

	// Reload with role
	var loaded usermodel.User
	if err := r.db.Preload("Role").First(&loaded, user.ID).Error; err != nil {
		return nil, err
	}

	return &loaded, nil
}

func (r *authRepo) UpdateLastLogin(id uint) error {
	now := time.Now()
	return r.db.Model(&usermodel.User{}).Where("id = ?", id).Update("last_login_at", &now).Error
}

func (r *authRepo) SaveRefreshToken(token *authmodel.RefreshToken) error {
	return r.db.Create(token).Error
}

func (r *authRepo) FindRefreshTokenByHash(hash string) (*authmodel.RefreshToken, error) {
	var token authmodel.RefreshToken
	err := r.db.Where("token_hash = ? AND expires_at > ?", hash, time.Now()).First(&token).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *authRepo) DeleteRefreshToken(id uint) error {
	return r.db.Delete(&authmodel.RefreshToken{}, id).Error
}

func (r *authRepo) DeleteUserRefreshTokens(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&authmodel.RefreshToken{}).Error
}

func (r *authRepo) SavePasswordResetToken(token *authmodel.PasswordResetToken) error {
	return r.db.Create(token).Error
}

func (r *authRepo) FindPasswordResetTokenByHash(hash string) (*authmodel.PasswordResetToken, error) {
	var token authmodel.PasswordResetToken
	err := r.db.Where("token_hash = ? AND expires_at > ? AND used_at IS NULL", hash, time.Now()).First(&token).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *authRepo) MarkResetTokenUsed(id uint) error {
	now := time.Now()
	return r.db.Model(&authmodel.PasswordResetToken{}).Where("id = ?", id).Update("used_at", &now).Error
}

func (r *authRepo) SaveActivationToken(token *authmodel.ActivationToken) error {
	return r.db.Create(token).Error
}

func (r *authRepo) FindActivationTokenByHash(hash string) (*authmodel.ActivationToken, error) {
	var token authmodel.ActivationToken
	err := r.db.Where("token_hash = ? AND expires_at > ? AND used_at IS NULL", hash, time.Now()).First(&token).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *authRepo) MarkActivationTokenUsed(id uint) error {
	now := time.Now()
	return r.db.Model(&authmodel.ActivationToken{}).Where("id = ?", id).Update("used_at", &now).Error
}

func (r *authRepo) ActivateUser(userID uint) error {
	now := time.Now()
	return r.db.Model(&usermodel.User{}).Where("id = ?", userID).Updates(map[string]any{
		"status":            "active",
		"email_verified_at": &now,
	}).Error
}

func (r *authRepo) SetUserPassword(userID uint, hashedPassword string) error {
	return r.db.Model(&usermodel.User{}).Where("id = ?", userID).Updates(map[string]any{
		"password":             hashedPassword,
		"must_change_password": false,
	}).Error
}

func (r *authRepo) FindRoleByID(id uint) (*rolemodel.Role, error) {
	var role rolemodel.Role
	err := r.db.First(&role, id).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *authRepo) FindRoleByName(name string) (*rolemodel.Role, error) {
	var role rolemodel.Role
	err := r.db.Where("name = ?", name).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}
