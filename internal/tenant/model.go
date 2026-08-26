package tenant

import (
	"time"
	"uuid"
)

type Organization struct {
	ID        uuid.UUID
	Name      string
	Slug      string
	CreatedAt time.Time
	UpdatedAt *time.Time
}

type membershipRole string

const (
	MembershipRoleOwner  membershipRole = "owner"
	MembershipRoleAdmin  membershipRole = "admin"
	MembershipRoleMember membershipRole = "member"
)

type membershipStatus string

const (
	MembershipStatusInvited   membershipStatus = "invited"
	MembershipStatusActive    membershipStatus = "active"
	MembershipStatusDeclined  membershipStatus = "declined"
	MembershipStatusSuspended membershipStatus = "suspended"
)

type Membership struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	UserID         uuid.UUID
	Role           membershipRole
	Permissions    []string
	Status         membershipStatus
	CreatedAt      time.Time
	UpdatedAt      *time.Time
}

type Invitation struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	UserID         uuid.UUID
	Role           membershipRole
	Permissions    []string
	ExpiresAt      time.Time
	CreatedAt      time.Time
}
