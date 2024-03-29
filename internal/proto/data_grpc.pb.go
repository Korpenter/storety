// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.6
// source: data.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// DataClient is the client API for Data service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
//go:generate mockery --name=DataClient -r --case underscore --with-expecter --structname DataClient --filename data_client.go
type DataClient interface {
	CreateData(ctx context.Context, in *CreateDataRequest, opts ...grpc.CallOption) (*CreateDataResponse, error)
	GetContent(ctx context.Context, in *GetContentRequest, opts ...grpc.CallOption) (*GetContentResponse, error)
	ListData(ctx context.Context, in *ListDataRequest, opts ...grpc.CallOption) (*ListDataResponse, error)
	DeleteData(ctx context.Context, in *DeleteDataRequest, opts ...grpc.CallOption) (*DeleteDataResponse, error)
	CreateBatchData(ctx context.Context, in *CreateBatchDataRequest, opts ...grpc.CallOption) (*CreateBatchResponse, error)
	UpdateBatchData(ctx context.Context, in *UpdateBatchDataRequest, opts ...grpc.CallOption) (*UpdateBatchResponse, error)
	SyncData(ctx context.Context, in *SyncRequest, opts ...grpc.CallOption) (*SyncResponse, error)
}

type dataClient struct {
	cc grpc.ClientConnInterface
}

func NewDataClient(cc grpc.ClientConnInterface) DataClient {
	return &dataClient{cc}
}

func (c *dataClient) CreateData(ctx context.Context, in *CreateDataRequest, opts ...grpc.CallOption) (*CreateDataResponse, error) {
	out := new(CreateDataResponse)
	err := c.cc.Invoke(ctx, "/proto.Data/CreateData", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dataClient) GetContent(ctx context.Context, in *GetContentRequest, opts ...grpc.CallOption) (*GetContentResponse, error) {
	out := new(GetContentResponse)
	err := c.cc.Invoke(ctx, "/proto.Data/GetContent", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dataClient) ListData(ctx context.Context, in *ListDataRequest, opts ...grpc.CallOption) (*ListDataResponse, error) {
	out := new(ListDataResponse)
	err := c.cc.Invoke(ctx, "/proto.Data/ListData", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dataClient) DeleteData(ctx context.Context, in *DeleteDataRequest, opts ...grpc.CallOption) (*DeleteDataResponse, error) {
	out := new(DeleteDataResponse)
	err := c.cc.Invoke(ctx, "/proto.Data/DeleteData", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dataClient) CreateBatchData(ctx context.Context, in *CreateBatchDataRequest, opts ...grpc.CallOption) (*CreateBatchResponse, error) {
	out := new(CreateBatchResponse)
	err := c.cc.Invoke(ctx, "/proto.Data/CreateBatchData", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dataClient) UpdateBatchData(ctx context.Context, in *UpdateBatchDataRequest, opts ...grpc.CallOption) (*UpdateBatchResponse, error) {
	out := new(UpdateBatchResponse)
	err := c.cc.Invoke(ctx, "/proto.Data/UpdateBatchData", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dataClient) SyncData(ctx context.Context, in *SyncRequest, opts ...grpc.CallOption) (*SyncResponse, error) {
	out := new(SyncResponse)
	err := c.cc.Invoke(ctx, "/proto.Data/SyncData", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DataServer is the server API for Data service.
// All implementations must embed UnimplementedDataServer
// for forward compatibility
type DataServer interface {
	CreateData(context.Context, *CreateDataRequest) (*CreateDataResponse, error)
	GetContent(context.Context, *GetContentRequest) (*GetContentResponse, error)
	ListData(context.Context, *ListDataRequest) (*ListDataResponse, error)
	DeleteData(context.Context, *DeleteDataRequest) (*DeleteDataResponse, error)
	CreateBatchData(context.Context, *CreateBatchDataRequest) (*CreateBatchResponse, error)
	UpdateBatchData(context.Context, *UpdateBatchDataRequest) (*UpdateBatchResponse, error)
	SyncData(context.Context, *SyncRequest) (*SyncResponse, error)
	mustEmbedUnimplementedDataServer()
}

// UnimplementedDataServer must be embedded to have forward compatible implementations.
type UnimplementedDataServer struct {
}

func (UnimplementedDataServer) CreateData(context.Context, *CreateDataRequest) (*CreateDataResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateData not implemented")
}
func (UnimplementedDataServer) GetContent(context.Context, *GetContentRequest) (*GetContentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetContent not implemented")
}
func (UnimplementedDataServer) ListData(context.Context, *ListDataRequest) (*ListDataResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListData not implemented")
}
func (UnimplementedDataServer) DeleteData(context.Context, *DeleteDataRequest) (*DeleteDataResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteData not implemented")
}
func (UnimplementedDataServer) CreateBatchData(context.Context, *CreateBatchDataRequest) (*CreateBatchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateBatchData not implemented")
}
func (UnimplementedDataServer) UpdateBatchData(context.Context, *UpdateBatchDataRequest) (*UpdateBatchResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateBatchData not implemented")
}
func (UnimplementedDataServer) SyncData(context.Context, *SyncRequest) (*SyncResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SyncData not implemented")
}
func (UnimplementedDataServer) mustEmbedUnimplementedDataServer() {}

// UnsafeDataServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DataServer will
// result in compilation errors.
type UnsafeDataServer interface {
	mustEmbedUnimplementedDataServer()
}

func RegisterDataServer(s grpc.ServiceRegistrar, srv DataServer) {
	s.RegisterService(&Data_ServiceDesc, srv)
}

func _Data_CreateData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateDataRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataServer).CreateData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Data/CreateData",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataServer).CreateData(ctx, req.(*CreateDataRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Data_GetContent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetContentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataServer).GetContent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Data/GetContent",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataServer).GetContent(ctx, req.(*GetContentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Data_ListData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListDataRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataServer).ListData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Data/ListData",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataServer).ListData(ctx, req.(*ListDataRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Data_DeleteData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteDataRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataServer).DeleteData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Data/DeleteData",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataServer).DeleteData(ctx, req.(*DeleteDataRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Data_CreateBatchData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateBatchDataRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataServer).CreateBatchData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Data/CreateBatchData",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataServer).CreateBatchData(ctx, req.(*CreateBatchDataRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Data_UpdateBatchData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateBatchDataRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataServer).UpdateBatchData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Data/UpdateBatchData",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataServer).UpdateBatchData(ctx, req.(*UpdateBatchDataRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Data_SyncData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SyncRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataServer).SyncData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.Data/SyncData",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataServer).SyncData(ctx, req.(*SyncRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Data_ServiceDesc is the grpc.ServiceDesc for Data service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Data_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Data",
	HandlerType: (*DataServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateData",
			Handler:    _Data_CreateData_Handler,
		},
		{
			MethodName: "GetContent",
			Handler:    _Data_GetContent_Handler,
		},
		{
			MethodName: "ListData",
			Handler:    _Data_ListData_Handler,
		},
		{
			MethodName: "DeleteData",
			Handler:    _Data_DeleteData_Handler,
		},
		{
			MethodName: "CreateBatchData",
			Handler:    _Data_CreateBatchData_Handler,
		},
		{
			MethodName: "UpdateBatchData",
			Handler:    _Data_UpdateBatchData_Handler,
		},
		{
			MethodName: "SyncData",
			Handler:    _Data_SyncData_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "data.proto",
}
