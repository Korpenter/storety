package validators

import (
	"errors"
	"github.com/Mldlr/storety/internal/server/models"
)

var (
	ErrEmptyUsername = errors.New("username cannot be empty")
	ErrEmptyPass     = errors.New("password cannot be empty")
)

func ValidateAuthorization(user *models.User) error {
	if user.Login == "" {
		return ErrEmptyUsername
	}
	if user.Password == "" {
		return ErrEmptyPass
	}
	return nil
}
