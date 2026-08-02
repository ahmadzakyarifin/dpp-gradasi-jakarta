package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/user/entity"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/user/mapper"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/user/model"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type UserRepo interface {
	Create(ctx context.Context, db *gorm.DB, user *entity.User) error
	CreateTx(ctx context.Context, user *entity.User) error
	FindAll(ctx context.Context, roleID uint) ([]entity.User, error)
	FindByID(ctx context.Context, id uint) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	FindByPhone(ctx context.Context, phone string) (*entity.User, error)
	FindByNIK(ctx context.Context, nik string) (*entity.User, error)
	Update(ctx context.Context, db *gorm.DB, user *entity.User) error
	UpdateTx(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id uint) error
	UpdatePassword(ctx context.Context, id uint, passwordHash string) error
	Activate(ctx context.Context, id uint) error
	ToggleStatus(ctx context.Context, id uint) error
	FindPaginated(ctx context.Context, page, limit int, search, roleIDStr, status, sort, relation string, trashed bool) ([]entity.User, int, error)
	CountActive(ctx context.Context) (int, error)
	CountActiveByPeriod(ctx context.Context, start, end *time.Time) (int, error)
	HasRolePermissions(ctx context.Context, roleID uint, permissions []string) (bool, error)
	HasAnyRolePermission(ctx context.Context, roleID uint, permissions []string) (bool, error)
	BulkDelete(ctx context.Context, ids []uint) error
	Restore(ctx context.Context, id uint) error
	BulkRestore(ctx context.Context, ids []uint) error
	CountStudentsByParent(ctx context.Context, parentID uint) (int, error)
}

type userRepo struct {
	db          *gorm.DB
	redisClient *redis.Client
}

func NewUserRepo(db *gorm.DB, redisClient *redis.Client) UserRepo {
	return &userRepo{db: db, redisClient: redisClient}
}

func (r *userRepo) Create(ctx context.Context, db *gorm.DB, u *entity.User) error {
	if db == nil {
		db = r.db
	}
	m := mapper.EntityToUserModel(u)
	if err := db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	u.ID = m.ID
	u.CreatedAt = m.CreatedAt
	u.UpdatedAt = m.UpdatedAt
	return nil
}

func (r *userRepo) FindAll(ctx context.Context, roleID uint) ([]entity.User, error) {
	q := r.db.WithContext(ctx).Preload("Role")
	if roleID != 0 {
		q = q.Where("role_id = ?", roleID)
	}
	var models []model.UserModel
	if err := q.Find(&models).Error; err != nil {
		return nil, err
	}
	users := make([]entity.User, len(models))
	for i := range models {
		users[i] = *mapper.ModelToUserEntity(&models[i])
	}
	return users, nil
}

func (r *userRepo) FindByID(ctx context.Context, id uint) (*entity.User, error) {
	var m model.UserModel
	err := r.db.WithContext(ctx).Unscoped().Preload("Role").Where("users.id = ?", id).First(&m).Error
	if err != nil {
		return nil, err
	}
	return mapper.ModelToUserEntity(&m), nil
}

func (r *userRepo) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	var m model.UserModel
	err := r.db.WithContext(ctx).Unscoped().Preload("Role").Where("email = ?", email).First(&m).Error
	if err != nil {
		return nil, err
	}
	return mapper.ModelToUserEntity(&m), nil
}

func (r *userRepo) FindByPhone(ctx context.Context, phone string) (*entity.User, error) {
	var m model.UserModel
	err := r.db.WithContext(ctx).Preload("Role").Where("phone = ?", phone).First(&m).Error
	if err != nil {
		return nil, err
	}
	return mapper.ModelToUserEntity(&m), nil
}

func (r *userRepo) FindByNIK(ctx context.Context, nik string) (*entity.User, error) {
	// NIK sudah dipindah ke tabel guardians.
	return nil, fmt.Errorf("nik moved to guardians table")
}

func (r *userRepo) Update(ctx context.Context, db *gorm.DB, u *entity.User) error {
	if db == nil {
		db = r.db
	}
	m := mapper.EntityToUserModel(u)
	updates := map[string]any{
		"role_id":       m.RoleID,
		"name":          m.Name,
		"email":         m.Email,
		"phone":         m.Phone,
		"photo_path":    m.PhotoPath,
		"status":        m.Status,
		"date_of_birth": m.DateOfBirth,
		"country_code":  m.CountryCode,
	}
	return db.WithContext(ctx).Model(&model.UserModel{}).
		Where("id = ?", m.ID).
		Updates(updates).Error
}

func (r *userRepo) CreateTx(ctx context.Context, u *entity.User) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return r.Create(ctx, tx, u)
	})
}

func (r *userRepo) UpdateTx(ctx context.Context, u *entity.User) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return r.Update(ctx, tx, u)
	})
}

func (r *userRepo) Delete(ctx context.Context, id uint) error {
	res := r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.UserModel{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("pengguna sudah terhapus atau tidak ditemukan")
	}
	return nil
}

func (r *userRepo) BulkDelete(ctx context.Context, ids []uint) error {
	res := r.db.WithContext(ctx).Where("id IN ?", ids).Delete(&model.UserModel{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("tidak ada data yang perlu dihapus")
	}
	return nil
}

func (r *userRepo) Restore(ctx context.Context, id uint) error {
	res := r.db.WithContext(ctx).Unscoped().Model(&model.UserModel{}).
		Where("id = ? AND deleted_at IS NOT NULL", id).
		Update("deleted_at", nil)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("data tidak ditemukan di tempat sampah atau sudah aktif")
	}
	return nil
}

