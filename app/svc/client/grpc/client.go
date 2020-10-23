// Code generated by truss. DO NOT EDIT.
// Rerunning truss will overwrite this file.
// Version:
// Version Date:

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
			pb.UserLoginReply{},
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

	return svc.Endpoints{
		LoginEndpoint:   loginEndpoint,
		GetCodeEndpoint: getcodeEndpoint,
	}, nil
}

// GRPC Client Decode

// DecodeGRPCLoginResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC login reply to a user-domain login response. Primarily useful in a client.
func DecodeGRPCLoginResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.UserLoginReply)
	return reply, nil
}

// DecodeGRPCGetCodeResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC getcode reply to a user-domain getcode response. Primarily useful in a client.
func DecodeGRPCGetCodeResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.GenericReply)
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
