package service

import (
	"context"
	"math"

	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/activitylog/dto"
	mapper "github.com/ahmadzakyarifin/schoolpay/backend/internal/module/activitylog/maper"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/activitylog/repository"
	"gorm.io/gorm"
)

type activityLogService struct {
	repo repository.ActivityLogRepository
}

func NewActivityLogService(
	repo repository.ActivityLogRepository,
) ActivityLogService {
	return &activityLogService{
		repo: repo,
	}
}

func (s *activityLogService) Log(
	ctx context.Context,
	db *gorm.DB,
	input *dto.ActivityLogInput,
) error {

	entity := mapper.InputToEntity(input)

	entity.RiskLevel = determineRisk(entity.Action)

	return s.repo.Create(ctx, db, entity)
}

func (s *activityLogService) List(
	ctx context.Context,
	req *dto.ActivityLogQueryReq,
) (*dto.ActivityLogListRes, error) {

	if req.Page <= 0 {
		req.Page = 1
	}

	if req.Limit <= 0 {
		req.Limit = 10
	}

	summary, err := s.repo.GetSummary(ctx, req)
	if err != nil {
		return nil, err
	}

	logs, total, err := s.repo.List(ctx, req)
	if err != nil {
		return nil, err
	}

	items := make([]dto.ActivityLogItemRes, len(logs))

	for i := range logs {
		items[i] = mapper.EntityToResponse(&logs[i])
	}

	totalPages := int(math.Ceil(float64(total) / float64(req.Limit)))

	if totalPages == 0 {
		totalPages = 1
	}

	return &dto.ActivityLogListRes{
		Summary: summary,
		Pagination: dto.ActivityLogPaginationRes{
			Page:       req.Page,
			Limit:      req.Limit,
			Total:      total,
			TotalPages: totalPages,
		},
		Items: items,
	}, nil
}

func (s *activityLogService) Detail(
	ctx context.Context,
	id uint64,
) (*dto.ActivityLogDetailRes, error) {

	log, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	res := mapper.EntityToDetailResponse(log)

	return &res, nil
}

func (s *activityLogService) EntityLogs(
	ctx context.Context,
	entityType string,
	entityID uint64,
) ([]dto.ActivityLogItemRes, error) {

	logs, err := s.repo.EntityLogs(ctx, entityType, entityID)
	if err != nil {
		return nil, err
	}

	res := make([]dto.ActivityLogItemRes, len(logs))
	for i := range logs {
		res[i] = mapper.EntityToResponse(&logs[i])
	}

	return res, nil
}
