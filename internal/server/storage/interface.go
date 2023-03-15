package storage

import (
	"context"
	"github.com/Mldlr/storety/internal/server/models"
	"github.com/google/uuid"
)

//go:generate mockery --name=Storage -r --case underscore --with-expecter --structname Storage --filename storage.go
type Storage interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetIdPassByName(ctx context.Context, username string) (uuid.UUID, string, error)
	GetSession(ctx context.Context, sessionID uuid.UUID, refreshToken string) (uuid.UUID, error)
	CreateSession(ctx context.Context, session, oldSession *models.Session) error
	CreateData(ctx context.Context, data *models.Data) error
	GetDataContentByName(ctx context.Context, userID uuid.UUID, name string) ([]byte, string, error)
	GetAllDataInfo(ctx context.Context, userID uuid.UUID) ([]models.DataInfo, error)
	DeleteDataByName(ctx context.Context, userID uuid.UUID, name string) error
}
