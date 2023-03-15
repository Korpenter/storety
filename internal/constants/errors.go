package constants

import "errors"

var (
	ErrUserExists          = errors.New("username taken")
	ErrUserNotFound        = errors.New("username not found")
	ErrSessionNotFound     = errors.New("session not found")
	ErrCreatingSession     = errors.New("unable to create session")
	ErrDeletingSession     = errors.New("unable to delete session")
	ErrCreatingData        = errors.New("unable to create data")
	ErrGettingData         = errors.New("unable to get data")
	ErrDeletingData        = errors.New("unable to delete data")
	ErrUpdatingVersion     = errors.New("unable to update version")
	ErrNoData              = errors.New("no data found")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrExpiredToken        = errors.New("token is expired")
)
