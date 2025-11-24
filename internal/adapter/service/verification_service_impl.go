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
	domainErrors "github.com/risk-place-angola/backend-risk-place/internal/domain/errors"
	domainrepository "github.com/risk-place-angola/backend-risk-place/internal/domain/repository"
	domainService "github.com/risk-place-angola/backend-risk-place/internal/domain/service"
)

const (
	verificationCodeTTL       = 10 * time.Minute
	verificationCodeMaxValue  = 1000000
	verificationCodeFormat    = "%06d"
	maxVerificationAttempts   = 5
	attemptLockoutDuration    = 15 * time.Minute
	resendCooldown            = 60 * time.Second
)

type verificationServiceImpl struct {
	cache              port.Cache
	smsNotifier        port.NotifierSMSService
	emailService       port.EmailService
	frontendURL        string
	translationService *TranslationService
	userRepo           domainrepository.UserRepository
}

func NewVerificationService(
	cache port.Cache,
	smsNotifier port.NotifierSMSService,
	emailService port.EmailService,
	frontendURL string,
	translationService *TranslationService,
	userRepo domainrepository.UserRepository,
) domainService.VerificationService {
	return &verificationServiceImpl{
		cache:              cache,
		smsNotifier:        smsNotifier,
		emailService:       emailService,
		frontendURL:        frontendURL,
		translationService: translationService,
		userRepo:           userRepo,
	}
}

func (s *verificationServiceImpl) SendCode(ctx context.Context, userID uuid.UUID, phone, email string) error {
	return s.sendCodeWithType(ctx, userID, phone, email, false)
}

func (s *verificationServiceImpl) SendPasswordResetCode(ctx context.Context, userID uuid.UUID, phone, email string) error {
	return s.sendCodeWithType(ctx, userID, phone, email, true)
}

func (s *verificationServiceImpl) sendCodeWithType(ctx context.Context, userID uuid.UUID, phone, email string, isPasswordReset bool) error {
	code := s.generateCode()
	key := fmt.Sprintf("verification:%s", userID.String())

	if err := s.cache.Set(ctx, key, code, verificationCodeTTL); err != nil {
		return fmt.Errorf("failed to store code: %w", err)
	}

	lang := s.getUserLanguage(ctx, userID)

	if err := s.sendViaSMS(ctx, phone, code, lang, isPasswordReset); err != nil {
		slog.Warn("SMS send failed, falling back to email", "user_id", userID, "error", err)
		if emailErr := s.sendViaEmail(ctx, email, code, lang, isPasswordReset); emailErr != nil {
			return fmt.Errorf("both SMS and email failed: SMS=%w, Email=%w", err, emailErr)
		}
		return domainErrors.ErrSentViaEmail
	}

	return nil
}

func (s *verificationServiceImpl) VerifyCode(ctx context.Context, userID uuid.UUID, code string) (bool, error) {
	key := fmt.Sprintf("verification:%s", userID.String())
	attemptsKey := fmt.Sprintf("verification:attempts:%s", userID.String())
	lockoutKey := fmt.Sprintf("verification:lockout:%s", userID.String())

	if _, err := s.cache.Get(ctx, lockoutKey); err == nil {
		return false, domainErrors.ErrVerificationLocked
	}

	storedCode, err := s.cache.Get(ctx, key)
	if err != nil {
		return false, fmt.Errorf("code not found or expired")
	}

	if storedCode != code {
		if err := s.incrementAttempts(ctx, attemptsKey, lockoutKey, userID); err != nil {
			slog.Error("Failed to increment verification attempts", "user_id", userID, "error", err)
		}
		return false, nil
	}

	if err := s.cache.Delete(ctx, key); err != nil {
		slog.Warn("Failed to delete verification code", "user_id", userID, "error", err)
	}
	if err := s.cache.Delete(ctx, attemptsKey); err != nil {
		slog.Warn("Failed to delete attempts counter", "user_id", userID, "error", err)
	}

	return true, nil
}

