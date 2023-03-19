package service

import (
	"context"
	"github.com/Mldlr/storety/internal/client/config"
	"github.com/Mldlr/storety/internal/client/models"
	"github.com/Mldlr/storety/internal/client/storage"
	pb "github.com/Mldlr/storety/internal/proto"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// DataService is an interface for the Data service.
type DataService interface {
	CreateData(n, t string, content []byte) error
	ListData() ([]models.DataInfo, error)
	GetData(n string) ([]byte, string, error)
	DeleteData(n string) error
	SyncData() error
	SetStorage(s storage.Storage)
}

// DataServiceImpl is a client for the Data service.
type DataServiceImpl struct {
	ctx          context.Context
	remoteClient pb.DataClient
	conn         *grpc.ClientConn
	storage      storage.Storage
	cfg          *config.Config
}

// NewDataService creates a new DataServiceImpl instance and returns a pointer to it.
// It takes a context, a gRPC client connection, and a configuration object as parameters.
func NewDataService(ctx context.Context, conn *grpc.ClientConn, cfg *config.Config) *DataServiceImpl {
	return &DataServiceImpl{
		ctx:          ctx,
		conn:         conn,
		remoteClient: pb.NewDataClient(conn),
		cfg:          cfg,
	}
}

// SetStorage sets the storage layer for the DataServiceImpl.
func (c *DataServiceImpl) SetStorage(s storage.Storage) {
	if c.storage != nil {
		c.storage.Close()
	}
	c.storage = s
}

// CreateData creates a new data entry locally.
func (c *DataServiceImpl) CreateData(name, typ string, content []byte) error {
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	data := &models.Data{
		ID:      id,
		Name:    name,
		Type:    typ,
		Content: content,
	}
	err = c.storage.CreateData(c.ctx, data)
	if err != nil {
		return err
	}
	return nil
}

// ListData gets list of data from local storage.
func (c *DataServiceImpl) ListData() ([]models.DataInfo, error) {
	return c.storage.GetAllDataInfo(c.ctx)
}

// GetData gets data from local storage.
func (c *DataServiceImpl) GetData(name string) ([]byte, string, error) {
	return c.storage.GetDataContentByName(c.ctx, name)
}

// DeleteData deletes data locally.
func (c *DataServiceImpl) DeleteData(name string) error {
	return c.storage.DeleteDataByName(c.ctx, name)
}

// SyncData get data from remote storage and syncs it with local storage.
func (c *DataServiceImpl) SyncData() error {
	newData, deleteData, lastSync, err := c.storage.GetSyncData(c.ctx)
	if err != nil {
		return err
	}
	req := &pb.SyncRequest{
		CreateData: make([]*pb.DataItem, len(newData)),
		DeleteData: make([]*pb.DataItem, len(deleteData)),
		LastSync:   timestamppb.New(lastSync),
	}
	for i, d := range newData {
		req.CreateData[i] = &pb.DataItem{
			Id:        d.ID.String(),
			Name:      d.Name,
			Type:      d.Type,
			Content:   d.Content,
			UpdatedAt: timestamppb.New(d.UpdatedAt),
			Deleted:   d.Deleted,
		}
	}
	for i, d := range deleteData {
		req.DeleteData[i] = &pb.DataItem{
			Id:        d.ID.String(),
			Name:      d.Name,
			Type:      d.Type,
			Content:   d.Content,
			UpdatedAt: timestamppb.New(d.UpdatedAt),
			Deleted:   d.Deleted,
		}
	}
	updates, err := c.remoteClient.SyncData(c.ctx, req)
	if err != nil {
		return err
	}
	updateData := make([]models.Data, len(updates.Data))
	for i, d := range updates.Data {
		updateData[i] = models.Data{
			ID:        uuid.MustParse(d.Id),
			Name:      d.Name,
			Type:      d.Type,
			Content:   d.Content,
			UpdatedAt: d.UpdatedAt.AsTime(),
			Deleted:   d.Deleted,
		}
	}
	err = c.storage.UpdateSyncData(c.ctx, newData, updateData)
	if err != nil {
		return err
	}
	return nil
}
