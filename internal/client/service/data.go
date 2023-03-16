package service

import (
	"context"
	"github.com/Mldlr/storety/internal/client/config"
	pb "github.com/Mldlr/storety/internal/proto"
	"google.golang.org/grpc"
)

// DataClient is a client for the Data service.
type DataClient struct {
	ctx        context.Context
	dataClient pb.DataClient
	cfg        *config.Config
}

// NewDataClient creates a new DataClient instance and returns a pointer to it.
// It takes a context, a gRPC client connection, and a configuration object as parameters.
func NewDataClient(ctx context.Context, conn *grpc.ClientConn, cfg *config.Config) *DataClient {
	return &DataClient{
		ctx:        ctx,
		dataClient: pb.NewDataClient(conn),
		cfg:        cfg,
	}
}

// CreateData makes a request to the CreateData RPC to create a new data.
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

// ListData makes a request to the ListData RPC to list all data.
func (c *DataClient) ListData() (*pb.ListDataResponse, error) {
	request := &pb.ListDataRequest{}

	resp, err := c.dataClient.ListData(c.ctx, request)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetData makes a request to the GetData RPC to get a data entry.
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

// DeleteData makes a request to the DeleteData RPC to delete a data entry.
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
