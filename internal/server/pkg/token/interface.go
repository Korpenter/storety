// Package token provides functionality for generating and verifying JWT tokens.
package token

import "github.com/google/uuid"

// TokenAuth is the interface for the token auth.
//
//go:generate mockery --name=TokenAuth -r --case underscore --with-expecter --structname TokenAuth --filename tokenAuth.go
type TokenAuth interface {
	GenerateTokenPair(id, sessionID uuid.UUID) (string, string, error)
	Verify(token string) (uuid.UUID, error)
}
