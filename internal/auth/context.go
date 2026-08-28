package auth

import (
	"context"
	"fmt"
	"log"
	"uuid"
)

type contextKey string

const userIDContextKey contextKey = "userID"

func WithUserID(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, userIDContextKey, userID)
}

func UserIDFomContext(ctx context.Context) (uuid.UUID, error) {
	value, ok := ctx.Value(userIDContextKey).(uuid.UUID)
	if !ok {
		return uuid.UUID{}, fmt.Errorf("user id not found in tontext")
	}
	return value, nil
}

func MustUserIDFomContext(ctx context.Context) uuid.UUID {
	value, ok := ctx.Value(userIDContextKey).(uuid.UUID)
	if !ok {
		log.Panicf("required value from context: %s", userIDContextKey)
	}
	return value
}
