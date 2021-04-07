// Code generated by truss. DO NOT EDIT.
// Rerunning truss will overwrite this file.
// Version: e12fd89529
// Version Date: 2021-03-04T06:59:01Z

package svc

// This file provides server-side bindings for the gRPC transport.
// It utilizes the transport/grpc.Server.

import (
	"context"
	"net/http"

	"google.golang.org/grpc/metadata"

	grpctransport "github.com/go-kit/kit/transport/grpc"

	// This Service
	pb "glab.tagtic.cn/ad_gains/kitty/proto"
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
		getinfo: grpctransport.NewServer(
			endpoints.GetInfoEndpoint,
			DecodeGRPCGetInfoRequest,
			EncodeGRPCGetInfoResponse,
			serverOptions...,
		),
		getinfobatch: grpctransport.NewServer(
			endpoints.GetInfoBatchEndpoint,
			DecodeGRPCGetInfoBatchRequest,
			EncodeGRPCGetInfoBatchResponse,
			serverOptions...,
		),
		updateinfo: grpctransport.NewServer(
			endpoints.UpdateInfoEndpoint,
			DecodeGRPCUpdateInfoRequest,
			EncodeGRPCUpdateInfoResponse,
			serverOptions...,
		),
		bind: grpctransport.NewServer(
			endpoints.BindEndpoint,
			DecodeGRPCBindRequest,
			EncodeGRPCBindResponse,
			serverOptions...,
		),
		bindad: grpctransport.NewServer(
			endpoints.BindAdEndpoint,
			DecodeGRPCBindAdRequest,
			EncodeGRPCBindAdResponse,
			serverOptions...,
		),
		unbind: grpctransport.NewServer(
			endpoints.UnbindEndpoint,
			DecodeGRPCUnbindRequest,
			EncodeGRPCUnbindResponse,
			serverOptions...,
		),
		refresh: grpctransport.NewServer(
			endpoints.RefreshEndpoint,
			DecodeGRPCRefreshRequest,
			EncodeGRPCRefreshResponse,
			serverOptions...,
		),
		softdelete: grpctransport.NewServer(
			endpoints.SoftDeleteEndpoint,
			DecodeGRPCSoftDeleteRequest,
			EncodeGRPCSoftDeleteResponse,
			serverOptions...,
		),
	}
}

// grpcServer implements the AppServer interface
type grpcServer struct {
	login        grpctransport.Handler
	getcode      grpctransport.Handler
	getinfo      grpctransport.Handler
	getinfobatch grpctransport.Handler
	updateinfo   grpctransport.Handler
	bind         grpctransport.Handler
	bindad       grpctransport.Handler
	unbind       grpctransport.Handler
	refresh      grpctransport.Handler
	softdelete   grpctransport.Handler
}

// Methods for grpcServer to implement AppServer interface

func (s *grpcServer) Login(ctx context.Context, req *pb.UserLoginRequest) (*pb.UserInfoReply, error) {
	_, rep, err := s.login.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.UserInfoReply), nil
}

func (s *grpcServer) GetCode(ctx context.Context, req *pb.GetCodeRequest) (*pb.GenericReply, error) {
	_, rep, err := s.getcode.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.GenericReply), nil
}

func (s *grpcServer) GetInfo(ctx context.Context, req *pb.UserInfoRequest) (*pb.UserInfoReply, error) {
	_, rep, err := s.getinfo.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.UserInfoReply), nil
}

func (s *grpcServer) GetInfoBatch(ctx context.Context, req *pb.UserInfoBatchRequest) (*pb.UserInfoBatchReply, error) {
	_, rep, err := s.getinfobatch.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.UserInfoBatchReply), nil
}

func (s *grpcServer) UpdateInfo(ctx context.Context, req *pb.UserInfoUpdateRequest) (*pb.UserInfoReply, error) {
	_, rep, err := s.updateinfo.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.UserInfoReply), nil
}

func (s *grpcServer) Bind(ctx context.Context, req *pb.UserBindRequest) (*pb.UserInfoReply, error) {
	_, rep, err := s.bind.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.UserInfoReply), nil
}

func (s *grpcServer) BindAd(ctx context.Context, req *pb.UserBindAdRequest) (*pb.GenericReply, error) {
	_, rep, err := s.bindad.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.GenericReply), nil
}

