package auth

import (
	"time"
	"uuid"
)

type User struct {
	ID        uuid.UUID
	Email     string
	GoogleID  string
	CreatedAt time.Time
	UpdatedAt *time.Time
}

type Session struct {
	ID     uuid.UUID
	UserID uuid.UUID
	// OrganizationID   uuid.UUID
	// OrganizationSlug string
	ExpiresAt time.Time
	CreatedAt time.Time
}
