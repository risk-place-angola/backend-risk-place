package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/model"
)

//nolint:interfacebloat // repository interface naturally has many methods
type UserRepository interface {
	GenericRepository[model.User]
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	FindByEmailOrPhone(ctx context.Context, identifier string) (*model.User, error)
	AddCodeToUser(ctx context.Context, userID uuid.UUID, code string, expiration time.Time) error
	MarkAccountVerified(ctx context.Context, userID uuid.UUID) error
	UpdateUserPassword(ctx context.Context, userID uuid.UUID, newPassword string) error
	UserHasPermission(ctx context.Context, userID uuid.UUID, permission string) (bool, error)
	ListAllDeviceTokensExceptUser(ctx context.Context, excludeUserID uuid.UUID) ([]string, error)
	UpdateUserDeviceInfo(ctx context.Context, userID uuid.UUID, fcmToken string, language string) error
	ListDeviceTokensByUserIDs(ctx context.Context, userIDs []uuid.UUID) ([]string, error)
	ListDeviceTokensForAlertNotification(ctx context.Context, userIDs []uuid.UUID, severityLevel string, distanceMeters int) ([]model.DeviceToken, error)
	ListDeviceTokensForReportNotification(ctx context.Context, userIDs []uuid.UUID, isVerified bool, distanceMeters int) ([]model.DeviceToken, error)
	UpdateSavedLocations(ctx context.Context, userID uuid.UUID, homeAddress, workAddress *model.SavedLocation) error
	UpdateNotificationPreferences(ctx context.Context, userID uuid.UUID, pushEnabled, smsEnabled bool) error
	GetNotificationPreferences(ctx context.Context, userID uuid.UUID) (pushEnabled, smsEnabled bool, err error)
	GetUserLanguageAndPhone(ctx context.Context, userID uuid.UUID) (language, phone string, err error)
}
