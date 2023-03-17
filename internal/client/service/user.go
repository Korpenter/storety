package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/Mldlr/storety/internal/client/config"
	"github.com/Mldlr/storety/internal/client/pkg/utils"
	"github.com/Mldlr/storety/internal/constants"
	pb "github.com/Mldlr/storety/internal/proto"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/pbkdf2"
	"google.golang.org/grpc"
	"os"
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
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return err
	}
	request := &pb.CreateUserRequest{
		Login:    username,
		Password: password,
		Salt:     base64.StdEncoding.EncodeToString(salt),
	}
	fmt.Println(request)
	result, err := c.userClient.CreateUser(c.ctx, request)
	if err != nil {
		return err
	}
	err = c.cfg.UpdateTokens(result.AuthToken, result.RefreshToken)
	if err != nil {
		return err
	}
	key := pbkdf2.Key([]byte(password), salt, 10000, 32, sha256.New)
	hashedKey, err := bcrypt.GenerateFromPassword(key, 14)
	if err != nil {
		return err
	}
	err = c.cfg.UpdateKey(key)
	if err != nil {
		return err
	}
	err = utils.SaveHashedKeyAndSalt(c.cfg.SaltsFile, username, hashedKey, salt)
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
	salt, err := base64.StdEncoding.DecodeString(result.Salt)
	if err != nil {
		return err
	}
	key := pbkdf2.Key([]byte(password), salt, 10000, 32, sha256.New)
	hashedKey, err := bcrypt.GenerateFromPassword(key, 14)
	if err != nil {
		return err
	}
	err = c.cfg.UpdateKey(key)
	if err != nil {
		return err
	}
	err = utils.SaveHashedKeyAndSalt(c.cfg.SaltsFile, username, hashedKey, salt)
	if err != nil {
		return err
	}
	return nil
}

// LocalLogin makes an attempt to authorize user locally.
func (c *UserClient) LocalLogin(username, password string) error {
	hashedKey, salt, err := utils.GetHashedKeyAndSalt(c.cfg.SaltsFile, username)
	if err != nil {
		if errors.Is(err, constants.ErrUserNotFound) {
			return fmt.Errorf("no local data found for %s, login/register remote first", username)
		} else if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("no local authorization data stored, login/register remote first")
		}
		return err
	}
	key := pbkdf2.Key([]byte(password), salt, 10000, 32, sha256.New)
	err = c.cfg.UpdateKey(key)
	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword(hashedKey, key)
	if err != nil {
		return constants.ErrInvalidCredentials
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
