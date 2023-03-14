package handler

import (
	"context"
	"errors"
	pb "github.com/Mldlr/storety/internal/proto"
	"github.com/Mldlr/storety/internal/server/models"
	"github.com/Mldlr/storety/internal/server/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *StoretyHandler) CreateData(ctx context.Context, request *pb.CreateDataRequest) (*pb.CreateDataResponse, error) {
	session := ctx.Value(models.SessionKey{}).(*models.Session)
	in := &models.Data{
		UserID:  session.UserID,
		Name:    request.Name,
		Type:    request.Type,
		Content: request.Content,
	}
	err := s.dataService.CreateData(ctx, in)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.CreateDataResponse{}, nil
}

func (s *StoretyHandler) GetContent(ctx context.Context, request *pb.GetContentRequest) (*pb.GetContentResponse, error) {
	session := ctx.Value(models.SessionKey{}).(*models.Session)
	content, err := s.dataService.GetDataContent(ctx, session.UserID, request.Name)
	if err != nil {
		if errors.Is(err, storage.ErrGettingData) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.GetContentResponse{Content: content}, nil
}

func (s *StoretyHandler) DeleteData(ctx context.Context, request *pb.DeleteDataRequest) (*pb.DeleteDataResponse, error) {
	session := ctx.Value(models.SessionKey{}).(*models.Session)
	err := s.dataService.DeleteData(ctx, session.UserID, request.Name)
	if err != nil {
		if errors.Is(err, storage.ErrDeletingData) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.DeleteDataResponse{}, nil
}

func (s *StoretyHandler) ListData(ctx context.Context, request *pb.ListDataRequest) (*pb.ListDataResponse, error) {
	session := ctx.Value(models.SessionKey{}).(*models.Session)
	list, err := s.dataService.ListData(ctx, session.UserID)
	if err != nil {
		if errors.Is(err, storage.ErrNoData) {
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
