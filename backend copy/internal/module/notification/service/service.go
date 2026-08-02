package service

import (
	"context"

	"github.com/ahmadzakyarifin/schoolpay/backend/internal/module/notification/dto"
)

// NotificationService adalah antarmuka pengiriman notifikasi.
// Caller bertanggung jawab merender template sebelum memanggil method ini.
type NotificationService interface {
	// SendEmail mengirim notifikasi email: simpan ke DB (Pending) + enqueue ke worker.
	SendEmail(ctx context.Context, req dto.SendEmailRequest) error
	// SendWhatsApp mengirim notifikasi WhatsApp: simpan ke DB (Pending) + enqueue ke worker.
	SendWhatsApp(ctx context.Context, req dto.SendWhatsAppRequest) error
}
