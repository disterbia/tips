// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.25.1
// source: email.proto

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

const (
	EmailService_SendEmail_FullMethodName                 = "/emailservice.EmailService/SendEmail"
	EmailService_SendCodeEmail_FullMethodName             = "/emailservice.EmailService/SendCodeEmail"
	EmailService_KldgaSendEmail_FullMethodName            = "/emailservice.EmailService/KldgaSendEmail"
	EmailService_KldgaSendCompetitionEmail_FullMethodName = "/emailservice.EmailService/KldgaSendCompetitionEmail"
	EmailService_AdapfitInquire_FullMethodName            = "/emailservice.EmailService/AdapfitInquire"
)

// EmailServiceClient is the client API for EmailService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EmailServiceClient interface {
	SendEmail(ctx context.Context, in *EmailRequest, opts ...grpc.CallOption) (*EmailResponse, error)
	SendCodeEmail(ctx context.Context, in *EmailCodeRequest, opts ...grpc.CallOption) (*EmailResponse, error)
	KldgaSendEmail(ctx context.Context, in *KldgaEmailRequest, opts ...grpc.CallOption) (*EmailResponse, error)
	KldgaSendCompetitionEmail(ctx context.Context, in *KldgaCompetitionRequest, opts ...grpc.CallOption) (*EmailResponse, error)
	AdapfitInquire(ctx context.Context, in *AdapfitReqeust, opts ...grpc.CallOption) (*EmailResponse, error)
}

type emailServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewEmailServiceClient(cc grpc.ClientConnInterface) EmailServiceClient {
	return &emailServiceClient{cc}
}

func (c *emailServiceClient) SendEmail(ctx context.Context, in *EmailRequest, opts ...grpc.CallOption) (*EmailResponse, error) {
	out := new(EmailResponse)
	err := c.cc.Invoke(ctx, EmailService_SendEmail_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *emailServiceClient) SendCodeEmail(ctx context.Context, in *EmailCodeRequest, opts ...grpc.CallOption) (*EmailResponse, error) {
	out := new(EmailResponse)
	err := c.cc.Invoke(ctx, EmailService_SendCodeEmail_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *emailServiceClient) KldgaSendEmail(ctx context.Context, in *KldgaEmailRequest, opts ...grpc.CallOption) (*EmailResponse, error) {
	out := new(EmailResponse)
	err := c.cc.Invoke(ctx, EmailService_KldgaSendEmail_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *emailServiceClient) KldgaSendCompetitionEmail(ctx context.Context, in *KldgaCompetitionRequest, opts ...grpc.CallOption) (*EmailResponse, error) {
	out := new(EmailResponse)
	err := c.cc.Invoke(ctx, EmailService_KldgaSendCompetitionEmail_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *emailServiceClient) AdapfitInquire(ctx context.Context, in *AdapfitReqeust, opts ...grpc.CallOption) (*EmailResponse, error) {
	out := new(EmailResponse)
	err := c.cc.Invoke(ctx, EmailService_AdapfitInquire_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EmailServiceServer is the server API for EmailService service.
// All implementations must embed UnimplementedEmailServiceServer
// for forward compatibility
type EmailServiceServer interface {
	SendEmail(context.Context, *EmailRequest) (*EmailResponse, error)
	SendCodeEmail(context.Context, *EmailCodeRequest) (*EmailResponse, error)
	KldgaSendEmail(context.Context, *KldgaEmailRequest) (*EmailResponse, error)
	KldgaSendCompetitionEmail(context.Context, *KldgaCompetitionRequest) (*EmailResponse, error)
	AdapfitInquire(context.Context, *AdapfitReqeust) (*EmailResponse, error)
	mustEmbedUnimplementedEmailServiceServer()
}

// UnimplementedEmailServiceServer must be embedded to have forward compatible implementations.
type UnimplementedEmailServiceServer struct {
}

func (UnimplementedEmailServiceServer) SendEmail(context.Context, *EmailRequest) (*EmailResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendEmail not implemented")
}
func (UnimplementedEmailServiceServer) SendCodeEmail(context.Context, *EmailCodeRequest) (*EmailResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendCodeEmail not implemented")
}
func (UnimplementedEmailServiceServer) KldgaSendEmail(context.Context, *KldgaEmailRequest) (*EmailResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method KldgaSendEmail not implemented")
}
func (UnimplementedEmailServiceServer) KldgaSendCompetitionEmail(context.Context, *KldgaCompetitionRequest) (*EmailResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method KldgaSendCompetitionEmail not implemented")
}
func (UnimplementedEmailServiceServer) AdapfitInquire(context.Context, *AdapfitReqeust) (*EmailResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AdapfitInquire not implemented")
}
func (UnimplementedEmailServiceServer) mustEmbedUnimplementedEmailServiceServer() {}

// UnsafeEmailServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EmailServiceServer will
// result in compilation errors.
type UnsafeEmailServiceServer interface {
	mustEmbedUnimplementedEmailServiceServer()
}

func RegisterEmailServiceServer(s grpc.ServiceRegistrar, srv EmailServiceServer) {
	s.RegisterService(&EmailService_ServiceDesc, srv)
}

func _EmailService_SendEmail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmailRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EmailServiceServer).SendEmail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EmailService_SendEmail_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EmailServiceServer).SendEmail(ctx, req.(*EmailRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EmailService_SendCodeEmail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmailCodeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EmailServiceServer).SendCodeEmail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EmailService_SendCodeEmail_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EmailServiceServer).SendCodeEmail(ctx, req.(*EmailCodeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EmailService_KldgaSendEmail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(KldgaEmailRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EmailServiceServer).KldgaSendEmail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EmailService_KldgaSendEmail_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EmailServiceServer).KldgaSendEmail(ctx, req.(*KldgaEmailRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EmailService_KldgaSendCompetitionEmail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(KldgaCompetitionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EmailServiceServer).KldgaSendCompetitionEmail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EmailService_KldgaSendCompetitionEmail_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EmailServiceServer).KldgaSendCompetitionEmail(ctx, req.(*KldgaCompetitionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EmailService_AdapfitInquire_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AdapfitReqeust)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EmailServiceServer).AdapfitInquire(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EmailService_AdapfitInquire_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EmailServiceServer).AdapfitInquire(ctx, req.(*AdapfitReqeust))
	}
	return interceptor(ctx, in, info, handler)
}

// EmailService_ServiceDesc is the grpc.ServiceDesc for EmailService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var EmailService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "emailservice.EmailService",
	HandlerType: (*EmailServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendEmail",
			Handler:    _EmailService_SendEmail_Handler,
		},
		{
			MethodName: "SendCodeEmail",
			Handler:    _EmailService_SendCodeEmail_Handler,
		},
		{
			MethodName: "KldgaSendEmail",
			Handler:    _EmailService_KldgaSendEmail_Handler,
		},
		{
			MethodName: "KldgaSendCompetitionEmail",
			Handler:    _EmailService_KldgaSendCompetitionEmail_Handler,
		},
		{
			MethodName: "AdapfitInquire",
			Handler:    _EmailService_AdapfitInquire_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "email.proto",
}
