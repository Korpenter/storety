package data

import (
	"context"
	"github.com/Mldlr/storety/internal/server/models"
	"github.com/Mldlr/storety/internal/server/storage"
	"github.com/google/uuid"
	"github.com/samber/do"
)

// ServiceImpl is the implementation of the data service.
type ServiceImpl struct {
	storage storage.Storage
}

// NewService creates a new data service.
func NewService(i *do.Injector) *ServiceImpl {
	repo := do.MustInvoke[storage.Storage](i)
	return &ServiceImpl{
		storage: repo,
	}
}

// CreateData creates a new data.
func (s *ServiceImpl) CreateData(ctx context.Context, data *models.Data) error {
	var err error
	data.ID, err = uuid.NewRandom()
	if err != nil {
		return err
	}
	return s.storage.CreateData(ctx, data)
}

// GetDataContent gets the content of a data.
func (s *ServiceImpl) GetDataContent(ctx context.Context, userID uuid.UUID, name string) ([]byte, string, error) {
	return s.storage.GetDataContentByName(ctx, userID, name)
}

// DeleteData deletes a data.
func (s *ServiceImpl) DeleteData(ctx context.Context, userID uuid.UUID, name string) error {
	return s.storage.DeleteDataByName(ctx, userID, name)
}

// ListData lists all data.
func (s *ServiceImpl) ListData(ctx context.Context, userID uuid.UUID) ([]models.DataInfo, error) {
	return s.storage.GetAllDataInfo(ctx, userID)
}
