package router

import (
	"github.com/risk-place-angola/backend-risk-place/internal/adapter/http/middleware"
	"net/http"
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

type RouteGroups struct {
	Public               *RouteGroup
	ProtectedJWT         *RouteGroup
	ProtectedAPIKey      *RouteGroup
	ProtectedAPIKeyLimit *RouteGroup
}

func NewGroups(mux *http.ServeMux, mw MWSet) RouteGroups {
	return RouteGroups{
		Public:               NewRouteGroup(mux, mw.Logging),
		ProtectedJWT:         NewRouteGroup(mux, mw.Logging, mw.JWT),
		ProtectedAPIKey:      NewRouteGroup(mux, mw.Logging, mw.APIKey),
		ProtectedAPIKeyLimit: NewRouteGroup(mux, mw.Logging, mw.APIKeyWithLimit),
	}
}
