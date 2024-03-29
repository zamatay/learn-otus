// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package v1

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

// EventsClient is the client API for Events service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EventsClient interface {
	AddEvent(ctx context.Context, in *EventRequest, opts ...grpc.CallOption) (*OkResponse, error)
	EditEvent(ctx context.Context, in *EventRequest, opts ...grpc.CallOption) (*OkResponse, error)
	RemoveEvent(ctx context.Context, in *IdRequest, opts ...grpc.CallOption) (*OkResponse, error)
	List(ctx context.Context, in *DateRequest, opts ...grpc.CallOption) (*EventDataSet, error)
	GetEvent(ctx context.Context, in *IdRequest, opts ...grpc.CallOption) (*EventRequest, error)
}

type eventsClient struct {
	cc grpc.ClientConnInterface
}

func NewEventsClient(cc grpc.ClientConnInterface) EventsClient {
	return &eventsClient{cc}
}

func (c *eventsClient) AddEvent(ctx context.Context, in *EventRequest, opts ...grpc.CallOption) (*OkResponse, error) {
	out := new(OkResponse)
	err := c.cc.Invoke(ctx, "/event.Events/AddEvent", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *eventsClient) EditEvent(ctx context.Context, in *EventRequest, opts ...grpc.CallOption) (*OkResponse, error) {
	out := new(OkResponse)
	err := c.cc.Invoke(ctx, "/event.Events/EditEvent", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *eventsClient) RemoveEvent(ctx context.Context, in *IdRequest, opts ...grpc.CallOption) (*OkResponse, error) {
	out := new(OkResponse)
	err := c.cc.Invoke(ctx, "/event.Events/RemoveEvent", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *eventsClient) List(ctx context.Context, in *DateRequest, opts ...grpc.CallOption) (*EventDataSet, error) {
	out := new(EventDataSet)
	err := c.cc.Invoke(ctx, "/event.Events/List", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *eventsClient) GetEvent(ctx context.Context, in *IdRequest, opts ...grpc.CallOption) (*EventRequest, error) {
	out := new(EventRequest)
	err := c.cc.Invoke(ctx, "/event.Events/GetEvent", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EventsServer is the server API for Events service.
// All implementations must embed UnimplementedEventsServer
// for forward compatibility
type EventsServer interface {
	AddEvent(context.Context, *EventRequest) (*OkResponse, error)
	EditEvent(context.Context, *EventRequest) (*OkResponse, error)
	RemoveEvent(context.Context, *IdRequest) (*OkResponse, error)
	List(context.Context, *DateRequest) (*EventDataSet, error)
	GetEvent(context.Context, *IdRequest) (*EventRequest, error)
	mustEmbedUnimplementedEventsServer()
}

// UnimplementedEventsServer must be embedded to have forward compatible implementations.
type UnimplementedEventsServer struct {
}

func (UnimplementedEventsServer) AddEvent(context.Context, *EventRequest) (*OkResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddEvent not implemented")
}
func (UnimplementedEventsServer) EditEvent(context.Context, *EventRequest) (*OkResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EditEvent not implemented")
}
func (UnimplementedEventsServer) RemoveEvent(context.Context, *IdRequest) (*OkResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveEvent not implemented")
}
func (UnimplementedEventsServer) List(context.Context, *DateRequest) (*EventDataSet, error) {
	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}
func (UnimplementedEventsServer) GetEvent(context.Context, *IdRequest) (*EventRequest, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetEvent not implemented")
}
func (UnimplementedEventsServer) mustEmbedUnimplementedEventsServer() {}

// UnsafeEventsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EventsServer will
// result in compilation errors.
type UnsafeEventsServer interface {
	mustEmbedUnimplementedEventsServer()
}

func RegisterEventsServer(s grpc.ServiceRegistrar, srv EventsServer) {
	s.RegisterService(&Events_ServiceDesc, srv)
}

func _Events_AddEvent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EventRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventsServer).AddEvent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/event.Events/AddEvent",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventsServer).AddEvent(ctx, req.(*EventRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Events_EditEvent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EventRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventsServer).EditEvent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/event.Events/EditEvent",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventsServer).EditEvent(ctx, req.(*EventRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Events_RemoveEvent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventsServer).RemoveEvent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/event.Events/RemoveEvent",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventsServer).RemoveEvent(ctx, req.(*IdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Events_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventsServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/event.Events/List",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventsServer).List(ctx, req.(*DateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Events_GetEvent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EventsServer).GetEvent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/event.Events/GetEvent",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EventsServer).GetEvent(ctx, req.(*IdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Events_ServiceDesc is the grpc.ServiceDesc for Events service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Events_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "event.Events",
	HandlerType: (*EventsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddEvent",
			Handler:    _Events_AddEvent_Handler,
		},
		{
			MethodName: "EditEvent",
			Handler:    _Events_EditEvent_Handler,
		},
		{
			MethodName: "RemoveEvent",
			Handler:    _Events_RemoveEvent_Handler,
		},
		{
			MethodName: "List",
			Handler:    _Events_List_Handler,
		},
		{
			MethodName: "GetEvent",
			Handler:    _Events_GetEvent_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "EventService.proto",
}
