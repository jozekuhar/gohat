package ctxutil

import (
	"context"
	"log"
)

type contextKey string

var userIDContextKey contextKey = "userID"

func WithUserID(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, userIDContextKey, userID)
}

func MustUserIDFomContext(ctx context.Context) string {
	value, ok := ctx.Value(userIDContextKey).(string)
	if !ok {
		log.Panicf("required value from context: %s", userIDContextKey)
	}
	return value
}
