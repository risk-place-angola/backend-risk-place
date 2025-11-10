package port

import "context"

type NotifierPushService interface {
	NotifyPush(ctx context.Context, deviceToken string, title string, message string, data map[string]string) error
	NotifyPushMulti(ctx context.Context, deviceTokens []string, title string, message string, data map[string]string) error
}

type NotifierSMSService interface {
	NotifySMS(ctx context.Context, phone string, message string) error
}

type NotifierHubService interface {
	BroadcastAlert(ctx context.Context, alertID, message string, lat, lon float64, radius float64) error
	BroadcastReport(ctx context.Context, reportID, message string, lat, lon, radius float64) error
}

type NotifierUserService interface {
	NotifyUser(ctx context.Context, userID string, message string, data map[string]string) error
}
