package data

import (
	"github.com/Mldlr/storety/internal/client/models"
	"github.com/Mldlr/storety/internal/client/storage"
)

// Service is an interface for the Data service.
type Service interface {
	// CreateData creates a new data entry locally.
	CreateData(n, t string, content []byte) error

	// ListData gets list of data from local storage.
	ListData() ([]models.DataInfo, error)

	// GetData gets data from local storage.
	GetData(n string) ([]byte, string, error)

	// DeleteData deletes data locally.
	DeleteData(n string) error

	// SyncData get data from remote storage and syncs it with local storage.
	SyncData() error

	// SetStorage sets the storage layer for the ServiceImpl.
	SetStorage(s storage.Storage)
}
