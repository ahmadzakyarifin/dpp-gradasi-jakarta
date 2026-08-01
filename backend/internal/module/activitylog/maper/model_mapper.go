package mapper

import (
	"fmt"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/activitylog/dto"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/activitylog/entity"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/activitylog/model"
)

// =========================================
// Request -> Entity
// =========================================

func InputToEntity(in *dto.ActivityLogInput) *entity.ActivityLog {
	if in == nil {
		return nil
	}

	metadata := in.Metadata
	if metadata == nil {
		metadata = map[string]any{}
	}

	return &entity.ActivityLog{
		ActorID:     in.ActorID,
		ActorName:   in.ActorName,
		ActorRole:   in.ActorRole,
		Action:      in.Action,
		EntityType:  in.EntityType,
		EntityID:    in.EntityID,
		EntityLabel: in.EntityLabel,
		Description: in.Description,
		IPAddress:   in.IPAddress,
		UserAgent:   in.UserAgent,
		Metadata:    metadata,
	}
}

// =========================================
// Model -> Entity
// =========================================

func ModelToEntity(m *model.ActivityLog) *entity.ActivityLog {
	if m == nil {
		return nil
	}

	metadata := m.Metadata
	if metadata == nil {
		metadata = map[string]any{}
	}

	return &entity.ActivityLog{
		ID:          m.ID,
		ActorID:     m.ActorID,
		ActorName:   m.ActorName,
		ActorRole:   m.ActorRole,
		Action:      m.Action,
		EntityType:  m.EntityType,
		EntityID:    m.EntityID,
		EntityLabel: m.EntityLabel,
		RiskLevel:   m.RiskLevel,
		Description: m.Description,
		IPAddress:   m.IPAddress,
		UserAgent:   m.UserAgent,
		Metadata:    map[string]any(metadata),
		CreatedAt:   m.CreatedAt,
	}
}

// =========================================
// Entity -> Model
// =========================================

func EntityToModel(e *entity.ActivityLog) *model.ActivityLog {
	if e == nil {
		return nil
	}

	m := &model.ActivityLog{
		ID:          e.ID,
		ActorID:     e.ActorID,
		ActorName:   e.ActorName,
		ActorRole:   e.ActorRole,
		Action:      e.Action,
		EntityType:  e.EntityType,
		EntityID:    e.EntityID,
		EntityLabel: e.EntityLabel,
		RiskLevel:   e.RiskLevel,
		Description: e.Description,
		IPAddress:   e.IPAddress,
		UserAgent:   e.UserAgent,
		Metadata:    model.JSONMap(e.Metadata),
		CreatedAt:   e.CreatedAt,
	}
	return m
}

// =========================================
// Entity -> Response
// =========================================

func EntityToResponse(e *entity.ActivityLog) dto.ActivityLogItemRes {
	return dto.ActivityLogItemRes{
		ID:          fmt.Sprintf("%d", e.ID),
		Time:        e.CreatedAt.Format("2006-01-02 15:04:05"),
		Actor:       e.ActorName,
		Role:        e.ActorRole,
		Action:      e.Action,
		Entity:      e.EntityType,
		EntityLabel: e.EntityLabel,
		Description: e.Description,
		IP:          e.IPAddress,
		Device:      e.UserAgent,
		Risk:        e.RiskLevel,
		Metadata:    e.Metadata,
	}
}

// ============================================
//  Entity -> Detail Response
// ============================================

func EntityToDetailResponse(e *entity.ActivityLog) dto.ActivityLogDetailRes {
	return dto.ActivityLogDetailRes{
		ID:          fmt.Sprintf("%d", e.ID),
		Time:        e.CreatedAt.Format("2006-01-02 15:04:05"),
		Actor:       e.ActorName,
		Role:        e.ActorRole,
		Action:      e.Action,
		Entity:      e.EntityType,
		EntityLabel: e.EntityLabel,
		Description: e.Description,
		IP:          e.IPAddress,
		Device:      e.UserAgent,
		Risk:        e.RiskLevel,
		Metadata:    e.Metadata,
	}
}
