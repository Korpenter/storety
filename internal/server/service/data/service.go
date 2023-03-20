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

// CreateData creates a new data entry for the specified user.
func (s *ServiceImpl) CreateData(ctx context.Context, userID uuid.UUID, data *models.Data) error {
	var err error
	data.ID, err = uuid.NewRandom()
	if err != nil {
		return err
	}
	return s.storage.CreateData(ctx, userID, data)
}

// GetDataContent retrieves the content and content type of specified data entry for a user.
func (s *ServiceImpl) GetDataContent(ctx context.Context, userID uuid.UUID, name string) ([]byte, string, error) {
	return s.storage.GetDataContentByName(ctx, userID, name)
}

// CreateBatch adds a new data batch in the database for the specified user.
func (s *ServiceImpl) CreateBatch(ctx context.Context, userID uuid.UUID, dataBatch []models.Data) error {
	if len(dataBatch) > 0 {
		return s.storage.CreateBatch(ctx, userID, dataBatch)
	}
	return nil
}

// UpdateBatch updates a data batch in the database for the specified user.
func (s *ServiceImpl) UpdateBatch(ctx context.Context, userID uuid.UUID, dataBatch []models.Data) error {
	if len(dataBatch) > 0 {
		return s.storage.UpdateBatch(ctx, userID, dataBatch)
	}
	return nil
}

// DeleteData removes a specified data entry for a user.
func (s *ServiceImpl) DeleteData(ctx context.Context, userID uuid.UUID, name string) error {
	return s.storage.DeleteDataByName(ctx, userID, name)
}

// ListData retrieves a list of all data entries associated with a user.
func (s *ServiceImpl) ListData(ctx context.Context, userID uuid.UUID) ([]models.DataInfo, error) {
	return s.storage.GetAllDataInfo(ctx, userID)
}

// GetSyncData gets new data and ids for client to send to server for update.
func (s *ServiceImpl) GetSyncData(ctx context.Context, userID uuid.UUID, syncData []models.SyncData) ([]models.Data, []string, error) {
	ids := make([]uuid.UUID, len(syncData))
	for i := range syncData {
		ids[i] = syncData[i].ID
	}
	if len(syncData) == 0 {
		newData, err := s.storage.GetNewData(ctx, userID, ids)
		if err != nil {
			return nil, nil, err
		}
		return newData, nil, nil
	}
	updatedData, requestID, err := s.storage.GetDataByUpdateAndHash(ctx, userID, syncData)
	if err != nil {
		return nil, nil, err
	}
	newData, err := s.storage.GetNewData(ctx, userID, ids)
	if err != nil {
		return nil, nil, err
	}
	return append(updatedData, newData...), requestID, nil
}
