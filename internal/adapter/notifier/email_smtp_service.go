package notifier

import (
	"context"
	"fmt"
	"github.com/risk-place-angola/backend-risk-place/internal/application/port"
	"github.com/risk-place-angola/backend-risk-place/internal/config"
	"net/smtp"
)

type SmtpEmailService struct {
	host     string
	port     string
	username string
	password string
	from     string
}

func NewSmtpEmailService(cfg config.Config) port.EmailService {
	return &SmtpEmailService{
		host:     cfg.EmailConfig.Host,
		port:     cfg.EmailConfig.Port,
		username: cfg.EmailConfig.User,
		password: cfg.EmailConfig.Pass,
		from:     cfg.EmailConfig.From,
	}
}

func (s *SmtpEmailService) SendEmail(ctx context.Context, to, subject, body string) error {
	auth := smtp.PlainAuth("", s.username, s.password, s.host)
	msg := buildPlainMessage(s.from, to, subject, body)

	addr := fmt.Sprintf("%s:%s", s.host, s.port)
	return smtp.SendMail(addr, auth, s.from, []string{to}, []byte(msg))
}

func (s *SmtpEmailService) SendHTMLEmail(ctx context.Context, to, subject, htmlBody string) error {
	auth := smtp.PlainAuth("", s.username, s.password, s.host)
	msg := buildHTMLMessage(s.from, to, subject, htmlBody)

	addr := fmt.Sprintf("%s:%s", s.host, s.port)
	return smtp.SendMail(addr, auth, s.from, []string{to}, []byte(msg))
}

func buildPlainMessage(from, to, subject, body string) string {
	return fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=\"utf-8\"\r\n\r\n%s",
		from, to, subject, body,
	)
}

func buildHTMLMessage(from, to, subject, htmlBody string) string {
	return fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=\"utf-8\"\r\n\r\n%s",
		from, to, subject, htmlBody,
	)
}
