package handler

import (
	"context"
	"errors"
	"fmt"
	pb "github.com/Mldlr/storety/internal/proto"
	"github.com/Mldlr/storety/internal/server/models"
	"github.com/Mldlr/storety/internal/server/pkg/util/validators"
	"github.com/Mldlr/storety/internal/server/service/user"
	"github.com/Mldlr/storety/internal/server/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *StoretyHandler) CreateUser(ctx context.Context, request *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	in := &models.User{
		Login:    request.Login,
		Password: request.Password,
	}
	if err := validators.ValidateAuthorization(in); err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("%v: %v", user.ErrInvalidCredentials, err))
	}
	session, err := s.userService.CreateUser(ctx, in)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, fmt.Sprintf("%v: %v", user.ErrInvalidCredentials, err))
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.CreateUserResponse{AuthToken: session.AuthToken, RefreshToken: session.RefreshToken}, nil
}

func (s *StoretyHandler) LogInUser(ctx context.Context, request *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	in := &models.User{
		Login:    request.Login,
		Password: request.Password,
	}
	if err := validators.ValidateAuthorization(in); err != nil {
		return nil, errors.Join(user.ErrInvalidCredentials, err)
	}
	session, err := s.userService.LogInUser(ctx, in)
	if err != nil {
		if errors.Is(err, user.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("%v: %v", user.ErrInvalidCredentials, err))
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.LoginUserResponse{AuthToken: session.AuthToken, RefreshToken: session.RefreshToken}, nil
}

func (s *StoretyHandler) RefreshUserSession(ctx context.Context, request *pb.RefreshUserSessionRequest) (*pb.RefreshUserSessionResponse, error) {
	session := ctx.Value(models.SessionKey{}).(*models.Session)
	session, err := s.userService.RefreshUserSession(ctx, session)
	if err != nil {
		if errors.Is(err, user.ErrInvalidRefreshToken) {
			return nil, status.Error(codes.PermissionDenied, fmt.Sprintf("%v: %v", user.ErrInvalidRefreshToken, err))
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.RefreshUserSessionResponse{AuthToken: session.AuthToken, RefreshToken: session.RefreshToken}, nil
}
