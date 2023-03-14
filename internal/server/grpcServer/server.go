package grpcServer

import (
	pb "github.com/Mldlr/storety/internal/proto"
	"github.com/Mldlr/storety/internal/server/config"
	"github.com/Mldlr/storety/internal/server/handler"
	"github.com/Mldlr/storety/internal/server/interceptors"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/samber/do"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"syscall"
)

type GRPCServer struct {
	srv *grpc.Server
	cfg *config.Config
	log *zap.Logger
}

func NewGRPCServer(i *do.Injector) *GRPCServer {
	cfg := do.MustInvoke[*config.Config](i)
	log := do.MustInvoke[*zap.Logger](i)
	authInterceptor := interceptors.NewAuthInterceptor(i)
	h := handler.NewStoretyHandler(i)
	srv := grpc.NewServer(grpc.ChainUnaryInterceptor(grpc_zap.UnaryServerInterceptor(log), authInterceptor.UnaryInterceptor))
	pb.RegisterDataServer(srv, h)
	pb.RegisterUserServer(srv, h)
	return &GRPCServer{
		srv: srv,
		cfg: cfg,
		log: log,
	}
}

func (s *GRPCServer) Run() {
	listener, err := net.Listen("tcp", s.cfg.ServiceAddress)
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
