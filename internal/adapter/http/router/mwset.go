package router

import (
	"net/http"

	"github.com/risk-place-angola/backend-risk-place/internal/adapter/http/middleware"
	"github.com/risk-place-angola/backend-risk-place/internal/infra/bootstrap"
)

// MWSet contains all available middleware for route groups.
// Pattern: Compose these middlewares in RouteGroups based on authentication needs.
type MWSet struct {
	Logging middleware.Middleware

	JWT          middleware.Middleware
	OptionalAuth middleware.Middleware

	APIKey          middleware.Middleware
	APIKeyWithLimit middleware.Middleware

	RequirePermission func(resource, action string) middleware.Middleware
}

func NewMWSet(c *bootstrap.Container) MWSet {
	return MWSet{
		Logging: middleware.Logging,
		JWT: func(next http.Handler) http.Handler {
			return c.AuthMiddleware.ValidateJWT(next)
		},
		OptionalAuth: func(next http.Handler) http.Handler {
			return c.OptionalAuthMiddleware.ValidateOptional(next)
		},
		RequirePermission: func(resource, action string) middleware.Middleware {
			return func(next http.Handler) http.Handler {
				return c.AuthorizationMiddleware.RequirePermission(resource, action)(next)
			}
		},
	}
}
