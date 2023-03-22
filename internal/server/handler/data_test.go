package handler

import (
	"context"
	"github.com/Mldlr/storety/internal/constants"
	pb "github.com/Mldlr/storety/internal/proto"
	"github.com/Mldlr/storety/internal/server/mocks"
	"github.com/Mldlr/storety/internal/server/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
	"time"
)

func TestCreateData(t *testing.T) {
	userID := uuid.New()
	tests := []struct {
		name    string
		setup   func(ctx context.Context, us *mocks.DataService)
		req     *pb.CreateDataRequest
		want    *pb.CreateDataResponse
		ctx     context.Context
		errCode codes.Code
	}{
		{
			name: "Create data successfully",
			setup: func(ctx context.Context, us *mocks.DataService) {
				us.EXPECT().CreateData(mock.AnythingOfType("*context.valueCtx"), userID, &models.Data{
					Name:    "password",
					Type:    "binary",
					Content: []byte("123"),
				}).Return(nil)
			},
			req: &pb.CreateDataRequest{
				Data: &pb.DataItem{
					Name:    "password",
					Type:    "binary",
					Content: []byte("123"),
				},
			},
			want:    &pb.CreateDataResponse{},
			ctx:     context.WithValue(context.Background(), models.SessionKey{}, &models.Session{UserID: userID}),
			errCode: codes.OK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var resp *pb.CreateDataResponse
			var err error

			ctx := context.Background()
			mockUserSrv := new(mocks.DataService)
			if tt.setup != nil {
				tt.setup(ctx, mockUserSrv)
			}
			mockDep := StoretyHandler{dataService: mockUserSrv}
			resp, err = mockDep.CreateData(tt.ctx, tt.req)
			require.EqualValues(t, tt.want, resp)
			if statusErr, ok := status.FromError(err); ok {
				require.Equal(t, tt.errCode.String(), statusErr.Code().String())
			}
		})
	}
}

func TestDeleteData(t *testing.T) {
	userID := uuid.New()
	tests := []struct {
		name    string
		setup   func(ctx context.Context, us *mocks.DataService)
		req     *pb.DeleteDataRequest
		ctx     context.Context
		errCode codes.Code
	}{
		{
			name: "Update data successfully",
			setup: func(ctx context.Context, us *mocks.DataService) {
				us.EXPECT().DeleteData(mock.AnythingOfType("*context.valueCtx"), userID, "password").
					Return(nil)
			},
			req: &pb.DeleteDataRequest{
				Name: "password",
			},
			ctx:     context.WithValue(context.Background(), models.SessionKey{}, &models.Session{UserID: userID}),
			errCode: codes.OK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error

			ctx := context.Background()
			mockUserSrv := new(mocks.DataService)
			if tt.setup != nil {
				tt.setup(ctx, mockUserSrv)
			}
			mockDep := StoretyHandler{dataService: mockUserSrv}
			_, err = mockDep.DeleteData(tt.ctx, tt.req)
			if statusErr, ok := status.FromError(err); ok {
				require.Equal(t, tt.errCode.String(), statusErr.Code().String())
			}
		})
	}
}

func TestListData(t *testing.T) {
	userID := uuid.New()
	tests := []struct {
		name    string
		setup   func(ctx context.Context, us *mocks.DataService)
		req     *pb.ListDataRequest
		resp    *pb.ListDataResponse
		ctx     context.Context
		errCode codes.Code
	}{
		{
			name: "List data successfully",
			setup: func(ctx context.Context, us *mocks.DataService) {
				us.EXPECT().ListData(mock.AnythingOfType("*context.valueCtx"), userID).
					Return([]models.DataInfo{{Name: "dataName", Type: "dataType"}}, nil)
			},
			req:     &pb.ListDataRequest{},
			resp:    &pb.ListDataResponse{Data: []*pb.DataInfo{{Name: "dataName", Type: "dataType"}}},
			ctx:     context.WithValue(context.Background(), models.SessionKey{}, &models.Session{UserID: userID}),
			errCode: codes.OK,
		},
		{
			name: "List data with no data to list",
			setup: func(ctx context.Context, us *mocks.DataService) {
				us.EXPECT().ListData(mock.AnythingOfType("*context.valueCtx"), userID).
					Return(nil, constants.ErrNoData)
			},
			req:     &pb.ListDataRequest{},
			ctx:     context.WithValue(context.Background(), models.SessionKey{}, &models.Session{UserID: userID}),
			errCode: codes.NotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			var resp *pb.ListDataResponse
			ctx := context.Background()
			mockUserSrv := new(mocks.DataService)
			if tt.setup != nil {
				tt.setup(ctx, mockUserSrv)
			}
			mockDep := StoretyHandler{dataService: mockUserSrv}
			resp, err = mockDep.ListData(tt.ctx, tt.req)
			if tt.resp != nil {
				require.EqualValues(t, tt.resp, resp)
			}
			if statusErr, ok := status.FromError(err); ok {
				require.Equal(t, tt.errCode.String(), statusErr.Code().String())
			}
		})
	}
}

