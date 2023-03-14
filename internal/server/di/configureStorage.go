package di

import (
	"context"
	"github.com/Mldlr/storety/internal/server/config"
	"github.com/Mldlr/storety/internal/server/migration"
	"github.com/Mldlr/storety/internal/server/storage"
	"github.com/Mldlr/storety/internal/server/storage/postgres"
	"github.com/samber/do"
	"go.uber.org/zap"
)

func configureStorage(i *do.Injector) {
	log := do.MustInvoke[*zap.Logger](i)
	cfg := do.MustInvoke[*config.Config](i)
	if cfg.PostgresURI != "" {
		d, err := postgres.NewDB(cfg.PostgresURI)
		if err != nil {
			log.Fatal("Error initiating postgres connection", zap.Error(err))
		}

		err = d.Ping(context.Background())
		if err != nil {
			log.Fatal("Error reaching db", zap.Error(err))
		}

		err = migration.RunMigrations(cfg.PostgresURI)
		if err != nil {
			log.Fatal("Error running migration", zap.Error(err))
		}

		do.Provide(
			i,
			func(i *do.Injector) (storage.Storage, error) {
				return d, nil
			},
		)
		return
	}

	log.Fatal("configuring storage")
}
