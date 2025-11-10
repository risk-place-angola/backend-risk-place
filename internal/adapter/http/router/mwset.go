package router

import (
	"net/http"

	"github.com/risk-place-angola/backend-risk-place/internal/adapter/http/middleware"
	"github.com/risk-place-angola/backend-risk-place/internal/infra/bootstrap"
)

type MWSet struct {
	Logging         middleware.Middleware
	JWT             middleware.Middleware
	APIKey          middleware.Middleware
	APIKeyWithLimit middleware.Middleware
}

func NewMWSet(c *bootstrap.Container) MWSet {
	return MWSet{
		Logging: middleware.Logging,
		JWT: func(next http.Handler) http.Handler {
			return c.AuthMiddleware.ValidateJWT(next)
		},
	}
}
