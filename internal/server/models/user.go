package models

import "github.com/google/uuid"

// User is the user model.
type User struct {
	ID       uuid.UUID
	Login    string
	Password string
}
