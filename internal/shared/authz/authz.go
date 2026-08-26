package authz

import (
	"slices"
	"uuid"
)

type role string

const (
	OwnerRole  role = "owner"
	AdminRole  role = "admin"
	MemberRole role = "member"
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
	string(OwnerRole):  {},
	string(AdminRole):  {MembershipRead, MembershipCreate, MembershipUpdate, MembershipDelete},
	string(MemberRole): {},
}

type Identity struct {
	OrganizationID   uuid.UUID
	OrganizationSlug string
	Role             string // map?
	Permissions      []string
}

func (a *Identity) HasPermission(perm permission) bool {
	// admin check?
	if a.Role == string(OwnerRole) {
		return true
	}

	if a.Role == string(AdminRole) && slices.Contains(rolePermissions[string(AdminRole)], perm) {
		return true
	}

	return slices.Contains(a.Permissions, string(perm))
}
