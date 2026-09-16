package identity

import (
	"context"
	"fmt"
	"log"
	"slices"
	"uuid"

	"mimokocke/internal/shared/permissions"
)

type contextKey string

const (
	authContextKey     contextKey = "authCtx"
	identityContextKey contextKey = "identityCtx"
)

type AuthCtx struct {
	UserID    uuid.UUID
	UserEmail string
}

func WithAuth(ctx context.Context, auth AuthCtx) context.Context {
	return context.WithValue(ctx, authContextKey, auth)
}

func AuthFromContext(ctx context.Context) (AuthCtx, error) {
	value, ok := ctx.Value(authContextKey).(AuthCtx)
	if !ok {
		return AuthCtx{}, fmt.Errorf("user not found in context")
	}
	return value, nil
}

func MustAuthFromContext(ctx context.Context) AuthCtx {
	value, ok := ctx.Value(authContextKey).(AuthCtx)
	if !ok {
		log.Panicf("required value from context: %s", authContextKey)
	}
	return value
}

type IdentityCtx struct {
	ID          uuid.UUID
	Email       string
	FirstName   string
	LastName    string
	OrgID       uuid.UUID
	OrgName     string
	OrgSlug     string
	Role        permissions.MembershipRole
	Permissions []permissions.MembershipPermission
}

func (a *IdentityCtx) HasPermission(perm permissions.MembershipPermission) bool {
	if a.Role == permissions.RoleOwner {
		return true
	}

	if a.Role == permissions.RoleAdmin &&
		slices.Contains(rolePermissions[permissions.RoleAdmin], perm) {
		return true
	}

	return slices.Contains(a.Permissions, perm)
}

func WithIdentity(ctx context.Context, identity IdentityCtx) context.Context {
	return context.WithValue(ctx, identityContextKey, identity)
}

func MustIdentityFromContext(ctx context.Context) IdentityCtx {
	value, ok := ctx.Value(identityContextKey).(IdentityCtx)
	if !ok {
		log.Panicf("required value from context: %s", identityContextKey)
	}
	return value
}

var rolePermissions = map[permissions.MembershipRole][]permissions.MembershipPermission{
	permissions.RoleOwner: {},
	permissions.RoleAdmin: {
		permissions.MembershipRead,
		permissions.MembershipCreate,
		permissions.MembershipUpdate,
		permissions.MembershipDelete,
	},
	permissions.RoleMember: {},
}
