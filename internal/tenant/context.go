package tenant

import (
	"context"
	"log"

	"mimokocke/internal/shared/authz"
)

type contextKey string

const identityContextKey contextKey = "identity"

func WithIdentity(ctx context.Context, identity authz.Identity) context.Context {
	return context.WithValue(ctx, identityContextKey, identity)
}

func MustIdentityFromContext(ctx context.Context) authz.Identity {
	value, ok := ctx.Value(identityContextKey).(authz.Identity)
	if !ok {
		log.Panicf("required value from context: %s", identityContextKey)
	}
	return value
}
