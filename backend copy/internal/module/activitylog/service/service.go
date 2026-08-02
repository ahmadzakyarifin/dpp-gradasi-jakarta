package service

import (
	"context"

	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/activitylog/dto"
	"gorm.io/gorm"
)

type ActivityLogService interface {

	// Digunakan module lain untuk mencatat activity.
	Log(ctx context.Context, db *gorm.DB, input *dto.ActivityLogInput) error

	// API
	List(ctx context.Context, req *dto.ActivityLogQueryReq) (*dto.ActivityLogListRes, error)

	Detail(ctx context.Context, id uint64) (*dto.ActivityLogDetailRes, error)

	EntityLogs(ctx context.Context, entityType string, entityID uint64) ([]dto.ActivityLogItemRes, error)
}
