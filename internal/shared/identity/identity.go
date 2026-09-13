package identity

import (
	"context"
	"log"
	"slices"
	"uuid"

	"mimokocke/internal/shared/permissions"
)

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

type Identity struct {
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

func (a *Identity) HasPermission(perm permissions.MembershipPermission) bool {
	if a.Role == permissions.RoleOwner {
		return true
	}

	if a.Role == permissions.RoleAdmin &&
		slices.Contains(rolePermissions[permissions.RoleAdmin], perm) {
		return true
	}

	return slices.Contains(a.Permissions, perm)
}

type contextKey string

const identityContextKey contextKey = "identity"

func WithContext(ctx context.Context, identity Identity) context.Context {
	return context.WithValue(ctx, identityContextKey, identity)
}

func MustFromContext(ctx context.Context) Identity {
	value, ok := ctx.Value(identityContextKey).(Identity)
	if !ok {
		log.Panicf("required value from context: %s", identityContextKey)
	}
	return value
}
