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
	s.db.Model(&beritamodel.Berita{}).Where("deleted_at IS NULL").Count(&totalBerita)
	res.TotalBerita = totalBerita

	// Kegiatan
	var totalKegiatan int64
	s.db.Model(&kegiatanmodel.Kegiatan{}).Where("deleted_at IS NULL").Count(&totalKegiatan)
	res.TotalKegiatan = totalKegiatan

	// Pengurus
	var totalPengurus int64
	s.db.Model(&pengurusmodel.Pengurus{}).Where("deleted_at IS NULL").Count(&totalPengurus)
	res.TotalPengurus = totalPengurus

	// Kontak
	var totalKontak, unreadKontak int64
	s.db.Model(&kontakmodel.PesanKontak{}).Where("deleted_at IS NULL").Count(&totalKontak)
	s.db.Model(&kontakmodel.PesanKontak{}).Where("deleted_at IS NULL AND is_read = ?", false).Count(&unreadKontak)
	res.TotalKontak = totalKontak
	res.UnreadKontak = unreadKontak

	// Admin users (non-super, tidak terhapus)
	var totalAdmin, pendingAdmin int64
	s.db.Model(&usermodel.User{}).Where("role_id <> ? AND deleted_at IS NULL", 1).Count(&totalAdmin)
	s.db.Model(&usermodel.User{}).Where("role_id <> ? AND deleted_at IS NULL AND status = ?", 1, "pending_activation").Count(&pendingAdmin)
	res.TotalAdmin = totalAdmin
	res.PendingAdmin = pendingAdmin

	// Activity log summary
	summary, err := s.activityLogSvc.Summary(ctx, &actdto.ActivityLogQueryReq{})
	if err == nil {
		res.ActivityLogs = summary.TotalLogs
		res.HighRiskLogs = summary.HighRisk
		res.FailedLogin = summary.FailedLogin
		res.CMSActions = summary.CMSAction
	}

	return res, nil
}
