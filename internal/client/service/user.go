package service

import (
	"context"
	"github.com/Mldlr/storety/internal/client/config"
	pb "github.com/Mldlr/storety/internal/proto"
	"google.golang.org/grpc"
)

// UserClient is a client for the User service.
type UserClient struct {
	ctx        context.Context
	userClient pb.UserClient
	cfg        *config.Config
}

// NewUserClient creates a new UserClient instance and returns a pointer to it.
// It takes a context, a gRPC client connection, and a configuration object as parameters.
func NewUserClient(ctx context.Context, conn *grpc.ClientConn, cfg *config.Config) *UserClient {
	return &UserClient{
		ctx:        ctx,
		userClient: pb.NewUserClient(conn),
		cfg:        cfg,
	}
}

// CreateUser makes a request to the CreateUser RPC to create a new user and updates the config.
func (c *UserClient) CreateUser(username, password string) error {
	request := &pb.CreateUserRequest{
		Login:    username,
		Password: password,
	}
	result, err := c.userClient.CreateUser(c.ctx, request)
	if err != nil {
		return err
	}
	err = c.cfg.UpdateTokens(result.AuthToken, result.RefreshToken)
	if err != nil {
		return err
	}
	err = c.cfg.UpdateKey(password)
	if err != nil {
		return err
	}
	return nil
}

// LogInUser makes a request to the LogInUser RPC to log in a user and updates the config.
func (c *UserClient) LogInUser(username, password string) error {
	request := &pb.LoginUserRequest{
		Login:    username,
		Password: password,
	}
	result, err := c.userClient.LogInUser(c.ctx, request)
	if err != nil {
		return err
	}
	err = c.cfg.UpdateTokens(result.AuthToken, result.RefreshToken)
	if err != nil {
		return err
	}
	err = c.cfg.UpdateKey(password)
	if err != nil {
		return err
	}
	return nil
}

// RefreshToken makes a request to the RefreshUserSession RPC to refresh the user's session and updates the config.
func (c *UserClient) RefreshToken() error {
	request := &pb.RefreshUserSessionRequest{}
	result, err := c.userClient.RefreshUserSession(c.ctx, request)
	if err != nil {
		return err
	}
	err = c.cfg.UpdateTokens(result.AuthToken, result.RefreshToken)
	if err != nil {
		return err
	}
	return nil
}
