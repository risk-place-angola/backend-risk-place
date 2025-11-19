package util

import "context"

// ContextKey is a type for context keys to avoid collisions.
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
