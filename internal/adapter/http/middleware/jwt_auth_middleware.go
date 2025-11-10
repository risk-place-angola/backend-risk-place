package middleware

import (
	"context"
	"fmt"
	httputil "github.com/risk-place-angola/backend-risk-place/internal/adapter/http/util"
	"github.com/risk-place-angola/backend-risk-place/internal/config"
	"log/slog"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/risk-place-angola/backend-risk-place/internal/adapter/http/util"
)

type AuthMiddleware struct {
	jwtSecret string
}

// NewAuthMiddleware creates a new instance of the authentication middleware.
func NewAuthMiddleware(cfg config.Config) *AuthMiddleware {
	return &AuthMiddleware{
		jwtSecret: cfg.JWTSecret,
	}
}

// ValidateJWT is a middleware that validates JWT tokens in incoming requests.
func (m *AuthMiddleware) ValidateJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sub, err := m.ValidateJWTFromRequest(r)
		if err != nil {
			slog.Error("JWT validation failed", slog.Any("error", err))
			httputil.Error(w, "Unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), util.UserIDCtxKey, sub)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *AuthMiddleware) ValidateJWTFromRequest(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")

	if authHeader == "" {
		return "", fmt.Errorf("missing or invalid authorization header")
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

	token, err := jwt.Parse(tokenStr, func(_ *jwt.Token) (interface{}, error) {
		return []byte(m.jwtSecret), nil
	})

	if err != nil || !token.Valid {
		slog.Error("Error parsing JWT token", slog.Any("error", err))
		return "", fmt.Errorf("invalid or expired token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid token claims")
	}

	sub, ok := claims["sub"].(string)
	if !ok || sub == "" {
		return "", fmt.Errorf("invalid token subject")
	}

	return sub, nil
}

// Chain applies a list of middleware functions to an http.Handler.
func Chain(h http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}

// ChainFunc applies a list of middleware functions to an http.HandlerFunc.
func ChainFunc(h http.HandlerFunc, middlewares ...func(http.Handler) http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		Chain(h, middlewares...).ServeHTTP(w, r)
	}
}
