package repository

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/auth/model"
	userentity "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/user/entity"
	usermapper "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/user/mapper"
	usermodel "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/user/model"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

const (
	TokenResetPassword = "reset_password"
)

type authRepo struct {
	db          *gorm.DB
	redisClient *redis.Client
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func NewAuthRepo(db *gorm.DB, redisClient *redis.Client) AuthRepo {
	return &authRepo{db: db, redisClient: redisClient}
}

func (r *authRepo) WithTx(tx *gorm.DB) AuthRepo {
	return &authRepo{db: tx, redisClient: r.redisClient}
}

func (r *authRepo) GetDB() *gorm.DB {
	return r.db
}

func (r *authRepo) GetPermissionsByRoleID(ctx context.Context, roleID uint) ([]string, error) {
	return getCachedRolePermissions(ctx, r.db, r.redisClient, roleID)
}

func (r *authRepo) FindByEmail(ctx context.Context, email string) (*userentity.User, error) {
	var u usermodel.UserModel
	err := r.db.WithContext(ctx).Preload("Role").Where("users.email = ?", email).First(&u).Error
	if err != nil {
		return nil, err
	}
	return usermapper.ModelToUserEntity(&u), nil
}

func (r *authRepo) FindUserByID(ctx context.Context, id uint) (*userentity.User, error) {
	var u usermodel.UserModel
	err := r.db.WithContext(ctx).Preload("Role").Where("users.id = ?", id).First(&u).Error
	if err != nil {
		return nil, err
	}
	return usermapper.ModelToUserEntity(&u), nil
}

func (r *authRepo) SaveRefreshToken(ctx context.Context, userID uint, token string, expiresAt time.Time, ip string, userAgent string, device string) error {
	hash := hashToken(token)
	m := &model.RefreshTokenModel{
		UserID:    userID,
		TokenHash: hash,
		ExpiresAt: expiresAt,
	}

	if ip != "" {
		m.IPAddress = &ip
	}
	if userAgent != "" {
		m.UserAgent = &userAgent
	}

	if device != "" {
		m.DeviceName = &device
	}

	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}

	return nil
}

func (r *authRepo) FindUserByRefreshToken(ctx context.Context, token string) (*userentity.User, time.Time, error) {
	hash := hashToken(token)
	var rt model.RefreshTokenModel
	err := r.db.WithContext(ctx).
		Where("token_hash = ?", hash).
		Where("revoked_at IS NULL").
		Where("expires_at > ?", time.Now()).
		First(&rt).Error
	if err != nil {
		return nil, time.Time{}, err
	}

	user, err := r.FindUserByID(
		ctx,
		rt.UserID,
	)

	if err != nil {
		return nil, time.Time{}, err
	}

	return user, rt.ExpiresAt, nil
}
func (r *authRepo) DeleteRefreshToken(ctx context.Context, token string) error {
	hashedToken := hashToken(token)
	now := time.Now()
	res := r.db.WithContext(ctx).Model(&model.RefreshTokenModel{}).
		Where("token_hash = ? AND revoked_at IS NULL", hashedToken).
		Update("revoked_at", now)
	return res.Error
}

func (r *authRepo) DeleteAllUserRefreshTokens(ctx context.Context, userID uint) error {
	now := time.Now()
	res := r.db.WithContext(ctx).Model(&model.RefreshTokenModel{}).
		Where("user_id = ? AND revoked_at IS NULL", userID).
		Update("revoked_at", now)
	return res.Error
}

func (r *authRepo) UpdatePassword(ctx context.Context, userID uint, hashedPassword string) error {
	res := r.db.WithContext(ctx).
		Model(&usermodel.UserModel{}).
		Where("id = ?", userID).
		Update("password_hash", hashedPassword)
	return res.Error
}

func (r *authRepo) UpdateUserContact(ctx context.Context, userID uint, field string, value string, verifiedAt time.Time) error {
	updates := map[string]any{
		field: value,
	}
	switch field {
	case "email":
		updates["email_verified_at"] = verifiedAt
	case "phone":
		updates["phone_verified_at"] = verifiedAt
	}

	return r.db.WithContext(ctx).
		Model(&usermodel.UserModel{}).
		Where("id = ?", userID).
		Updates(updates).Error
}

func (r *authRepo) SaveAuthToken(ctx context.Context, userID uint, token string, tokenType string, expiresAt time.Time) error {
	hash := hashToken(token)

	switch tokenType {

	case TokenResetPassword:

		m := &model.PasswordResetTokenModel{
			UserID:    userID,
			TokenHash: hash,
			ExpiresAt: expiresAt,
		}

		if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
			return err
		}

		return nil

	default:
		return fmt.Errorf(
			"unknown auth token type: %s",
			tokenType,
		)
	}
}

func (r *authRepo) FindAuthToken(ctx context.Context, token string, tokenType string) (uint, error) {
	hash := hashToken(token)

	switch tokenType {
	case "reset_password":
		var row model.PasswordResetTokenModel
		err := r.db.WithContext(ctx).
			Where("token_hash=?", hash).
			Where("used_at IS NULL").
			Where("expires_at > ?", time.Now()).
			First(&row).Error

		if err != nil {
			return 0, err
		}

		return row.UserID, nil
	default:

		return 0, fmt.Errorf(
			"unknown token type",
		)
	}
}

func (r *authRepo) DeleteAuthToken(ctx context.Context, token string, tokenType string) error {
	hash := hashToken(token)
	now := time.Now()

	switch tokenType {
	case "reset_password":
		res := r.db.WithContext(ctx).
			Model(&model.PasswordResetTokenModel{}).
			Where("token_hash=?", hash).
			Update("used_at", now)
		return res.Error

	default:
		return fmt.Errorf(
			"unknown token type",
		)
	}
}

func getCachedRolePermissions(ctx context.Context, db *gorm.DB, redisClient *redis.Client, roleID uint) ([]string, error) {
	if redisClient == nil {
		return fetchRolePermissionsFromDB(ctx, db, roleID)
	}

	cacheKey := fmt.Sprintf("role_permissions:%d", roleID)
	cached, err := redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var perms []string
		if err := json.Unmarshal([]byte(cached), &perms); err == nil {
			return perms, nil
		}
	}

	perms, err := fetchRolePermissionsFromDB(ctx, db, roleID)
	if err != nil {
		return nil, err
	}

	if b, err := json.Marshal(perms); err == nil {
		redisClient.Set(ctx, cacheKey, b, 24*time.Hour)
	}

	return perms, nil
}

func fetchRolePermissionsFromDB(ctx context.Context, db *gorm.DB, roleID uint) ([]string, error) {
	var names []string
	err := db.WithContext(ctx).
		Table("role_permissions AS rp").
		Select("p.name").
		Joins("JOIN permissions AS p ON rp.permission_id = p.id").
		Where("rp.role_id = ?", roleID).
		Scan(&names).Error
	return names, err
}
