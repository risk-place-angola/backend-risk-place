package router

import (
	"net/http"

	"github.com/risk-place-angola/backend-risk-place/internal/adapter/http/middleware"
)

type RouteGroup struct {
	mux         *http.ServeMux
	middlewares []middleware.Middleware
}

func NewRouteGroup(mux *http.ServeMux, middlewares ...middleware.Middleware) *RouteGroup {
	return &RouteGroup{
		mux:         mux,
		middlewares: middlewares,
	}
}

func (rg *RouteGroup) Handle(pattern string, handler http.Handler) {
	final := middleware.Use(handler, rg.middlewares...)
	rg.mux.Handle(pattern, final)
}

func (g *RouteGroup) HandleFunc(pattern string, hf http.HandlerFunc) {
	g.Handle(pattern, hf)
}

// RouteGroups organizes routes by authentication requirements.
// Pattern: Choose the appropriate group based on business rules.
type RouteGroups struct {
	// Public: No authentication required (e.g., login, signup, health)
	Public *RouteGroup

	// OptionalAuth: Accepts JWT OR Device-Id OR Anonymous
	// Use for: Read-only public data that can be enhanced with user context
	OptionalAuth *RouteGroup

	// ProtectedJWT: Requires valid JWT token (authenticated users only)
	// Use for: User-specific operations, writes, sensitive data
	ProtectedJWT *RouteGroup

	// ProtectedAPIKey: Requires API Key (for service-to-service)
	ProtectedAPIKey *RouteGroup

	// ProtectedAPIKeyLimit: API Key with rate limiting
	ProtectedAPIKeyLimit *RouteGroup
}

func NewGroups(mux *http.ServeMux, mw MWSet) RouteGroups {
	return RouteGroups{
		Public:               NewRouteGroup(mux, mw.Logging),
		OptionalAuth:         NewRouteGroup(mux, mw.Logging, mw.OptionalAuth),
		ProtectedJWT:         NewRouteGroup(mux, mw.Logging, mw.JWT),
		ProtectedAPIKey:      NewRouteGroup(mux, mw.Logging, mw.APIKey),
		ProtectedAPIKeyLimit: NewRouteGroup(mux, mw.Logging, mw.APIKeyWithLimit),
	}
}
