package repository

import (
	"context"

	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/activitylog/dto"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/activitylog/entity"
	mapper "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/activitylog/maper"
	"github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/activitylog/model"
	"gorm.io/gorm"
)

type activityLogRepository struct {
	db *gorm.DB
}

func NewActivityLogRepository(
	db *gorm.DB,
) ActivityLogRepository {

	return &activityLogRepository{
		db: db,
	}
}

// Create menyimpan activity log.
func (r *activityLogRepository) Create(
	ctx context.Context,
	db *gorm.DB,
	log *entity.ActivityLog,
) error {

	if db == nil {
		db = r.db
	}

	modelLog := mapper.EntityToModel(log)

	return db.WithContext(ctx).Create(modelLog).Error
}

// List mengambil data activity log.
func (r *activityLogRepository) List(
	ctx context.Context,
	req *dto.ActivityLogQueryReq,
) ([]entity.ActivityLog, int64, error) {

	var models []model.ActivityLog

	q := r.buildQuery(req)

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	page := req.Page
	if page <= 0 {
		page = 1
	}

	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}

	q = q.
		Order("created_at DESC").
		Order("id DESC").
		Limit(limit).
		Offset((page - 1) * limit)

	if err := q.Find(&models).Error; err != nil {
		return nil, 0, err
	}

	result := make(
		[]entity.ActivityLog,
		len(models),
	)

	for i := range models {
		result[i] = *mapper.ModelToEntity(&models[i])
	}

	return result, int64(total), nil
}

// buildQuery membuat query dinamis.
func (r *activityLogRepository) buildQuery(
	req *dto.ActivityLogQueryReq,
) *gorm.DB {

	q := r.db.
		WithContext(context.Background()).
		Model(&model.ActivityLog{})

	if req.Action != "" &&
		req.Action != "all" {

		q = q.Where(
			"action = ?",
			req.Action,
		)
	}

	if req.Entity != "" &&
		req.Entity != "all" {

		q = q.Where(
			"entity_type = ?",
			req.Entity,
		)
	}

	if req.Role != "" &&
		req.Role != "all" {

		q = q.Where(
			"actor_role = ?",
			req.Role,
		)
	}

	if req.Risk != "" &&
		req.Risk != "all" {

		q = q.Where(
			"risk_level = ?",
			req.Risk,
		)
	}

	if req.Search != "" {

		search := "%" + req.Search + "%"

		q = q.Where(
			"(actor_name LIKE ? OR actor_role LIKE ? OR action LIKE ? OR entity_type LIKE ? OR entity_label LIKE ?)",
			search, search, search, search, search,
		)
	}

	return q
}

// GetSummary mengambil statistik dashboard.
func (r *activityLogRepository) GetSummary(
	ctx context.Context,
	req *dto.ActivityLogQueryReq,
) (dto.ActivityLogSummaryRes, error) {

	var summary dto.ActivityLogSummaryRes

	total, err := r.countQuery(ctx, r.buildQuery(req))
	if err != nil {
		return summary, err
	}

	summary.TotalLogs = total

	highRisk, err := r.countQuery(ctx, r.buildQuery(req).
		Where("risk_level = ?", "high"))
	if err != nil {
		return summary, err
	}

	summary.HighRisk = highRisk

	failedLogin, err := r.countQuery(ctx, r.buildQuery(req).
		Where("action = ?", "auth.failed_login"))
	if err != nil {
		return summary, err
	}

	summary.FailedLogin = failedLogin

	finance, err := r.countQuery(ctx, r.buildQuery(req).
		Where("entity_type IN ?", []string{
			"payment",
			"invoice",
			"billing",
			"payment_gateway",
		}))
	if err != nil {
		return summary, err
	}

	summary.FinanceAction = finance

	return summary, nil
}

func (r *activityLogRepository) countQuery(
	ctx context.Context,
	q *gorm.DB,
) (int64, error) {

	var count int64
	if err := q.WithContext(ctx).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// FindByID mengambil detail activity log berdasarkan id.
func (r *activityLogRepository) FindByID(
	ctx context.Context,
	id uint64,
) (*entity.ActivityLog, error) {

	var m model.ActivityLog

	err := r.db.
		WithContext(ctx).
		Where("id = ?", id).
		First(&m).Error

	if err != nil {
		return nil, err
	}

	return mapper.ModelToEntity(&m), nil
}

// EntityLogs mengambil history sebuah entity.
func (r *activityLogRepository) EntityLogs(
	ctx context.Context,
	entityType string,
	entityID uint64,
) ([]entity.ActivityLog, error) {

	var models []model.ActivityLog

	err := r.db.
		WithContext(ctx).
		Where(
			"entity_type = ? AND entity_id = ?",
			entityType,
			entityID,
		).
		Order("created_at DESC").
		Order("id DESC").
		Find(&models).Error

	if err != nil {
		return nil, err
	}

	result := make(
		[]entity.ActivityLog,
		len(models),
	)

	for i := range models {
		result[i] = *mapper.ModelToEntity(&models[i])
	}

	return result, nil
}
