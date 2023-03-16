// Package data provides the interface for the data service.
package data

import (
	"context"
	"github.com/Mldlr/storety/internal/server/models"
	"github.com/google/uuid"
)

// Service is the interface for the data service.
//
// Service provides methods for creating, retrieving, listing, and deleting data
// associated with a user.
//
//go:generate mockery --name=Service -r --case underscore --with-expecter --structname DataService --filename data_service.go
type Service interface {
	// CreateData adds a new data entry in the database for the specified user.
	CreateData(ctx context.Context, data *models.Data) error

	// GetDataContent retrieves the content and content type of specified data entry for a user.
	GetDataContent(ctx context.Context, userID uuid.UUID, name string) ([]byte, string, error)

	// DeleteData removes a specified data entry for a user.
	DeleteData(ctx context.Context, userID uuid.UUID, name string) error

	// ListData retrieves a list of all data entries associated with a user.
	ListData(ctx context.Context, userID uuid.UUID) ([]models.DataInfo, error)
}
