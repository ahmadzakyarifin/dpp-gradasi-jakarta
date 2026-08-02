package service

import (
	"context"
	"fmt"

	jobs "github.com/ahmadzakyarifin/schoolpay/backend/internal/job"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/notification"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/notification/dto"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/notification/model"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/notification/repository"
	"gorm.io/gorm"
)

type notificationService struct {
	db        *gorm.DB
	repo      repository.NotificationRepo
	jobClient *jobs.Client
}

func NewNotificationService(db *gorm.DB, repo repository.NotificationRepo, jobClient *jobs.Client) NotificationService {
	return &notificationService{
		db:        db,
		repo:      repo,
		jobClient: jobClient,
	}
}

func (s *notificationService) SendEmail(ctx context.Context, req dto.SendEmailRequest) error {
	notif := &model.Notification{
		Channel:         notification.ChannelEmail,
		Destination:     req.To,
		Subject:         req.Subject,
		Message:         req.HTML,
		Status:          notification.StatusPending,
		RecipientUserID: req.UserID,
	}

	if err := s.repo.Create(ctx, notif); err != nil {
		return fmt.Errorf("notification: simpan rekor email: %w", err)
	}

	if err := s.jobClient.EnqueueEmail(jobs.EmailJob{
		NotificationID: &notif.ID,
		To:             req.To,
		Subject:        req.Subject,
		HTML:           req.HTML,
	}); err != nil {
		return fmt.Errorf("notification: enqueue email: %w", err)
	}

	return nil
}

func (s *notificationService) SendWhatsApp(ctx context.Context, req dto.SendWhatsAppRequest) error {
	notif := &model.Notification{
		Channel:         notification.ChannelWhatsApp,
		Destination:     req.To,
		Message:         req.Text,
		Status:          notification.StatusPending,
		RecipientUserID: req.UserID,
	}

	if err := s.repo.Create(ctx, notif); err != nil {
		return fmt.Errorf("notification: simpan rekor whatsapp: %w", err)
	}

	if err := s.jobClient.EnqueueWhatsApp(jobs.WhatsAppJob{
		NotificationID: &notif.ID,
		ChatID:         req.To,
		Text:           req.Text,
	}); err != nil {
		return fmt.Errorf("notification: enqueue whatsapp: %w", err)
	}

	return nil
}
