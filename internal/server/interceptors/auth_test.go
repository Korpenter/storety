package interceptors

import (
	"context"
	"github.com/Mldlr/storety/internal/constants"
	pb "github.com/Mldlr/storety/internal/proto"
	"github.com/Mldlr/storety/internal/server/mocks"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"testing"
)

func TestAuthInterceptor_UnaryInterceptor_Unprotected(t *testing.T) {
	tests := []struct {
		name         string
		setup        func(ctx context.Context, ta *mocks.TokenAuth)
		req          interface{}
		ctx          context.Context
		unaryInfo    *grpc.UnaryServerInfo
		wantedErrMsg string
		errCode      codes.Code
	}{
		{
			name: "Unprotected route request",
			unaryInfo: &grpc.UnaryServerInfo{
				FullMethod: "/proto.User/CreateUser",
			},
			ctx: context.Background(),
			req: &pb.CreateUserRequest{
				Login:    "username",
				Password: "password",
			},
			errCode: codes.OK,
		},
		{
			name: "Refresh route successful request",
			setup: func(ctx context.Context, ta *mocks.TokenAuth) {
				ta.EXPECT().Verify("refreshToken").
					Return(uuid.New(), nil)
			},
			req: &pb.RefreshUserSessionRequest{},
			ctx: metadata.NewIncomingContext(context.Background(),
				metadata.New(map[string]string{"refresh_token": "refreshToken"})),
			unaryInfo: &grpc.UnaryServerInfo{
				FullMethod: "/proto.User/RefreshUserSession",
			},
			errCode: codes.OK,
		},
		{
			name: "Refresh route unsuccessful request with no token",

			req: &pb.RefreshUserSessionRequest{},
			ctx: context.Background(),
			unaryInfo: &grpc.UnaryServerInfo{
				FullMethod: "/proto.User/RefreshUserSession",
			},
			wantedErrMsg: "missing refresh_token",
			errCode:      codes.PermissionDenied,
		},
		{
			name: "Refresh route unsuccessful request with expired token",
			setup: func(ctx context.Context, ta *mocks.TokenAuth) {
				ta.EXPECT().Verify("expiredToken").
					Return(uuid.Nil, jwt.ErrTokenExpired)
			},
			req: &pb.RefreshUserSessionRequest{},
			ctx: metadata.NewIncomingContext(context.Background(),
				metadata.New(map[string]string{"refresh_token": "expiredToken"})),
			unaryInfo: &grpc.UnaryServerInfo{
				FullMethod: "/proto.User/RefreshUserSession",
			},
			wantedErrMsg: constants.ErrExpiredToken.Error(),
			errCode:      codes.PermissionDenied,
		},
	}

	unaryHandler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, nil
	}
	mockAuth := new(mocks.TokenAuth)
	interceptor := AuthServerInterceptor{
		tokenAuth: mockAuth,
		unprotectedRoutes: map[string]struct{}{
			"/proto.User/CreateUser": struct{}{},
			"/proto.User/LogInUser":  struct{}{},
		},
		refreshRoute: map[string]struct{}{
			"/proto.User/RefreshUserSession": struct{}{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.setup != nil {
				tt.setup(tt.ctx, mockAuth)
			}
			_, err := interceptor.UnaryInterceptor(tt.ctx, tt.req, tt.unaryInfo, unaryHandler)
			if statusErr, ok := status.FromError(err); ok {
				require.Equal(t, tt.errCode.String(), statusErr.Code().String())
				require.Equal(t, statusErr.Message(), tt.wantedErrMsg)
			}
		})
	}
}
