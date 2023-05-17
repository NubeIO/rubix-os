// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.21.12
// source: module.proto

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

const (
	Module_Init_FullMethodName         = "/proto.Module/Init"
	Module_GetUrlPrefix_FullMethodName = "/proto.Module/GetUrlPrefix"
	Module_Get_FullMethodName          = "/proto.Module/Get"
	Module_Post_FullMethodName         = "/proto.Module/Post"
	Module_Put_FullMethodName          = "/proto.Module/Put"
	Module_Patch_FullMethodName        = "/proto.Module/Patch"
	Module_Delete_FullMethodName       = "/proto.Module/Delete"
)

// ModuleClient is the client API for Module service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ModuleClient interface {
	Init(ctx context.Context, in *InitRequest, opts ...grpc.CallOption) (*Empty, error)
	GetUrlPrefix(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*GetUrlPrefixResponse, error)
	Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*Response, error)
	Post(ctx context.Context, in *PostRequest, opts ...grpc.CallOption) (*Response, error)
	Put(ctx context.Context, in *PutRequest, opts ...grpc.CallOption) (*Response, error)
	Patch(ctx context.Context, in *PatchRequest, opts ...grpc.CallOption) (*Response, error)
	Delete(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*Response, error)
}

type moduleClient struct {
	cc grpc.ClientConnInterface
}

func NewModuleClient(cc grpc.ClientConnInterface) ModuleClient {
	return &moduleClient{cc}
}

func (c *moduleClient) Init(ctx context.Context, in *InitRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, Module_Init_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *moduleClient) GetUrlPrefix(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*GetUrlPrefixResponse, error) {
	out := new(GetUrlPrefixResponse)
	err := c.cc.Invoke(ctx, Module_GetUrlPrefix_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *moduleClient) Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, Module_Get_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *moduleClient) Post(ctx context.Context, in *PostRequest, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, Module_Post_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *moduleClient) Put(ctx context.Context, in *PutRequest, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, Module_Put_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *moduleClient) Patch(ctx context.Context, in *PatchRequest, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, Module_Patch_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *moduleClient) Delete(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, Module_Delete_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ModuleServer is the server API for Module service.
// All implementations should embed UnimplementedModuleServer
// for forward compatibility
type ModuleServer interface {
	Init(context.Context, *InitRequest) (*Empty, error)
	GetUrlPrefix(context.Context, *Empty) (*GetUrlPrefixResponse, error)
	Get(context.Context, *GetRequest) (*Response, error)
	Post(context.Context, *PostRequest) (*Response, error)
	Put(context.Context, *PutRequest) (*Response, error)
	Patch(context.Context, *PatchRequest) (*Response, error)
	Delete(context.Context, *DeleteRequest) (*Response, error)
}

// UnimplementedModuleServer should be embedded to have forward compatible implementations.
type UnimplementedModuleServer struct {
}

func (UnimplementedModuleServer) Init(context.Context, *InitRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Init not implemented")
}
func (UnimplementedModuleServer) GetUrlPrefix(context.Context, *Empty) (*GetUrlPrefixResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUrlPrefix not implemented")
}
func (UnimplementedModuleServer) Get(context.Context, *GetRequest) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedModuleServer) Post(context.Context, *PostRequest) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Post not implemented")
}
func (UnimplementedModuleServer) Put(context.Context, *PutRequest) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Put not implemented")
}
func (UnimplementedModuleServer) Patch(context.Context, *PatchRequest) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Patch not implemented")
}
func (UnimplementedModuleServer) Delete(context.Context, *DeleteRequest) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}

// UnsafeModuleServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ModuleServer will
// result in compilation errors.
type UnsafeModuleServer interface {
	mustEmbedUnimplementedModuleServer()
}

func RegisterModuleServer(s grpc.ServiceRegistrar, srv ModuleServer) {
	s.RegisterService(&Module_ServiceDesc, srv)
}

