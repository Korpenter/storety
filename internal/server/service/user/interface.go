package user

import (
	"context"
	"github.com/Mldlr/storety/internal/server/models"
)

//go:generate mockery --name=Service -r --case underscore --with-expecter --structname UserService --filename user_service.go
type Service interface {
	CreateUser(ctx context.Context, user *models.User) (*models.Session, error)
	LogInUser(ctx context.Context, user *models.User) (*models.Session, error)
	RefreshUserSession(ctx context.Context, oldSession *models.Session) (*models.Session, error)
}
