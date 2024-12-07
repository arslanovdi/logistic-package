// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: logistic_package_api.proto

package logistic_package_api_v1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	LogisticPackageApiService_CreateV1_FullMethodName = "/logistic_package_api.v1.LogisticPackageApiService/CreateV1"
	LogisticPackageApiService_DeleteV1_FullMethodName = "/logistic_package_api.v1.LogisticPackageApiService/DeleteV1"
	LogisticPackageApiService_GetV1_FullMethodName    = "/logistic_package_api.v1.LogisticPackageApiService/GetV1"
	LogisticPackageApiService_ListV1_FullMethodName   = "/logistic_package_api.v1.LogisticPackageApiService/ListV1"
	LogisticPackageApiService_UpdateV1_FullMethodName = "/logistic_package_api.v1.LogisticPackageApiService/UpdateV1"
)

// LogisticPackageApiServiceClient is the client API for LogisticPackageApiService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LogisticPackageApiServiceClient interface {
	CreateV1(ctx context.Context, in *CreateRequestV1, opts ...grpc.CallOption) (*CreateResponseV1, error)
	DeleteV1(ctx context.Context, in *DeleteV1Request, opts ...grpc.CallOption) (*DeleteV1Response, error)
	GetV1(ctx context.Context, in *GetV1Request, opts ...grpc.CallOption) (*GetV1Response, error)
	ListV1(ctx context.Context, in *ListV1Request, opts ...grpc.CallOption) (*ListV1Response, error)
	UpdateV1(ctx context.Context, in *UpdateV1Request, opts ...grpc.CallOption) (*UpdateV1Response, error)
}

type logisticPackageApiServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewLogisticPackageApiServiceClient(cc grpc.ClientConnInterface) LogisticPackageApiServiceClient {
	return &logisticPackageApiServiceClient{cc}
}

func (c *logisticPackageApiServiceClient) CreateV1(ctx context.Context, in *CreateRequestV1, opts ...grpc.CallOption) (*CreateResponseV1, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateResponseV1)
	err := c.cc.Invoke(ctx, LogisticPackageApiService_CreateV1_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *logisticPackageApiServiceClient) DeleteV1(ctx context.Context, in *DeleteV1Request, opts ...grpc.CallOption) (*DeleteV1Response, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteV1Response)
	err := c.cc.Invoke(ctx, LogisticPackageApiService_DeleteV1_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *logisticPackageApiServiceClient) GetV1(ctx context.Context, in *GetV1Request, opts ...grpc.CallOption) (*GetV1Response, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetV1Response)
	err := c.cc.Invoke(ctx, LogisticPackageApiService_GetV1_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *logisticPackageApiServiceClient) ListV1(ctx context.Context, in *ListV1Request, opts ...grpc.CallOption) (*ListV1Response, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListV1Response)
	err := c.cc.Invoke(ctx, LogisticPackageApiService_ListV1_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *logisticPackageApiServiceClient) UpdateV1(ctx context.Context, in *UpdateV1Request, opts ...grpc.CallOption) (*UpdateV1Response, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateV1Response)
	err := c.cc.Invoke(ctx, LogisticPackageApiService_UpdateV1_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LogisticPackageApiServiceServer is the server API for LogisticPackageApiService service.
// All implementations must embed UnimplementedLogisticPackageApiServiceServer
// for forward compatibility.
type LogisticPackageApiServiceServer interface {
	CreateV1(context.Context, *CreateRequestV1) (*CreateResponseV1, error)
	DeleteV1(context.Context, *DeleteV1Request) (*DeleteV1Response, error)
	GetV1(context.Context, *GetV1Request) (*GetV1Response, error)
	ListV1(context.Context, *ListV1Request) (*ListV1Response, error)
	UpdateV1(context.Context, *UpdateV1Request) (*UpdateV1Response, error)
	mustEmbedUnimplementedLogisticPackageApiServiceServer()
}

// UnimplementedLogisticPackageApiServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedLogisticPackageApiServiceServer struct{}

func (UnimplementedLogisticPackageApiServiceServer) CreateV1(context.Context, *CreateRequestV1) (*CreateResponseV1, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateV1 not implemented")
}
func (UnimplementedLogisticPackageApiServiceServer) DeleteV1(context.Context, *DeleteV1Request) (*DeleteV1Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteV1 not implemented")
}
func (UnimplementedLogisticPackageApiServiceServer) GetV1(context.Context, *GetV1Request) (*GetV1Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetV1 not implemented")
}
func (UnimplementedLogisticPackageApiServiceServer) ListV1(context.Context, *ListV1Request) (*ListV1Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListV1 not implemented")
}
func (UnimplementedLogisticPackageApiServiceServer) UpdateV1(context.Context, *UpdateV1Request) (*UpdateV1Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateV1 not implemented")
}
func (UnimplementedLogisticPackageApiServiceServer) mustEmbedUnimplementedLogisticPackageApiServiceServer() {
}
func (UnimplementedLogisticPackageApiServiceServer) testEmbeddedByValue() {}

// UnsafeLogisticPackageApiServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LogisticPackageApiServiceServer will
// result in compilation errors.
type UnsafeLogisticPackageApiServiceServer interface {
	mustEmbedUnimplementedLogisticPackageApiServiceServer()
}

func RegisterLogisticPackageApiServiceServer(s grpc.ServiceRegistrar, srv LogisticPackageApiServiceServer) {
	// If the following call pancis, it indicates UnimplementedLogisticPackageApiServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&LogisticPackageApiService_ServiceDesc, srv)
}

func _LogisticPackageApiService_CreateV1_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateRequestV1)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LogisticPackageApiServiceServer).CreateV1(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LogisticPackageApiService_CreateV1_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogisticPackageApiServiceServer).CreateV1(ctx, req.(*CreateRequestV1))
	}
	return interceptor(ctx, in, info, handler)
}

