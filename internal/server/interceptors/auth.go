package interceptors

import (
	"context"
	"errors"
	"fmt"
	"github.com/Mldlr/storety/internal/server/models"
	"github.com/Mldlr/storety/internal/server/pkg/token"
	"github.com/Mldlr/storety/internal/server/pkg/util/helpers"
	"github.com/golang-jwt/jwt/v5"
	"github.com/samber/do"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Auth is an authentication interceptor.
type AuthInterceptor struct {
	tokenAuth         token.TokenAuth
	unprotectedRoutes map[string]struct{}
	refreshRoute      string
}

func NewAuthInterceptor(i *do.Injector) *AuthInterceptor {
	tokenAuth := do.MustInvoke[token.TokenAuth](i)
	return &AuthInterceptor{
		tokenAuth: tokenAuth,
		unprotectedRoutes: map[string]struct{}{
			"/proto.User/CreateUser": struct{}{},
			"/proto.User/LogInUser":  struct{}{},
		},
		refreshRoute: "/proto.User/RefreshUserSession",
	}
}

func (a *AuthInterceptor) UnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if _, ok := a.unprotectedRoutes[info.FullMethod]; ok {
		return handler(ctx, req)
	}
	var tokenName string
	if info.FullMethod == a.refreshRoute {
		tokenName = "refresh_token"
	} else {
		tokenName = "auth_token"
	}
	tokenMD, ok := helpers.CheckMDValue(ctx, tokenName)
	if !ok {
		return nil, status.Error(codes.PermissionDenied, fmt.Sprintf("missing %s", tokenName))
	}
	id, err := a.tokenAuth.Verify(tokenMD)
	fmt.Println(err)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, status.Error(codes.PermissionDenied, jwt.ErrTokenExpired.Error())
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
