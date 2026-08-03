package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/helper"
	activitylogdto "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/activitylog/dto"
	activitylogservice "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/activitylog/service"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/role/dto"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/role/entity"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/role/mapper"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/role/repository"
	"gorm.io/gorm"
)

type RoleService interface {
	GetAllRoles(ctx context.Context, req dto.RoleQueryReq) ([]entity.Role, int, error)
	GetRoleByID(ctx context.Context, id uint) (*entity.Role, error)
	CreateRole(ctx context.Context, req dto.RoleCreateReq) (*entity.Role, error)
	UpdateRole(ctx context.Context, id uint, req dto.RoleUpdateReq) (*entity.Role, error)
	DeleteRole(ctx context.Context, id uint) error
	RestoreRole(ctx context.Context, id uint) error
	UpdateStatus(ctx context.Context, id uint, req dto.RoleStatusReq) error
	BulkDelete(ctx context.Context, ids []uint) error
	BulkRestore(ctx context.Context, ids []uint) error
	GetDependencyInfo(ctx context.Context, id uint) (map[string]interface{}, error)
	CheckUnique(ctx context.Context, field string, value string, excludeID uint) (bool, error)
}

type roleService struct {
	db    *gorm.DB
	repo  repository.RoleRepo
	audit activitylogservice.ActivityLogService
}

func NewRoleService(db *gorm.DB, repo repository.RoleRepo, audit activitylogservice.ActivityLogService) RoleService {
	return &roleService{db: db, repo: repo, audit: audit}
}

func (s *roleService) log(ctx context.Context, db *gorm.DB, input *activitylogdto.ActivityLogInput) {
	if s.audit == nil {
		return
	}
	userID, userName, role, ipAddress, userAgent := helper.GetAuditMeta(ctx)
	var uID *uint
	if userID > 0 {
		uID = &userID
	}
	input.ActorID = uID
	input.ActorName = userName
	input.ActorRole = role
	input.IPAddress = ipAddress
	input.UserAgent = userAgent

	_ = s.audit.Log(ctx, db, input)
}

func roleDependencyMessage(userCount int) string {
	var messages []string
	if userCount > 0 {
		messages = append(messages, fmt.Sprintf("%d pengguna aktif", userCount))
	}
	return strings.Join(messages, ", ")
}

func (s *roleService) GetAllRoles(ctx context.Context, req dto.RoleQueryReq) ([]entity.Role, int, error) {
	req.Normalize()
	roles, total, err := s.repo.GetAllRoles(ctx, req.Page, req.Limit, req.Search, req.Status, req.Sort)
	if err != nil {
		return nil, 0, err
	}

	// Populate user_count per role (N+1 diterima — data role sedikit).
	for i := range roles {
		count, err := s.repo.CountUsers(ctx, roles[i].ID)
		if err != nil {
			return nil, 0, err
		}
		roles[i].UserCount = count
	}

	return roles, total, nil
}

func (s *roleService) GetRoleByID(ctx context.Context, id uint) (*entity.Role, error) {
	existing, err := s.repo.GetRoleByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, helper.NewNotFoundError("Role tidak ditemukan")
	}
	count, err := s.repo.CountUsers(ctx, existing.ID)
	if err != nil {
		return nil, err
	}
	existing.UserCount = count
	return existing, nil
}

func (s *roleService) CreateRole(ctx context.Context, req dto.RoleCreateReq) (*entity.Role, error) {
	v := helper.NewValidationError()

	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		v.Add("name", "Kode role wajib diisi")
	}

	req.DisplayName = strings.TrimSpace(req.DisplayName)
	if req.DisplayName == "" {
		v.Add("display_name", "Nama tampilan role wajib diisi")
	}

	// Validasi duplikat name/display_name
	existingRole, err := s.repo.GetRoleByNameWithDeleted(ctx, req.Name)
	if err == nil && existingRole != nil && existingRole.DeletedAt.Valid {
		v.Add("name", "Kode role ini sudah pernah digunakan dan saat ini berada di keranjang sampah. Silakan pulihkan data tersebut.")
	} else if existingRole != nil {
		v.Add("name", "Kode role ini sudah digunakan")
	}

	existingDisplay, err := s.repo.GetRoleByDisplayNameWithDeleted(ctx, req.DisplayName)
	if err == nil && existingDisplay != nil && existingDisplay.DeletedAt.Valid {
		v.Add("display_name", "Nama tampilan ini sudah pernah digunakan dan saat ini berada di keranjang sampah. Silakan pulihkan data tersebut.")
	} else if existingDisplay != nil {
		v.Add("display_name", "Nama tampilan ini sudah digunakan")
	}

	if len(v.Errors) > 0 {
		return nil, v
	}

	role := mapper.CreateReqToEntity(&req)

	err = s.repo.CreateRole(ctx, role)
	if err != nil {
		return nil, err
	}

	s.log(ctx, s.db, &activitylogdto.ActivityLogInput{
		Action:      "roles.create",
		EntityType:  "roles",
		EntityID:    &role.ID,
		EntityLabel: role.DisplayName,
		Description: fmt.Sprintf("Membuat role baru: %s", role.DisplayName),
		Metadata: map[string]any{
			"new_values": map[string]interface{}{
				"name":         role.Name,
				"display_name": role.DisplayName,
				"is_active":    role.IsActive,
			},
		},
	})

	return role, nil
}

