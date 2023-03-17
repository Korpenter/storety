// Package user provides the interface for the user service.
package user

import (
	"context"
	"github.com/Mldlr/storety/internal/server/models"
)

// Service is the interface for the user service.
//
// Service provides methods for creating a new user, logging in a user,
// and refreshing a user session.

//go:generate mockery --name=Service -r --case underscore --with-expecter --structname UserService --filename user_service.go
type Service interface {
	// CreateUser creates a new user and returns a new session for the created user, or an error if any occurs.
	CreateUser(ctx context.Context, user *models.User) (*models.Session, error)

	// LogInUser logs in a user and returns a new session and stored salt for the logged-in user, or an error if any occurs.
	LogInUser(ctx context.Context, user *models.User) (*models.Session, string, error)

	// RefreshUserSession refreshes a user session and returns a new session for the user, or an error if any occurs.
	RefreshUserSession(ctx context.Context, oldSession *models.Session) (*models.Session, error)
}
