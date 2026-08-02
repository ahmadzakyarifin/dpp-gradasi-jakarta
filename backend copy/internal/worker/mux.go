package worker

import (
	"github.com/hibiken/asynq"

	"github.com/ahmadzakyarifin/schoolpay/backend/internal/infrastructure"
	jobs "github.com/ahmadzakyarifin/schoolpay/backend/internal/job"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/notification/repository"
	"gorm.io/gorm"
)

func NewMux(
	db *gorm.DB,
	mail *infrastructure.Mail,
	waha *infrastructure.WAHA,
) *asynq.ServeMux {

	mux := asynq.NewServeMux()

	notifRepo := repository.NewNotificationRepo(db)

	emailWorker := NewEmailWorker(mail, notifRepo)
	whatsAppWorker := NewWhatsAppWorker(waha, notifRepo)

	mux.HandleFunc(
		jobs.TaskEmail,
		emailWorker.ProcessTask,
	)

	mux.HandleFunc(
		jobs.TaskWhatsApp,
		whatsAppWorker.ProcessTask,
	)

	return mux
}
