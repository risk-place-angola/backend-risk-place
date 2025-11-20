package port

import "context"

type NotifierPushService interface {
	NotifyPush(ctx context.Context, deviceToken string, title string, message string, data map[string]string) error
	NotifyPushMulti(ctx context.Context, deviceTokens []string, title string, message string, data map[string]string) error
}

type NotifierSMSService interface {
	NotifySMS(ctx context.Context, phone string, message string) error
}

type NotificationService interface {
	SendNotificationWithFallback(ctx context.Context, userID, deviceID, language, riskType, eventKey string, data map[string]string) error
	SendNotificationToMultiple(ctx context.Context, userIDs []string, deviceTokens []string, language, riskType, eventKey string, data map[string]string) error
}

type NotifierHubService interface {
	BroadcastAlert(ctx context.Context, alertID, message string, lat, lon float64, radius float64) error
	BroadcastReport(ctx context.Context, reportID, message string, lat, lon, radius float64) error
}

type NotifierUserService interface {
	NotifyUser(ctx context.Context, userID string, message string, data map[string]string) error
}

type NearbyUsersService interface {
	UpdateUserLocation(ctx context.Context, userID string, deviceID string, lat, lon, speed, heading float64, isAnonymous bool) error
	GetNearbyUsers(ctx context.Context, userID string, lat, lon, radiusMeters float64) ([]NearbyUser, error)
	CleanupStaleLocations(ctx context.Context) error
}

type NearbyUser struct {
	UserID      string  `json:"user_id"`
	AnonymousID string  `json:"anonymous_id"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	AvatarID    string  `json:"avatar_id"`
	Color       string  `json:"color"`
	Speed       float64 `json:"speed"`
	Heading     float64 `json:"heading"`
}
