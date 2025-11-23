package port

import "context"

type EmailService interface {
	SendEmail(ctx context.Context, to string, subject string, body string) error
	SendHTMLEmail(ctx context.Context, to string, subject string, htmlBody string) error
}
