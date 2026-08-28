package auth

import (
	"time"
	"uuid"
)

type User struct {
	ID        uuid.UUID
	Email     string
	CreatedAt time.Time
	UpdatedAt *time.Time
}

type AuthProvider string

const (
	AuthProviderPassword = "password"
	AuthProviderGoogle   = "google"
)

type Authentication struct {
	ID           uuid.UUID
	UserID       uuid.UUID
	Provider     AuthProvider
	ProviderID   *string
	PasswordHash *string
	CreatedAt    time.Time
	UpdatedAt    *time.Time
}

type Session struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	ExpiresAt time.Time
	CreatedAt time.Time
}
