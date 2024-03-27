// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: pcache.proto

package __

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

// PcacheClient is the client API for Pcache service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PcacheClient interface {
	Get(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error)
}

type pcacheClient struct {
	cc grpc.ClientConnInterface
}

func NewPcacheClient(cc grpc.ClientConnInterface) PcacheClient {
	return &pcacheClient{cc}
}

func (c *pcacheClient) Get(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/pcachepb.Pcache/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PcacheServer is the server API for Pcache service.
// All implementations must embed UnimplementedPcacheServer
// for forward compatibility
type PcacheServer interface {
	Get(context.Context, *Request) (*Response, error)
	mustEmbedUnimplementedPcacheServer()
}

// UnimplementedPcacheServer must be embedded to have forward compatible implementations.
type UnimplementedPcacheServer struct {
}

func (UnimplementedPcacheServer) Get(context.Context, *Request) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedPcacheServer) mustEmbedUnimplementedPcacheServer() {}

// UnsafePcacheServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PcacheServer will
// result in compilation errors.
type UnsafePcacheServer interface {
	mustEmbedUnimplementedPcacheServer()
}

func RegisterPcacheServer(s grpc.ServiceRegistrar, srv PcacheServer) {
	s.RegisterService(&Pcache_ServiceDesc, srv)
}

func _Pcache_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PcacheServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pcachepb.Pcache/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PcacheServer).Get(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

// Pcache_ServiceDesc is the grpc.ServiceDesc for Pcache service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Pcache_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pcachepb.Pcache",
	HandlerType: (*PcacheServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Get",
			Handler:    _Pcache_Get_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pcache.proto",
}
