// Code generated by truss. DO NOT EDIT.
// Rerunning truss will overwrite this file.
// Version: b158a9d285
// Version Date: 2020-10-26T02:16:40Z

// Package grpc provides a gRPC client for the App service.
package grpc

import (
	"context"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"

	// This Service
	"github.com/Reasno/kitty/app/svc"
	pb "github.com/Reasno/kitty/proto"
)

// New returns an service backed by a gRPC client connection. It is the
// responsibility of the caller to dial, and later close, the connection.
func New(conn *grpc.ClientConn, options ...ClientOption) (pb.AppServer, error) {
	var cc clientConfig

	for _, f := range options {
		err := f(&cc)
		if err != nil {
			return nil, errors.Wrap(err, "cannot apply option")
		}
	}

	clientOptions := []grpctransport.ClientOption{
		grpctransport.ClientBefore(
			contextValuesToGRPCMetadata(cc.headers)),
	}
	var loginEndpoint endpoint.Endpoint
	{
		loginEndpoint = grpctransport.NewClient(
			conn,
			"kitty.App",
			"Login",
			EncodeGRPCLoginRequest,
			DecodeGRPCLoginResponse,
			pb.UserInfoReply{},
			clientOptions...,
		).Endpoint()
	}

	var getcodeEndpoint endpoint.Endpoint
	{
		getcodeEndpoint = grpctransport.NewClient(
			conn,
			"kitty.App",
			"GetCode",
			EncodeGRPCGetCodeRequest,
			DecodeGRPCGetCodeResponse,
			pb.GenericReply{},
			clientOptions...,
		).Endpoint()
	}

	var getinfoEndpoint endpoint.Endpoint
	{
		getinfoEndpoint = grpctransport.NewClient(
			conn,
			"kitty.App",
			"GetInfo",
			EncodeGRPCGetInfoRequest,
			DecodeGRPCGetInfoResponse,
			pb.UserInfoReply{},
			clientOptions...,
		).Endpoint()
	}

	var updateinfoEndpoint endpoint.Endpoint
	{
		updateinfoEndpoint = grpctransport.NewClient(
			conn,
			"kitty.App",
			"UpdateInfo",
			EncodeGRPCUpdateInfoRequest,
			DecodeGRPCUpdateInfoResponse,
			pb.UserInfoReply{},
			clientOptions...,
		).Endpoint()
	}

	var bindEndpoint endpoint.Endpoint
	{
		bindEndpoint = grpctransport.NewClient(
			conn,
			"kitty.App",
			"Bind",
			EncodeGRPCBindRequest,
			DecodeGRPCBindResponse,
			pb.UserInfoReply{},
			clientOptions...,
		).Endpoint()
	}

	var unbindEndpoint endpoint.Endpoint
	{
		unbindEndpoint = grpctransport.NewClient(
			conn,
			"kitty.App",
			"Unbind",
			EncodeGRPCUnbindRequest,
			DecodeGRPCUnbindResponse,
			pb.UserInfoReply{},
			clientOptions...,
		).Endpoint()
	}

	var refreshEndpoint endpoint.Endpoint
	{
		refreshEndpoint = grpctransport.NewClient(
			conn,
			"kitty.App",
			"Refresh",
			EncodeGRPCRefreshRequest,
			DecodeGRPCRefreshResponse,
			pb.UserInfoReply{},
			clientOptions...,
		).Endpoint()
	}

	return svc.Endpoints{
		LoginEndpoint:      loginEndpoint,
		GetCodeEndpoint:    getcodeEndpoint,
		GetInfoEndpoint:    getinfoEndpoint,
		UpdateInfoEndpoint: updateinfoEndpoint,
		BindEndpoint:       bindEndpoint,
		UnbindEndpoint:     unbindEndpoint,
		RefreshEndpoint:    refreshEndpoint,
	}, nil
}

// GRPC Client Decode

// DecodeGRPCLoginResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC login reply to a user-domain login response. Primarily useful in a client.
func DecodeGRPCLoginResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.UserInfoReply)
	return reply, nil
}

// DecodeGRPCGetCodeResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC getcode reply to a user-domain getcode response. Primarily useful in a client.
func DecodeGRPCGetCodeResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.GenericReply)
	return reply, nil
}

// DecodeGRPCGetInfoResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC getinfo reply to a user-domain getinfo response. Primarily useful in a client.
func DecodeGRPCGetInfoResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.UserInfoReply)
	return reply, nil
}

// DecodeGRPCUpdateInfoResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC updateinfo reply to a user-domain updateinfo response. Primarily useful in a client.
func DecodeGRPCUpdateInfoResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.UserInfoReply)
	return reply, nil
}

// DecodeGRPCBindResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC bind reply to a user-domain bind response. Primarily useful in a client.
func DecodeGRPCBindResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.UserInfoReply)
	return reply, nil
}

// DecodeGRPCUnbindResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC unbind reply to a user-domain unbind response. Primarily useful in a client.
func DecodeGRPCUnbindResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.UserInfoReply)
	return reply, nil
}

// DecodeGRPCRefreshResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC refresh reply to a user-domain refresh response. Primarily useful in a client.
func DecodeGRPCRefreshResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.UserInfoReply)
	return reply, nil
}

// GRPC Client Encode

// EncodeGRPCLoginRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain login request to a gRPC login request. Primarily useful in a client.
func EncodeGRPCLoginRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.UserLoginRequest)
	return req, nil
}

// EncodeGRPCGetCodeRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain getcode request to a gRPC getcode request. Primarily useful in a client.
func EncodeGRPCGetCodeRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.GetCodeRequest)
	return req, nil
}

// EncodeGRPCGetInfoRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain getinfo request to a gRPC getinfo request. Primarily useful in a client.
func EncodeGRPCGetInfoRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.UserInfoRequest)
	return req, nil
}

// EncodeGRPCUpdateInfoRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain updateinfo request to a gRPC updateinfo request. Primarily useful in a client.
func EncodeGRPCUpdateInfoRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.UserInfoUpdateRequest)
	return req, nil
}

// EncodeGRPCBindRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain bind request to a gRPC bind request. Primarily useful in a client.
func EncodeGRPCBindRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.UserBindRequest)
	return req, nil
}

// EncodeGRPCUnbindRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain unbind request to a gRPC unbind request. Primarily useful in a client.
func EncodeGRPCUnbindRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.UserUnbindRequest)
	return req, nil
}

// EncodeGRPCRefreshRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain refresh request to a gRPC refresh request. Primarily useful in a client.
func EncodeGRPCRefreshRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.UserRefreshRequest)
	return req, nil
}

type clientConfig struct {
	headers []string
}

// ClientOption is a function that modifies the client config
type ClientOption func(*clientConfig) error

func CtxValuesToSend(keys ...string) ClientOption {
	return func(o *clientConfig) error {
		o.headers = keys
		return nil
	}
}

func contextValuesToGRPCMetadata(keys []string) grpctransport.ClientRequestFunc {
	return func(ctx context.Context, md *metadata.MD) context.Context {
		var pairs []string
		for _, k := range keys {
			if v, ok := ctx.Value(k).(string); ok {
				pairs = append(pairs, k, v)
			}
		}

		if pairs != nil {
			*md = metadata.Join(*md, metadata.Pairs(pairs...))
		}

		return ctx
	}
}
