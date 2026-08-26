package auth

import (
	"context"
	"log"
	"uuid"
)

type contextKey string

const userIDContextKey contextKey = "userID"

func WithUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, userIDContextKey, userID)
}

func MustUserIDFomContext(ctx context.Context) uuid.UUID {
	value, ok := ctx.Value(userIDContextKey).(uuid.UUID)
	if !ok {
		log.Panicf("required value from context: %s", userIDContextKey)
	}
	return value
}
