// Package grpcServer provides the gRPC server implementation for the Storety service.
package grpcServer

import (
	"crypto/tls"
	pb "github.com/Mldlr/storety/internal/proto"
	"github.com/Mldlr/storety/internal/server/config"
	"github.com/Mldlr/storety/internal/server/handler"
	"github.com/Mldlr/storety/internal/server/interceptors"
	pkgTls "github.com/Mldlr/storety/internal/server/pkg/tls"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/samber/do"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"os"
	"os/signal"
	"syscall"
)

// GRPCServer is the gRPC server for the Storety service.
type GRPCServer struct {
	srv *grpc.Server
	cfg *config.Config
	log *zap.Logger
}

// NewGRPCServer creates a new GRPCServer with the provided dependency injector.
// It returns a pointer to the newly created GRPCServer.
func NewGRPCServer(i *do.Injector) *GRPCServer {
	cfg := do.MustInvoke[*config.Config](i)
	log := do.MustInvoke[*zap.Logger](i)
	authInterceptor := interceptors.NewAuthInterceptor(i)
	h := handler.NewStoretyHandler(i)

	certFiles := []string{cfg.CertFile, cfg.KeyFile}
	for _, file := range certFiles {
		if _, err := os.Stat(file); err != nil {
			err = pkgTls.GenerateCert(cfg)
			if err != nil {
				log.Fatal("failed to generate certificate", zap.Error(err))
			}
			break
		}
	}

	srv := grpc.NewServer(grpc.ChainUnaryInterceptor(grpc_zap.UnaryServerInterceptor(log), authInterceptor.UnaryInterceptor))
	pb.RegisterDataServer(srv, h)
	pb.RegisterUserServer(srv, h)
	return &GRPCServer{
		srv: srv,
		cfg: cfg,
		log: log,
	}
}

// Run starts the gRPC server and listens for incoming connections.
// It also handles graceful shutdown on receiving termination signals.
func (s *GRPCServer) Run() {
	cert, err := tls.LoadX509KeyPair(s.cfg.CertFile, s.cfg.KeyFile)
	if err != nil {
		s.log.Fatal("failed to load certificate", zap.Error(err))
	}

	listener, err := tls.Listen("tcp", s.cfg.ServiceAddress, &tls.Config{Certificates: []tls.Certificate{cert}})
	if err != nil {
		s.log.Fatal("failed to listen", zap.String("address", s.cfg.ServiceAddress))
	}

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	go func() {
		if err = s.srv.Serve(listener); err != nil {
			s.log.Fatal("failed to serve", zap.Error(err))
		}
	}()
	<-sigint
	s.log.Info("shutting down")
	s.srv.GracefulStop()
	listener.Close()
}
