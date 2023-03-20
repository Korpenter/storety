// Package storage defines the interface for the storage layer defining the methods
// for handling user sessions and data storage.
package storage

import (
	"context"
	"github.com/Mldlr/storety/internal/client/models"
	"github.com/google/uuid"
)

// Storage is the interface for the storage layer, which defines the methods for handling user sessions and data storage.
type Storage interface {
	// CreateData creates a new data entry in the storage for a user.
	CreateData(ctx context.Context, data *models.Data) error

	// GetDataContentByName retrieves the content and type of data entry by name.
	GetDataContentByName(ctx context.Context, name string) ([]byte, string, error)

	// GetAllDataInfo retrieves the list of all data entries' information.
	GetAllDataInfo(ctx context.Context) ([]models.DataInfo, error)

	// DeleteDataByName deletes a data entry by name.
	DeleteDataByName(ctx context.Context, name string) error

	GetNewData(ctx context.Context) ([]models.Data, error)

	GetSyncData(ctx context.Context) ([]models.SyncData, error)

	SyncBatch(ctx context.Context, syncBatch []models.Data) error

	GetBatch(ctx context.Context, ids []uuid.UUID) ([]models.Data, error)

	//SetSyncedStatus sets synced status to true for new data that was sent to the server.
	SetSyncedStatus(ctx context.Context, newData []models.Data) error

	// Close closes the db connection.
	Close() error
}
