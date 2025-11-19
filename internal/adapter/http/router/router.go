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

	// ========================================
	// PUBLIC ROUTES (No authentication)
	// ========================================
	g.Public.HandleFunc("GET /api/v1/health", healthCheckHandler)

	// Anonymous Device Registration
	g.Public.HandleFunc("POST /api/v1/devices/register", container.DeviceHandler.RegisterDevice)
	g.Public.HandleFunc("PUT /api/v1/devices/location", container.DeviceHandler.UpdateDeviceLocation)

	// Authentication
	g.Public.HandleFunc("POST /api/v1/auth/signup", container.UserHandler.Signup)
	g.Public.HandleFunc("POST /api/v1/auth/login", container.UserHandler.Login)
	g.Public.HandleFunc("POST /api/v1/auth/confirm", container.UserHandler.ConfirmSignup)
	g.Public.HandleFunc("POST /api/v1/auth/password/forgot", container.UserHandler.ForgotPassword)
	g.Public.HandleFunc("POST /api/v1/auth/password/reset", container.UserHandler.ResetPassword)

	// Risks
	g.OptionalAuth.HandleFunc("GET /api/v1/risks/types", container.RiskHandler.ListRiskTypes)
	g.OptionalAuth.HandleFunc("GET /api/v1/risks/types/{id}", container.RiskHandler.GetRiskType)
	g.OptionalAuth.HandleFunc("GET /api/v1/risks/topics", container.RiskHandler.ListRiskTopics)
	g.OptionalAuth.HandleFunc("GET /api/v1/risks/topics/{id}", container.RiskHandler.GetRiskTopic)

	// User Management
	g.ProtectedJWT.HandleFunc("GET /api/v1/users/me", container.UserHandler.Me)
	g.ProtectedJWT.HandleFunc("PUT /api/v1/users/profile", container.UserHandler.UpdateProfile)

	// Alerts
	g.ProtectedJWT.HandleFunc("POST /api/v1/alerts", container.AlertHandler.CreateAlert)

	// Location Sharing
	g.OptionalAuth.HandleFunc("POST /api/v1/location-sharing", container.LocationSharingHandler.CreateLocationSharing)
	g.OptionalAuth.HandleFunc("PUT /api/v1/location-sharing/{id}/location", container.LocationSharingHandler.UpdateLocationSharing)
	g.OptionalAuth.HandleFunc("DELETE /api/v1/location-sharing/{id}", container.LocationSharingHandler.DeleteLocationSharing)
	g.Public.HandleFunc("GET /share/{token}", container.LocationSharingHandler.GetPublicLocationSharing)

	// Safe Routes
	g.OptionalAuth.HandleFunc("POST /api/v1/routes/safe-route", container.SafeRouteHandler.CalculateSafeRoute)
	g.OptionalAuth.HandleFunc("POST /api/v1/routes/incidents-heatmap", container.SafeRouteHandler.GetIncidentsHeatmap)
	g.ProtectedJWT.HandleFunc("POST /api/v1/routes/navigate-home", container.SafeRouteHandler.NavigateToHome)
	g.ProtectedJWT.HandleFunc("POST /api/v1/routes/navigate-work", container.SafeRouteHandler.NavigateToWork)

	// Reports
	g.OptionalAuth.HandleFunc("GET /api/v1/reports", container.ReportHandler.List)
	g.ProtectedJWT.HandleFunc("POST /api/v1/reports", container.ReportHandler.Create)
	g.OptionalAuth.HandleFunc("GET /api/v1/reports/nearby", container.ReportHandler.ListNearby)
	g.OptionalAuth.HandleFunc("PUT /api/v1/reports/{id}/location", container.ReportHandler.UpdateLocation)
	g.ProtectedJWT.HandleFunc("POST /api/v1/reports/{id}/verify", container.ReportHandler.Verify)
	g.ProtectedJWT.HandleFunc("POST /api/v1/reports/{id}/resolve", container.ReportHandler.Resolve)

	// WebSocket connection
	mux.HandleFunc("/ws/alerts", container.WSHandler.HandleWebSocket)

	// Swagger documentation
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
