package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/Mldlr/storety/cmd/client/cmd"
	"github.com/Mldlr/storety/internal/client/config"
	interceptors "github.com/Mldlr/storety/internal/client/interceptor"
	"github.com/Mldlr/storety/internal/client/service"
	"github.com/Mldlr/storety/internal/client/service/crypto"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/samber/do"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
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

	injector := do.New()
	cfg := config.NewConfig()

	do.Provide(
		injector,
		func(i *do.Injector) (*config.Config, error) {
			return cfg, nil
		},
	)

	ctx := context.Background()
	cert, err := tls.LoadX509KeyPair(cfg.CertFile, cfg.KeyFile)
	if err != nil {
		log.Fatal("Failed to load certificate:", err)
	}
	cred := credentials.NewTLS(&tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
	})
	keepaliveParams := keepalive.ClientParameters{
		Time:                10 * time.Second,
		Timeout:             2 * time.Second,
		PermitWithoutStream: true,
	}
	authInterceptor := interceptors.NewAuthClientInterceptor(cfg)
	retryInterceptor := interceptors.NewRetryClientInterceptor(cfg, 10, 5*time.Second)
	opts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithTransportCredentials(cred),
		grpc.WithKeepaliveParams(keepaliveParams),
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(
			authInterceptor.UnaryInterceptor,
			retryInterceptor.UnaryInterceptor,
		)),
	}
	var conn *grpc.ClientConn
	conn, err = grpc.DialContext(ctx, cfg.ServiceAddress, opts...)
	if err != nil {
		log.Println("Failed to connect to server:", err)
	}
	do.Provide(
		injector,
		func(i *do.Injector) (*grpc.ClientConn, error) {
			return conn, nil
		},
	)
	cryptoSvc := crypto.NewCrypto(cfg)
	userService := service.NewUserService(ctx, conn, cfg)
	dataService := service.NewDataService(ctx, conn, cfg)
	do.Provide(
		injector,
		func(i *do.Injector) (crypto.Crypto, error) {
			return *cryptoSvc, nil
		},
	)
	do.Provide(
		injector,
		func(i *do.Injector) (service.UserService, error) {
			return userService, nil
		},
	)
	do.Provide(
		injector,
		func(i *do.Injector) (service.DataService, error) {
			return dataService, nil
		},
	)

	defer conn.Close()
	cmd.Execute(injector)
}
