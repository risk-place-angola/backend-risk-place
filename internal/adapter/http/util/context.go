package util

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type ContextKey string

const (
	UserIDCtxKey       ContextKey = "user_id"
	IdentifierCtxKey   ContextKey = "identifier"
	IsAuthenticatedKey ContextKey = "is_authenticated"
)

func GetUserIDFromContext(ctx context.Context) (string, bool) {
	userIDRaw := ctx.Value(UserIDCtxKey)
	userID, ok := userIDRaw.(string)
	if !ok || userID == "" {
		return "", false
	}
	return userID, true
}

func ExtractAndValidateUserID(w http.ResponseWriter, r *http.Request) (uuid.UUID, bool) {
	userID, ok := r.Context().Value("user_id").(string)
	if !ok || userID == "" {
		Error(w, "unauthorized", http.StatusUnauthorized)
		return uuid.Nil, false
	}

	uid, err := uuid.Parse(userID)
	if err != nil {
		Error(w, "invalid user ID", http.StatusBadRequest)
		return uuid.Nil, false
	}

	return uid, true
}

func ExtractAndValidatePathID(w http.ResponseWriter, r *http.Request, paramName, entityName string) (uuid.UUID, bool) {
	id := r.PathValue(paramName)
	if id == "" {
		Error(w, entityName+" ID is required", http.StatusBadRequest)
		return uuid.Nil, false
	}

	parsed, err := uuid.Parse(id)
	if err != nil {
		Error(w, "invalid "+entityName+" ID", http.StatusBadRequest)
		return uuid.Nil, false
	}

	return parsed, true
}

func GetIdentifierFromContext(ctx context.Context) (string, bool) {
	identifierRaw := ctx.Value(IdentifierCtxKey)
	identifier, ok := identifierRaw.(string)
	if !ok || identifier == "" {
		return "", false
	}

	isAuthRaw := ctx.Value(IsAuthenticatedKey)
	isAuth, ok := isAuthRaw.(bool)
	if !ok {
		isAuth = false
	}

	return identifier, isAuth
}

// UserIdentifier holds information about the current user or anonymous session
type UserIdentifier struct {
	UserID          string
	DeviceID        string
	IsAuthenticated bool
}

// ExtractUserIdentifier extracts user ID or device ID from the request context and headers
// This is the unified function that handles both authenticated and anonymous users
func ExtractUserIdentifier(r *http.Request) (*UserIdentifier, bool) {
	ctx := r.Context()

	// Try to get authenticated user ID from context
	userID, hasUserID := GetUserIDFromContext(ctx)

	// Check authentication status
	isAuthRaw := ctx.Value(IsAuthenticatedKey)
	isAuthenticated, _ := isAuthRaw.(bool)

	// Try to get device ID from headers
	deviceID := r.Header.Get("X-Device-Id")
	if deviceID == "" {
		deviceID = r.Header.Get("Device-Id")
	}

	identifier := &UserIdentifier{}

	// Case 1: Authenticated user
	if hasUserID && isAuthenticated {
		identifier.UserID = userID
		identifier.IsAuthenticated = true
		// Device ID is optional for authenticated users
		if deviceID == "" {
			identifier.DeviceID = userID
		} else {
			identifier.DeviceID = deviceID
		}
		return identifier, true
	}

	// Case 2: Anonymous user with device ID
	if deviceID != "" {
		identifier.DeviceID = deviceID
		identifier.UserID = deviceID
		identifier.IsAuthenticated = false
		return identifier, true
	}

	// Case 3: Anonymous user with device ID in context (from middleware)
	if hasUserID && !isAuthenticated {
		identifier.DeviceID = userID
		identifier.UserID = userID
		identifier.IsAuthenticated = false
		return identifier, true
	}

	// No valid identifier found
	return nil, false
}

// ExtractUserIdentifierOrError is a helper that returns error response if extraction fails
func ExtractUserIdentifierOrError(w http.ResponseWriter, r *http.Request) (*UserIdentifier, bool) {
	identifier, ok := ExtractUserIdentifier(r)
	if !ok {
		Error(w, "user ID or device ID required", http.StatusUnauthorized)
		return nil, false
	}
	return identifier, true
}
