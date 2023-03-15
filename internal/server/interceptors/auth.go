package interceptors

import (
	"context"
	"errors"
	"fmt"
	"github.com/Mldlr/storety/internal/constants"
	"github.com/Mldlr/storety/internal/server/models"
	"github.com/Mldlr/storety/internal/server/pkg/token"
	"github.com/Mldlr/storety/internal/server/pkg/util/helpers"
	"github.com/golang-jwt/jwt/v5"
	"github.com/samber/do"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// AuthServerInterceptor implements a gRPC server interceptor for authentication.
type AuthServerInterceptor struct {
	tokenAuth         token.TokenAuth
	unprotectedRoutes map[string]struct{}
	refreshRoute      map[string]struct{}
}

// NewAuthInterceptor returns a new authentication interceptor.
func NewAuthInterceptor(i *do.Injector) *AuthServerInterceptor {
	tokenAuth := do.MustInvoke[token.TokenAuth](i)
	return &AuthServerInterceptor{
		tokenAuth: tokenAuth,
		unprotectedRoutes: map[string]struct{}{
			"/proto.User/CreateUser": struct{}{},
			"/proto.User/LogInUser":  struct{}{},
		},
		refreshRoute: map[string]struct{}{
			"/proto.User/RefreshUserSession": struct{}{},
		},
	}
}

// UnaryInterceptor implements the UnaryInterceptor method of the grpc.UnaryServerInterceptor interface.
func (a *AuthServerInterceptor) UnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if _, ok := a.unprotectedRoutes[info.FullMethod]; ok {
		return handler(ctx, req)
	}
	var tokenName string
	if _, ok := a.refreshRoute[info.FullMethod]; ok {
		tokenName = "refresh_token"
	} else {
		tokenName = "auth_token"
	}
	tokenMD, ok := helpers.CheckMDValue(ctx, tokenName)
	if !ok {
		return nil, status.Error(codes.PermissionDenied, fmt.Sprintf("missing %s", tokenName))
	}
	id, err := a.tokenAuth.Verify(tokenMD)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, status.Error(codes.PermissionDenied, constants.ErrExpiredToken.Error())
		}
		return nil, status.Error(codes.PermissionDenied, fmt.Sprintf("invalid %s: %s", tokenName, err.Error()))
	}
	session := &models.Session{}
	if tokenName == "refresh_token" {
		session.RefreshToken = tokenMD
		session.ID = id
	} else {
		session.AuthToken = tokenMD
		session.UserID = id
	}
	ctxNew := context.WithValue(ctx, models.SessionKey{}, session)
	return handler(ctxNew, req)
}
