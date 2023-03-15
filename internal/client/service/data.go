package service

import (
	"context"
	"github.com/Mldlr/storety/internal/client/config"
	pb "github.com/Mldlr/storety/internal/proto"
	"google.golang.org/grpc"
)

type DataClient struct {
	ctx        context.Context
	dataClient pb.DataClient
	cfg        *config.Config
}

func NewDataClient(ctx context.Context, conn *grpc.ClientConn, cfg *config.Config) *DataClient {
	return &DataClient{
		ctx:        ctx,
		dataClient: pb.NewDataClient(conn),
		cfg:        cfg,
	}
}

func (c *DataClient) CreateData(n, t string, content []byte) error {
	request := &pb.CreateDataRequest{
		Name:    n,
		Type:    t,
		Content: content,
	}
	_, err := c.dataClient.CreateData(c.ctx, request)
	if err != nil {
		return err
	}
	return nil
}

func (c *DataClient) ListData() (*pb.ListDataResponse, error) {
	request := &pb.ListDataRequest{}

	resp, err := c.dataClient.ListData(c.ctx, request)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *DataClient) GetData(name string) (*pb.GetContentResponse, error) {
	request := &pb.GetContentRequest{
		Name: name,
	}

	result, err := c.dataClient.GetContent(c.ctx, request)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *DataClient) DeleteData(name string) error {
	request := &pb.DeleteDataRequest{
		Name: name,
	}

	_, err := c.dataClient.DeleteData(c.ctx, request)
	if err != nil {
		return err
	}
	return nil
}