func _Module_Init_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InitRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ModuleServer).Init(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Module_Init_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ModuleServer).Init(ctx, req.(*InitRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Module_GetUrlPrefix_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ModuleServer).GetUrlPrefix(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Module_GetUrlPrefix_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ModuleServer).GetUrlPrefix(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Module_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ModuleServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Module_Get_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ModuleServer).Get(ctx, req.(*GetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Module_Post_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PostRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ModuleServer).Post(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Module_Post_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ModuleServer).Post(ctx, req.(*PostRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Module_Put_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PutRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ModuleServer).Put(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Module_Put_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ModuleServer).Put(ctx, req.(*PutRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Module_Patch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PatchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ModuleServer).Patch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Module_Patch_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ModuleServer).Patch(ctx, req.(*PatchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Module_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ModuleServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Module_Delete_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ModuleServer).Delete(ctx, req.(*DeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Module_ServiceDesc is the grpc.ServiceDesc for Module service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Module_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.Module",
	HandlerType: (*ModuleServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Init",
			Handler:    _Module_Init_Handler,
		},
		{
			MethodName: "GetUrlPrefix",
			Handler:    _Module_GetUrlPrefix_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _Module_Get_Handler,
		},
		{
			MethodName: "Post",
			Handler:    _Module_Post_Handler,
		},
		{
			MethodName: "Put",
			Handler:    _Module_Put_Handler,
		},
		{
			MethodName: "Patch",
			Handler:    _Module_Patch_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _Module_Delete_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "module.proto",
}

const (
	DBHelper_GetList_FullMethodName = "/proto.DBHelper/GetList"
	DBHelper_Get_FullMethodName     = "/proto.DBHelper/Get"
	DBHelper_Post_FullMethodName    = "/proto.DBHelper/Post"
	DBHelper_Put_FullMethodName     = "/proto.DBHelper/Put"
	DBHelper_Patch_FullMethodName   = "/proto.DBHelper/Patch"
	DBHelper_Delete_FullMethodName  = "/proto.DBHelper/Delete"
)

// DBHelperClient is the client API for DBHelper service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DBHelperClient interface {
	GetList(ctx context.Context, in *GetListRequest, opts ...grpc.CallOption) (*Response, error)
	Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*Response, error)
	Post(ctx context.Context, in *PostRequest, opts ...grpc.CallOption) (*Response, error)
	Put(ctx context.Context, in *PutRequest, opts ...grpc.CallOption) (*Response, error)
	Patch(ctx context.Context, in *PatchRequest, opts ...grpc.CallOption) (*Response, error)
	Delete(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*Response, error)
}

type dBHelperClient struct {
	cc grpc.ClientConnInterface
}

func NewDBHelperClient(cc grpc.ClientConnInterface) DBHelperClient {
	return &dBHelperClient{cc}
}

func (c *dBHelperClient) GetList(ctx context.Context, in *GetListRequest, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, DBHelper_GetList_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dBHelperClient) Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, DBHelper_Get_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dBHelperClient) Post(ctx context.Context, in *PostRequest, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, DBHelper_Post_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dBHelperClient) Put(ctx context.Context, in *PutRequest, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, DBHelper_Put_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dBHelperClient) Patch(ctx context.Context, in *PatchRequest, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, DBHelper_Patch_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dBHelperClient) Delete(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, DBHelper_Delete_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DBHelperServer is the server API for DBHelper service.
// All implementations should embed UnimplementedDBHelperServer
// for forward compatibility
type DBHelperServer interface {
	GetList(context.Context, *GetListRequest) (*Response, error)
	Get(context.Context, *GetRequest) (*Response, error)
	Post(context.Context, *PostRequest) (*Response, error)
	Put(context.Context, *PutRequest) (*Response, error)
	Patch(context.Context, *PatchRequest) (*Response, error)
	Delete(context.Context, *DeleteRequest) (*Response, error)
}

// UnimplementedDBHelperServer should be embedded to have forward compatible implementations.
type UnimplementedDBHelperServer struct {
}

func (UnimplementedDBHelperServer) GetList(context.Context, *GetListRequest) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetList not implemented")
}
func (UnimplementedDBHelperServer) Get(context.Context, *GetRequest) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedDBHelperServer) Post(context.Context, *PostRequest) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Post not implemented")
}
func (UnimplementedDBHelperServer) Put(context.Context, *PutRequest) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Put not implemented")
}
func (UnimplementedDBHelperServer) Patch(context.Context, *PatchRequest) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Patch not implemented")
}
func (UnimplementedDBHelperServer) Delete(context.Context, *DeleteRequest) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}

// UnsafeDBHelperServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DBHelperServer will
// result in compilation errors.
type UnsafeDBHelperServer interface {
	mustEmbedUnimplementedDBHelperServer()
}

func RegisterDBHelperServer(s grpc.ServiceRegistrar, srv DBHelperServer) {
	s.RegisterService(&DBHelper_ServiceDesc, srv)
}

func _DBHelper_GetList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DBHelperServer).GetList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DBHelper_GetList_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DBHelperServer).GetList(ctx, req.(*GetListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DBHelper_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DBHelperServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DBHelper_Get_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DBHelperServer).Get(ctx, req.(*GetRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DBHelper_Post_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PostRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DBHelperServer).Post(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DBHelper_Post_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DBHelperServer).Post(ctx, req.(*PostRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DBHelper_Put_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PutRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DBHelperServer).Put(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DBHelper_Put_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DBHelperServer).Put(ctx, req.(*PutRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DBHelper_Patch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PatchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DBHelperServer).Patch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DBHelper_Patch_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DBHelperServer).Patch(ctx, req.(*PatchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DBHelper_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DBHelperServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DBHelper_Delete_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DBHelperServer).Delete(ctx, req.(*DeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// DBHelper_ServiceDesc is the grpc.ServiceDesc for DBHelper service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DBHelper_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.DBHelper",
	HandlerType: (*DBHelperServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetList",
			Handler:    _DBHelper_GetList_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _DBHelper_Get_Handler,
		},
		{
			MethodName: "Post",
			Handler:    _DBHelper_Post_Handler,
		},
		{
			MethodName: "Put",
			Handler:    _DBHelper_Put_Handler,
		},
		{
			MethodName: "Patch",
			Handler:    _DBHelper_Patch_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _DBHelper_Delete_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "module.proto",
}
