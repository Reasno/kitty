// Code generated by truss. DO NOT EDIT.
// Rerunning truss will overwrite this file.
// Version: 831b290599
// Version Date: 2020-11-16T05:27:36Z

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

// MakeGRPCServer makes a set of endpoints available as a gRPC ShareServer.
func MakeGRPCServer(endpoints Endpoints, options ...grpctransport.ServerOption) pb.ShareServer {
	serverOptions := []grpctransport.ServerOption{
		grpctransport.ServerBefore(metadataToContext),
	}
	serverOptions = append(serverOptions, options...)
	return &grpcServer{
		// share

		invitebyurl: grpctransport.NewServer(
			endpoints.InviteByUrlEndpoint,
			DecodeGRPCInviteByUrlRequest,
			EncodeGRPCInviteByUrlResponse,
			serverOptions...,
		),
		invitebytoken: grpctransport.NewServer(
			endpoints.InviteByTokenEndpoint,
			DecodeGRPCInviteByTokenRequest,
			EncodeGRPCInviteByTokenResponse,
			serverOptions...,
		),
		addinvitationcode: grpctransport.NewServer(
			endpoints.AddInvitationCodeEndpoint,
			DecodeGRPCAddInvitationCodeRequest,
			EncodeGRPCAddInvitationCodeResponse,
			serverOptions...,
		),
		listfriend: grpctransport.NewServer(
			endpoints.ListFriendEndpoint,
			DecodeGRPCListFriendRequest,
			EncodeGRPCListFriendResponse,
			serverOptions...,
		),
		claimreward: grpctransport.NewServer(
			endpoints.ClaimRewardEndpoint,
			DecodeGRPCClaimRewardRequest,
			EncodeGRPCClaimRewardResponse,
			serverOptions...,
		),
	}
}

// grpcServer implements the ShareServer interface
type grpcServer struct {
	invitebyurl       grpctransport.Handler
	invitebytoken     grpctransport.Handler
	addinvitationcode grpctransport.Handler
	listfriend        grpctransport.Handler
	claimreward       grpctransport.Handler
}

// Methods for grpcServer to implement ShareServer interface

func (s *grpcServer) InviteByUrl(ctx context.Context, req *pb.ShareEmptyRequest) (*pb.ShareDataReply, error) {
	_, rep, err := s.invitebyurl.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.ShareDataReply), nil
}

func (s *grpcServer) InviteByToken(ctx context.Context, req *pb.ShareEmptyRequest) (*pb.ShareDataReply, error) {
	_, rep, err := s.invitebytoken.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.ShareDataReply), nil
}

func (s *grpcServer) AddInvitationCode(ctx context.Context, req *pb.ShareAddInvitationRequest) (*pb.ShareGenericReply, error) {
	_, rep, err := s.addinvitationcode.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.ShareGenericReply), nil
}

func (s *grpcServer) ListFriend(ctx context.Context, req *pb.ShareListFriendRequest) (*pb.ShareListFriendReply, error) {
	_, rep, err := s.listfriend.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.ShareListFriendReply), nil
}

func (s *grpcServer) ClaimReward(ctx context.Context, req *pb.ShareClaimRewardRequest) (*pb.ShareGenericReply, error) {
	_, rep, err := s.claimreward.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.ShareGenericReply), nil
}

// Server Decode

// DecodeGRPCInviteByUrlRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC invitebyurl request to a user-domain invitebyurl request. Primarily useful in a server.
func DecodeGRPCInviteByUrlRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.ShareEmptyRequest)
	return req, nil
}

// DecodeGRPCInviteByTokenRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC invitebytoken request to a user-domain invitebytoken request. Primarily useful in a server.
func DecodeGRPCInviteByTokenRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.ShareEmptyRequest)
	return req, nil
}

// DecodeGRPCAddInvitationCodeRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC addinvitationcode request to a user-domain addinvitationcode request. Primarily useful in a server.
func DecodeGRPCAddInvitationCodeRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.ShareAddInvitationRequest)
	return req, nil
}

// DecodeGRPCListFriendRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC listfriend request to a user-domain listfriend request. Primarily useful in a server.
func DecodeGRPCListFriendRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.ShareListFriendRequest)
	return req, nil
}

// DecodeGRPCClaimRewardRequest is a transport/grpc.DecodeRequestFunc that converts a
// gRPC claimreward request to a user-domain claimreward request. Primarily useful in a server.
func DecodeGRPCClaimRewardRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.ShareClaimRewardRequest)
	return req, nil
}

// Server Encode

// EncodeGRPCInviteByUrlResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain invitebyurl response to a gRPC invitebyurl reply. Primarily useful in a server.
func EncodeGRPCInviteByUrlResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*pb.ShareDataReply)
	return resp, nil
}

// EncodeGRPCInviteByTokenResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain invitebytoken response to a gRPC invitebytoken reply. Primarily useful in a server.
func EncodeGRPCInviteByTokenResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*pb.ShareDataReply)
	return resp, nil
}

// EncodeGRPCAddInvitationCodeResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain addinvitationcode response to a gRPC addinvitationcode reply. Primarily useful in a server.
func EncodeGRPCAddInvitationCodeResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*pb.ShareGenericReply)
	return resp, nil
}

// EncodeGRPCListFriendResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain listfriend response to a gRPC listfriend reply. Primarily useful in a server.
func EncodeGRPCListFriendResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*pb.ShareListFriendReply)
	return resp, nil
}

// EncodeGRPCClaimRewardResponse is a transport/grpc.EncodeResponseFunc that converts a
// user-domain claimreward response to a gRPC claimreward reply. Primarily useful in a server.
func EncodeGRPCClaimRewardResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(*pb.ShareGenericReply)
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
