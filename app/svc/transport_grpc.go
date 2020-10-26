// Code generated by truss. DO NOT EDIT.
// Rerunning truss will overwrite this file.
// Version: b158a9d285
// Version Date: 2020-10-26T02:16:40Z

package svc

// This file provides server-side bindings for the gRPC transport.
// It utilizes the transport/grpc.Server.

import (
	"context"
	"net/http"

	"google.golang.org/grpc/metadata"

	grpctransport "github.com/go-kit/kit/transport/grpc"

	// This Service
	pb "github.com/Reasno/kitty/proto"
)

// MakeGRPCServer makes a set of endpoints available as a gRPC AppServer.
func MakeGRPCServer(endpoints Endpoints, options ...grpctransport.ServerOption) pb.AppServer {
	serverOptions := []grpctransport.ServerOption{
		grpctransport.ServerBefore(metadataToContext),
	}
	serverOptions = append(serverOptions, options...)
	return &grpcServer{
		// app

		login: grpctransport.NewServer(
			endpoints.LoginEndpoint,
			DecodeGRPCLoginRequest,
			EncodeGRPCLoginResponse,
			serverOptions...,
		),
		getcode: grpctransport.NewServer(
			endpoints.GetCodeEndpoint,
			DecodeGRPCGetCodeRequest,
			EncodeGRPCGetCodeResponse,
			serverOptions...,
		),
	}
}

// grpcServer implements the AppServer interface
type grpcServer struct {
	login   grpctransport.Handler
	getcode grpctransport.Handler
}

// Methods for grpcServer to implement AppServer interface

func (s *grpcServer) Login(ctx context.Context, req *pb.UserLoginRequest) (*pb.UserLoginReply, error) {
	_, rep, err := s.login.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.UserLoginReply), nil
}

func (s *grpcServer) GetCode(ctx context.Context, req *pb.GetCodeRequest) (*pb.GenericReply, error) {
	_, rep, err := s.getcode.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.GenericReply), nil
}

// Server Decode

// DecodeGRPCLoginRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC login request to a user-domain login request. Primarily useful in a server.
func DecodeGRPCLoginRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.UserLoginRequest)
	return req, nil
}

// DecodeGRPCGetCodeRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC getcode request to a user-domain getcode request. Primarily useful in a server.
func DecodeGRPCGetCodeRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.GetCodeRequest)
	return req, nil
}

// Server Encode

// EncodeGRPCLoginResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain login response to a gRPC login reply. Primarily useful in a server.
func EncodeGRPCLoginResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*pb.UserLoginReply)
	return resp, nil
}

// EncodeGRPCGetCodeResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain getcode response to a gRPC getcode reply. Primarily useful in a server.
func EncodeGRPCGetCodeResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*pb.GenericReply)
	return resp, nil
}

// Helpers

func metadataToContext(ctx context.Context, md metadata.MD) context.Context {
	for k, v := range md {
		if v != nil {
			// The key is added both in metadata format (k) which is all lower
			// and the http.CanonicalHeaderKey of the key so that it can be
			// accessed in either format
			ctx = context.WithValue(ctx, k, v[0])
			ctx = context.WithValue(ctx, http.CanonicalHeaderKey(k), v[0])
		}
	}

	return ctx
}
