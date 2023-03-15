package constants

import "errors"

// Errors returned by the database and services.
var (
	// ErrUserExists is returned when a user already exists.
	ErrUserExists = errors.New("username taken")
	// ErrUserNotFound is returned when a user does not exist.
	ErrUserNotFound = errors.New("username not found")
	// ErrSessionNotFound is returned when a session does not exist.
	ErrSessionNotFound = errors.New("session not found")
	// ErrCreatingSession is returned when a session could not be created.
	ErrCreatingSession = errors.New("unable to create session")
	// ErrDeletingSession is returned when a session could not be deleted.
	ErrDeletingSession = errors.New("unable to delete session")
	// ErrCreatingData is returned when data could not be created.
	ErrCreatingData = errors.New("unable to create data")
	// ErrGettingData is returned when data could not be retrieved.
	ErrGettingData = errors.New("unable to get data")
	// ErrUpdatingData is returned when data could not be updated.
	ErrDeletingData = errors.New("unable to delete data")
	// ErrUpdatingVersion is returned when the version could not be updated.
	ErrUpdatingVersion = errors.New("unable to update version")
	// ErrNoData is returned when no data is found.
	ErrNoData = errors.New("no data found")
	// ErrInvalidCredentials is returned when the credentials are invalid.
	ErrInvalidCredentials = errors.New("invalid credentials")
	// ErrInvalidRefreshToken is returned when the refresh token is invalid.
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	// ErrExpiredToken is returned when the token is expired.
	ErrExpiredToken = errors.New("token is expired")
)
