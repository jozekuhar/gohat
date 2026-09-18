package sqlc

import "time"

type MembershipStatus string

func (s MembershipStatus) String() string {
	return string(s)
}

const (
	MemberStatusActive   MembershipStatus = "active"
	MemberStatusInactive MembershipStatus = "inactive"
)

type InvitationStatus string

func (s InvitationStatus) String() string {
	return string(s)
}

const (
	InvitationStatusPending  InvitationStatus = "pending"
	InvitationStatusAccepted InvitationStatus = "accepted"
	InvitationStatusDeclined InvitationStatus = "declined"
	InvitationStatusCanceled InvitationStatus = "canceled"
	InvitationStatusExpired  InvitationStatus = "expired"
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
	if time.Now().After(i.ExpiresAt) {
		return InvitationStatusExpired
	}
	return InvitationStatusPending
}

func (i *Invitation) IsPending() bool {
	return i.Status() == InvitationStatusPending
}
