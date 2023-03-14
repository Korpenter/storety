package models

import "github.com/google/uuid"

type SessionKey struct{}

type Session struct {
	AuthToken    string
	RefreshToken string
	UserID       uuid.UUID
	ID           uuid.UUID
}
