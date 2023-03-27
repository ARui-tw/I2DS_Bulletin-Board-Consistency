// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.9
// source: BBC/BBC.proto

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

// BulletinClient is the client API for Bulletin service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BulletinClient interface {
	Post(ctx context.Context, in *Content, opts ...grpc.CallOption) (*ACK, error)
	Read(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*ReadResult, error)
	Choose(ctx context.Context, in *ID, opts ...grpc.CallOption) (*Content, error)
	Reply(ctx context.Context, in *Node, opts ...grpc.CallOption) (*ACK, error)
	Update(ctx context.Context, in *Node, opts ...grpc.CallOption) (*ACK, error)
}

type bulletinClient struct {
	cc grpc.ClientConnInterface
}

func NewBulletinClient(cc grpc.ClientConnInterface) BulletinClient {
	return &bulletinClient{cc}
}

func (c *bulletinClient) Post(ctx context.Context, in *Content, opts ...grpc.CallOption) (*ACK, error) {
	out := new(ACK)
	err := c.cc.Invoke(ctx, "/BBC.Bulletin/Post", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bulletinClient) Read(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*ReadResult, error) {
	out := new(ReadResult)
	err := c.cc.Invoke(ctx, "/BBC.Bulletin/Read", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bulletinClient) Choose(ctx context.Context, in *ID, opts ...grpc.CallOption) (*Content, error) {
	out := new(Content)
	err := c.cc.Invoke(ctx, "/BBC.Bulletin/Choose", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bulletinClient) Reply(ctx context.Context, in *Node, opts ...grpc.CallOption) (*ACK, error) {
	out := new(ACK)
	err := c.cc.Invoke(ctx, "/BBC.Bulletin/Reply", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bulletinClient) Update(ctx context.Context, in *Node, opts ...grpc.CallOption) (*ACK, error) {
	out := new(ACK)
	err := c.cc.Invoke(ctx, "/BBC.Bulletin/Update", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BulletinServer is the server API for Bulletin service.
// All implementations must embed UnimplementedBulletinServer
// for forward compatibility
type BulletinServer interface {
	Post(context.Context, *Content) (*ACK, error)
	Read(context.Context, *Empty) (*ReadResult, error)
	Choose(context.Context, *ID) (*Content, error)
	Reply(context.Context, *Node) (*ACK, error)
	Update(context.Context, *Node) (*ACK, error)
	mustEmbedUnimplementedBulletinServer()
}

// UnimplementedBulletinServer must be embedded to have forward compatible implementations.
type UnimplementedBulletinServer struct {
}

func (UnimplementedBulletinServer) Post(context.Context, *Content) (*ACK, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Post not implemented")
}
func (UnimplementedBulletinServer) Read(context.Context, *Empty) (*ReadResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Read not implemented")
}
func (UnimplementedBulletinServer) Choose(context.Context, *ID) (*Content, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Choose not implemented")
}
func (UnimplementedBulletinServer) Reply(context.Context, *Node) (*ACK, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Reply not implemented")
}
func (UnimplementedBulletinServer) Update(context.Context, *Node) (*ACK, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedBulletinServer) mustEmbedUnimplementedBulletinServer() {}

// UnsafeBulletinServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BulletinServer will
// result in compilation errors.
type UnsafeBulletinServer interface {
	mustEmbedUnimplementedBulletinServer()
}

func RegisterBulletinServer(s grpc.ServiceRegistrar, srv BulletinServer) {
	s.RegisterService(&Bulletin_ServiceDesc, srv)
}

func _Bulletin_Post_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Content)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BulletinServer).Post(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/BBC.Bulletin/Post",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BulletinServer).Post(ctx, req.(*Content))
	}
	return interceptor(ctx, in, info, handler)
}

func _Bulletin_Read_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BulletinServer).Read(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/BBC.Bulletin/Read",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BulletinServer).Read(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Bulletin_Choose_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BulletinServer).Choose(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/BBC.Bulletin/Choose",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BulletinServer).Choose(ctx, req.(*ID))
	}
	return interceptor(ctx, in, info, handler)
}

