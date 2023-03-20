package service

import (
	"context"
	"fmt"
	"github.com/Mldlr/storety/internal/client/config"
	"github.com/Mldlr/storety/internal/client/models"
	"github.com/Mldlr/storety/internal/client/storage"
	pb "github.com/Mldlr/storety/internal/proto"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
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
	newData, err := c.storage.GetNewData(c.ctx)
	if err != nil {
		return err
	}
	if len(newData) > 0 {
		req := &pb.CreateBatchDataRequest{
			Data: make([]*pb.DataItem, len(newData)),
		}
		for i, d := range newData {
			req.Data[i] = &pb.DataItem{
				Id:        d.ID.String(),
				Name:      d.Name,
				Type:      d.Type,
				Content:   d.Content,
				UpdatedAt: timestamppb.New(d.UpdatedAt),
				Deleted:   d.Deleted,
			}
		}
		_, err = c.remoteClient.CreateBatchData(c.ctx, req)
		if err != nil {
			return err
		}
		err = c.storage.SetSyncedStatus(c.ctx, newData)
		if err != nil {
			return err
		}
	}
	syncData, err := c.storage.GetSyncData(c.ctx)
	if err != nil {
		return err
	}
	syncReq := &pb.SyncRequest{SyncInfo: make([]*pb.SyncDataItem, len(syncData))}
	for i, d := range syncData {
		syncReq.SyncInfo[i] = &pb.SyncDataItem{
			Id:        d.ID.String(),
			Hash:      d.Hash,
			UpdatedAt: timestamppb.New(d.UpdatedAt),
		}
	}
	syncResp, err := c.remoteClient.SyncData(c.ctx, syncReq)
	if syncResp != nil {
		if syncResp.UpdateData != nil || len(syncResp.UpdateData) > 0 {
			serverUpdates := make([]models.Data, len(syncResp.UpdateData))
			for i, d := range syncResp.UpdateData {
				serverUpdates[i] = models.Data{
					ID:        uuid.MustParse(d.Id),
					Name:      d.Name,
					Type:      d.Type,
					Content:   d.Content,
					UpdatedAt: d.UpdatedAt.AsTime(),
					Deleted:   d.Deleted,
				}
			}
			err = c.storage.SyncBatch(c.ctx, serverUpdates)
			if err != nil {
				return err
			}

		}
		if syncResp.RequestedUpdates != nil || len(syncResp.RequestedUpdates) > 0 {
			requestedIDs := make([]uuid.UUID, len(syncResp.RequestedUpdates))
			for i, d := range syncResp.RequestedUpdates {
				requestedIDs[i] = uuid.MustParse(d)
			}
			requestedData, err := c.storage.GetBatch(c.ctx, requestedIDs)
			if err != nil {
				return err
			}
			updReq := &pb.UpdateBatchDataRequest{Data: make([]*pb.DataItem, len(requestedData))}
			for i, d := range requestedData {
				updReq.Data[i] = &pb.DataItem{
					Id:        d.ID.String(),
					Name:      d.Name,
					Type:      d.Type,
					Content:   d.Content,
					UpdatedAt: timestamppb.New(d.UpdatedAt),
					Deleted:   d.Deleted,
				}
			}
			_, err = c.remoteClient.UpdateBatchData(c.ctx, updReq)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// StartSyncData starts a goroutine that syncs data regularly.
func (c *DataServiceImpl) StartSyncData() {
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-c.ctx.Done():
				return
			case <-ticker.C:
				if c.cfg.EncryptionKey != nil && c.conn != nil {
					err := c.SyncData()
					if err != nil {
						fmt.Printf("Error syncing data: %v\n", err)
					}
				}
			}
		}
	}()
}
