// Package token provides functionality for generating and verifying JWT tokens.
package token

import "github.com/google/uuid"

// TokenAuth is the interface for the token auth.
//
//go:generate mockery --name=TokenAuth -r --case underscore --with-expecter --structname TokenAuth --filename tokenAuth.go
type TokenAuth interface {
	// GenerateTokenPair generates a new token pair for the specified user and session.
	GenerateTokenPair(id, sessionID uuid.UUID) (string, string, error)

	// Verify verifies the specified token and returns the user ID/session ID, or an error if any occurs.
	Verify(token string) (uuid.UUID, error)
}
