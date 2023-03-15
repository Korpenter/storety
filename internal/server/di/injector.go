package di

import (
	"github.com/Mldlr/storety/internal/server/config"
	"github.com/samber/do"
	"go.uber.org/zap"
)

// ConfigureDependencies configures the dependencies for the Storety server.
func ConfigureDependencies(cfg *config.Config, log *zap.Logger) *do.Injector {
	injector := do.New()

	do.Provide(
		injector,
		func(i *do.Injector) (*config.Config, error) {
			return cfg, nil
		},
	)

	do.Provide(
		injector,
		func(i *do.Injector) (*zap.Logger, error) {
			return log, nil
		},
	)

	configureStorage(injector)
	configureServices(injector)

	return injector
}
