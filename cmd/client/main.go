package main

import (
	"context"
	"github.com/Mldlr/storety/cmd/client/cmd"
	"github.com/Mldlr/storety/internal/client/config"
	interceptors "github.com/Mldlr/storety/internal/client/interceptor"
	"github.com/Mldlr/storety/internal/client/service"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

func main() {
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
	crypto := service.NewCrypto(cfg)
	userClient := service.NewUserClient(ctx, conn, cfg)
	dataClient := service.NewDataClient(ctx, conn, cfg)
	defer conn.Close()
	cmd.Execute(userClient, dataClient, crypto)
}
