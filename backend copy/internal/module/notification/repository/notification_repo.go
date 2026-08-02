package repository

import (
	"context"
	"time"

	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/notification/model"
	"gorm.io/gorm"
)

// NotificationRepo mendefinisikan operasi database untuk notifikasi.
type NotificationRepo interface {
	Create(ctx context.Context, notif *model.Notification) error
	UpdateStatus(ctx context.Context, id uint, status string, providerMsgID string, errorMsg string) error
	FindPaginated(ctx context.Context, channel string, status string, search string, page int, pageSize int) ([]model.NotificationRow, int64, error)
	GetDB() *gorm.DB
}

type notificationRepo struct {
	db *gorm.DB
}

func NewNotificationRepo(db *gorm.DB) NotificationRepo {
	return &notificationRepo{db: db}
}

func (r *notificationRepo) GetDB() *gorm.DB {
	return r.db
}

func (r *notificationRepo) Create(ctx context.Context, notif *model.Notification) error {
	return r.db.WithContext(ctx).Create(notif).Error
}

func (r *notificationRepo) UpdateStatus(ctx context.Context, id uint, status string, providerMsgID string, errorMsg string) error {
	updates := map[string]any{
		"status": status,
	}
	switch status {
	case "sent":
		now := time.Now()
		updates["sent_at"] = &now
		if providerMsgID != "" {
			updates["provider_message_id"] = providerMsgID
		}
	case "failed":
		if errorMsg != "" {
			updates["error_message"] = errorMsg
		}
	}

	return r.db.WithContext(ctx).
		Model(&model.Notification{}).
		Where("id = ?", id).
		Updates(updates).Error
}

func (r *notificationRepo) FindPaginated(ctx context.Context, channel string, status string, search string, page int, pageSize int) ([]model.NotificationRow, int64, error) {
	db := r.db.WithContext(ctx).
		Table("notifications").
		Select("notifications.*, users.name as recipient_name").
		Joins("LEFT JOIN users ON users.id = notifications.recipient_user_id")

	if channel != "" && channel != "all" {
		db = db.Where("channel = ?", channel)
	}
	if status != "" && status != "all" {
		db = db.Where("status = ?", status)
	}
	if search != "" {
		q := "%" + search + "%"
		db = db.Where("destination LIKE ? OR subject LIKE ? OR message LIKE ? OR error_message LIKE ?", q, q, q, q)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var rows []model.NotificationRow
	if err := db.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&rows).Error; err != nil {
		return nil, 0, err
	}

	return rows, total, nil
}
