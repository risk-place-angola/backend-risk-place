package eventlistener

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/risk-place-angola/backend-risk-place/internal/adapter/websocket"
	"github.com/risk-place-angola/backend-risk-place/internal/application/port"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/event"
	domainrepository "github.com/risk-place-angola/backend-risk-place/internal/domain/repository"
)

func RegisterEventListeners(
	dispatcher port.EventDispatcher,
	hub *websocket.Hub,
	userRepo domainrepository.UserRepository,
	anonymousSessionRepo domainrepository.AnonymousSessionRepository,
	notifierPush port.NotifierPushService,
	notifierSMS port.NotifierSMSService,
) {
	registerBroadcastHandler[event.AlertCreatedEvent](
		dispatcher,
		hub,
		userRepo,
		anonymousSessionRepo,
		notifierPush,
		"AlertCreated",
		"🚨 Alerta de Risco",
		func(ctx context.Context, h *websocket.Hub, ev event.AlertCreatedEvent) {
			h.BroadcastAlert(ctx, ev.AlertID.String(), ev.Message, ev.Latitude, ev.Longitude, ev.Radius)
		},
		"alert_id",
	)

	registerBroadcastHandler[event.ReportCreatedEvent](
		dispatcher,
		hub,
		userRepo,
		anonymousSessionRepo,
		notifierPush,
		"ReportCreated",
		"📍 Novo Relato de Risco",
		func(ctx context.Context, h *websocket.Hub, ev event.ReportCreatedEvent) {
			h.BroadcastReport(ctx, ev.ReportID.String(), ev.Message, ev.Latitude, ev.Longitude, ev.Radius)
		},
		"report_id",
	)

	dispatcher.Register("ReportResolved", func(e event.Event) {
		ev, ok := e.(event.ReportResolvedEvent)
		if !ok {
			slog.Error("failed to cast event to ReportResolvedEvent")
			return
		}

		for _, uid := range ev.UserIDs {
			hub.NotifyUser(uid, "report_resolved", map[string]string{
				"report_id": ev.ReportID.String(),
				"message":   ev.Message,
			})
		}
	})

	dispatcher.Register("ReportVerified", func(e event.Event) {
		ev, ok := e.(event.ReportVerifiedEvent)
		if !ok {
			slog.Error("failed to cast event to ReportVerifiedEvent")
			return
		}

		hub.NotifyUser(ev.UserID.String(), "report_verified", map[string]string{
			"report_id": ev.ReportID.String(),
			"message":   "Seu relatório foi verificado.",
		})
	})
}

func registerBroadcastHandler[T any](
	dispatcher port.EventDispatcher,
	hub *websocket.Hub,
	userRepo domainrepository.UserRepository,
	anonymousSessionRepo domainrepository.AnonymousSessionRepository,
	notifierPush port.NotifierPushService,
	eventName, title string,
	broadcast func(context.Context, *websocket.Hub, T),
	idKey string,
) {
	registerEventHandlers(dispatcher, eventName, func(e event.Event) {
		ev, ok := e.(T)
		if !ok {
			slog.Error("failed to cast event", "event_name", eventName)
			return
		}

		ctx := context.Background()
		broadcast(ctx, hub, ev)

		var userID []uuid.UUID
		var lat, lon, radius float64

		switch v := any(ev).(type) {
		case event.AlertCreatedEvent:
			userID = v.UserID
			lat = v.Latitude
			lon = v.Longitude
			radius = v.Radius
		case event.ReportCreatedEvent:
			userID = v.UserID
			lat = v.Latitude
			lon = v.Longitude
			radius = v.Radius
		}

		// Get FCM tokens from authenticated users
		deviceTokens, err := userRepo.ListDeviceTokensByUserIDs(ctx, userID)
		if err != nil {
			slog.Error("failed to list device tokens", "event_name", eventName, "error", err)
		}

		// Get FCM tokens from anonymous sessions in radius
		anonymousTokens, err := anonymousSessionRepo.GetFCMTokensInRadius(ctx, lat, lon, radius)
		if err != nil {
			slog.Error("failed to list anonymous tokens", "event_name", eventName, "error", err)
		}

		// Combine all tokens
		allTokens := make([]string, 0, len(deviceTokens)+len(anonymousTokens))
		allTokens = append(allTokens, deviceTokens...)
		allTokens = append(allTokens, anonymousTokens...)

		if len(allTokens) > 0 {
			var id string
			switch v := any(ev).(type) {
			case event.AlertCreatedEvent:
				id = v.AlertID.String()
			case event.ReportCreatedEvent:
				id = v.ReportID.String()
			}

			slog.Info("sending push notifications",
				slog.String("event", eventName),
				slog.Int("authenticated_users", len(deviceTokens)),
				slog.Int("anonymous_sessions", len(anonymousTokens)),
				slog.Int("total", len(allTokens)))

			err = notifierPush.NotifyPushMulti(ctx, allTokens, title, getMessage(ev), map[string]string{
				idKey: id,
			})
			if err != nil {
				slog.Error("failed to send push notification", "event_name", eventName, "error", err)
			}
		}
	})
}

func getMessage(e any) string {
	switch v := e.(type) {
	case event.AlertCreatedEvent:
		return v.Message
	case event.ReportCreatedEvent:
		return v.Message
	default:
		return ""
	}
}

func registerEventHandlers(dispatcher port.EventDispatcher, eventName string, handler func(e event.Event)) {
	dispatcher.Register(eventName, handler)
}