func (s *verificationServiceImpl) ResendCode(ctx context.Context, userID uuid.UUID, phone, email string) error {
	return s.resendCodeWithType(ctx, userID, phone, email, false)
}

func (s *verificationServiceImpl) ResendPasswordResetCode(ctx context.Context, userID uuid.UUID, phone, email string) error {
	return s.resendCodeWithType(ctx, userID, phone, email, true)
}

func (s *verificationServiceImpl) resendCodeWithType(ctx context.Context, userID uuid.UUID, phone, email string, isPasswordReset bool) error {
	cooldownKey := fmt.Sprintf("verification:resend:%s", userID.String())
	lockoutKey := fmt.Sprintf("verification:lockout:%s", userID.String())

	if _, err := s.cache.Get(ctx, lockoutKey); err == nil {
		return domainErrors.ErrVerificationLocked
	}

	if _, err := s.cache.Get(ctx, cooldownKey); err == nil {
		return domainErrors.ErrVerificationCooldown
	}

	if err := s.cache.Set(ctx, cooldownKey, "1", resendCooldown); err != nil {
		slog.Warn("Failed to set resend cooldown", "user_id", userID, "error", err)
	}

	return s.sendCodeWithType(ctx, userID, phone, email, isPasswordReset)
}

func (s *verificationServiceImpl) sendViaSMS(ctx context.Context, phone, code string, lang Language, isPasswordReset bool) error {
	msgKey := "verification_code_sms"
	if isPasswordReset {
		msgKey = "password_reset_sms"
	}
	msg := s.translationService.GetMessage(msgKey, lang, "")
	message := fmt.Sprintf("%s: %s. %s", msg.Title, code, msg.Body)
	return s.smsNotifier.NotifySMS(ctx, phone, message)
}

func (s *verificationServiceImpl) sendViaEmail(ctx context.Context, email, code string, lang Language, isPasswordReset bool) error {
	msgKey := "verification_code_email"
	if isPasswordReset {
		msgKey = "password_reset_email"
	}
	msg := s.translationService.GetMessage(msgKey, lang, "")
	htmlBody := fmt.Sprintf(`
		<h2>%s</h2>
		<p>%s: <strong>%s</strong></p>
	`, msg.Title, msg.Body, code)

	return s.emailService.SendHTMLEmail(ctx, email, msg.Title, htmlBody)
}

func (s *verificationServiceImpl) incrementAttempts(ctx context.Context, attemptsKey, lockoutKey string, userID uuid.UUID) error {
	attemptsStr, err := s.cache.Get(ctx, attemptsKey)
	attempts := 1
	if err == nil {
		if _, parseErr := fmt.Sscanf(attemptsStr, "%d", &attempts); parseErr == nil {
			attempts++
		}
	}

	if attempts >= maxVerificationAttempts {
		if err := s.cache.Set(ctx, lockoutKey, "1", attemptLockoutDuration); err != nil {
			return err
		}
		slog.Warn("Verification attempts exceeded, user locked out", "user_id", userID, "attempts", attempts)
		return nil
	}

	return s.cache.Set(ctx, attemptsKey, fmt.Sprintf("%d", attempts), verificationCodeTTL)
}

func (s *verificationServiceImpl) getUserLanguage(ctx context.Context, userID uuid.UUID) Language {
	language, _, err := s.userRepo.GetUserLanguageAndPhone(ctx, userID)
	if err != nil {
		slog.Warn("Failed to get user language, using default", "user_id", userID, "error", err)
		return LanguagePortuguese
	}
	return s.translationService.ParseLanguage(language)
}

func (s *verificationServiceImpl) generateCode() string {
	code, _ := rand.Int(rand.Reader, big.NewInt(verificationCodeMaxValue))
	return fmt.Sprintf(verificationCodeFormat, code.Int64())
}
