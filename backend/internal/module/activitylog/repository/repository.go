package repository

import (
	"context"

	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/activitylog/dto"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/activitylog/entity"
	"gorm.io/gorm"
)

type ActivityLogRepository interface {

	// Create menyimpan activity log baru.
	Create(
		ctx context.Context,
		db *gorm.DB,
		log *entity.ActivityLog,
	) error

	// List mengambil daftar activity log dengan filter dan pagination.
	List(
		ctx context.Context,
		req *dto.ActivityLogQueryReq,
	) ([]entity.ActivityLog, int64, error)

	// GetSummary mengambil data summary untuk dashboard audit.
	GetSummary(
		ctx context.Context,
		req *dto.ActivityLogQueryReq,
	) (dto.ActivityLogSummaryRes, error)

	// FindByID mengambil satu activity log berdasarkan id.
	FindByID(
		ctx context.Context,
		id uint64,
	) (*entity.ActivityLog, error)

	// EntityLogs mengambil history berdasarkan entity tertentu.
	EntityLogs(
		ctx context.Context,
		entityType string,
		entityID uint64,
	) ([]entity.ActivityLog, error)
}
