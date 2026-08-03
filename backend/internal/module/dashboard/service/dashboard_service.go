package service

import (
	"context"

	actdto "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/activitylog/dto"
	activitylogservice "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/activitylog/service"
	beritamodel "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/berita/model"
	dashdto "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/dashboard/dto"
	kegiatanmodel "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/kegiatan/model"
	kontakmodel "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/kontak/model"
	pengurusmodel "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/pengurus/model"
	usermodel "github.com/ahmadzakyarifin/dpp-gradasi/backend/internal/module/user/model"
	"gorm.io/gorm"
)

type DashboardService interface {
	GetSummary(ctx context.Context) (*dashdto.DashboardSummaryRes, error)
}

type dashboardService struct {
	db             *gorm.DB
	activityLogSvc activitylogservice.ActivityLogService
}

func NewDashboardService(db *gorm.DB, activityLogSvc activitylogservice.ActivityLogService) DashboardService {
	return &dashboardService{db: db, activityLogSvc: activityLogSvc}
}

func (s *dashboardService) GetSummary(ctx context.Context) (*dashdto.DashboardSummaryRes, error) {
	res := &dashdto.DashboardSummaryRes{}

	// Berita (tidak terhapus)
	var totalBerita int64
	if err := s.db.Model(&beritamodel.Berita{}).Where("deleted_at IS NULL").Count(&totalBerita).Error; err != nil {
		return nil, err
	}
	res.TotalBerita = totalBerita

	// Kegiatan
	var totalKegiatan int64
	if err := s.db.Model(&kegiatanmodel.Kegiatan{}).Where("deleted_at IS NULL").Count(&totalKegiatan).Error; err != nil {
		return nil, err
	}
	res.TotalKegiatan = totalKegiatan

	// Pengurus
	var totalPengurus int64
	if err := s.db.Model(&pengurusmodel.Pengurus{}).Where("deleted_at IS NULL").Count(&totalPengurus).Error; err != nil {
		return nil, err
	}
	res.TotalPengurus = totalPengurus

	// Kontak
	var totalKontak, unreadKontak int64
	if err := s.db.Model(&kontakmodel.PesanKontak{}).Where("deleted_at IS NULL").Count(&totalKontak).Error; err != nil {
		return nil, err
	}
	if err := s.db.Model(&kontakmodel.PesanKontak{}).Where("deleted_at IS NULL AND is_read = ?", false).Count(&unreadKontak).Error; err != nil {
		return nil, err
	}
	res.TotalKontak = totalKontak
	res.UnreadKontak = unreadKontak

	// Admin users (non-super, tidak terhapus)
	var totalAdmin, pendingAdmin int64
	if err := s.db.Model(&usermodel.UserModel{}).Where("role_id <> ? AND deleted_at IS NULL", 1).Count(&totalAdmin).Error; err != nil {
		return nil, err
	}
	if err := s.db.Model(&usermodel.UserModel{}).Where("role_id <> ? AND deleted_at IS NULL AND status = ?", 1, "inactive").Count(&pendingAdmin).Error; err != nil {
		return nil, err
	}
	res.TotalAdmin = totalAdmin
	res.PendingAdmin = pendingAdmin

	// Activity log summary
	summary, err := s.activityLogSvc.Summary(ctx, &actdto.ActivityLogQueryReq{})
	if err == nil {
		res.ActivityLogs = summary.TotalLogs
		res.HighRiskLogs = summary.HighRisk
		res.FailedLogin = summary.FailedLogin
		res.CMSActions = summary.FinanceAction
	}

	return res, nil
}
