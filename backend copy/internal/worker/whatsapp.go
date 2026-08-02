package worker

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/hibiken/asynq"

	"github.com/ahmadzakyarifin/schoolpay/backend/internal/infrastructure"
	jobs "github.com/ahmadzakyarifin/schoolpay/backend/internal/job"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/notification"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/notification/repository"
)

type WhatsAppWorker struct {
	waha      *infrastructure.WAHA
	notifRepo repository.NotificationRepo
}

func NewWhatsAppWorker(waha *infrastructure.WAHA, notifRepo repository.NotificationRepo) *WhatsAppWorker {
	return &WhatsAppWorker{
		waha:      waha,
		notifRepo: notifRepo,
	}
}

func (w *WhatsAppWorker) ProcessTask(ctx context.Context, task *asynq.Task) error {
	if len(task.Payload()) == 0 {
		return errors.New("empty whatsapp payload")
	}

	var job jobs.WhatsAppJob
	if err := json.Unmarshal(task.Payload(), &job); err != nil {
		return err
	}

	err := w.waha.SendText(ctx, infrastructure.SendTextRequest{
		ChatID: job.ChatID,
		Text:   job.Text,
	})

	if job.NotificationID != nil {
		status := notification.StatusSent
		var errMsg string
		if err != nil {
			status = notification.StatusFailed
			errMsg = err.Error()
		}
		_ = w.notifRepo.UpdateStatus(ctx, *job.NotificationID, status, "", errMsg)
	}

	return err
}
