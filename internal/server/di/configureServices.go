package di

import (
	"github.com/Mldlr/storety/internal/server/pkg/token"
	data2 "github.com/Mldlr/storety/internal/server/service/data"
	user2 "github.com/Mldlr/storety/internal/server/service/user"
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
	dataService := data2.NewService(i)
	do.Provide(
		i,
		func(i *do.Injector) (data2.Service, error) {
			return dataService, nil
		},
	)
	userService := user2.NewService(i)
	do.Provide(
		i,
		func(i *do.Injector) (user2.Service, error) {
			return userService, nil
		},
	)
}
