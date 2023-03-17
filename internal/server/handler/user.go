// Package handler provides the main Storety gRPC server handler.
package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/Mldlr/storety/internal/constants"
	pb "github.com/Mldlr/storety/internal/proto"
	"github.com/Mldlr/storety/internal/server/models"
	"github.com/Mldlr/storety/internal/server/pkg/util/validators"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CreateUser creates a new user account.
func (s *StoretyHandler) CreateUser(ctx context.Context, request *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	in := &models.User{
		Login:    request.Login,
		Password: request.Password,
		Salt:     request.Salt,
	}
	if err := validators.ValidateAuthorization(in); err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("%v: %v", constants.ErrInvalidCredentials, err))
	}
	session, err := s.userService.CreateUser(ctx, in)
	if err != nil {
		if errors.Is(err, constants.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, fmt.Sprintf("%v: %v", constants.ErrInvalidCredentials, err))
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.CreateUserResponse{AuthToken: session.AuthToken, RefreshToken: session.RefreshToken}, nil
}

// LogInUser authenticates a user and logs them in.
func (s *StoretyHandler) LogInUser(ctx context.Context, request *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	in := &models.User{
		Login:    request.Login,
		Password: request.Password,
	}
	if err := validators.ValidateAuthorization(in); err != nil {
		return nil, errors.Join(constants.ErrInvalidCredentials, err)
	}
	session, salt, err := s.userService.LogInUser(ctx, in)
	if err != nil {
		if errors.Is(err, constants.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("%v: %v", constants.ErrInvalidCredentials, err))
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.LoginUserResponse{AuthToken: session.AuthToken, RefreshToken: session.RefreshToken, Salt: salt}, nil
}

// RefreshUserSession refreshes the user's authentication and refresh tokens.
func (s *StoretyHandler) RefreshUserSession(ctx context.Context, request *pb.RefreshUserSessionRequest) (*pb.RefreshUserSessionResponse, error) {
	session := ctx.Value(models.SessionKey{}).(*models.Session)
	session, err := s.userService.RefreshUserSession(ctx, session)
	if err != nil {
		if errors.Is(err, constants.ErrInvalidRefreshToken) {
			return nil, status.Error(codes.PermissionDenied, fmt.Sprintf("%v: %v", constants.ErrInvalidRefreshToken, err))
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.RefreshUserSessionResponse{AuthToken: session.AuthToken, RefreshToken: session.RefreshToken}, nil
}
