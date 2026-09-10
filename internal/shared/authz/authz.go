package authz

import (
	"errors"
	"slices"
	"uuid"
)

type Role string

func (r Role) String() string {
	return string(r)
}

var ErrPermissionDenied = errors.New("permission denied")

const (
	RoleOwner  Role = "owner"
	RoleAdmin  Role = "admin"
	RoleMember Role = "member"
)

type Permission string

const (
	PermMembershipRead   Permission = "membership:read"
	PermMembershipCreate Permission = "membership:create"
	PermMembershipUpdate Permission = "membership:update"
	PermMembershipDelete Permission = "membership:delete"
	PermChannelRead      Permission = "channel:read"
	PermChannelCreate    Permission = "channel:create"
	PermChannelUpdate    Permission = "channel:update"
	PermChannelDelete    Permission = "channel:delete"
)

var rolePermissions = map[Role][]Permission{
	RoleOwner: {},
	RoleAdmin: {
		PermMembershipRead,
		PermMembershipCreate,
		PermMembershipUpdate,
		PermMembershipDelete,
	},
	RoleMember: {},
}

type Identity struct {
	ID          uuid.UUID
	Email       string
	FirstName   string
	LastName    string
	OrgID       uuid.UUID
	OrgName     string
	OrgSlug     string
	Role        Role
	Permissions []Permission
}

func (a *Identity) HasPermission(perm Permission) bool {
	if a.Role == RoleOwner {
		return true
	}

	if a.Role == RoleAdmin && slices.Contains(rolePermissions[RoleAdmin], perm) {
		return true
	}

	return slices.Contains(a.Permissions, perm)
}
