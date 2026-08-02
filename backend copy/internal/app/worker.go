package app

import (
	"context"
	"log"

	"github.com/ahmadzakyarifin/schoolpay/backend/config"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/infrastructure"
	jobs "github.com/ahmadzakyarifin/schoolpay/backend/internal/job"
	"github.com/ahmadzakyarifin/schoolpay/backend/internal/worker"
	"gorm.io/gorm"
)

// startBackgroundJobs menjalankan worker asynq di goroutine terpisah.
// Worker membaca job dari Redis (email/whatsapp) lalu mengirim via SMTP/WAHA.
func startBackgroundJobs(
	ctx context.Context,
	db *gorm.DB,
	appConfig *config.Config,
	jobClient *jobs.Client,
) {
	if jobClient == nil {
		log.Println("worker: job client nil, skip")
		return
	}

	// Inisialisasi sender infrastructure — baca dari env var saja.
	mail, err := infrastructure.NewMail(appConfig)
	if err != nil {
		log.Printf("worker: mail gagal diinisialisasi: %v (email worker nonaktif)", err)
	}
	waha := infrastructure.NewWAHA(appConfig)

	// Setup worker server & mux
	srv := jobs.NewServer(appConfig)
	mux := worker.NewMux(db, mail, waha)

	// Jalankan worker di background
	go func() {
		log.Println("worker: asynq server starting...")
		if err := srv.Run(mux); err != nil {
			log.Printf("worker: asynq server error: %v", err)
		}
	}()

	log.Println("worker: started")
}
