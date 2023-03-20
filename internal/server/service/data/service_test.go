package data

import (
	"context"
	"github.com/Mldlr/storety/internal/server/mocks"
	"github.com/Mldlr/storety/internal/server/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestServiceImpl_GetSyncData(t *testing.T) {
	userID := uuid.New()
	tests := []struct {
		name      string
		setup     func(ctx context.Context, s *mocks.Storage)
		userID    uuid.UUID
		syncData  []models.SyncData
		wantUpd   []models.Data
		wantIds   []string
		wantedErr error
	}{
		{
			name: "Create Sync Data with no data from client",
			setup: func(ctx context.Context, s *mocks.Storage) {
				s.EXPECT().GetNewData(ctx, userID, mock.AnythingOfType("[]uuid.UUID")).
					Return([]models.Data{{Name: "Test"}}, nil)
			},
			userID:    userID,
			wantUpd:   []models.Data{{Name: "Test"}},
			wantIds:   nil,
			wantedErr: nil,
		},
		{
			name: "Create Sync Data with data from client",
			setup: func(ctx context.Context, s *mocks.Storage) {
				s.EXPECT().GetDataByUpdateAndHash(ctx, userID, []models.SyncData{{Hash: "1"}}).
					Return([]models.Data{{Name: "Test1"}}, []string{"1", "2", "3"}, nil)
				s.EXPECT().GetNewData(ctx, userID, mock.AnythingOfType("[]uuid.UUID")).
					Return([]models.Data{{Name: "Test2"}}, nil)

			},
			userID:    userID,
			syncData:  []models.SyncData{{Hash: "1"}},
			wantUpd:   []models.Data{{Name: "Test1"}, {Name: "Test2"}},
			wantIds:   []string{"1", "2", "3"},
			wantedErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			mockStorage := new(mocks.Storage)
			if tt.setup != nil {
				tt.setup(ctx, mockStorage)
			}
			mockService := ServiceImpl{storage: mockStorage}
			updates, requestedIDs, err := mockService.GetSyncData(ctx, tt.userID, tt.syncData)
			if tt.wantedErr != nil {
				require.ErrorIs(t, err, tt.wantedErr)
				require.Nil(t, updates, requestedIDs)
				return
			}
			require.EqualValues(t, tt.wantUpd, updates)
			require.EqualValues(t, tt.wantIds, requestedIDs)
		})
	}
}
