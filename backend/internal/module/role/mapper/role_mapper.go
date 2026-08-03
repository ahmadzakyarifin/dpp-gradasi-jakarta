package mapper

import (
	"time"

	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/role/dto"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/role/entity"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/role/model"
	"gorm.io/gorm"
)

func toGormDeletedAt(t *time.Time) gorm.DeletedAt {
	if t == nil {
		return gorm.DeletedAt{}
	}
	return gorm.DeletedAt{Time: *t, Valid: true}
}

//
// =========================================
// Request -> Entity
// =========================================
//

func CreateReqToEntity(req *dto.RoleCreateReq) *entity.Role {
	if req == nil {
		return nil
	}
	return &entity.Role{
		Name:        req.Name,
		DisplayName: req.DisplayName,
		IsActive:    req.IsActive,
	}
}

func UpdateReqToEntity(req *dto.RoleUpdateReq, role *entity.Role) {
	if req == nil || role == nil {
		return
	}
	role.Name = req.Name
	role.DisplayName = req.DisplayName
	role.IsActive = req.IsActive
}

func StatusReqToEntity(req *dto.RoleStatusReq, role *entity.Role) {
	if req == nil || role == nil {
		return
	}
	role.IsActive = req.IsActive
}

//
// =========================================
// Entity -> Model
// =========================================
//

func EntityToModel(e *entity.Role) *model.RoleModel {
	if e == nil {
		return nil
	}

	return &model.RoleModel{
		ID:          e.ID,
		Name:        e.Name,
		DisplayName: e.DisplayName,
		IsActive:    e.IsActive,
		IsSystem:    e.IsSystem,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
		DeletedAt:   toGormDeletedAt(e.DeletedAt),
	}
}

//
// =========================================
// Model -> Entity
// =========================================
//

func ModelToEntity(m *model.RoleModel) *entity.Role {
	if m == nil {
		return nil
	}

	var deletedAt *time.Time
	if m.DeletedAt.Valid {
		deletedAt = &m.DeletedAt.Time
	}

	return &entity.Role{
		ID:          m.ID,
		Name:        m.Name,
		DisplayName: m.DisplayName,
		IsActive:    m.IsActive,
		IsSystem:    m.IsSystem,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		DeletedAt:   deletedAt,
	}
}

func ModelListToEntity(models []model.RoleModel) []entity.Role {
	list := make([]entity.Role, 0, len(models))

	for i := range models {
		list = append(list, *ModelToEntity(&models[i]))
	}

	return list
}

//
// =========================================
// Entity -> Response DTO
// =========================================
//

func EntityToResponse(e *entity.Role) dto.RoleRes {
	if e == nil {
		return dto.RoleRes{}
	}

	return dto.RoleRes{
		ID:          e.ID,
		Name:        e.Name,
		DisplayName: e.DisplayName,
		IsSystem:    e.IsSystem,
		IsActive:    e.IsActive,
		UserCount:   e.UserCount,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
		DeletedAt:   e.DeletedAt,
	}
}

func EntityListToResponse(entities []entity.Role) []dto.RoleRes {
	list := make([]dto.RoleRes, 0, len(entities))

	for i := range entities {
		list = append(list, EntityToResponse(&entities[i]))
	}

	return list
}
