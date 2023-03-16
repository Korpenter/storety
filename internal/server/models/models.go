// Package models provides the data models used by the Storety server.
package models

import "github.com/google/uuid"

// User is the user model.
type User struct {
	ID       uuid.UUID
	Login    string
	Password string
}

// SessionKey for retrieval of the session from the context.
type SessionKey struct{}

// Session is the session model.
type Session struct {
	AuthToken    string
	RefreshToken string
	UserID       uuid.UUID
	ID           uuid.UUID
}

// Data is the data model.
type Data struct {
	ID      uuid.UUID
	UserID  uuid.UUID
	Name    string
	Type    string
	Content []byte
}

// DataInfo is the data info model.
type DataInfo struct {
	Name string
	Type string
}
