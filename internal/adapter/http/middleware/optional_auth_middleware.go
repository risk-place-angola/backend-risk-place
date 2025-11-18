package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	httputil "github.com/risk-place-angola/backend-risk-place/internal/adapter/http/util"
)

type OptionalAuthMiddleware struct {
	authMiddleware *AuthMiddleware
}

func NewOptionalAuthMiddleware(authMiddleware *AuthMiddleware) *OptionalAuthMiddleware {
	return &OptionalAuthMiddleware{
		authMiddleware: authMiddleware,
	}
}

func (m *OptionalAuthMiddleware) ExtractIdentifier(r *http.Request) (string, bool, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
		userID, err := m.authMiddleware.ValidateJWTFromRequest(r)
		if err == nil {
			return userID, true, nil
		}
		slog.Debug("JWT validation failed, trying device_id", slog.Any("error", err))
	}

	deviceID := r.Header.Get("X-Device-Id")
	if deviceID == "" {
		deviceID = r.Header.Get("Device-Id")
	}

	const minDeviceIDLength = 16
	if deviceID != "" {
		if len(deviceID) >= minDeviceIDLength {
			return deviceID, false, nil
		}
		slog.Debug("Invalid device_id format", slog.String("device_id", deviceID))
	}

	return "", false, http.ErrNoLocation
}

func (m *OptionalAuthMiddleware) ValidateOptional(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		identifier, isAuthenticated, err := m.ExtractIdentifier(r)

		type contextKey string
		const isAuthenticatedKey contextKey = "is_authenticated"

		ctx := r.Context()

		if err != nil {
			// Allow anonymous access - no identifier required
			slog.Debug("No identifier found, allowing anonymous access")
			ctx = context.WithValue(ctx, httputil.UserIDCtxKey, "")
			ctx = context.WithValue(ctx, isAuthenticatedKey, false)
		} else {
			ctx = context.WithValue(ctx, httputil.UserIDCtxKey, identifier)
			ctx = context.WithValue(ctx, isAuthenticatedKey, isAuthenticated)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *OptionalAuthMiddleware) RequireDeviceID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		identifier, _, err := m.ExtractIdentifier(r)
		if err != nil {
			httputil.Error(w, "device_id or authentication required", http.StatusUnauthorized)
			return
		}

		type contextKey string
		const deviceIDKey contextKey = "device_id"

		ctx := context.WithValue(r.Context(), deviceIDKey, identifier)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
