// Package handler provides the main Storety gRPC server handler.
package handler

import (
	pb "github.com/Mldlr/storety/internal/proto"
	"github.com/Mldlr/storety/internal/server/config"
	"github.com/Mldlr/storety/internal/server/service/data"
	"github.com/Mldlr/storety/internal/server/service/user"
	"github.com/samber/do"
	"go.uber.org/zap"
)

// StoretyHandler is the handler for the Storety gRPC server.
// It embeds both UnimplementedDataServer and UnimplementedUserServer interfaces.
type StoretyHandler struct {
	pb.UnimplementedDataServer
	pb.UnimplementedUserServer
	userService user.Service
	dataService data.Service
	cfg         *config.Config
	log         *zap.Logger
}

// NewStoretyHandler creates a new StoretyHandler with the provided dependencies.
// It takes a *do.Injector as input and initializes the required dependencies.
func NewStoretyHandler(i *do.Injector) *StoretyHandler {
	userService := do.MustInvoke[user.Service](i)
	dataService := do.MustInvoke[data.Service](i)
	cfg := do.MustInvoke[*config.Config](i)
	logger := do.MustInvoke[*zap.Logger](i)
	return &StoretyHandler{
		userService: userService,
		dataService: dataService,
		cfg:         cfg,
		log:         logger,
	}
}
