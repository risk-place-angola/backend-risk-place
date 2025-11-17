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

	// Health
	g.Public.HandleFunc("GET /api/v1/health", healthCheckHandler)

	// Auth
	g.Public.HandleFunc("POST /api/v1/auth/signup", container.UserHandler.Signup)
	g.Public.HandleFunc("POST /api/v1/auth/login", container.UserHandler.Login)
	g.Public.HandleFunc("POST /api/v1/auth/confirm", container.UserHandler.ConfirmSignup)
	g.Public.HandleFunc("POST /api/v1/auth/password/forgot", container.UserHandler.ForgotPassword)
	g.Public.HandleFunc("POST /api/v1/auth/password/reset", container.UserHandler.ResetPassword)

	// User protected routes
	g.ProtectedJWT.HandleFunc("GET /api/v1/users/me", container.UserHandler.Me)

	// WebSocket connection
	mux.HandleFunc("/ws/alerts", container.WSHandler.HandleWebSocket)

	// Risks
	g.ProtectedJWT.HandleFunc("GET /api/v1/risks/types", container.RiskHandler.ListRiskTypes)
	g.ProtectedJWT.HandleFunc("GET /api/v1/risks/topics", container.RiskHandler.ListRiskTopics)

	// Alert
	g.ProtectedJWT.HandleFunc("POST /api/v1/alerts", container.AlertHandler.CreateAlert)

	// Reports
	g.ProtectedJWT.HandleFunc("POST /api/v1/reports", container.ReportHandler.Create)
	g.ProtectedJWT.HandleFunc("GET /api/v1/reports/nearby", container.ReportHandler.ListNearby)
	g.ProtectedJWT.HandleFunc("POST /api/v1/reports/{id}/verify", container.ReportHandler.Verify)
	g.ProtectedJWT.HandleFunc("POST /api/v1/reports/{id}/resolve", container.ReportHandler.Resolve)

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
