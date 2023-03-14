package user

import (
	"context"
	mocks2 "github.com/Mldlr/storety/internal/server/mocks"
	"github.com/Mldlr/storety/internal/server/models"
	"github.com/Mldlr/storety/internal/server/storage"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestService_CreateUser(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(ctx context.Context, ta *mocks2.TokenAuth, s *mocks2.Storage)
		user      *models.User
		want      *models.Session
		wantedErr error
	}{
		{
			name: "Create user successfully",
			setup: func(ctx context.Context, ta *mocks2.TokenAuth, s *mocks2.Storage) {
				var nilSession *models.Session
				s.EXPECT().CreateUser(ctx, mock.AnythingOfType("*models.User")).Return(nil)
				ta.EXPECT().GenerateTokenPair(mock.AnythingOfType("uuid.UUID"), mock.AnythingOfType("uuid.UUID")).
					Return("auth_token", "refresh_token", nil)
				s.EXPECT().CreateSession(ctx, mock.AnythingOfType("*models.Session"), nilSession).
					Return(nil)
			},
			user: &models.User{
				Login:    "username",
				Password: "password",
			},
			want: &models.Session{
				AuthToken:    "auth_token",
				RefreshToken: "refresh_token",
			},
			wantedErr: nil,
		},
		{
			name: "Fail to create user with duplicated username",
			setup: func(ctx context.Context, ta *mocks2.TokenAuth, s *mocks2.Storage) {
				s.EXPECT().CreateUser(ctx, mock.AnythingOfType("*models.User")).Return(storage.ErrUserExists)
			},
			user: &models.User{
				Login:    "username",
				Password: "",
			},
			want: &models.Session{
				AuthToken:    "auth_token",
				RefreshToken: "refresh_token",
			},
			wantedErr: storage.ErrUserExists,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			mockTokenAuth := new(mocks2.TokenAuth)
			mockStorage := new(mocks2.Storage)
			if tt.setup != nil {
				tt.setup(ctx, mockTokenAuth, mockStorage)
			}
			mockService := ServiceImpl{tokenAuth: mockTokenAuth, storage: mockStorage}
			session, err := mockService.CreateUser(ctx, tt.user)
			if tt.wantedErr != nil {
				require.ErrorIs(t, err, tt.wantedErr)
				require.Nil(t, session)
				return
			}
			require.Equal(t, tt.want.AuthToken, session.AuthToken)
			require.Equal(t, tt.want.RefreshToken, session.RefreshToken)
		})
	}
}

func TestService_LogInUser(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(ctx context.Context, ta *mocks2.TokenAuth, s *mocks2.Storage)
		user      *models.User
		want      *models.Session
		wantedErr error
	}{
		{
			name: "Invalid credentials",
			setup: func(ctx context.Context, ta *mocks2.TokenAuth, s *mocks2.Storage) {
				s.EXPECT().GetIdPassByName(ctx, "username").
					Return(uuid.New(), "password", nil)
			},
			user: &models.User{
				Login:    "username",
				Password: "password",
			},
			wantedErr: ErrInvalidCredentials,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			mockTokenAuth := new(mocks2.TokenAuth)
			mockStorage := new(mocks2.Storage)
			if tt.setup != nil {
				tt.setup(ctx, mockTokenAuth, mockStorage)
			}
			mockService := ServiceImpl{tokenAuth: mockTokenAuth, storage: mockStorage}
			_, err := mockService.LogInUser(ctx, tt.user)
			if tt.wantedErr != nil {
				require.ErrorIs(t, err, tt.wantedErr)
				return
			}
		})
	}
}

func TestService_RefreshUserSession(t *testing.T) {
	id, err := uuid.NewRandom()
	assert.NoError(t, err)
	uid, err := uuid.NewRandom()
	assert.NoError(t, err)
	tests := []struct {
		name      string
		setup     func(ctx context.Context, ta *mocks2.TokenAuth, s *mocks2.Storage)
		session   *models.Session
		want      *models.Session
		wantedErr error
	}{
		{
			name: "Refresh session successfully",
			setup: func(ctx context.Context, ta *mocks2.TokenAuth, s *mocks2.Storage) {
				s.EXPECT().GetSession(ctx, id, "OldRefreshToken").
					Return(uid, nil)
				ta.EXPECT().GenerateTokenPair(uid, mock.AnythingOfType("uuid.UUID")).
					Return("auth_token", "refresh_token", nil)
				s.EXPECT().
					CreateSession(ctx, mock.AnythingOfType("*models.Session"), mock.AnythingOfType("*models.Session")).
					Return(nil)
			},
			session: &models.Session{
				ID:           id,
				RefreshToken: "OldRefreshToken",
			},
			want: &models.Session{
				AuthToken:    "auth_token",
				RefreshToken: "refresh_token",
				UserID:       uid,
			},
			wantedErr: nil,
		},
		{
			name: "Fail to refresh session with session not found",
			setup: func(ctx context.Context, ta *mocks2.TokenAuth, s *mocks2.Storage) {
				s.EXPECT().GetSession(ctx, id, "OldRefreshToken").
					Return(uuid.Nil, storage.ErrSessionNotFound)
			},
			session: &models.Session{
				ID:           id,
				RefreshToken: "OldRefreshToken",
			},
			want:      nil,
			wantedErr: storage.ErrSessionNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			mockTokenAuth := new(mocks2.TokenAuth)
			mockStorage := new(mocks2.Storage)
			if tt.setup != nil {
				tt.setup(ctx, mockTokenAuth, mockStorage)
			}
			mockService := ServiceImpl{tokenAuth: mockTokenAuth, storage: mockStorage}
			session, err := mockService.RefreshUserSession(ctx, tt.session)
			if tt.wantedErr != nil {
				require.ErrorIs(t, err, tt.wantedErr)
				require.Nil(t, session)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want.AuthToken, session.AuthToken)
			require.Equal(t, tt.want.RefreshToken, session.RefreshToken)
			require.Equal(t, tt.want.UserID, session.UserID)
		})
	}
}
