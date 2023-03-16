package main

import (
	"context"
	"fmt"
	"github.com/Mldlr/storety/cmd/client/cmd"
	"github.com/Mldlr/storety/internal/client/config"
	interceptors "github.com/Mldlr/storety/internal/client/interceptor"
	"github.com/Mldlr/storety/internal/client/service"
	"github.com/Mldlr/storety/internal/client/service/crypto"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
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

	ctx := context.Background()
	cfg := config.NewConfig()
	authInterceptor := interceptors.NewAuthClientInterceptor(cfg)
	retryInterceptor := interceptors.NewRetryClientInterceptor(cfg, 10, 5*time.Second)
	opts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(
			authInterceptor.UnaryInterceptor,
			retryInterceptor.UnaryInterceptor,
		)),
	}
	conn, err := grpc.Dial(cfg.ServiceAddress, opts...)
	if err != nil {
		log.Fatal("Failed to dial server:", err)
	}
	crypto := crypto.NewCrypto(cfg)
	userClient := service.NewUserClient(ctx, conn, cfg)
	dataClient := service.NewDataClient(ctx, conn, cfg)
	defer conn.Close()
	cmd.Execute(userClient, dataClient, crypto)
}
