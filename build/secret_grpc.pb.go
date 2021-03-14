// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package genproto

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

// SecretAppClient is the client API for SecretApp service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SecretAppClient interface {
	CreateSecret(ctx context.Context, in *CreateSecretRequest, opts ...grpc.CallOption) (*CreateSecretResponse, error)
	SeeSecret(ctx context.Context, in *SeeSecretRequest, opts ...grpc.CallOption) (*SeeSecretResponse, error)
}

type secretAppClient struct {
	cc grpc.ClientConnInterface
}

func NewSecretAppClient(cc grpc.ClientConnInterface) SecretAppClient {
	return &secretAppClient{cc}
}

func (c *secretAppClient) CreateSecret(ctx context.Context, in *CreateSecretRequest, opts ...grpc.CallOption) (*CreateSecretResponse, error) {
	out := new(CreateSecretResponse)
	err := c.cc.Invoke(ctx, "/sharesecret.SecretApp/CreateSecret", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *secretAppClient) SeeSecret(ctx context.Context, in *SeeSecretRequest, opts ...grpc.CallOption) (*SeeSecretResponse, error) {
	out := new(SeeSecretResponse)
	err := c.cc.Invoke(ctx, "/sharesecret.SecretApp/SeeSecret", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SecretAppServer is the server API for SecretApp service.
// All implementations should embed UnimplementedSecretAppServer
// for forward compatibility
type SecretAppServer interface {
	CreateSecret(context.Context, *CreateSecretRequest) (*CreateSecretResponse, error)
	SeeSecret(context.Context, *SeeSecretRequest) (*SeeSecretResponse, error)
}

// UnimplementedSecretAppServer should be embedded to have forward compatible implementations.
type UnimplementedSecretAppServer struct {
}

func (UnimplementedSecretAppServer) CreateSecret(context.Context, *CreateSecretRequest) (*CreateSecretResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateSecret not implemented")
}
func (UnimplementedSecretAppServer) SeeSecret(context.Context, *SeeSecretRequest) (*SeeSecretResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SeeSecret not implemented")
}

// UnsafeSecretAppServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SecretAppServer will
// result in compilation errors.
type UnsafeSecretAppServer interface {
	mustEmbedUnimplementedSecretAppServer()
}

func RegisterSecretAppServer(s grpc.ServiceRegistrar, srv SecretAppServer) {
	s.RegisterService(&SecretApp_ServiceDesc, srv)
}

func _SecretApp_CreateSecret_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateSecretRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SecretAppServer).CreateSecret(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sharesecret.SecretApp/CreateSecret",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SecretAppServer).CreateSecret(ctx, req.(*CreateSecretRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SecretApp_SeeSecret_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SeeSecretRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SecretAppServer).SeeSecret(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/sharesecret.SecretApp/SeeSecret",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SecretAppServer).SeeSecret(ctx, req.(*SeeSecretRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// SecretApp_ServiceDesc is the grpc.ServiceDesc for SecretApp service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SecretApp_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "sharesecret.SecretApp",
	HandlerType: (*SecretAppServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateSecret",
			Handler:    _SecretApp_CreateSecret_Handler,
		},
		{
			MethodName: "SeeSecret",
			Handler:    _SecretApp_SeeSecret_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "secret.proto",
}
