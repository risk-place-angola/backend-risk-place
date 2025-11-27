package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	httputil "github.com/risk-place-angola/backend-risk-place/internal/adapter/http/util"
	"github.com/risk-place-angola/backend-risk-place/internal/domain/service"
)

type contextKey string

const permissionCheckedKey contextKey = "permission_checked"

type AuthorizationMiddleware struct {
	authzService *service.AuthorizationService
}

func NewAuthorizationMiddleware(authzService *service.AuthorizationService) *AuthorizationMiddleware {
	return &AuthorizationMiddleware{
		authzService: authzService,
	}
}

func (m *AuthorizationMiddleware) RequirePermission(resource, action string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userIDStr, ok := r.Context().Value(httputil.UserIDCtxKey).(string)
			if !ok || userIDStr == "" {
				slog.Warn("authorization check failed: missing user ID in context")
				httputil.Error(w, "forbidden", http.StatusForbidden)
				return
			}

			userID, err := uuid.Parse(userIDStr)
			if err != nil {
				slog.Error("authorization check failed: invalid user ID", "error", err)
				httputil.Error(w, "forbidden", http.StatusForbidden)
				return
			}

			hasPermission, err := m.authzService.HasPermission(r.Context(), userID, resource, action)
			if err != nil {
				slog.Error("authorization check failed", "error", err, "user_id", userID, "resource", resource, "action", action)
				httputil.Error(w, "forbidden", http.StatusForbidden)
				return
			}

			if !hasPermission {
				slog.Warn("permission denied", "user_id", userID, "resource", resource, "action", action)
				httputil.Error(w, "forbidden: insufficient permissions", http.StatusForbidden)
				return
			}

			ctx := context.WithValue(r.Context(), permissionCheckedKey, true)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
