// Package constants contains errors returned by the database and services.
package constants

import "errors"

var (
	// ErrUserExists is returned when a user already exists.
	ErrUserExists = errors.New("username taken")

	// ErrUserNotFound is returned when a user does not exist.
	ErrUserNotFound = errors.New("user not found")

	// ErrSessionNotFound is returned when a session does not exist.
	ErrSessionNotFound = errors.New("session not found")

	// ErrCreateSession is returned when a session could not be created.
	ErrCreateSession = errors.New("unable to create session")

	// ErrDeleteSession is returned when a session could not be deleted.
	ErrDeleteSession = errors.New("unable to delete session")

	// ErrCreateData is returned when data could not be created.
	ErrCreateData = errors.New("unable to create data")

	// ErrGetData is returned when data could not be retrieved.
	ErrGetData = errors.New("unable to get data")

	// ErrDeleteData is returned when data could not be deleted.
	ErrDeleteData = errors.New("unable to delete data")

	// ErrUpdateData is returned when data could not be updated.
	ErrUpdateData = errors.New("unable to update data")

	// ErrNoData is returned when no data is found.
	ErrNoData = errors.New("no data found")

	// ErrInvalidCredentials is returned when the credentials are invalid.
	ErrInvalidCredentials = errors.New("invalid credentials")

	// ErrInvalidRefreshToken is returned when the refresh token is invalid.
	ErrInvalidRefreshToken = errors.New("invalid refresh token")

	// ErrExpiredToken is returned when the token is expired.
	ErrExpiredToken = errors.New("token is expired")
)