func (s *roleService) UpdateRole(ctx context.Context, id uint, req dto.RoleUpdateReq) (*entity.Role, error) {
	existing, err := s.repo.GetRoleByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, helper.NewNotFoundError("Role tidak ditemukan")
	}

	if existing.IsSystem {
		v := helper.NewValidationError()
		v.Add("general", "Role bawaan sistem tidak dapat diubah")
		return nil, v
	}

	v := helper.NewValidationError()

	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		v.Add("name", "Kode role wajib diisi")
	}

	req.DisplayName = strings.TrimSpace(req.DisplayName)
	if req.DisplayName == "" {
		v.Add("display_name", "Nama tampilan role wajib diisi")
	}

	existingRole, err := s.repo.GetRoleByNameWithDeleted(ctx, req.Name)
	if err == nil && existingRole != nil && existingRole.ID != id {
		if existingRole.DeletedAt.Valid {
			v.Add("name", "Kode role ini sudah pernah digunakan dan saat ini berada di keranjang sampah. Silakan pulihkan data tersebut.")
		} else {
			v.Add("name", "Kode role ini sudah digunakan")
		}
	}

	existingDisplay, err := s.repo.GetRoleByDisplayNameWithDeleted(ctx, req.DisplayName)
	if err == nil && existingDisplay != nil && existingDisplay.ID != id {
		if existingDisplay.DeletedAt.Valid {
			v.Add("display_name", "Nama tampilan ini sudah pernah digunakan dan saat ini berada di keranjang sampah. Silakan pulihkan data tersebut.")
		} else {
			v.Add("display_name", "Nama tampilan ini sudah digunakan")
		}
	}

	if len(v.Errors) > 0 {
		return nil, v
	}

	oldVals := map[string]interface{}{
		"name":         existing.Name,
		"display_name": existing.DisplayName,
		"is_active":    existing.IsActive,
	}

	mapper.UpdateReqToEntity(&req, existing)

	err = s.repo.UpdateRole(ctx, existing)
	if err != nil {
		return nil, err
	}

	s.log(ctx, s.db, &activitylogdto.ActivityLogInput{
		Action:      "roles.update",
		EntityType:  "roles",
		EntityID:    &existing.ID,
		EntityLabel: existing.DisplayName,
		Description: fmt.Sprintf("Memperbarui role: %s", existing.DisplayName),
		Metadata: map[string]any{
			"old_values": oldVals,
			"new_values": map[string]interface{}{
				"name":         existing.Name,
				"display_name": existing.DisplayName,
				"is_active":    existing.IsActive,
			},
		},
	})

	return existing, nil
}

func (s *roleService) DeleteRole(ctx context.Context, id uint) error {
	existing, err := s.repo.GetRoleByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return helper.NewNotFoundError("Role tidak ditemukan")
	}

	if existing.IsSystem {
		v := helper.NewValidationError()
		v.Add("general", "Role bawaan sistem tidak dapat dihapus")
		return v
	}

	hasUsers, err := s.repo.HasUsers(ctx, id)
	if err != nil {
		return err
	}
	if hasUsers {
		v := helper.NewValidationError()
		v.Add("general", "Role ini sedang digunakan oleh pengguna aktif dan tidak dapat dihapus")
		return v
	}

	err = s.repo.DeleteRole(ctx, id)
	if err != nil {
		return err
	}

	s.log(ctx, s.db, &activitylogdto.ActivityLogInput{
		Action:      "roles.delete",
		EntityType:  "roles",
		EntityID:    &id,
		EntityLabel: existing.DisplayName,
		Description: fmt.Sprintf("Menghapus role: %s", existing.DisplayName),
		Metadata: map[string]any{
			"old_values": map[string]interface{}{
				"is_active": existing.IsActive,
				"status":    "active",
			},
			"new_values": map[string]interface{}{
				"is_active": false,
				"status":    "deleted",
			},
		},
	})

	return nil
}

