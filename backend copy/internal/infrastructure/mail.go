package infrastructure

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ahmadzakyarifin/schoolpay/backend/config"
	"github.com/wneessen/go-mail"
)

type Mail struct {
	cfg *config.Config
}

type MailRequest struct {
	To      string
	Subject string
	Text    string
	HTML    string
}

func NewMail(cfg *config.Config) (*Mail, error) {
	return &Mail{
		cfg: cfg,
	}, nil
}

func (m *Mail) Send(ctx context.Context, req MailRequest) error {
	req.To = strings.TrimSpace(req.To)
	req.Subject = strings.TrimSpace(req.Subject)
	req.Text = strings.TrimSpace(req.Text)
	req.HTML = strings.TrimSpace(req.HTML)

	if req.To == "" {
		return errors.New("mail: recipient wajib diisi")
	}

	if req.Subject == "" {
		return errors.New("mail: subject wajib diisi")
	}

	if req.Text == "" && req.HTML == "" {
		return errors.New("mail: message wajib diisi")
	}

	smtpHost := m.cfg.SMTP.Host
	smtpPort := m.cfg.SMTP.Port
	smtpEmail := m.cfg.SMTP.Email
	smtpPass := m.cfg.SMTP.Pass
	fromName := m.cfg.SMTP.FromName

	client, err := mail.NewClient(
		smtpHost,
		mail.WithPort(smtpPort),
		mail.WithUsername(smtpEmail),
		mail.WithPassword(smtpPass),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithTimeout(time.Duration(m.cfg.SMTP.TimeoutSec)*time.Second),
	)
	if err != nil {
		return fmt.Errorf("smtp client dibuat: %w", err)
	}

	msg := mail.NewMsg()

	if err := msg.FromFormat(fromName, smtpEmail); err != nil {
		return fmt.Errorf("mail: pengirim: %w", err)
	}

	if err := msg.To(req.To); err != nil {
		return fmt.Errorf("mail: penerima: %w", err)
	}

	msg.Subject(req.Subject)

	switch {
	case req.Text != "" && req.HTML != "":
		msg.SetBodyString(mail.TypeTextPlain, req.Text)
		msg.AddAlternativeString(mail.TypeTextHTML, req.HTML)

	case req.HTML != "":
		msg.SetBodyString(mail.TypeTextHTML, req.HTML)

	default:
		msg.SetBodyString(mail.TypeTextPlain, req.Text)
	}

	if err := client.DialAndSendWithContext(ctx, msg); err != nil {
		return fmt.Errorf("mail: email dikirim: %w", err)
	}

	return nil
}
