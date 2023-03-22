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

// CreateBatchData creates a batch of data items.
func (s *StoretyHandler) CreateBatchData(ctx context.Context, request *pb.CreateBatchDataRequest) (*pb.CreateBatchResponse, error) {
	session := ctx.Value(models.SessionKey{}).(*models.Session)
	createItems := make([]models.Data, len(request.Data))
	for i, d := range request.Data {
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
	err := s.dataService.CreateBatch(ctx, session.UserID, createItems)
	if err != nil {
		if errors.Is(err, constants.ErrGetData) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.CreateBatchResponse{}, nil
}

// UpdateBatchData updates a batch of data items.
func (s *StoretyHandler) UpdateBatchData(ctx context.Context, request *pb.UpdateBatchDataRequest) (*pb.UpdateBatchResponse, error) {
	session := ctx.Value(models.SessionKey{}).(*models.Session)
	updateItems := make([]models.Data, len(request.Data))
	for i, d := range request.Data {
		id, err := uuid.Parse(d.Id)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		updateItems[i] = models.Data{
			ID:        id,
			Name:      d.Name,
			Type:      d.Type,
			Content:   d.Content,
			UpdatedAt: d.UpdatedAt.AsTime(),
			Deleted:   d.Deleted,
		}
	}
	err := s.dataService.UpdateBatch(ctx, session.UserID, updateItems)
	if err != nil {
		if errors.Is(err, constants.ErrGetData) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.UpdateBatchResponse{}, nil
}

// SyncData accepts data to update on the server and sends updates to user client.
func (s *StoretyHandler) SyncData(ctx context.Context, request *pb.SyncRequest) (*pb.SyncResponse, error) {
	session := ctx.Value(models.SessionKey{}).(*models.Session)
	syncData := make([]models.SyncData, len(request.SyncInfo))
	for i, d := range request.SyncInfo {
		id, err := uuid.Parse(d.Id)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		syncData[i] = models.SyncData{
			ID:        id,
			Hash:      d.Hash,
			UpdatedAt: d.UpdatedAt.AsTime(),
		}
	}
	updates, requestedUpdates, err := s.dataService.GetSyncData(ctx, session.UserID, syncData)
	if err != nil {
		if errors.Is(err, constants.ErrGetData) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	resp := &pb.SyncResponse{
		UpdateData:       make([]*pb.DataItem, len(updates)),
		RequestedUpdates: requestedUpdates,
	}
	for i, v := range updates {
		resp.UpdateData[i] = &pb.DataItem{
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
