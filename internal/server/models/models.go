// Package models provides the data models used by the Storety server.
package models

import (
	"github.com/google/uuid"
	"time"
)

// User is the user model.
type User struct {
	ID       uuid.UUID
	Login    string
	Password string
	Salt     string
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
	ID        uuid.UUID
	Name      string
	Type      string
	Content   []byte
	UpdatedAt time.Time
	Synced    bool
	Deleted   bool
}

// DataInfo is the data info model.
type DataInfo struct {
	Name string
	Type string
}

// SyncData is the data sync model for syncing client db with server.
type SyncData struct {
	CreateData []Data
	DeleteData []Data
	LastSync   time.Time
}