func (r *userRepo) BulkRestore(ctx context.Context, ids []uint) error {
	res := r.db.WithContext(ctx).Unscoped().Model(&model.UserModel{}).
		Where("id IN ? AND deleted_at IS NOT NULL", ids).
		Update("deleted_at", nil)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("tidak ada data yang perlu dipulihkan")
	}
	return nil
}

func (r *userRepo) CountStudentsByParent(ctx context.Context, parentID uint) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.UserModel{}).
		Where("user_id = ? AND deleted_at IS NULL", parentID).
		Count(&count).Error
	return int(count), err
}

func (r *userRepo) UpdatePassword(ctx context.Context, id uint, hash string) error {
	res := r.db.WithContext(ctx).Model(&model.UserModel{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"password_hash": hash,
			"status":        "active",
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("pengguna tidak ditemukan")
	}
	return nil
}

func (r *userRepo) Activate(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&model.UserModel{}).
		Where("id = ?", id).
		Update("status", "active").Error
}

func (r *userRepo) ToggleStatus(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Model(&model.UserModel{}).
		Where("id = ?", id).
		Update("status", gorm.Expr("CASE WHEN status = 'active' THEN 'inactive' ELSE 'active' END")).Error
}

func (r *userRepo) FindPaginated(ctx context.Context, page, limit int, search, roleIDStr, status, sort, relation string, trashed bool) ([]entity.User, int, error) {
	studentCountSub := "(SELECT COUNT(id) FROM students WHERE user_id = users.id AND deleted_at IS NULL)"

	q := r.db.WithContext(ctx).Model(&model.UserModel{})

	if trashed {
		q = q.Unscoped().Where("users.deleted_at IS NOT NULL")
	} else {
		switch status {
		case "active":
			q = q.Where("users.status = ?", "active")
		case "inactive":
			q = q.Where("users.status = ?", "inactive")
		}
	}

	switch relation {
	case "no_child":
		q = q.Where(studentCountSub + " = 0")
	case "has_child":
		q = q.Where(studentCountSub + " > 0")
	}

	if search != "" {
		s := "%" + search + "%"
		q = q.Where("users.name LIKE ? OR users.email LIKE ? OR users.phone LIKE ?", s, s, s)
	}

	if roleIDStr != "" {
		q = q.Where("users.role_id = ?", roleIDStr)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	switch sort {
	case "name_asc":
		q = q.Order("users.name ASC")
	case "name_desc":
		q = q.Order("users.name DESC")
	case "created_asc":
		q = q.Order("users.created_at ASC")
	case "created_desc":
		q = q.Order("users.created_at DESC")
	default:
		q = q.Order("users.created_at DESC")
	}

	studentNamesSub := `(SELECT GROUP_CONCAT(CONCAT_WS('::', CAST(s.id AS CHAR), s.name) ORDER BY s.name ASC SEPARATOR '||') FROM students s WHERE s.user_id = users.id AND s.deleted_at IS NULL) AS student_names`

	var models []model.UserModel
	err := q.
		Select("users.*", studentCountSub+" AS student_count", studentNamesSub).
		Preload("Role").
		Limit(limit).
		Offset((page - 1) * limit).
		Find(&models).Error
	if err != nil {
		return nil, 0, err
	}

	users := make([]entity.User, len(models))
	for i := range models {
		users[i] = *mapper.ModelToUserEntity(&models[i])
	}

	return users, int(total), nil
}

func (r *userRepo) CountActive(ctx context.Context) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.UserModel{}).Where("status = ?", "active").Count(&count).Error
	return int(count), err
}

func (r *userRepo) CountActiveByPeriod(ctx context.Context, start, end *time.Time) (int, error) {
	q := r.db.WithContext(ctx).Model(&model.UserModel{})
	if start != nil {
		q = q.Where("created_at >= ?", start)
	}
	if end != nil {
		q = q.Where("created_at <= ?", end)
	}
	var count int64
	err := q.Count(&count).Error
	return int(count), err
}

func (r *userRepo) HasRolePermissions(ctx context.Context, roleID uint, permissions []string) (bool, error) {
	if len(permissions) == 0 {
		return false, nil
	}
	cachedPerms, err := getCachedRolePermissions(ctx, r.db, r.redisClient, roleID)
	if err != nil {
		return false, err
	}
	permMap := make(map[string]bool)
	for _, p := range cachedPerms {
		permMap[p] = true
	}
	for _, p := range permissions {
		if !permMap[p] {
			return false, nil
		}
	}
	return true, nil
}

func (r *userRepo) HasAnyRolePermission(ctx context.Context, roleID uint, permissions []string) (bool, error) {
	if len(permissions) == 0 {
		return false, nil
	}
	cachedPerms, err := getCachedRolePermissions(ctx, r.db, r.redisClient, roleID)
	if err != nil {
		return false, err
	}
	permMap := make(map[string]bool)
	for _, p := range cachedPerms {
		permMap[p] = true
	}
	for _, p := range permissions {
		if permMap[p] {
			return true, nil
		}
	}
	return false, nil
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