func _LogisticPackageApiService_DeleteV1_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteV1Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LogisticPackageApiServiceServer).DeleteV1(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LogisticPackageApiService_DeleteV1_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogisticPackageApiServiceServer).DeleteV1(ctx, req.(*DeleteV1Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _LogisticPackageApiService_GetV1_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetV1Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LogisticPackageApiServiceServer).GetV1(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LogisticPackageApiService_GetV1_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogisticPackageApiServiceServer).GetV1(ctx, req.(*GetV1Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _LogisticPackageApiService_ListV1_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListV1Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LogisticPackageApiServiceServer).ListV1(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LogisticPackageApiService_ListV1_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogisticPackageApiServiceServer).ListV1(ctx, req.(*ListV1Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _LogisticPackageApiService_UpdateV1_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateV1Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LogisticPackageApiServiceServer).UpdateV1(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LogisticPackageApiService_UpdateV1_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogisticPackageApiServiceServer).UpdateV1(ctx, req.(*UpdateV1Request))
	}
	return interceptor(ctx, in, info, handler)
}

// LogisticPackageApiService_ServiceDesc is the grpc.ServiceDesc for LogisticPackageApiService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var LogisticPackageApiService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "logistic_package_api.v1.LogisticPackageApiService",
	HandlerType: (*LogisticPackageApiServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateV1",
			Handler:    _LogisticPackageApiService_CreateV1_Handler,
		},
		{
			MethodName: "DeleteV1",
			Handler:    _LogisticPackageApiService_DeleteV1_Handler,
		},
		{
			MethodName: "GetV1",
			Handler:    _LogisticPackageApiService_GetV1_Handler,
		},
		{
			MethodName: "ListV1",
			Handler:    _LogisticPackageApiService_ListV1_Handler,
		},
		{
			MethodName: "UpdateV1",
			Handler:    _LogisticPackageApiService_UpdateV1_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "logistic_package_api.proto",
}