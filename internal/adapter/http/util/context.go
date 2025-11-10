package util

import "context"

// ContextKey is a type for context keys to avoid collisions.
type ContextKey string

const UserIDCtxKey ContextKey = "user_id"

func GetUserIDFromContext(ctx context.Context) (string, bool) {
	userIDRaw := ctx.Value(UserIDCtxKey)
	userID, ok := userIDRaw.(string)
	if !ok || userID == "" {
		return "", false
	}
	return userID, true
}
