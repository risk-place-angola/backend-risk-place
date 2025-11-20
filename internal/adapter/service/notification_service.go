package service

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/risk-place-angola/backend-risk-place/internal/application/port"
	domainrepository "github.com/risk-place-angola/backend-risk-place/internal/domain/repository"
)

type NotificationService struct {
	translationService   *TranslationService
	pushService          port.NotifierPushService
	smsService           port.NotifierSMSService
	userRepo             domainrepository.UserRepository
	anonymousSessionRepo domainrepository.AnonymousSessionRepository
}

func NewNotificationService(
	translationService *TranslationService,
	pushService port.NotifierPushService,
	smsService port.NotifierSMSService,
	userRepo domainrepository.UserRepository,
	anonymousSessionRepo domainrepository.AnonymousSessionRepository,
) *NotificationService {
	return &NotificationService{
		translationService:   translationService,
		pushService:          pushService,
		smsService:           smsService,
		userRepo:             userRepo,
		anonymousSessionRepo: anonymousSessionRepo,
	}
}

func (s *NotificationService) SendNotificationWithFallback(ctx context.Context, userID, deviceID, language, riskType, eventKey string, data map[string]string) error {
	lang := s.translationService.ParseLanguage(language)
	msg := s.translationService.GetMessage(eventKey, lang, riskType)

	var pushEnabled, smsEnabled bool
	var phone, fcmToken string
	var err error

	if userID != "" {
		uid, parseErr := uuid.Parse(userID)
		if parseErr != nil {
			slog.Error("invalid user ID", "error", parseErr)
			return parseErr
		}

		pushEnabled, smsEnabled, err = s.userRepo.GetNotificationPreferences(ctx, uid)
		if err != nil {
			slog.Error("failed to get user notification preferences", "error", err)
			pushEnabled = true
			smsEnabled = false
		}

		language, phone, err = s.userRepo.GetUserLanguageAndPhone(ctx, uid)
		if err != nil {
			slog.Error("failed to get user language and phone", "error", err)
		}

		user, err := s.userRepo.FindByID(ctx, uid)
		if err == nil && user != nil {
			fcmToken = user.DeviceToken
		}
	} else if deviceID != "" {
		pushEnabled, smsEnabled, err = s.anonymousSessionRepo.GetNotificationPreferences(ctx, deviceID)
		if err != nil {
			slog.Error("failed to get anonymous notification preferences", "error", err)
			pushEnabled = true
			smsEnabled = false
		}

		session, err := s.anonymousSessionRepo.FindByDeviceID(ctx, deviceID)
		if err == nil && session != nil {
			fcmToken = session.DeviceFCMToken
		}
	}

	if pushEnabled && fcmToken != "" {
		err = s.pushService.NotifyPush(ctx, fcmToken, msg.Title, msg.Body, data)
		if err == nil {
			return nil
		}
		slog.Error("push notification failed, falling back to SMS", "error", err)
	}

	if smsEnabled && phone != "" {
		err = s.smsService.NotifySMS(ctx, phone, msg.Body)
		if err != nil {
			slog.Error("SMS notification failed", "error", err)
			return err
		}
	}

	return nil
}

func (s *NotificationService) SendNotificationToMultiple(ctx context.Context, userIDs []uuid.UUID, deviceTokens []string, language, riskType, eventKey string, data map[string]string) error {
	lang := s.translationService.ParseLanguage(language)
	msg := s.translationService.GetMessage(eventKey, lang, riskType)

	var tokens []string
	var phones []string

	for _, uid := range userIDs {
		pushEnabled, smsEnabled, err := s.userRepo.GetNotificationPreferences(ctx, uid)
		if err != nil {
			slog.Error("failed to get notification preferences", "user_id", uid, "error", err)
			continue
		}

		if pushEnabled {
			user, err := s.userRepo.FindByID(ctx, uid)
			if err == nil && user != nil && user.DeviceToken != "" {
				tokens = append(tokens, user.DeviceToken)
			}
		}

		if smsEnabled {
			_, phone, err := s.userRepo.GetUserLanguageAndPhone(ctx, uid)
			if err == nil && phone != "" {
				phones = append(phones, phone)
			}
		}
	}

	tokens = append(tokens, deviceTokens...)

	if len(tokens) > 0 {
		err := s.pushService.NotifyPushMulti(ctx, tokens, msg.Title, msg.Body, data)
		if err != nil {
			slog.Error("push notification to multiple failed", "error", err)

			if len(phones) > 0 {
				for _, phone := range phones {
					if err := s.smsService.NotifySMS(ctx, phone, msg.Body); err != nil {
						slog.Error("SMS fallback failed", "phone", phone, "error", err)
					}
				}
			}
		}
	}

	return nil
}
