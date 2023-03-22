// Package validators contains the validators for the server.
package validators

import (
	"errors"
	"github.com/Mldlr/storety/internal/server/models"
)

var (
	// ErrEmptyUsername is returned when the username is empty.
	ErrEmptyUsername = errors.New("username cannot be empty")
	// ErrEmptyPass is returned when the password is empty.
	ErrEmptyPass = errors.New("password cannot be empty")
)

// ValidateAuthorization validates the user login and password.
func ValidateAuthorization(user *models.User) error {
	if user.Login == "" {
		return ErrEmptyUsername
	}
	if user.Password == "" {
		return ErrEmptyPass
	}
	return nil
}
