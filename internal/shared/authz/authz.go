package authz

import (
	"slices"
	"uuid"
)

type role string

const (
	RoleOwner  role = "owner"
	RoleAdmin  role = "admin"
	RoleMember role = "member"
)

type permission string

const (
	MembershipRead   permission = "membership:read"
	MembershipCreate permission = "membership:create"
	MembershipUpdate permission = "membership:update"
	MembershipDelete permission = "membership:delete"
	ChannelRead      permission = "channel:read"
	ChannelCreate    permission = "channel:create"
	ChannelUpdate    permission = "channel:update"
	ChannelDelete    permission = "channel:delete"

	// TODO(jozekuhar): billing is future feature
	BillingRead   permission = "billing:read"
	BillingCreate permission = "billing:create"
	BillingUpdate permission = "billing:update"
	BillingDelete permission = "billing:write"
)

var rolePermissions = map[string][]permission{
	string(RoleOwner):  {},
	string(RoleAdmin):  {MembershipRead, MembershipCreate, MembershipUpdate, MembershipDelete},
	string(RoleMember): {},
}

type Identity struct {
	OrganizationID   uuid.UUID
	OrganizationSlug string
	Role             string // map?
	Permissions      []string
}

func (a *Identity) HasPermission(perm permission) bool {
	// admin check?
	if a.Role == string(RoleOwner) {
		return true
	}

	if a.Role == string(RoleAdmin) && slices.Contains(rolePermissions[string(RoleAdmin)], perm) {
		return true
	}

	return slices.Contains(a.Permissions, string(perm))
}