func TestCreateBatchData(t *testing.T) {
	userID := uuid.New()
	tests := []struct {
		name    string
		setup   func(ctx context.Context, us *mocks.DataService)
		req     *pb.CreateBatchDataRequest
		want    *pb.CreateBatchResponse
		ctx     context.Context
		errCode codes.Code
	}{
		{
			name: "Create data batch successfully",
			setup: func(ctx context.Context, us *mocks.DataService) {
				us.EXPECT().CreateBatch(mock.AnythingOfType("*context.valueCtx"), userID,
					[]models.Data{{ID: userID, Name: "testName", Type: "binary", UpdatedAt: time.Unix(0, 0).UTC(), Content: []byte("123")}}).
					Return(nil)
			},
			req: &pb.CreateBatchDataRequest{
				Data: []*pb.DataItem{{Id: userID.String(), Name: "testName", Type: "binary", Content: []byte("123")}},
			},
			want:    &pb.CreateBatchResponse{},
			ctx:     context.WithValue(context.Background(), models.SessionKey{}, &models.Session{UserID: userID}),
			errCode: codes.OK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var resp *pb.CreateBatchResponse
			var err error

			ctx := context.Background()
			mockDataSrv := new(mocks.DataService)
			if tt.setup != nil {
				tt.setup(ctx, mockDataSrv)
			}
			mockDep := StoretyHandler{dataService: mockDataSrv}
			resp, err = mockDep.CreateBatchData(tt.ctx, tt.req)
			require.EqualValues(t, tt.want, resp)
			if statusErr, ok := status.FromError(err); ok {
				require.Equal(t, tt.errCode.String(), statusErr.Code().String())
			}
		})
	}
}

func TestUpdateBatchData(t *testing.T) {
	userID := uuid.New()
	tests := []struct {
		name    string
		setup   func(ctx context.Context, us *mocks.DataService)
		req     *pb.UpdateBatchDataRequest
		want    *pb.UpdateBatchResponse
		ctx     context.Context
		errCode codes.Code
	}{
		{
			name: "Update data successfully",
			setup: func(ctx context.Context, us *mocks.DataService) {
				us.EXPECT().UpdateBatch(mock.AnythingOfType("*context.valueCtx"), userID,
					[]models.Data{{ID: userID, Name: "testName", Type: "binary", UpdatedAt: time.Unix(0, 0).UTC(), Content: []byte("123")}}).
					Return(nil)
			},
			req: &pb.UpdateBatchDataRequest{
				Data: []*pb.DataItem{{Id: userID.String(), Name: "testName", Type: "binary", Content: []byte("123")}},
			},
			want:    &pb.UpdateBatchResponse{},
			ctx:     context.WithValue(context.Background(), models.SessionKey{}, &models.Session{UserID: userID}),
			errCode: codes.OK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var resp *pb.UpdateBatchResponse
			var err error

			ctx := context.Background()
			mockDataSrv := new(mocks.DataService)
			if tt.setup != nil {
				tt.setup(ctx, mockDataSrv)
			}
			mockDep := StoretyHandler{dataService: mockDataSrv}
			resp, err = mockDep.UpdateBatchData(tt.ctx, tt.req)
			require.EqualValues(t, tt.want, resp)
			if statusErr, ok := status.FromError(err); ok {
				require.Equal(t, tt.errCode.String(), statusErr.Code().String())
			}
		})
	}
}

func TestSyncData(t *testing.T) {
	userID := uuid.New()
	id := uuid.New()
	tests := []struct {
		name    string
		setup   func(ctx context.Context, us *mocks.DataService)
		req     *pb.SyncRequest
		want    *pb.SyncResponse
		ctx     context.Context
		errCode codes.Code
	}{
		{
			name: "SyncData successfully",
			setup: func(ctx context.Context, us *mocks.DataService) {
				us.EXPECT().GetSyncData(mock.AnythingOfType("*context.valueCtx"), userID,
					[]models.SyncData{{ID: userID, Hash: "testName", UpdatedAt: time.Unix(0, 0).UTC()}}).
					Return([]models.Data{{ID: id}}, []string{"1", "2"}, nil)
			},
			req: &pb.SyncRequest{
				SyncInfo: []*pb.SyncDataItem{{Id: userID.String(), Hash: "testName"}},
			},
			want: &pb.SyncResponse{
				UpdateData: []*pb.DataItem{
					{
						Id:        id.String(),
						UpdatedAt: timestamppb.New(time.Time{}),
					},
				},
				RequestedUpdates: []string{"1", "2"},
			},
			ctx:     context.WithValue(context.Background(), models.SessionKey{}, &models.Session{UserID: userID}),
			errCode: codes.OK,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var resp *pb.SyncResponse
			var err error

			ctx := context.Background()
			mockDataSrv := new(mocks.DataService)
			if tt.setup != nil {
				tt.setup(ctx, mockDataSrv)
			}
			mockDep := StoretyHandler{dataService: mockDataSrv}
			resp, err = mockDep.SyncData(tt.ctx, tt.req)
			require.EqualValues(t, tt.want, resp)
			if statusErr, ok := status.FromError(err); ok {
				require.Equal(t, tt.errCode.String(), statusErr.Code().String())
			}
		})
	}
}
