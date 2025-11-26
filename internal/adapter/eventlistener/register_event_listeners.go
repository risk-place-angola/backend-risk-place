package eventlistener

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/risk-place-angola/backend-risk-place/internal/adapter/service"
	"github.com/risk-place-angola/backend-risk-place/internal/adapter/websocket"
	"github.com/risk-place-angola/backend-risk-place/internal/application/port"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/event"
	domainrepository "github.com/risk-place-angola/backend-risk-place/internal/domain/repository"
	domainservice "github.com/risk-place-angola/backend-risk-place/internal/domain/service"
)

func RegisterEventListeners(
	dispatcher port.EventDispatcher,
	hub *websocket.Hub,
	userRepo domainrepository.UserRepository,
	anonymousSessionRepo domainrepository.AnonymousSessionRepository,
	settingsRepo domainrepository.SafetySettingsRepository,
	notifierPush port.NotifierPushService,
	notifierSMS port.NotifierSMSService,
	translationService *service.TranslationService,
) {
	settingsChecker := domainservice.NewSettingsChecker(settingsRepo, anonymousSessionRepo)

	registerBroadcastHandler[event.AlertCreatedEvent](
		dispatcher,
		hub,
		userRepo,
		anonymousSessionRepo,
		settingsChecker,
		notifierPush,
		notifierSMS,
		translationService,
		"AlertCreated",
		func(ctx context.Context, h *websocket.Hub, ev event.AlertCreatedEvent) {
			h.BroadcastAlert(ctx, ev.AlertID.String(), ev.Message, ev.Latitude, ev.Longitude, ev.Radius, ev.Severity)
		},
		"alert_id",
	)

	registerBroadcastHandler[event.ReportCreatedEvent](
		dispatcher,
		hub,
		userRepo,
		anonymousSessionRepo,
		settingsChecker,
		notifierPush,
		notifierSMS,
		translationService,
		"ReportCreated",
		func(ctx context.Context, h *websocket.Hub, ev event.ReportCreatedEvent) {
			h.BroadcastReport(ctx, ev.ReportID.String(), ev.Message, ev.Latitude, ev.Longitude, ev.Radius, ev.IsVerified)
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
	_ domainservice.SettingsChecker,
	notifierPush port.NotifierPushService,
	_ port.NotifierSMSService,
	translationService *service.TranslationService,
	eventName string,
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
		var riskType string
		var id string
		var deviceTokens []string
		var anonymousTokens []string

		switch v := any(ev).(type) {
		case event.AlertCreatedEvent:
			userID = v.UserID
			lat = v.Latitude
			lon = v.Longitude
			radius = v.Radius
			riskType = v.RiskType
			id = v.AlertID.String()

			distanceMeters := int(radius)

			authTokens, err := userRepo.ListDeviceTokensForAlertNotification(ctx, userID, v.Severity, distanceMeters)
			if err != nil {
				slog.Error("failed to list device tokens for alert", "error", err)
			} else {
				for _, token := range authTokens {
					deviceTokens = append(deviceTokens, token.FCMToken)
				}
			}

			anonTokens, err := anonymousSessionRepo.GetFCMTokensForAlertNotification(ctx, lat, lon, radius, v.Severity)
			if err != nil {
				slog.Error("failed to list anonymous tokens for alert", "error", err)
			} else {
				for _, token := range anonTokens {
					anonymousTokens = append(anonymousTokens, token.FCMToken)
				}
			}

		case event.ReportCreatedEvent:
			userID = v.UserID
			lat = v.Latitude
			lon = v.Longitude
			radius = v.Radius
			riskType = v.RiskType
			id = v.ReportID.String()

			distanceMeters := int(radius)

			authTokens, err := userRepo.ListDeviceTokensForReportNotification(ctx, userID, v.IsVerified, distanceMeters)
			if err != nil {
				slog.Error("failed to list device tokens for report", "error", err)
			} else {
				for _, token := range authTokens {
					deviceTokens = append(deviceTokens, token.FCMToken)
				}
			}

			anonTokens, err := anonymousSessionRepo.GetFCMTokensForReportNotification(ctx, lat, lon, radius, v.IsVerified)
			if err != nil {
				slog.Error("failed to list anonymous tokens for report", "error", err)
			} else {
				for _, token := range anonTokens {
					anonymousTokens = append(anonymousTokens, token.FCMToken)
				}
			}
		}

		allTokens := make([]string, 0, len(deviceTokens)+len(anonymousTokens))
		allTokens = append(allTokens, deviceTokens...)
		allTokens = append(allTokens, anonymousTokens...)

		if len(allTokens) > 0 {
			slog.Info("sending push notifications with settings filters",
				slog.String("event", eventName),
				slog.Int("authenticated_users", len(deviceTokens)),
				slog.Int("anonymous_sessions", len(anonymousTokens)),
				slog.Int("total", len(allTokens)))

			eventKey := "alert_created"
			if eventName == "ReportCreated" {
				eventKey = "report_created"
			}

			msg := translationService.GetMessage(eventKey, service.LanguagePortuguese, riskType)

			err := notifierPush.NotifyPushMulti(ctx, allTokens, msg.Title, msg.Body, map[string]string{
				idKey: id,
			})
			if err != nil {
				slog.Error("failed to send push notification", "event_name", eventName, "error", err)
			}
		} else {
			slog.Debug("no users eligible for notification after settings filter", "event_name", eventName)
		}
	})
}
func registerEventHandlers(dispatcher port.EventDispatcher, eventName string, handler func(e event.Event)) {
	dispatcher.Register(eventName, handler)
}
