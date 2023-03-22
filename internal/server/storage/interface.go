// Package storage defines the interface for the storage layer defining the methods
// for handling user sessions and data storage.
package storage

import (
	"context"
	"github.com/Mldlr/storety/internal/server/models"
	"github.com/google/uuid"
)

// Storage is the interface for the storage layer, which defines the methods for handling user sessions and data storage.
//
//go:generate mockery --name=Storage -r --case underscore --with-expecter --structname Storage --filename storage.go
type Storage interface {
	// CreateUser creates a new user in the storage.
	CreateUser(ctx context.Context, user *models.User) error

	// GetUserDataByName retrieves the UUID, password and salt for a user with the given username.
	GetUserDataByName(ctx context.Context, username string) (uuid.UUID, string, string, error)

	// GetSession retrieves the user's UUID associated with the given session ID and refresh token.
	GetSession(ctx context.Context, sessionID uuid.UUID, refreshToken string) (uuid.UUID, error)

	// CreateSession creates a new session or updates an existing one.
	CreateSession(ctx context.Context, session, oldSession *models.Session) error

	// CreateData creates a new data entry in the storage for a user.
	CreateData(ctx context.Context, userID uuid.UUID, data *models.Data) error

	// GetDataContentByName retrieves the content and type of data entry by name for the given user's UUID.
	GetDataContentByName(ctx context.Context, userID uuid.UUID, name string) ([]byte, string, error)

	// GetAllDataInfo retrieves the list of all data entries' information for the given user's UUID.
	GetAllDataInfo(ctx context.Context, userID uuid.UUID) ([]models.DataInfo, error)

	// DeleteDataByName deletes a data entry by name for the given user's UUID.
	DeleteDataByName(ctx context.Context, userID uuid.UUID, name string) error

	// GetNewData retrieves all data entries that were created after the last sync.
	GetNewData(ctx context.Context, userID uuid.UUID, ids []uuid.UUID) ([]models.Data, error)

	// CreateBatch creates a new batch of data entries in the storage for a user.
	CreateBatch(ctx context.Context, userID uuid.UUID, dataBatch []models.Data) error

	// UpdateBatch updates a batch of data entries in the storage for a user.
	UpdateBatch(ctx context.Context, userID uuid.UUID, dataBatch []models.Data) error

	// GetDataByUpdateAndHash retrieves data entries that were created after the last sync and have a different hash
	// and IDs of entries that were updated locally but not synced.
	GetDataByUpdateAndHash(ctx context.Context, userID uuid.UUID, syncData []models.SyncData) ([]models.Data, []string, error)
}
