package interceptors

import (
	"context"
	"github.com/Mldlr/storety/internal/client/config"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"testing"
)

func TestAuthClientInterceptor(t *testing.T) {
	cfg := &config.Config{
		JWTAuthToken:    "test_auth_token",
		JWTRefreshToken: "test_refresh_token",
	}
	tests := []struct {
		name    string
		method  string
		token   string
		invoker func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error
	}{
		{
			name:   "Unprotected method",
			method: "/proto.User/CreateUser",
			token:  "",
			invoker: func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
				return nil
			},
		},
		{
			name:   "Protected method",
			method: "/proto.User/SomeProtectedMethod",
			token:  "test_auth_token",
			invoker: func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
				md, _ := metadata.FromOutgoingContext(ctx)
				assert.Equal(t, []string{cfg.JWTAuthToken}, md["auth_token"])
				return nil
			},
		},
		{
			name:   "Refresh method",
			method: "/proto.User/RefreshUserSession",
			token:  "test_refresh_token",
			invoker: func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, opts ...grpc.CallOption) error {
				md, _ := metadata.FromOutgoingContext(ctx)
				assert.Equal(t, []string{cfg.JWTRefreshToken}, md["refresh_token"])
				return nil
			},
		},
	}

	interceptor := NewAuthClientInterceptor(cfg)
	cc := &grpc.ClientConn{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := interceptor.UnaryInterceptor(context.Background(), tt.method, nil, nil, cc, tt.invoker, nil)
			assert.NoError(t, err)
		})
	}
}
