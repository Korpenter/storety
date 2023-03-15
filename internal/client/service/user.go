package service

import (
	"context"
	"github.com/Mldlr/storety/internal/client/config"
	pb "github.com/Mldlr/storety/internal/proto"
	"google.golang.org/grpc"
)

type UserClient struct {
	ctx        context.Context
	userClient pb.UserClient
	cfg        *config.Config
}

func NewUserClient(ctx context.Context, conn *grpc.ClientConn, cfg *config.Config) *UserClient {
	return &UserClient{
		ctx:        ctx,
		userClient: pb.NewUserClient(conn),
		cfg:        cfg,
	}
}

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
