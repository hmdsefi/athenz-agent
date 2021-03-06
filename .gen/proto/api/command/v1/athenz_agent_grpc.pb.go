// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package v1

import (
	context "context"
	v1 "github.com/hamed-yousefi/athenz-agent/.gen/proto/api/message/v1"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// AthenzAgentClient is the client API for AthenzAgent service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AthenzAgentClient interface {
	CheckAccessWithToken(ctx context.Context, in *v1.AccessCheckRequest, opts ...grpc.CallOption) (*v1.AccessCheckResponse, error)
	GetServiceToken(ctx context.Context, in *v1.ServiceTokenRequest, opts ...grpc.CallOption) (*v1.ServiceTokenResponse, error)
}

type athenzAgentClient struct {
	cc grpc.ClientConnInterface
}

func NewAthenzAgentClient(cc grpc.ClientConnInterface) AthenzAgentClient {
	return &athenzAgentClient{cc}
}

func (c *athenzAgentClient) CheckAccessWithToken(ctx context.Context, in *v1.AccessCheckRequest, opts ...grpc.CallOption) (*v1.AccessCheckResponse, error) {
	out := new(v1.AccessCheckResponse)
	err := c.cc.Invoke(ctx, "/athenz.agent.api.command.v1.AthenzAgent/CheckAccessWithToken", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *athenzAgentClient) GetServiceToken(ctx context.Context, in *v1.ServiceTokenRequest, opts ...grpc.CallOption) (*v1.ServiceTokenResponse, error) {
	out := new(v1.ServiceTokenResponse)
	err := c.cc.Invoke(ctx, "/athenz.agent.api.command.v1.AthenzAgent/GetServiceToken", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AthenzAgentServer is the server API for AthenzAgent service.
// All implementations should embed UnimplementedAthenzAgentServer
// for forward compatibility
type AthenzAgentServer interface {
	CheckAccessWithToken(context.Context, *v1.AccessCheckRequest) (*v1.AccessCheckResponse, error)
	GetServiceToken(context.Context, *v1.ServiceTokenRequest) (*v1.ServiceTokenResponse, error)
}

// UnimplementedAthenzAgentServer should be embedded to have forward compatible implementations.
type UnimplementedAthenzAgentServer struct {
}

func (UnimplementedAthenzAgentServer) CheckAccessWithToken(context.Context, *v1.AccessCheckRequest) (*v1.AccessCheckResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckAccessWithToken not implemented")
}
func (UnimplementedAthenzAgentServer) GetServiceToken(context.Context, *v1.ServiceTokenRequest) (*v1.ServiceTokenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetServiceToken not implemented")
}

// UnsafeAthenzAgentServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AthenzAgentServer will
// result in compilation errors.
type UnsafeAthenzAgentServer interface {
	mustEmbedUnimplementedAthenzAgentServer()
}

func RegisterAthenzAgentServer(s grpc.ServiceRegistrar, srv AthenzAgentServer) {
	s.RegisterService(&AthenzAgent_ServiceDesc, srv)
}

func _AthenzAgent_CheckAccessWithToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(v1.AccessCheckRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AthenzAgentServer).CheckAccessWithToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/athenz.agent.api.command.v1.AthenzAgent/CheckAccessWithToken",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AthenzAgentServer).CheckAccessWithToken(ctx, req.(*v1.AccessCheckRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AthenzAgent_GetServiceToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(v1.ServiceTokenRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AthenzAgentServer).GetServiceToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/athenz.agent.api.command.v1.AthenzAgent/GetServiceToken",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AthenzAgentServer).GetServiceToken(ctx, req.(*v1.ServiceTokenRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AthenzAgent_ServiceDesc is the grpc.ServiceDesc for AthenzAgent service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AthenzAgent_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "athenz.agent.api.command.v1.AthenzAgent",
	HandlerType: (*AthenzAgentServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CheckAccessWithToken",
			Handler:    _AthenzAgent_CheckAccessWithToken_Handler,
		},
		{
			MethodName: "GetServiceToken",
			Handler:    _AthenzAgent_GetServiceToken_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/athenz/agent/api/command/v1/athenz_agent.proto",
}
