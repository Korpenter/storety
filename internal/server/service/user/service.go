package user

import (
	"context"
	"errors"
	"github.com/Mldlr/storety/internal/server/models"
	"github.com/Mldlr/storety/internal/server/pkg/token"
	"github.com/Mldlr/storety/internal/server/storage"
	"github.com/google/uuid"
	"github.com/samber/do"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type ServiceImpl struct {
	storage   storage.Storage
	tokenAuth token.TokenAuth
	log       *zap.Logger
}

func NewService(i *do.Injector) *ServiceImpl {
	repo := do.MustInvoke[storage.Storage](i)
	tokenAuth := do.MustInvoke[token.TokenAuth](i)
	log := do.MustInvoke[*zap.Logger](i)
	return &ServiceImpl{
		storage:   repo,
		tokenAuth: tokenAuth,
		log:       log,
	}
}

func (s *ServiceImpl) CreateUser(ctx context.Context, user *models.User) (*models.Session, error) {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		return nil, err
	}
	user.Password = string(hashBytes)
	user.ID, err = uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	err = s.storage.CreateUser(ctx, user)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			return nil, errors.Join(ErrInvalidCredentials, err)
		}
		return nil, err
	}
	session := &models.Session{UserID: user.ID}
	session.ID, err = uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	session.AuthToken, session.RefreshToken, err = s.tokenAuth.GenerateTokenPair(session.UserID, session.ID)
	if err != nil {
		return nil, err
	}
	err = s.storage.CreateSession(ctx, session, nil)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (s *ServiceImpl) LogInUser(ctx context.Context, user *models.User) (*models.Session, error) {
	uid, hash, err := s.storage.GetIdPassByName(ctx, user.Login)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return nil, errors.Join(ErrInvalidCredentials, err)
		}
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(user.Password))
	if err != nil {
		return nil, errors.Join(ErrInvalidCredentials, err)
	}
	session := &models.Session{UserID: uid}
	session.ID, err = uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	session.AuthToken, session.RefreshToken, err = s.tokenAuth.GenerateTokenPair(session.UserID, session.ID)
	if err != nil {
		return nil, err
	}
	err = s.storage.CreateSession(ctx, session, nil)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (s *ServiceImpl) RefreshUserSession(ctx context.Context, oldSession *models.Session) (*models.Session, error) {
	var err error
	session := &models.Session{}
	session.UserID, err = s.storage.GetSession(ctx, oldSession.ID, oldSession.RefreshToken)
	if err != nil {
		if errors.Is(err, storage.ErrSessionNotFound) {
			return nil, errors.Join(ErrInvalidRefreshToken, err)
		}
		return nil, err
	}
	session.ID, err = uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	session.AuthToken, session.RefreshToken, err = s.tokenAuth.GenerateTokenPair(session.UserID, session.ID)
	if err != nil {
		return nil, err
	}
	err = s.storage.CreateSession(ctx, session, oldSession)
	if err != nil {
		return nil, err
	}
	return session, nil
}
