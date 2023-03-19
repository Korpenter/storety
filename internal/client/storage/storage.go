// Package storage defines the interface for the storage layer defining the methods
// for handling user sessions and data storage.
package storage

import (
	"context"
	"github.com/Mldlr/storety/internal/client/models"
	"time"
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

	// GetSyncData get sync data for sync with remote and last sync timestamp.
	GetSyncData(ctx context.Context) ([]models.Data, []models.Data, time.Time, error)

	//UpdateSyncData updates data and last synced time in db.
	UpdateSyncData(ctx context.Context, syncedNewData []models.Data, updatedData []models.Data) error

	// Close closes the db connection.
	Close() error
}
