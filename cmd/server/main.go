package main

import (
	"fmt"
	"github.com/Mldlr/storety/internal/server/config"
	"github.com/Mldlr/storety/internal/server/di"
	"github.com/Mldlr/storety/internal/server/grpcServer"
	"go.uber.org/zap"
)

// Build info
var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

// NA is the string output if build info is not set
const NA string = "N/A"

func main() {
	if len(buildVersion) == 0 {
		buildVersion = NA
	}
	if len(buildDate) == 0 {
		buildDate = NA
	}
	if len(buildCommit) == 0 {
		buildCommit = NA
	}
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)
	cfg := config.NewConfig()
	log, _ := zap.NewProduction()
	injector := di.ConfigureDependencies(cfg, log)
	log.Info("starting with cfg:",
		zap.String("Storety Address:", cfg.ServiceAddress),
	)
	srv := grpcServer.NewGRPCServer(injector)
	srv.Run()
}
