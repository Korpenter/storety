package data

import (
	"context"
	"github.com/Mldlr/storety/internal/server/models"
	"github.com/google/uuid"
)

// Service is the interface for the data service.
//
//go:generate mockery --name=Service -r --case underscore --with-expecter --structname DataService --filename data_service.go
type Service interface {
	CreateData(ctx context.Context, data *models.Data) error
	GetDataContent(ctx context.Context, userID uuid.UUID, name string) ([]byte, string, error)
	DeleteData(ctx context.Context, userID uuid.UUID, name string) error
	ListData(ctx context.Context, userID uuid.UUID) ([]models.DataInfo, error)
}
