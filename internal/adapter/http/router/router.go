package router

import (
	"net/http"

	_ "github.com/risk-place-angola/backend-risk-place/api"
	"github.com/risk-place-angola/backend-risk-place/internal/infra/bootstrap"
	httpSwagger "github.com/swaggo/http-swagger"
)

func SetupRoutes(container *bootstrap.Container) *http.ServeMux {
	mux := http.NewServeMux()

	mw := NewMWSet(container)
	g := NewGroups(mux, mw)

	g.Public.HandleFunc("GET health", healthCheckHandler)
	g.Public.HandleFunc("POST /api/v1/devices/register", container.DeviceHandler.RegisterDevice)
	g.Public.HandleFunc("PUT /api/v1/devices/location", container.DeviceHandler.UpdateDeviceLocation)

	g.Public.HandleFunc("POST /api/v1/auth/signup", container.UserHandler.Signup)
	g.Public.HandleFunc("POST /api/v1/auth/login", container.UserHandler.Login)
	g.Public.HandleFunc("POST /api/v1/auth/refresh", container.UserHandler.RefreshToken)
	g.Public.HandleFunc("POST /api/v1/auth/confirm", container.UserHandler.ConfirmSignup)
	g.Public.HandleFunc("POST /api/v1/auth/resend-code", container.UserHandler.ResendCode)
	g.Public.HandleFunc("POST /api/v1/auth/password/forgot", container.UserHandler.ForgotPassword)
	g.Public.HandleFunc("POST /api/v1/auth/password/reset", container.UserHandler.ResetPassword)
	g.ProtectedJWT.HandleFunc("POST /api/v1/auth/logout", container.UserHandler.Logout)

	g.OptionalAuth.HandleFunc("GET /api/v1/risks/types", container.RiskHandler.ListRiskTypes)
	g.OptionalAuth.HandleFunc("GET /api/v1/risks/types/{id}", container.RiskHandler.GetRiskType)
	g.OptionalAuth.HandleFunc("GET /api/v1/risks/topics", container.RiskHandler.ListRiskTopics)
	g.OptionalAuth.HandleFunc("GET /api/v1/risks/topics/{id}", container.RiskHandler.GetRiskTopic)

	g.ProtectedJWT.HandleFunc("GET /api/v1/users/me", container.UserHandler.Me)
	g.ProtectedJWT.HandleFunc("PUT /api/v1/users/profile", container.UserHandler.UpdateProfile)
	g.ProtectedJWT.HandleFunc("PUT /api/v1/users/me/device", container.NotificationHandler.UpdateDeviceInfo)
	g.OptionalAuth.HandleFunc("PUT /api/v1/users/me/notifications/preferences", container.NotificationHandler.UpdateNotificationPreferences)
	g.OptionalAuth.HandleFunc("GET /api/v1/users/me/notifications/preferences", container.NotificationHandler.GetNotificationPreferences)

	g.ProtectedJWT.HandleFunc("GET /api/v1/users/me/emergency-contacts", container.EmergencyContactHandler.GetEmergencyContacts)
	g.ProtectedJWT.HandleFunc("POST /api/v1/users/me/emergency-contacts", container.EmergencyContactHandler.CreateEmergencyContact)
	g.ProtectedJWT.HandleFunc("PUT /api/v1/users/me/emergency-contacts/{id}", container.EmergencyContactHandler.UpdateEmergencyContact)
	g.ProtectedJWT.HandleFunc("DELETE /api/v1/users/me/emergency-contacts/{id}", container.EmergencyContactHandler.DeleteEmergencyContact)
	g.ProtectedJWT.HandleFunc("POST /api/v1/emergency/alert", container.EmergencyContactHandler.SendEmergencyAlert)

	g.OptionalAuth.HandleFunc("GET /api/v1/users/me/settings", container.SafetySettingsHandler.GetSettings)
	g.OptionalAuth.HandleFunc("PUT /api/v1/users/me/settings", container.SafetySettingsHandler.UpdateSettings)

	g.OptionalAuth.HandleFunc("POST /api/v1/alerts", container.AlertHandler.CreateAlert)
	g.OptionalAuth.HandleFunc("POST /api/v1/alerts/{id}/subscribe", container.MyAlertsHandler.SubscribeToAlert)
	g.OptionalAuth.HandleFunc("DELETE /api/v1/alerts/{id}/unsubscribe", container.MyAlertsHandler.UnsubscribeFromAlert)
	g.OptionalAuth.HandleFunc("GET /api/v1/users/me/alerts/created", container.MyAlertsHandler.GetMyCreatedAlerts)
	g.OptionalAuth.HandleFunc("GET /api/v1/users/me/alerts/subscribed", container.MyAlertsHandler.GetMySubscribedAlerts)
	g.ProtectedJWT.HandleFunc("PUT /api/v1/alerts/{id}", container.MyAlertsHandler.UpdateAlert)
	g.ProtectedJWT.HandleFunc("DELETE /api/v1/alerts/{id}", container.MyAlertsHandler.DeleteAlert)

	g.OptionalAuth.HandleFunc("POST /api/v1/location-sharing", container.LocationSharingHandler.CreateLocationSharing)
	g.OptionalAuth.HandleFunc("PUT /api/v1/location-sharing/{id}/location", container.LocationSharingHandler.UpdateLocationSharing)
	g.OptionalAuth.HandleFunc("DELETE /api/v1/location-sharing/{id}", container.LocationSharingHandler.DeleteLocationSharing)
	g.Public.HandleFunc("GET /share/{token}", container.LocationSharingHandler.GetPublicLocationSharing)

	g.OptionalAuth.HandleFunc("POST /api/v1/routes/safe-route", container.SafeRouteHandler.CalculateSafeRoute)
	g.OptionalAuth.HandleFunc("POST /api/v1/routes/incidents-heatmap", container.SafeRouteHandler.GetIncidentsHeatmap)
	g.ProtectedJWT.HandleFunc("POST /api/v1/routes/navigate-home", container.SafeRouteHandler.NavigateToHome)
	g.ProtectedJWT.HandleFunc("POST /api/v1/routes/navigate-work", container.SafeRouteHandler.NavigateToWork)

	g.OptionalAuth.HandleFunc("GET /api/v1/reports", container.ReportHandler.List)
	g.ProtectedJWT.HandleFunc("POST /api/v1/reports", container.ReportHandler.Create)
	g.OptionalAuth.HandleFunc("GET /api/v1/reports/nearby", container.ReportHandler.ListNearby)
	g.OptionalAuth.HandleFunc("PUT /api/v1/reports/{id}/location", container.ReportHandler.UpdateLocation)
	g.OptionalAuth.HandleFunc("POST /api/v1/reports/{id}/vote", container.ReportHandler.VoteReport)
	g.ProtectedJWT.HandleFunc("POST /api/v1/reports/{id}/verify", container.ReportHandler.Verify)
	g.ProtectedJWT.HandleFunc("POST /api/v1/reports/{id}/resolve", container.ReportHandler.Resolve)

	g.ProtectedJWT.HandleFunc("POST /api/v1/upload/risk-type-icon", container.StorageHandler.UploadRiskTypeIcon)
	g.ProtectedJWT.HandleFunc("POST /api/v1/upload/risk-topic-icon", container.StorageHandler.UploadRiskTopicIcon)
	g.Public.HandleFunc("GET /api/v1/storage/{path...}", container.StorageHandler.ServeFile)

	g.OptionalAuth.HandleFunc("POST /api/v1/users/location", container.NearbyUsersHandler.UpdateLocation)
	g.OptionalAuth.HandleFunc("POST /api/v1/users/nearby", container.NearbyUsersHandler.GetNearbyUsers)

	mux.HandleFunc("/ws/alerts", container.WSHandler.HandleWebSocket)
	mux.HandleFunc("/docs/", httpSwagger.WrapHandler)

	return mux
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}

func RoutesDEV(container *bootstrap.Container) {
	if !container.Cfg.IsDevelopment() {
		return
	}
}
