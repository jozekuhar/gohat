package permissions

import "errors"

var ErrDenied = errors.New("permission denied")

type MembershipRole string

func (r MembershipRole) String() string {
	return string(r)
}

const (
	RoleOwner  MembershipRole = "owner"
	RoleAdmin  MembershipRole = "admin"
	RoleMember MembershipRole = "member"
)

type MembershipPermission string

func (p MembershipPermission) String() string {
	return string(p)
}

const (
	MembershipRead   MembershipPermission = "membership:read"
	MembershipCreate MembershipPermission = "membership:create"
	MembershipUpdate MembershipPermission = "membership:update"
	MembershipDelete MembershipPermission = "membership:delete"
	ChannelRead      MembershipPermission = "channel:read"
	ChannelCreate    MembershipPermission = "channel:create"
	ChannelUpdate    MembershipPermission = "channel:update"
	ChannelDelete    MembershipPermission = "channel:delete"
	OrderRead        MembershipPermission = "order:read"
)
