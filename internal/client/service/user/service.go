package user

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
	"github.com/samber/do"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/pbkdf2"
	"google.golang.org/grpc"
	"log"
	"os"
)

// ServiceImpl is a client for the User service.
type ServiceImpl struct {
	ctx          context.Context
	conn         *grpc.ClientConn
	remoteClient pb.UserClient
	cfg          *config.Config
}

// NewServiceImpl creates a new ServiceImpl instance and returns a pointer to it.
// It takes a context, a gRPC client connection, and a configuration object as parameters.
func NewServiceImpl(i *do.Injector) *ServiceImpl {
	conn := do.MustInvoke[*grpc.ClientConn](i)
	cfg := do.MustInvoke[*config.Config](i)
	return &ServiceImpl{
		ctx:          context.Background(),
		conn:         conn,
		remoteClient: pb.NewUserClient(conn),
		cfg:          cfg,
	}
}

// CreateUser implements the CreateUser method of the Service interface.
func (c *ServiceImpl) CreateUser(username, password string) error {
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
	c.cfg.UpdateTokens(result.AuthToken, result.RefreshToken)
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

// LogInUser implements the LogInUser method of the Service interface.
func (c *ServiceImpl) LogInUser(username, password string) error {
	err := c.remoteLogInUser(username, password)
	if err != nil {
		log.Println("Failed to log in on remote server, attempting local login")
		fmt.Println(err)
	} else {
		log.Println("Successful remote log in")
		return nil
	}
	log.Println("attempting local login")
	err = c.localLogin(username, password)
	if err != nil {
		return fmt.Errorf("failed local login: %v", err)
	}
	return nil
}

// remoteLogInUser makes a request to the LogInUser RPC to log in a user and updates the config.
func (c *ServiceImpl) remoteLogInUser(username, password string) error {
	request := &pb.LoginUserRequest{
		Login:    username,
		Password: password,
	}
	result, err := c.remoteClient.LogInUser(c.ctx, request)
	if err != nil {
		return err
	}
	c.cfg.UpdateTokens(result.AuthToken, result.RefreshToken)
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
	fmt.Println("salt", salt)
	err = utils.SaveAuthData(c.cfg.SaltsFile, username, hashedKey, salt, result.AuthToken, result.RefreshToken)
	if err != nil {
		return err
	}

	return nil
}

// localLogin makes an attempt to authorize user locally.
func (c *ServiceImpl) localLogin(username, password string) error {
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

// RefreshToken implements the Service interface method RefreshToken.
func (c *ServiceImpl) RefreshToken() error {
	request := &pb.RefreshUserSessionRequest{}
	result, err := c.remoteClient.RefreshUserSession(c.ctx, request)
	if err != nil {
		return fmt.Errorf("failed to refresh token: %v", err)
	}
	c.cfg.UpdateTokens(result.AuthToken, result.RefreshToken)
	return nil
}
