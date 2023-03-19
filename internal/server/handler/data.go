// Package handler provides the main Storety gRPC server handler.
package handler

import (
	"context"
	"errors"
	"github.com/Mldlr/storety/internal/constants"
	pb "github.com/Mldlr/storety/internal/proto"
	"github.com/Mldlr/storety/internal/server/models"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// CreateData creates a new data item and stores it.
func (s *StoretyHandler) CreateData(ctx context.Context, request *pb.CreateDataRequest) (*pb.CreateDataResponse, error) {
	session := ctx.Value(models.SessionKey{}).(*models.Session)
	in := &models.Data{
		Name:    request.Data.Name,
		Type:    request.Data.Type,
		Content: request.Data.Content,
	}
	err := s.dataService.CreateData(ctx, session.UserID, in)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.CreateDataResponse{}, nil
}

// GetContent retrieves the content of a data item.
func (s *StoretyHandler) GetContent(ctx context.Context, request *pb.GetContentRequest) (*pb.GetContentResponse, error) {
	session := ctx.Value(models.SessionKey{}).(*models.Session)
	content, contentType, err := s.dataService.GetDataContent(ctx, session.UserID, request.Name)
	if err != nil {
		if errors.Is(err, constants.ErrGetData) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.GetContentResponse{Content: content, Type: contentType}, nil
}

// DeleteData removes a data item.
func (s *StoretyHandler) DeleteData(ctx context.Context, request *pb.DeleteDataRequest) (*pb.DeleteDataResponse, error) {
	session := ctx.Value(models.SessionKey{}).(*models.Session)
	err := s.dataService.DeleteData(ctx, session.UserID, request.Name)
	if err != nil {
		if errors.Is(err, constants.ErrDeleteData) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.DeleteDataResponse{}, nil
}

// ListData returns a list of all data items for a user.
func (s *StoretyHandler) ListData(ctx context.Context, request *pb.ListDataRequest) (*pb.ListDataResponse, error) {
	session := ctx.Value(models.SessionKey{}).(*models.Session)
	list, err := s.dataService.ListData(ctx, session.UserID)
	if err != nil {
		if errors.Is(err, constants.ErrNoData) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	var response []*pb.DataInfo
	for _, data := range list {
		item := &pb.DataInfo{
			Name: data.Name,
			Type: data.Type,
		}
		response = append(response, item)
	}
	return &pb.ListDataResponse{Data: response}, nil
}

// SyncData accepts data to update on the server and sends updates to user client.
func (s *StoretyHandler) SyncData(ctx context.Context, request *pb.SyncRequest) (*pb.SyncResponse, error) {
	session := ctx.Value(models.SessionKey{}).(*models.Session)
	createItems := make([]models.Data, len(request.CreateData))
	deleteItems := make([]models.Data, len(request.DeleteData))
	for i, d := range request.CreateData {
		id, err := uuid.Parse(d.Id)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		createItems[i] = models.Data{
			ID:        id,
			Name:      d.Name,
			Type:      d.Type,
			Content:   d.Content,
			UpdatedAt: d.UpdatedAt.AsTime(),
			Deleted:   d.Deleted,
		}
	}
	for i, v := range request.DeleteData {
		id, err := uuid.Parse(v.Id)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		deleteItems[i] = models.Data{
			ID:        id,
			Name:      v.Name,
			Type:      v.Type,
			Content:   v.Content,
			UpdatedAt: v.UpdatedAt.AsTime(),
			Deleted:   v.Deleted,
		}
	}
	syncData := models.SyncData{
		CreateData: createItems,
		DeleteData: deleteItems,
		LastSync:   request.LastSync.AsTime(),
	}
	updates, err := s.dataService.SyncData(ctx, session.UserID, syncData)
	if err != nil {
		if errors.Is(err, constants.ErrGetData) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	resp := &pb.SyncResponse{Data: make([]*pb.DataItem, len(updates))}
	for i, v := range updates {
		resp.Data[i] = &pb.DataItem{
			Id:        v.ID.String(),
			Name:      v.Name,
			Type:      v.Type,
			Content:   v.Content,
			UpdatedAt: timestamppb.New(v.UpdatedAt),
			Deleted:   v.Deleted,
		}
	}
	return resp, nil
}
