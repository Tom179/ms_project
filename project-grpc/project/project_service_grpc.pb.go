// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.22.3
// source: project_service.proto

package project_service_v1

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

// ProjectServiceClient is the client API for ProjectService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ProjectServiceClient interface {
	Index(ctx context.Context, in *IndexMessage, opts ...grpc.CallOption) (*IndexResponse, error)
	FindProjectByMemId(ctx context.Context, in *ProjectRpcMessage, opts ...grpc.CallOption) (*MyProjectResponse, error)
	FindProjectTemplate(ctx context.Context, in *ProjectRpcMessage, opts ...grpc.CallOption) (*ProjectTemplateResponse, error)
}

type projectServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewProjectServiceClient(cc grpc.ClientConnInterface) ProjectServiceClient {
	return &projectServiceClient{cc}
}

func (c *projectServiceClient) Index(ctx context.Context, in *IndexMessage, opts ...grpc.CallOption) (*IndexResponse, error) {
	out := new(IndexResponse)
	err := c.cc.Invoke(ctx, "/project.service.v1.ProjectService/Index", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *projectServiceClient) FindProjectByMemId(ctx context.Context, in *ProjectRpcMessage, opts ...grpc.CallOption) (*MyProjectResponse, error) {
	out := new(MyProjectResponse)
	err := c.cc.Invoke(ctx, "/project.service.v1.ProjectService/FindProjectByMemId", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *projectServiceClient) FindProjectTemplate(ctx context.Context, in *ProjectRpcMessage, opts ...grpc.CallOption) (*ProjectTemplateResponse, error) {
	out := new(ProjectTemplateResponse)
	err := c.cc.Invoke(ctx, "/project.service.v1.ProjectService/FindProjectTemplate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ProjectServiceServer is the server API for ProjectService service.
// All implementations must embed UnimplementedProjectServiceServer
// for forward compatibility
type ProjectServiceServer interface {
	Index(context.Context, *IndexMessage) (*IndexResponse, error)
	FindProjectByMemId(context.Context, *ProjectRpcMessage) (*MyProjectResponse, error)
	FindProjectTemplate(context.Context, *ProjectRpcMessage) (*ProjectTemplateResponse, error)
	mustEmbedUnimplementedProjectServiceServer()
}

// UnimplementedProjectServiceServer must be embedded to have forward compatible implementations.
type UnimplementedProjectServiceServer struct {
}

func (UnimplementedProjectServiceServer) Index(context.Context, *IndexMessage) (*IndexResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Index not implemented")
}
func (UnimplementedProjectServiceServer) FindProjectByMemId(context.Context, *ProjectRpcMessage) (*MyProjectResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindProjectByMemId not implemented")
}
func (UnimplementedProjectServiceServer) FindProjectTemplate(context.Context, *ProjectRpcMessage) (*ProjectTemplateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindProjectTemplate not implemented")
}
func (UnimplementedProjectServiceServer) mustEmbedUnimplementedProjectServiceServer() {}

// UnsafeProjectServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ProjectServiceServer will
// result in compilation errors.
type UnsafeProjectServiceServer interface {
	mustEmbedUnimplementedProjectServiceServer()
}

func RegisterProjectServiceServer(s grpc.ServiceRegistrar, srv ProjectServiceServer) {
	s.RegisterService(&ProjectService_ServiceDesc, srv)
}

func _ProjectService_Index_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IndexMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectServiceServer).Index(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/project.service.v1.ProjectService/Index",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectServiceServer).Index(ctx, req.(*IndexMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProjectService_FindProjectByMemId_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProjectRpcMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectServiceServer).FindProjectByMemId(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/project.service.v1.ProjectService/FindProjectByMemId",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectServiceServer).FindProjectByMemId(ctx, req.(*ProjectRpcMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProjectService_FindProjectTemplate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProjectRpcMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectServiceServer).FindProjectTemplate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/project.service.v1.ProjectService/FindProjectTemplate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectServiceServer).FindProjectTemplate(ctx, req.(*ProjectRpcMessage))
	}
	return interceptor(ctx, in, info, handler)
}

// ProjectService_ServiceDesc is the grpc.ServiceDesc for ProjectService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ProjectService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "project.service.v1.ProjectService",
	HandlerType: (*ProjectServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Index",
			Handler:    _ProjectService_Index_Handler,
		},
		{
			MethodName: "FindProjectByMemId",
			Handler:    _ProjectService_FindProjectByMemId_Handler,
		},
		{
			MethodName: "FindProjectTemplate",
			Handler:    _ProjectService_FindProjectTemplate_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "project_service.proto",
}
