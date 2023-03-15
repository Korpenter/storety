package data

import (
	"context"
	"github.com/Mldlr/storety/internal/server/models"
	"github.com/Mldlr/storety/internal/server/storage"
	"github.com/google/uuid"
	"github.com/samber/do"
)

type ServiceImpl struct {
	storage storage.Storage
}

func NewService(i *do.Injector) *ServiceImpl {
	repo := do.MustInvoke[storage.Storage](i)
	return &ServiceImpl{
		storage: repo,
	}
}

func (s *ServiceImpl) CreateData(ctx context.Context, data *models.Data) error {
	var err error
	data.ID, err = uuid.NewRandom()
	if err != nil {
		return err
	}
	return s.storage.CreateData(ctx, data)
}

func (s *ServiceImpl) GetDataContent(ctx context.Context, userID uuid.UUID, name string) ([]byte, string, error) {
	return s.storage.GetDataContentByName(ctx, userID, name)
}

func (s *ServiceImpl) DeleteData(ctx context.Context, userID uuid.UUID, name string) error {
	return s.storage.DeleteDataByName(ctx, userID, name)
}

func (s *ServiceImpl) ListData(ctx context.Context, userID uuid.UUID) ([]models.DataInfo, error) {
	return s.storage.GetAllDataInfo(ctx, userID)
}