func (s *grpcServer) Unbind(ctx context.Context, req *pb.UserUnbindRequest) (*pb.UserInfoReply, error) {
	_, rep, err := s.unbind.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.UserInfoReply), nil
}

func (s *grpcServer) Refresh(ctx context.Context, req *pb.UserRefreshRequest) (*pb.UserInfoReply, error) {
	_, rep, err := s.refresh.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.UserInfoReply), nil
}

func (s *grpcServer) SoftDelete(ctx context.Context, req *pb.UserSoftDeleteRequest) (*pb.UserInfoReply, error) {
	_, rep, err := s.softdelete.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.UserInfoReply), nil
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

// DecodeGRPCGetInfoRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC getinfo request to a user-domain getinfo request. Primarily useful in a server.
func DecodeGRPCGetInfoRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.UserInfoRequest)
	return req, nil
}

// DecodeGRPCGetInfoBatchRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC getinfobatch request to a user-domain getinfobatch request. Primarily useful in a server.
func DecodeGRPCGetInfoBatchRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.UserInfoBatchRequest)
	return req, nil
}

// DecodeGRPCUpdateInfoRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC updateinfo request to a user-domain updateinfo request. Primarily useful in a server.
func DecodeGRPCUpdateInfoRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.UserInfoUpdateRequest)
	return req, nil
}

// DecodeGRPCBindRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC bind request to a user-domain bind request. Primarily useful in a server.
func DecodeGRPCBindRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.UserBindRequest)
	return req, nil
}

// DecodeGRPCBindAdRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC bindad request to a user-domain bindad request. Primarily useful in a server.
func DecodeGRPCBindAdRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.UserBindAdRequest)
	return req, nil
}

// DecodeGRPCUnbindRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC unbind request to a user-domain unbind request. Primarily useful in a server.
func DecodeGRPCUnbindRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.UserUnbindRequest)
	return req, nil
}

// DecodeGRPCRefreshRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC refresh request to a user-domain refresh request. Primarily useful in a server.
func DecodeGRPCRefreshRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.UserRefreshRequest)
	return req, nil
}

// DecodeGRPCSoftDeleteRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC softdelete request to a user-domain softdelete request. Primarily useful in a server.
func DecodeGRPCSoftDeleteRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.UserSoftDeleteRequest)
	return req, nil
}

// Server Encode

// EncodeGRPCLoginResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain login response to a gRPC login reply. Primarily useful in a server.
func EncodeGRPCLoginResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*pb.UserInfoReply)
	return resp, nil
}

// EncodeGRPCGetCodeResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain getcode response to a gRPC getcode reply. Primarily useful in a server.
func EncodeGRPCGetCodeResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*pb.GenericReply)
	return resp, nil
}

// EncodeGRPCGetInfoResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain getinfo response to a gRPC getinfo reply. Primarily useful in a server.
func EncodeGRPCGetInfoResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*pb.UserInfoReply)
	return resp, nil
}

// EncodeGRPCGetInfoBatchResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain getinfobatch response to a gRPC getinfobatch reply. Primarily useful in a server.
func EncodeGRPCGetInfoBatchResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*pb.UserInfoBatchReply)
	return resp, nil
}

// EncodeGRPCUpdateInfoResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain updateinfo response to a gRPC updateinfo reply. Primarily useful in a server.
func EncodeGRPCUpdateInfoResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*pb.UserInfoReply)
	return resp, nil
}

// EncodeGRPCBindResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain bind response to a gRPC bind reply. Primarily useful in a server.
func EncodeGRPCBindResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*pb.UserInfoReply)
	return resp, nil
}

// EncodeGRPCBindAdResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain bindad response to a gRPC bindad reply. Primarily useful in a server.
func EncodeGRPCBindAdResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*pb.GenericReply)
	return resp, nil
}

// EncodeGRPCUnbindResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain unbind response to a gRPC unbind reply. Primarily useful in a server.
func EncodeGRPCUnbindResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*pb.UserInfoReply)
	return resp, nil
}

// EncodeGRPCRefreshResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain refresh response to a gRPC refresh reply. Primarily useful in a server.
func EncodeGRPCRefreshResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*pb.UserInfoReply)
	return resp, nil
}

// EncodeGRPCSoftDeleteResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain softdelete response to a gRPC softdelete reply. Primarily useful in a server.
func EncodeGRPCSoftDeleteResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*pb.UserInfoReply)
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
