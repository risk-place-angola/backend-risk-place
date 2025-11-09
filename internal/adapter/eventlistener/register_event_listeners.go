package eventlistener

import (
	"context"
	"log/slog"

	"github.com/risk-place-angola/backend-risk-place/internal/adapter/websocket"
	"github.com/risk-place-angola/backend-risk-place/internal/application/port"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/event"
	domainrepository "github.com/risk-place-angola/backend-risk-place/internal/domain/repository"
)

func RegisterEventListeners(
	dispatcher port.EventDispatcher,
	hub *websocket.Hub,
	userRepo domainrepository.UserRepository,
	notifierPush port.NotifierPushService,
	notifierSMS port.NotifierSMSService,
) {
	dispatcher.Register("AlertCreated", func(e event.Event) {
		ev := e.(event.AlertCreatedEvent)
		ctx := context.Background()

		hub.BroadcastAlert(ctx, ev.AlertID.String(), ev.Message, ev.Latitude, ev.Longitude, ev.Radius)

		deviceTokens, err := userRepo.ListDeviceTokensByUserIDs(ctx, ev.UserID)
		if err != nil {
			slog.Error("failed to list device tokens", "error", err)
			return
		}

		if len(deviceTokens) > 0 {
			err = notifierPush.NotifyPushMulti(ctx, deviceTokens, "🚨 Alerta de Risco", ev.Message, map[string]string{
				"alert_id": ev.AlertID.String(),
			})
			if err != nil {
				slog.Error("failed to send push notification for alert", "error", err)
			}
		}
	})

	dispatcher.Register("ReportCreated", func(e event.Event) {
		ev := e.(event.ReportCreatedEvent)
		ctx := context.Background()

		hub.BroadcastReport(ctx, ev.ReportID.String(), ev.Message, ev.Latitude, ev.Longitude, ev.Radius)

		deviceTokens, err := userRepo.ListDeviceTokensByUserIDs(ctx, ev.UserID)
		if err != nil {
			slog.Error("failed to list device tokens", "error", err)
			return
		}

		if len(deviceTokens) > 0 {
			err = notifierPush.NotifyPushMulti(ctx, deviceTokens, "📍 Novo Relato de Risco", ev.Message, map[string]string{
				"report_id": ev.ReportID.String(),
			})
			if err != nil {
				slog.Error("failed to send push notification for report", "error", err)
			}
		}
	})

	dispatcher.Register("ReportResolved", func(e event.Event) {
		ev := e.(event.ReportResolvedEvent)
		for _, uid := range ev.UserIDs {
			hub.NotifyUser(uid, "report_resolved", map[string]string{
				"report_id": ev.ReportID.String(),
				"message":   ev.Message,
			})
		}
	})

	dispatcher.Register("ReportVerified", func(e event.Event) {
		ev := e.(event.ReportVerifiedEvent)
		hub.NotifyUser(ev.UserID.String(), "report_verified", map[string]string{
			"report_id": ev.ReportID.String(),
			"message":   "Seu relatório foi verificado.",
		})
	})
}
