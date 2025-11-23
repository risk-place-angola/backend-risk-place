package notifier

import (
	"context"
	"log/slog"

	"firebase.google.com/go/v4/messaging"
)

const MaxBatchSize = 500
const ThrottleMs = 400

type FCMNotifier struct {
	messageClient *messaging.Client
}

func NewFCMNotifier(messageClient *messaging.Client) *FCMNotifier {
	return &FCMNotifier{
		messageClient: messageClient,
	}
}

// NotifyPush sends a push notification via FCM
func (f *FCMNotifier) NotifyPush(ctx context.Context, deviceToken string, title string, message string, data map[string]string) error {
	msg := &messaging.Message{
		Token: deviceToken,
		Notification: &messaging.Notification{
			Title: title,
			Body:  message,
		},
		Data: data,
	}

	_, err := f.messageClient.Send(ctx, msg)
	if err != nil {
		slog.Error("Error sending FCM push notification", "error", err)
		return err
	}

	return nil
}

// NotifyPushMulti sends push notifications to multiple device tokens via FCM
func (f *FCMNotifier) NotifyPushMulti(ctx context.Context, deviceTokens []string, title string, message string, data map[string]string) error {
	for i := 0; i < len(deviceTokens); i += MaxBatchSize {
		end := i + MaxBatchSize
		if end > len(deviceTokens) {
			end = len(deviceTokens)
		}
		batch := deviceTokens[i:end]

		msg := &messaging.MulticastMessage{
			Tokens: batch,
			Notification: &messaging.Notification{
				Title: title,
				Body:  message,
			},
			Data: data,
		}

		res, err := f.messageClient.SendEachForMulticast(ctx, msg)
		if err != nil {
			slog.Error("Error sending FCM multicast push notification", "error", err)
			return err
		}

		slog.Info("FCM multicast push notification sent", "successCount", res.SuccessCount, "failureCount", res.FailureCount)
	}

	return nil
}