func (s *roleService) RestoreRole(ctx context.Context, id uint) error {
	existing, err := s.repo.GetRoleByIDUnscoped(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return helper.NewNotFoundError("Role tidak ditemukan")
	}

	err = s.repo.Restore(ctx, id)
	if err != nil {
		return err
	}

	s.log(ctx, s.db, &activitylogdto.ActivityLogInput{
		Action:      "roles.restore",
		EntityType:  "roles",
		EntityID:    &id,
		EntityLabel: existing.DisplayName,
		Description: fmt.Sprintf("Memulihkan role: %s", existing.DisplayName),
		Metadata: map[string]any{
			"old_values": map[string]interface{}{
				"is_active": existing.IsActive,
				"status":    "deleted",
			},
			"new_values": map[string]interface{}{
				"is_active": true,
				"status":    "active",
			},
		},
	})

	return nil
}

func (s *roleService) UpdateStatus(ctx context.Context, id uint, req dto.RoleStatusReq) error {
	existing, err := s.repo.GetRoleByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return helper.NewNotFoundError("Role tidak ditemukan")
	}

	if existing.IsSystem {
		v := helper.NewValidationError()
		v.Add("general", "Status role bawaan sistem tidak boleh diubah")
		return v
	}

	if existing.IsActive {
		hasUsers, _ := s.repo.HasUsers(ctx, id)
		if hasUsers {
			v := helper.NewValidationError()
			v.Add("general", "Role masih digunakan oleh pengguna dan tidak dapat dinonaktifkan")
			return v
		}
	}

	mapper.StatusReqToEntity(&req, existing)

	err = s.repo.UpdateRole(ctx, existing)
	if err != nil {
		return err
	}

	s.log(ctx, s.db, &activitylogdto.ActivityLogInput{
		Action:      "roles.update",
		EntityType:  "roles",
		EntityID:    &id,
		EntityLabel: existing.DisplayName,
		Description: fmt.Sprintf("Mengubah status aktif/nonaktif role: %s", existing.DisplayName),
		Metadata: map[string]any{
			"old_values": map[string]interface{}{"is_active": !req.IsActive},
			"new_values": map[string]interface{}{"is_active": req.IsActive},
		},
	})

	return nil
}

func (s *roleService) BulkDelete(ctx context.Context, ids []uint) error {
	for _, id := range ids {
		hasUsers, _ := s.repo.HasUsers(ctx, id)
		if hasUsers {
			v := helper.NewValidationError()
			v.Add("general", "Salah satu role masih digunakan oleh pengguna dan tidak dapat dihapus massal")
			return v
		}
	}

	err := s.repo.BulkDelete(ctx, ids)
	if err != nil {
		return err
	}

	s.log(ctx, s.db, &activitylogdto.ActivityLogInput{
		Action:      "roles.delete",
		EntityType:  "roles",
		Description: fmt.Sprintf("Menghapus massal role dengan ID: %v", ids),
	})

	return nil
}

func (s *roleService) BulkRestore(ctx context.Context, ids []uint) error {
	err := s.repo.BulkRestore(ctx, ids)
	if err != nil {
		return err
	}

	s.log(ctx, s.db, &activitylogdto.ActivityLogInput{
		Action:      "roles.restore",
		EntityType:  "roles",
		Description: fmt.Sprintf("Memulihkan massal role dengan ID: %v", ids),
	})

	return nil
}

func (s *roleService) GetDependencyInfo(ctx context.Context, id uint) (map[string]interface{}, error) {
	count, err := s.repo.CountUsers(ctx, id)
	if err != nil {
		return nil, err
	}

	message := roleDependencyMessage(count)

	return map[string]interface{}{
		"has_dependencies": message != "",
		"message":          message,
		"user_count":       count,
	}, nil
}

func (s *roleService) CheckUnique(ctx context.Context, field string, value string, excludeID uint) (bool, error) {
	switch field {
	case "name":
		return s.repo.ExistsByName(ctx, value, excludeID)
	case "display_name":
		return s.repo.ExistsByDisplayName(ctx, value, excludeID)
	default:
		v := helper.NewValidationError()
		v.Add("field", "field harus name atau display_name")
		return false, v
	}
}
