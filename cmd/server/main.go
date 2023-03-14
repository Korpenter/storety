package main

import (
	"github.com/Mldlr/storety/internal/server/config"
	"github.com/Mldlr/storety/internal/server/di"
	"github.com/Mldlr/storety/internal/server/grpcServer"
	"go.uber.org/zap"
)

func main() {
	cfg := config.NewConfig()
	log, _ := zap.NewProduction()
	injector := di.ConfigureDependencies(cfg, log)
	log.Info("starting with cfg:",
		zap.String("Storety Address:", cfg.ServiceAddress),
	)
	srv := grpcServer.NewGRPCServer(injector)
	srv.Run()
}
