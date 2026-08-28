package tenant

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

type membershipStatus string

const (
	membershipStatusActive   membershipStatus = "active"
	membershipStatusInactive membershipStatus = "inactive"
)

type Membership struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	UserID         uuid.UUID
	FirstName      string
	LastName       string
	Role           authz.Role
	Permissions    []authz.Permission
	Status         membershipStatus
	CreatedAt      time.Time
	UpdatedAt      *time.Time
}

type Invitation struct {
	ID             uuid.UUID
	OrganizationID uuid.UUID
	Email          string
	FirstName      string
	LastName       string
	Role           authz.Role
	Permissions    []authz.Permission
	TokenHash      string
	ExpiresAt      time.Time
	CreatedAt      time.Time
}
