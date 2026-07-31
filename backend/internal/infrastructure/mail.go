package infrastructure

import (
	"fmt"
	"log"
	"strconv"

	"github.com/ahmadzakyarifin/dpp-gradasi/backend/config"
	gomail "gopkg.in/gomail.v2"
)

type Mailer struct {
	fromName string
	fromAddr string
	dialer   *gomail.Dialer
}

func NewMailer(cfg *config.Config) *Mailer {
	d := gomail.NewDialer(
		cfg.SMTP.Host,
		cfg.SMTP.Port,
		cfg.SMTP.Email,
		cfg.SMTP.Pass,
	)

	log.Printf("mailer siap: %s:%d", cfg.SMTP.Host, cfg.SMTP.Port)

	return &Mailer{
		fromName: cfg.SMTP.FromName,
		fromAddr: cfg.SMTP.Email,
		dialer:   d,
	}
}

func (m *Mailer) Send(to, subject, body string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", m.fromAddr)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", body)

	if err := m.dialer.DialAndSend(msg); err != nil {
		return fmt.Errorf("gagal kirim email: %w", err)
	}
	return nil
}

func (m *Mailer) SendAsync(to, subject, body string) {
	go func() {
		if err := m.Send(to, subject, body); err != nil {
			log.Printf("gagal kirim email ke %s: %v", to, err)
		}
	}()
}

func (m *Mailer) DSN() string {
	return m.fromAddr + " (via " + m.dialer.Host + ":" + strconv.Itoa(m.dialer.Port) + ")"
}
