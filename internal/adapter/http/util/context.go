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
