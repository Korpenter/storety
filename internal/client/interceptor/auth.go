package interceptors

import (
	"context"
	"github.com/Mldlr/storety/internal/client/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type AuthClientInterceptor struct {
	cfg               *config.Config
	unprotectedRoutes map[string]struct{}
	refreshRoute      map[string]struct{}
}

func NewAuthClientInterceptor(cfg *config.Config) *AuthClientInterceptor {
	return &AuthClientInterceptor{
		cfg: cfg,
		unprotectedRoutes: map[string]struct{}{
			"/proto.User/CreateUser": struct{}{},
			"/proto.User/LogInUser":  struct{}{},
		},
		refreshRoute: map[string]struct{}{
			"/proto.User/RefreshUserSession": struct{}{},
		},
	}
}

func (a *AuthClientInterceptor) UnaryInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, callOpts ...grpc.CallOption) error {
	if _, ok := a.unprotectedRoutes[method]; ok {
		return invoker(ctx, method, req, reply, cc, callOpts...)
	}
	var outCtx context.Context
	if _, ok := a.refreshRoute[method]; ok {
		outCtx = metadata.NewOutgoingContext(ctx, metadata.Pairs("refresh_token", a.cfg.JWTRefreshToken))
	} else {
		outCtx = metadata.NewOutgoingContext(ctx, metadata.Pairs("auth_token", a.cfg.JWTAuthToken))
	}
	return invoker(outCtx, method, req, reply, cc, callOpts...)
}
