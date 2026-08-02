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

type EmailWorker struct {
	mail      *infrastructure.Mail
	notifRepo repository.NotificationRepo
}

func NewEmailWorker(mail *infrastructure.Mail, notifRepo repository.NotificationRepo) *EmailWorker {
	return &EmailWorker{
		mail:      mail,
		notifRepo: notifRepo,
	}
}

func (w *EmailWorker) ProcessTask(ctx context.Context, task *asynq.Task) error {
	if len(task.Payload()) == 0 {
		return errors.New("empty email payload")
	}

	var job jobs.EmailJob
	if err := json.Unmarshal(task.Payload(), &job); err != nil {
		return err
	}

	err := w.mail.Send(ctx, infrastructure.MailRequest{
		To:      job.To,
		Subject: job.Subject,
		HTML:    job.HTML,
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
