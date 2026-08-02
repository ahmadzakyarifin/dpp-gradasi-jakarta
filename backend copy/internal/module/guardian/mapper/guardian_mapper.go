package mapper

import (
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/guardian/entity"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/guardian/model"
	"time"
)

func ModelToEntity(m *model.GuardianModel) *entity.Guardian {
	var deletedAt *time.Time
	if m.DeletedAt.Valid {
		deletedAt = &m.DeletedAt.Time
	}
	return &entity.Guardian{
		ID:          m.ID,
		UserID:      m.UserID,
		Name:        m.Name,
		Phone:       m.Phone,
		Email:       m.Email,
		NIK:         m.NIK,
		Education:   m.Education,
		Occupation:  m.Occupation,
		IncomeRange: m.IncomeRange,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
		DeletedAt:   deletedAt,
	}
}

func EntityToModel(e *entity.Guardian) *model.GuardianModel {
	return &model.GuardianModel{
		ID:          e.ID,
		UserID:      e.UserID,
		Name:        e.Name,
		Phone:       e.Phone,
		Email:       e.Email,
		NIK:         e.NIK,
		Education:   e.Education,
		Occupation:  e.Occupation,
		IncomeRange: e.IncomeRange,
	}
}
