package handler

import (
	"context"
	"errors"
	"github.com/Mldlr/storety/internal/constants"
	pb "github.com/Mldlr/storety/internal/proto"
	"github.com/Mldlr/storety/internal/server/mocks"
	"github.com/Mldlr/storety/internal/server/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(ctx context.Context, us *mocks.UserService)
		req     *pb.CreateUserRequest
		want    *pb.CreateUserResponse
		errCode codes.Code
	}{
		{
			name: "Create user successfully",
			setup: func(ctx context.Context, us *mocks.UserService) {
				us.EXPECT().CreateUser(mock.AnythingOfType("*context.emptyCtx"), &models.User{
					Login:    "username",
					Password: "password",
				}).Return(&models.Session{AuthToken: "auth_token", RefreshToken: "refresh_token"}, nil)
			},
			req: &pb.CreateUserRequest{
				Login:    "username",
				Password: "password",
			},
			want: &pb.CreateUserResponse{
				AuthToken:    "auth_token",
				RefreshToken: "refresh_token",
			},
			errCode: codes.OK,
		},
		{
			name: "Fail to create duplicate user",
			setup: func(ctx context.Context, us *mocks.UserService) {
				us.EXPECT().CreateUser(mock.AnythingOfType("*context.emptyCtx"), &models.User{
					Login:    "username",
					Password: "password",
				}).Return(nil, errors.Join(constants.ErrInvalidCredentials, constants.ErrUserExists))
			},
			req: &pb.CreateUserRequest{
				Login:    "username",
				Password: "password",
			},
			want:    nil,
			errCode: codes.AlreadyExists,
		},
		{
			name: "Fail to create user with invalid credentials",
			req: &pb.CreateUserRequest{
				Login:    "username",
				Password: "",
			},
			want:    nil,
			errCode: codes.InvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var resp *pb.CreateUserResponse
			var err error

			ctx := context.Background()
			mockUserSrv := new(mocks.UserService)
			if tt.setup != nil {
				tt.setup(ctx, mockUserSrv)
			}
			mockDep := StoretyHandler{userService: mockUserSrv}
			resp, err = mockDep.CreateUser(ctx, tt.req)
			require.EqualValues(t, tt.want, resp)
			if statusErr, ok := status.FromError(err); ok {
				require.Equal(t, tt.errCode.String(), statusErr.Code().String())
			}
		})
	}
}

func TestLoginUser(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(ctx context.Context, us *mocks.UserService)
		req     *pb.LoginUserRequest
		want    *pb.LoginUserResponse
		errCode codes.Code
	}{
		{
			name: "Login user successfully",
			setup: func(ctx context.Context, us *mocks.UserService) {
				us.EXPECT().LogInUser(mock.AnythingOfType("*context.emptyCtx"), &models.User{
					Login:    "username",
					Password: "password",
				}).Return(&models.Session{AuthToken: "auth_token", RefreshToken: "refresh_token"}, "salt", nil)
			},
			req: &pb.LoginUserRequest{
				Login:    "username",
				Password: "password",
			},
			want: &pb.LoginUserResponse{
				AuthToken:    "auth_token",
				RefreshToken: "refresh_token",
				Salt:         "salt",
			},
			errCode: codes.OK,
		},
		{
			name: "Attempt login with non existent username",
			setup: func(ctx context.Context, us *mocks.UserService) {
				us.EXPECT().LogInUser(mock.AnythingOfType("*context.emptyCtx"), &models.User{
					Login:    "username",
					Password: "password",
				}).Return(nil, "", errors.Join(constants.ErrInvalidCredentials, constants.ErrUserNotFound))
			},
			req: &pb.LoginUserRequest{
				Login:    "username",
				Password: "password",
			},
			want:    nil,
			errCode: codes.InvalidArgument,
		},
		{
			name: "Fail to login user with invalid credentials",
			req: &pb.LoginUserRequest{
				Login:    "username",
				Password: "",
			},
			want:    nil,
			errCode: codes.InvalidArgument,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var resp *pb.LoginUserResponse
			var err error

			ctx := context.Background()
			mockUserSrv := new(mocks.UserService)
			if tt.setup != nil {
				tt.setup(ctx, mockUserSrv)
			}
			mockDep := StoretyHandler{userService: mockUserSrv}
			resp, err = mockDep.LogInUser(ctx, tt.req)
			require.EqualValues(t, tt.want, resp)
			if statusErr, ok := status.FromError(err); ok {
				require.Equal(t, tt.errCode.String(), statusErr.Code().String())
			}
		})
	}
}

func TestRefreshUserSession(t *testing.T) {
	id, err := uuid.NewRandom()
	assert.NoError(t, err)
	tests := []struct {
		name    string
		setup   func(ctx context.Context, us *mocks.UserService)
		req     *pb.RefreshUserSessionRequest
		want    *pb.RefreshUserSessionResponse
		errCode codes.Code
	}{
		{
			name: "Refresh successfully",
			setup: func(ctx context.Context, us *mocks.UserService) {
				us.EXPECT().RefreshUserSession(mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("*models.Session")).
					Return(&models.Session{AuthToken: "AuthNew", RefreshToken: "refreshNew"}, nil)
			},
			req: &pb.RefreshUserSessionRequest{},
			want: &pb.RefreshUserSessionResponse{
				AuthToken:    "AuthNew",
				RefreshToken: "refreshNew",
			},
			errCode: codes.OK,
		},
		{
			name: "Refresh unsuccessfully with invalid refresh token",
			setup: func(ctx context.Context, us *mocks.UserService) {
				us.EXPECT().RefreshUserSession(mock.AnythingOfType("*context.valueCtx"), mock.AnythingOfType("*models.Session")).
					Return(nil, constants.ErrInvalidRefreshToken)
			},
			req:     &pb.RefreshUserSessionRequest{},
			want:    nil,
			errCode: codes.PermissionDenied,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var resp *pb.RefreshUserSessionResponse
			var err error

			ctx := context.Background()
			mockUserSrv := new(mocks.UserService)
			if tt.setup != nil {
				tt.setup(ctx, mockUserSrv)
			}
			incCtx := context.WithValue(ctx,
				models.SessionKey{},
				&models.Session{ID: id, RefreshToken: "refresh_token"},
			)
			mockHandler := StoretyHandler{userService: mockUserSrv}
			resp, err = mockHandler.RefreshUserSession(incCtx, tt.req)
			require.EqualValues(t, tt.want, resp)
			if statusErr, ok := status.FromError(err); ok {
				require.Equal(t, tt.errCode.String(), statusErr.Code().String())
			}
		})
	}
}
