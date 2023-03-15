package models

import "github.com/google/uuid"

// SessionKey for retrival of the session from the context.
type SessionKey struct{}

// Session is the session model.
type Session struct {
	AuthToken    string
	RefreshToken string
	UserID       uuid.UUID
	ID           uuid.UUID
}
