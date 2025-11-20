package repository

import (
	"context"

	"github.com/risk-place-angola/backend-risk-place/internal/domain/model"
)

type AnonymousSessionRepository interface {
	Create(ctx context.Context, session *model.AnonymousSession) error
	FindByDeviceID(ctx context.Context, deviceID string) (*model.AnonymousSession, error)
	Update(ctx context.Context, session *model.AnonymousSession) error
	UpdateLocation(ctx context.Context, deviceID string, lat, lon float64) error
	UpdateFCMToken(ctx context.Context, deviceID string, fcmToken string) error
	GetFCMTokensInRadius(ctx context.Context, lat, lon, radiusMeters float64) ([]string, error)
	Delete(ctx context.Context, deviceID string) error
	CleanupOldSessions(ctx context.Context, daysOld int) error
	TouchLastSeen(ctx context.Context, deviceID string) error
	UpdateNotificationPreferences(ctx context.Context, deviceID string, pushEnabled, smsEnabled bool) error
	GetNotificationPreferences(ctx context.Context, deviceID string) (pushEnabled, smsEnabled bool, err error)
}
