package service

import (
	"context"
	"crypto/rand"
	"fmt"
	"log/slog"
	"math/big"
	"time"

	"github.com/google/uuid"
	"github.com/risk-place-angola/backend-risk-place/internal/application/port"
	domainService "github.com/risk-place-angola/backend-risk-place/internal/domain/service"
)

const (
	verificationCodeTTL      = 10 * time.Minute
	verificationCodeMaxValue = 1000000
	verificationCodeFormat   = "%06d"
)

type verificationServiceImpl struct {
	cache        port.Cache
	smsNotifier  port.NotifierSMSService
	emailService port.EmailService
	frontendURL  string
}

func NewVerificationService(
	cache port.Cache,
	smsNotifier port.NotifierSMSService,
	emailService port.EmailService,
	frontendURL string,
) domainService.VerificationService {
	return &verificationServiceImpl{
		cache:        cache,
		smsNotifier:  smsNotifier,
		emailService: emailService,
		frontendURL:  frontendURL,
	}
}

func (s *verificationServiceImpl) SendCode(ctx context.Context, userID uuid.UUID, phone, email string) error {
	code := s.generateCode()
	key := fmt.Sprintf("verification:%s", userID.String())

	if err := s.cache.Set(ctx, key, code, verificationCodeTTL); err != nil {
		return fmt.Errorf("failed to store code: %w", err)
	}

	if err := s.sendViaSMS(ctx, phone, code); err != nil {
		slog.Warn("SMS send failed, falling back to email", "user_id", userID, "error", err)
		if emailErr := s.sendViaEmail(ctx, email, code); emailErr != nil {
			return fmt.Errorf("both SMS and email failed: SMS=%w, Email=%w", err, emailErr)
		}
		return nil
	}

	return nil
}

func (s *verificationServiceImpl) VerifyCode(ctx context.Context, userID uuid.UUID, code string) (bool, error) {
	key := fmt.Sprintf("verification:%s", userID.String())

	storedCode, err := s.cache.Get(ctx, key)
	if err != nil {
		return false, fmt.Errorf("code not found or expired")
	}

	if storedCode != code {
		return false, nil
	}

	if err := s.cache.Delete(ctx, key); err != nil {
		slog.Warn("Failed to delete verification code", "user_id", userID, "error", err)
	}

	return true, nil
}

func (s *verificationServiceImpl) ResendCode(ctx context.Context, userID uuid.UUID, phone, email string) error {
	key := fmt.Sprintf("verification:%s", userID.String())

	if _, err := s.cache.Get(ctx, key); err == nil {
		return fmt.Errorf("code already sent, please wait")
	}

	return s.SendCode(ctx, userID, phone, email)
}

func (s *verificationServiceImpl) sendViaSMS(ctx context.Context, phone, code string) error {
	message := fmt.Sprintf("Seu código de verificação Risk Place: %s. Válido por 10 minutos.", code)
	return s.smsNotifier.NotifySMS(ctx, phone, message)
}

func (s *verificationServiceImpl) sendViaEmail(ctx context.Context, email, code string) error {
	htmlBody := fmt.Sprintf(`
		<h2>Verificação de Conta - Risk Place Angola</h2>
		<p>Seu código de verificação: <strong>%s</strong></p>
		<p>O código expira em 10 minutos.</p>
	`, code)

	return s.emailService.SendHTMLEmail(ctx, email, "Código de Verificação", htmlBody)
}

func (s *verificationServiceImpl) generateCode() string {
	code, _ := rand.Int(rand.Reader, big.NewInt(verificationCodeMaxValue))
	return fmt.Sprintf(verificationCodeFormat, code.Int64())
}
