package di

import (
	"github.com/Mldlr/storety/internal/server/pkg/token"
	"github.com/Mldlr/storety/internal/server/service/data"
	"github.com/Mldlr/storety/internal/server/service/user"
	"github.com/samber/do"
)

// configureServices configures the services for the Storety server.
func configureServices(i *do.Injector) {
	tokenAuth := token.NewJwtAuth(i)
	do.Provide(
		i,
		func(i *do.Injector) (token.TokenAuth, error) {
			return tokenAuth, nil
		},
	)
	dataService := data.NewService(i)
	do.Provide(
		i,
		func(i *do.Injector) (data.Service, error) {
			return dataService, nil
		},
	)
	userService := user.NewService(i)
	do.Provide(
		i,
		func(i *do.Injector) (user.Service, error) {
			return userService, nil
		},
	)
}
