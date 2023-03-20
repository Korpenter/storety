package data

import (
	"context"
	"github.com/Mldlr/storety/internal/client/mocks"
	"github.com/Mldlr/storety/internal/client/models"
	pb "github.com/Mldlr/storety/internal/proto"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestCreateData(t *testing.T) {
	ctx := context.Background()
	storageMock := new(mocks.Storage)
	dataService := ServiceImpl{
		ctx:     ctx,
		storage: storageMock,
	}
	name := "Test Data"
	typ := "test"
	content := []byte("test content")

	storageMock.On("CreateData", ctx, mock.AnythingOfType("*models.Data")).Return(nil)
	err := dataService.CreateData(name, typ, content)
	assert.NoError(t, err)

	storageMock.AssertCalled(t, "CreateData", ctx, mock.AnythingOfType("*models.Data"))
	storageMock.AssertNumberOfCalls(t, "CreateData", 1)
	actualData := storageMock.Calls[0].Arguments.Get(1).(*models.Data)
	assert.Equal(t, name, actualData.Name)
	assert.Equal(t, typ, actualData.Type)
	assert.Equal(t, content, actualData.Content)
}

func TestSyncData(t *testing.T) {
	ctx := context.Background()
	storageMock := new(mocks.Storage)
	remoteClientMock := new(mocks.DataClient)

	dataService := ServiceImpl{
		ctx:          ctx,
		storage:      storageMock,
		remoteClient: remoteClientMock,
	}

	id := uuid.New()
	now := time.Now()
	dataItem := models.Data{
		ID:        id,
		Name:      "Test Data",
		Type:      "test",
		Content:   []byte("test content"),
		UpdatedAt: now,
		Deleted:   false,
	}

	storageMock.On("GetNewData", ctx).Return([]models.Data{dataItem}, nil)
	remoteClientMock.On("CreateBatchData", ctx, mock.AnythingOfType("*pb.CreateBatchDataRequest")).Return(&pb.CreateBatchResponse{}, nil)
	storageMock.On("SetSyncedStatus", ctx, []models.Data{dataItem}).Return(nil)
	storageMock.On("GetSyncData", ctx).Return([]models.Data{dataItem}, nil)
	syncDataItem := &pb.SyncDataItem{
		Id:        id.String(),
		Hash:      "",
		UpdatedAt: &timestamp.Timestamp{Seconds: now.Unix()},
	}
	remoteClientMock.On("SyncData", ctx, &pb.SyncRequest{SyncInfo: []*pb.SyncDataItem{syncDataItem}}).Return(&pb.SyncResponse{}, nil)

	err := dataService.SyncData()
	assert.NoError(t, err)

	storageMock.AssertExpectations(t)
	remoteClientMock.AssertExpectations(t)
}
