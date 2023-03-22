package user

import (
	"context"
	"errors"
	"github.com/Mldlr/storety/internal/client/config"
	"github.com/Mldlr/storety/internal/client/mocks"
	"github.com/Mldlr/storety/internal/client/pkg/utils"
	pb "github.com/Mldlr/storety/internal/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"os"
	"testing"
)

func TestCreateUser(t *testing.T) {
	ctx := context.Background()
	remoteClientMock := new(mocks.UserClient)

	saltsFile, err := os.CreateTemp("", "salts-test.json")
	assert.NoError(t, err)

	cfg := &config.Config{
		SaltsFile: saltsFile.Name(),
	}

	userService := &ServiceImpl{
		ctx:          ctx,
		remoteClient: remoteClientMock,
		cfg:          cfg,
	}
	username := "testuser"
	password := "testpassword"
	remoteClientMock.On("CreateUser", ctx, mock.AnythingOfType("*proto.CreateUserRequest")).
		Return(&pb.CreateUserResponse{
			AuthToken:    "test-auth-token",
			RefreshToken: "test-refresh-token",
		}, nil)

	err = userService.CreateUser(username, password)
	assert.NoError(t, err)
	_, _, authToken, refreshToken, err := utils.GetAuthData(saltsFile.Name(), username)
	assert.NoError(t, err)
	assert.Equal(t, "test-auth-token", authToken)
	assert.Equal(t, "test-refresh-token", refreshToken)
	remoteClientMock.AssertCalled(t, "CreateUser", ctx, mock.AnythingOfType("*proto.CreateUserRequest"))
	remoteClientMock.AssertNumberOfCalls(t, "CreateUser", 1)
}

func TestRefreshToken(t *testing.T) {
	ctx := context.Background()
	remoteClientMock := new(mocks.UserClient)

	cfg := &config.Config{
		JWTAuthToken:    "old-auth-token",
		JWTRefreshToken: "old-refresh-token",
	}

	service := ServiceImpl{
		ctx:          ctx,
		remoteClient: remoteClientMock,
		cfg:          cfg,
	}

	remoteClientMock.On("RefreshUserSession", ctx, mock.AnythingOfType("*proto.RefreshUserSessionRequest")).Return(&pb.RefreshUserSessionResponse{
		AuthToken:    "new-auth-token",
		RefreshToken: "new-refresh-token",
	}, nil)

	err := service.RefreshToken()
	assert.NoError(t, err)

	remoteClientMock.AssertCalled(t, "RefreshUserSession", ctx, mock.AnythingOfType("*proto.RefreshUserSessionRequest"))
	remoteClientMock.AssertNumberOfCalls(t, "RefreshUserSession", 1)

	assert.Equal(t, "new-auth-token", cfg.JWTAuthToken)
	assert.Equal(t, "new-refresh-token", cfg.JWTRefreshToken)
}

func TestLogInUser(t *testing.T) {
	ctx := context.Background()
	remoteClientMock := new(mocks.UserClient)

	saltsFile, err := os.CreateTemp("", "salts-test.json")
	assert.NoError(t, err)

	cfg := &config.Config{
		JWTAuthToken:    "old-auth-token",
		JWTRefreshToken: "old-refresh-token",
		SaltsFile:       saltsFile.Name(),
	}

	service := ServiceImpl{
		ctx:          ctx,
		remoteClient: remoteClientMock,
		cfg:          cfg,
	}

	tests := []struct {
		name                 string
		remoteClientResponse *pb.LoginUserResponse
		remoteClientError    error
		expectedError        error
	}{
		{
			name: "Successful remote login",
			remoteClientResponse: &pb.LoginUserResponse{
				AuthToken:    "new-auth-token",
				RefreshToken: "new-refresh-token",
				Salt:         "salt",
			},
			remoteClientError: nil,
			expectedError:     nil,
		},
		{
			name:                 "Failed remote login, successful local login",
			remoteClientResponse: nil,
			remoteClientError:    errors.New("random remote error"),
			expectedError:        nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			remoteClientMock.On("LogInUser", ctx, mock.AnythingOfType("*proto.LoginUserRequest")).Return(tt.remoteClientResponse, tt.remoteClientError)

			err := service.LogInUser("username", "password")
			assert.Equal(t, tt.expectedError, err)
			if tt.expectedError == nil {
				_, _, authToken, refreshToken, err := utils.GetAuthData(saltsFile.Name(), "username")
				assert.NoError(t, err)
				assert.Equal(t, "new-auth-token", authToken)
				assert.Equal(t, "new-refresh-token", refreshToken)

				remoteClientMock.AssertCalled(t, "LogInUser", ctx, mock.AnythingOfType("*proto.LoginUserRequest"))
				remoteClientMock.AssertNumberOfCalls(t, "LogInUser", 1)
				remoteClientMock.ExpectedCalls = []*mock.Call{}
				remoteClientMock.Calls = []mock.Call{}
			}
		})
	}
}
