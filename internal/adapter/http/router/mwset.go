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

	// JWT: Validates JWT token (Authorization: Bearer <token>) - REQUIRED
	JWT middleware.Middleware

	// OptionalAuth: Validates JWT OR Device-Id OR allows anonymous
	// Priority: JWT > Device-Id > Anonymous
	OptionalAuth middleware.Middleware

	APIKey          middleware.Middleware
	APIKeyWithLimit middleware.Middleware
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
	}
}
