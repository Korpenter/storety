package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/Mldlr/storety/internal/client/config"
	"github.com/Mldlr/storety/internal/client/pkg/helpers"
	"github.com/Mldlr/storety/internal/client/pkg/utils"
	"github.com/Mldlr/storety/internal/constants"
	pb "github.com/Mldlr/storety/internal/proto"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/pbkdf2"
	"google.golang.org/grpc"
	"log"
	"os"
)

// UserService is the interface for the user service.
type UserService interface {
	CreateUser(username, password string) error
	LogInUser(username, password string) error
	RefreshToken() error
}

// UserServiceImpl is a client for the User service.
type UserServiceImpl struct {
	ctx          context.Context
	conn         *grpc.ClientConn
	remoteClient pb.UserClient
	cfg          *config.Config
}

// NewUserService creates a new UserServiceImpl instance and returns a pointer to it.
// It takes a context, a gRPC client connection, and a configuration object as parameters.
func NewUserService(ctx context.Context, conn *grpc.ClientConn, cfg *config.Config) *UserServiceImpl {
	return &UserServiceImpl{
		ctx:          ctx,
		conn:         conn,
		remoteClient: pb.NewUserClient(conn),
		cfg:          cfg,
	}
}

// CreateUser makes a request to the CreateUser RPC to create a new user and updates the config.
func (c *UserServiceImpl) CreateUser(username, password string) error {
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
	result, err := c.remoteClient.CreateUser(c.ctx, request)
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
	c.cfg.UpdateKey(key)
	err = utils.SaveAuthData(c.cfg.SaltsFile, username, hashedKey, salt, result.AuthToken, result.RefreshToken)
	if err != nil {
		return err
	}

	return nil
}

// LogInUser attempts user authorization on the remote server and if it fails, attempts local authorization.
func (c *UserServiceImpl) LogInUser(username, password string) error {
	switch {
	case c.conn != nil:
		err := c.remoteLogInUser(username, password)
		if err != nil {
			log.Println("Failed to log in on remote server, attempting local login")
		} else {
			log.Println("Successful remote log in")
			return nil
		}
		fallthrough
	default:
		log.Println("attempting local login")
		err := c.localLogin(username, password)
		if err != nil {
			return helpers.LogError(fmt.Errorf("failed local login: %v", err))
		}
		return nil
	}
}

// remoteLogInUser makes a request to the LogInUser RPC to log in a user and updates the config.
func (c *UserServiceImpl) remoteLogInUser(username, password string) error {
	request := &pb.LoginUserRequest{
		Login:    username,
		Password: password,
	}

	result, err := c.remoteClient.LogInUser(c.ctx, request)
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
	c.cfg.UpdateKey(key)
	err = utils.SaveAuthData(c.cfg.SaltsFile, username, hashedKey, salt, result.AuthToken, result.RefreshToken)
	if err != nil {
		return err
	}

	return nil
}

// localLogin makes an attempt to authorize user locally.
func (c *UserServiceImpl) localLogin(username, password string) error {
	hashedKey, salt, authToken, refreshToken, err := utils.GetAuthData(c.cfg.SaltsFile, username)
	if err != nil {
		if errors.Is(err, constants.ErrUserNotFound) {
			return fmt.Errorf("no local data found for %s, login/register remote first", username)
		} else if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("no local authorization data stored, login/register remote first")
		}
		return err
	}
	key := pbkdf2.Key([]byte(password), salt, 10000, 32, sha256.New)
	c.cfg.UpdateKey(key)
	c.cfg.UpdateTokens(authToken, refreshToken)
	err = bcrypt.CompareHashAndPassword(hashedKey, key)
	if err != nil {
		return constants.ErrInvalidCredentials
	}
	return nil
}

// RefreshToken makes a request to the RefreshUserSession RPC to refresh the user's session and updates the config.
func (c *UserServiceImpl) RefreshToken() error {
	request := &pb.RefreshUserSessionRequest{}
	result, err := c.remoteClient.RefreshUserSession(c.ctx, request)
	if err != nil {
		return err
	}
	err = c.cfg.UpdateTokens(result.AuthToken, result.RefreshToken)
	if err != nil {
		return err
	}
	return nil
}
