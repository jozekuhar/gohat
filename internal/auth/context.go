package auth

import (
	"context"
	"fmt"
	"log"
	"uuid"
)

type contextKey string

const (
	userContextKey contextKey = "userContext"
)

type AuthContext struct {
	UserID    uuid.UUID
	UserEmail string
}

func WithAuthContext(ctx context.Context, u AuthContext) context.Context {
	return context.WithValue(ctx, userContextKey, u)
}

func AuthFromContext(ctx context.Context) (AuthContext, error) {
	value, ok := ctx.Value(userContextKey).(AuthContext)
	if !ok {
		return AuthContext{}, fmt.Errorf("user not found in context")
	}
	return value, nil
}

func MustAuthFromContext(ctx context.Context) AuthContext {
	value, ok := ctx.Value(userContextKey).(AuthContext)
	if !ok {
		log.Panicf("required value from context: %s", userContextKey)
	}
	return value
}