func _Bulletin_Reply_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Node)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BulletinServer).Reply(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/BBC.Bulletin/Reply",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BulletinServer).Reply(ctx, req.(*Node))
	}
	return interceptor(ctx, in, info, handler)
}

func _Bulletin_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Node)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BulletinServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/BBC.Bulletin/Update",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BulletinServer).Update(ctx, req.(*Node))
	}
	return interceptor(ctx, in, info, handler)
}

// Bulletin_ServiceDesc is the grpc.ServiceDesc for Bulletin service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Bulletin_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "BBC.Bulletin",
	HandlerType: (*BulletinServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Post",
			Handler:    _Bulletin_Post_Handler,
		},
		{
			MethodName: "Read",
			Handler:    _Bulletin_Read_Handler,
		},
		{
			MethodName: "Choose",
			Handler:    _Bulletin_Choose_Handler,
		},
		{
			MethodName: "Reply",
			Handler:    _Bulletin_Reply_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _Bulletin_Update_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "BBC/BBC.proto",
}

// PrimaryClient is the client API for Primary service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PrimaryClient interface {
	Post(ctx context.Context, in *Content, opts ...grpc.CallOption) (*ID, error)
	Reply(ctx context.Context, in *Node, opts ...grpc.CallOption) (*ID, error)
}

type primaryClient struct {
	cc grpc.ClientConnInterface
}

func NewPrimaryClient(cc grpc.ClientConnInterface) PrimaryClient {
	return &primaryClient{cc}
}

func (c *primaryClient) Post(ctx context.Context, in *Content, opts ...grpc.CallOption) (*ID, error) {
	out := new(ID)
	err := c.cc.Invoke(ctx, "/BBC.Primary/Post", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *primaryClient) Reply(ctx context.Context, in *Node, opts ...grpc.CallOption) (*ID, error) {
	out := new(ID)
	err := c.cc.Invoke(ctx, "/BBC.Primary/Reply", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PrimaryServer is the server API for Primary service.
// All implementations must embed UnimplementedPrimaryServer
// for forward compatibility
type PrimaryServer interface {
	Post(context.Context, *Content) (*ID, error)
	Reply(context.Context, *Node) (*ID, error)
	mustEmbedUnimplementedPrimaryServer()
}

// UnimplementedPrimaryServer must be embedded to have forward compatible implementations.
type UnimplementedPrimaryServer struct {
}

func (UnimplementedPrimaryServer) Post(context.Context, *Content) (*ID, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Post not implemented")
}
func (UnimplementedPrimaryServer) Reply(context.Context, *Node) (*ID, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Reply not implemented")
}
func (UnimplementedPrimaryServer) mustEmbedUnimplementedPrimaryServer() {}

// UnsafePrimaryServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PrimaryServer will
// result in compilation errors.
type UnsafePrimaryServer interface {
	mustEmbedUnimplementedPrimaryServer()
}

func RegisterPrimaryServer(s grpc.ServiceRegistrar, srv PrimaryServer) {
	s.RegisterService(&Primary_ServiceDesc, srv)
}

func _Primary_Post_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Content)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PrimaryServer).Post(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/BBC.Primary/Post",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PrimaryServer).Post(ctx, req.(*Content))
	}
	return interceptor(ctx, in, info, handler)
}

func _Primary_Reply_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Node)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PrimaryServer).Reply(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/BBC.Primary/Reply",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PrimaryServer).Reply(ctx, req.(*Node))
	}
	return interceptor(ctx, in, info, handler)
}

// Primary_ServiceDesc is the grpc.ServiceDesc for Primary service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Primary_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "BBC.Primary",
	HandlerType: (*PrimaryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Post",
			Handler:    _Primary_Post_Handler,
		},
		{
			MethodName: "Reply",
			Handler:    _Primary_Reply_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "BBC/BBC.proto",
}
