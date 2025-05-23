// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v6.30.2
// source: sso/sso.proto

package sso

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
	Auth_Register_FullMethodName = "/auth.Auth/Register"
	Auth_Login_FullMethodName    = "/auth.Auth/Login"
)

// AuthClient is the client API for Auth service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AuthClient interface {
	Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error)
	Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error)
}

type authClient struct {
	cc grpc.ClientConnInterface
}

func NewAuthClient(cc grpc.ClientConnInterface) AuthClient {
	return &authClient{cc}
}

func (c *authClient) Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RegisterResponse)
	err := c.cc.Invoke(ctx, Auth_Register_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authClient) Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LoginResponse)
	err := c.cc.Invoke(ctx, Auth_Login_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AuthServer is the server API for Auth service.
// All implementations must embed UnimplementedAuthServer
// for forward compatibility.
type AuthServer interface {
	Register(context.Context, *RegisterRequest) (*RegisterResponse, error)
	Login(context.Context, *LoginRequest) (*LoginResponse, error)
	mustEmbedUnimplementedAuthServer()
}

// UnimplementedAuthServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedAuthServer struct{}

func (UnimplementedAuthServer) Register(context.Context, *RegisterRequest) (*RegisterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}
func (UnimplementedAuthServer) Login(context.Context, *LoginRequest) (*LoginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (UnimplementedAuthServer) mustEmbedUnimplementedAuthServer() {}
func (UnimplementedAuthServer) testEmbeddedByValue()              {}

// UnsafeAuthServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AuthServer will
// result in compilation errors.
type UnsafeAuthServer interface {
	mustEmbedUnimplementedAuthServer()
}

func RegisterAuthServer(s grpc.ServiceRegistrar, srv AuthServer) {
	// If the following call pancis, it indicates UnimplementedAuthServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Auth_ServiceDesc, srv)
}

func _Auth_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Auth_Register_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServer).Register(ctx, req.(*RegisterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Auth_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Auth_Login_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServer).Login(ctx, req.(*LoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Auth_ServiceDesc is the grpc.ServiceDesc for Auth service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Auth_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "auth.Auth",
	HandlerType: (*AuthServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Register",
			Handler:    _Auth_Register_Handler,
		},
		{
			MethodName: "Login",
			Handler:    _Auth_Login_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "sso/sso.proto",
}

const (
	Withdrawals_TopUp_FullMethodName       = "/auth.Withdrawals/TopUp"
	Withdrawals_Balance_FullMethodName     = "/auth.Withdrawals/Balance"
	Withdrawals_Withdraw_FullMethodName    = "/auth.Withdrawals/Withdraw"
	Withdrawals_Withdrawals_FullMethodName = "/auth.Withdrawals/Withdrawals"
)

// WithdrawalsClient is the client API for Withdrawals service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type WithdrawalsClient interface {
	TopUp(ctx context.Context, in *TopUpRequest, opts ...grpc.CallOption) (*TopUpResponse, error)
	Balance(ctx context.Context, in *BalanceRequest, opts ...grpc.CallOption) (*BalanceResponse, error)
	Withdraw(ctx context.Context, in *WithdrawRequest, opts ...grpc.CallOption) (*WithdrawResponse, error)
	Withdrawals(ctx context.Context, in *WithdrawalsRequest, opts ...grpc.CallOption) (*WithdrawalsResponse, error)
}

type withdrawalsClient struct {
	cc grpc.ClientConnInterface
}

func NewWithdrawalsClient(cc grpc.ClientConnInterface) WithdrawalsClient {
	return &withdrawalsClient{cc}
}

func (c *withdrawalsClient) TopUp(ctx context.Context, in *TopUpRequest, opts ...grpc.CallOption) (*TopUpResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(TopUpResponse)
	err := c.cc.Invoke(ctx, Withdrawals_TopUp_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *withdrawalsClient) Balance(ctx context.Context, in *BalanceRequest, opts ...grpc.CallOption) (*BalanceResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(BalanceResponse)
	err := c.cc.Invoke(ctx, Withdrawals_Balance_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *withdrawalsClient) Withdraw(ctx context.Context, in *WithdrawRequest, opts ...grpc.CallOption) (*WithdrawResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(WithdrawResponse)
	err := c.cc.Invoke(ctx, Withdrawals_Withdraw_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *withdrawalsClient) Withdrawals(ctx context.Context, in *WithdrawalsRequest, opts ...grpc.CallOption) (*WithdrawalsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(WithdrawalsResponse)
	err := c.cc.Invoke(ctx, Withdrawals_Withdrawals_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// WithdrawalsServer is the server API for Withdrawals service.
// All implementations must embed UnimplementedWithdrawalsServer
// for forward compatibility.
type WithdrawalsServer interface {
	TopUp(context.Context, *TopUpRequest) (*TopUpResponse, error)
	Balance(context.Context, *BalanceRequest) (*BalanceResponse, error)
	Withdraw(context.Context, *WithdrawRequest) (*WithdrawResponse, error)
	Withdrawals(context.Context, *WithdrawalsRequest) (*WithdrawalsResponse, error)
	mustEmbedUnimplementedWithdrawalsServer()
}

// UnimplementedWithdrawalsServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedWithdrawalsServer struct{}

func (UnimplementedWithdrawalsServer) TopUp(context.Context, *TopUpRequest) (*TopUpResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TopUp not implemented")
}
func (UnimplementedWithdrawalsServer) Balance(context.Context, *BalanceRequest) (*BalanceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Balance not implemented")
}
func (UnimplementedWithdrawalsServer) Withdraw(context.Context, *WithdrawRequest) (*WithdrawResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Withdraw not implemented")
}
func (UnimplementedWithdrawalsServer) Withdrawals(context.Context, *WithdrawalsRequest) (*WithdrawalsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Withdrawals not implemented")
}
func (UnimplementedWithdrawalsServer) mustEmbedUnimplementedWithdrawalsServer() {}
func (UnimplementedWithdrawalsServer) testEmbeddedByValue()                     {}

// UnsafeWithdrawalsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to WithdrawalsServer will
// result in compilation errors.
type UnsafeWithdrawalsServer interface {
	mustEmbedUnimplementedWithdrawalsServer()
}

func RegisterWithdrawalsServer(s grpc.ServiceRegistrar, srv WithdrawalsServer) {
	// If the following call pancis, it indicates UnimplementedWithdrawalsServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Withdrawals_ServiceDesc, srv)
}

func _Withdrawals_TopUp_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TopUpRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WithdrawalsServer).TopUp(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Withdrawals_TopUp_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WithdrawalsServer).TopUp(ctx, req.(*TopUpRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Withdrawals_Balance_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BalanceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WithdrawalsServer).Balance(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Withdrawals_Balance_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WithdrawalsServer).Balance(ctx, req.(*BalanceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Withdrawals_Withdraw_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WithdrawRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WithdrawalsServer).Withdraw(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Withdrawals_Withdraw_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WithdrawalsServer).Withdraw(ctx, req.(*WithdrawRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Withdrawals_Withdrawals_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WithdrawalsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WithdrawalsServer).Withdrawals(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Withdrawals_Withdrawals_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WithdrawalsServer).Withdrawals(ctx, req.(*WithdrawalsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Withdrawals_ServiceDesc is the grpc.ServiceDesc for Withdrawals service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Withdrawals_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "auth.Withdrawals",
	HandlerType: (*WithdrawalsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "TopUp",
			Handler:    _Withdrawals_TopUp_Handler,
		},
		{
			MethodName: "Balance",
			Handler:    _Withdrawals_Balance_Handler,
		},
		{
			MethodName: "Withdraw",
			Handler:    _Withdrawals_Withdraw_Handler,
		},
		{
			MethodName: "Withdrawals",
			Handler:    _Withdrawals_Withdrawals_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "sso/sso.proto",
}
