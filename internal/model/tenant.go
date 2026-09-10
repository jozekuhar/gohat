package model

import (
	"time"
	"uuid"

	"mimokocke/internal/shared/authz"
)

type Organization struct {
	ID        uuid.UUID
	Name      string
	Slug      string
	CreatedAt time.Time
	UpdatedAt *time.Time
}

type MembershipStatus string

func (s MembershipStatus) String() string {
	return string(s)
}

const (
	MembershipStatusActive   MembershipStatus = "active"
	MembershipStatusInactive MembershipStatus = "inactive"
)

type Membership struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	UserID         uuid.UUID
	FirstName      string
	LastName       string
	Role           authz.Role
	Permissions    []authz.Permission
	Status         MembershipStatus
	CanceledAt     *time.Time
	CanceledByID   *uuid.UUID
	CreatedAt      time.Time
	UpdatedAt      *time.Time

	Organization *Organization `db:"-"`
	User         *User         `db:"-"`
	CanceledBy   *User         `db:"-"`
}

type Invitation struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	InviterID      uuid.UUID
	Email          string
	FirstName      string
	LastName       string
	Role           authz.Role
	Permissions    []authz.Permission
	TokenHash      string
	ExpiresAt      time.Time
	AcceptedAt     *time.Time
	DeclinedAt     *time.Time
	CanceledAt     *time.Time
	CanceledByID   *uuid.UUID
	CreatedAt      time.Time

	Organization *Organization `db:"-"`
	Inviter      *User         `db:"-"`
	CanceledBy   *User         `db:"-"`
}

type InvitationStatus string

const (
	InvitationStatusAccepted InvitationStatus = "accepted"
	InvitationStatusExpired  InvitationStatus = "expired"
	InvitationStatusDeclined InvitationStatus = "declined"
	InvitationStatusCanceled InvitationStatus = "canceled"
	InvitationStatusPending  InvitationStatus = "pending"
)

func (i *Invitation) Status() InvitationStatus {
	if i.AcceptedAt != nil {
		return InvitationStatusAccepted
	}
	if i.DeclinedAt != nil {
		return InvitationStatusDeclined
	}
	if i.CanceledAt != nil {
		return InvitationStatusCanceled
	}
	if time.Now().Before(i.ExpiresAt) {
		return InvitationStatusExpired
	}
	return InvitationStatusPending
}

func (i *Invitation) IsPending() bool {
	if i.AcceptedAt != nil || i.DeclinedAt != nil || i.CanceledAt != nil {
		return false
	}
	if time.Now().Before(i.ExpiresAt) {
		return false
	}
	return true
}
