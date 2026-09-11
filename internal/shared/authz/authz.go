package authz

import (
	"errors"
	"slices"
	"uuid"

	"mimokocke/internal/provider/db"
)

var ErrPermissionDenied = errors.New("permission denied")

var rolePermissions = map[db.MembershipRole][]db.MembershipPermission{
	db.RoleOwner: {},
	db.RoleAdmin: {
		db.PermMembershipRead,
		db.PermMembershipCreate,
		db.PermMembershipUpdate,
		db.PermMembershipDelete,
	},
	db.RoleMember: {},
}

type Identity struct {
	ID          uuid.UUID
	Email       string
	FirstName   string
	LastName    string
	OrgID       uuid.UUID
	OrgName     string
	OrgSlug     string
	Role        db.MembershipRole
	Permissions []db.MembershipPermission
}

func (a *Identity) HasPermission(perm db.MembershipPermission) bool {
	if a.Role == db.RoleOwner {
		return true
	}

	if a.Role == db.RoleAdmin && slices.Contains(rolePermissions[db.RoleAdmin], perm) {
		return true
	}

	return slices.Contains(a.Permissions, perm)
}
