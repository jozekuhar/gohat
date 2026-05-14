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

func MustUserIDFomContext(ctx context.Context) int64 {
	value, ok := ctx.Value(userIDContextKey).(int64)
	if !ok {
		log.Panicf("required value from context: %s", userIDContextKey)
	}
	return value
}

func UserIDFomContext(ctx context.Context) (int64, bool) {
	val, ok := ctx.Value(userIDContextKey).(int64)
	return val, ok
}
