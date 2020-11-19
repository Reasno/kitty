// Code generated by truss. DO NOT EDIT.
// Rerunning truss will overwrite this file.
// Version: 831b290599
// Version Date: 2020-11-16T05:27:36Z

// Package grpc provides a gRPC client for the Share service.
package grpc

import (
	"context"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"

	// This Service
	pb "glab.tagtic.cn/ad_gains/kitty/proto"
	"glab.tagtic.cn/ad_gains/kitty/share/svc"
)

// New returns an service backed by a gRPC client connection. It is the
// responsibility of the caller to dial, and later close, the connection.
func New(conn *grpc.ClientConn, options ...ClientOption) (pb.ShareServer, error) {
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
	var invitebyurlEndpoint endpoint.Endpoint
	{
		invitebyurlEndpoint = grpctransport.NewClient(
			conn,
			"kitty.Share",
			"InviteByUrl",
			EncodeGRPCInviteByUrlRequest,
			DecodeGRPCInviteByUrlResponse,
			pb.ShareDataReply{},
			clientOptions...,
		).Endpoint()
	}

	var invitebytokenEndpoint endpoint.Endpoint
	{
		invitebytokenEndpoint = grpctransport.NewClient(
			conn,
			"kitty.Share",
			"InviteByToken",
			EncodeGRPCInviteByTokenRequest,
			DecodeGRPCInviteByTokenResponse,
			pb.ShareDataReply{},
			clientOptions...,
		).Endpoint()
	}

	var addinvitationcodeEndpoint endpoint.Endpoint
	{
		addinvitationcodeEndpoint = grpctransport.NewClient(
			conn,
			"kitty.Share",
			"AddInvitationCode",
			EncodeGRPCAddInvitationCodeRequest,
			DecodeGRPCAddInvitationCodeResponse,
			pb.ShareGenericReply{},
			clientOptions...,
		).Endpoint()
	}

	var listfriendEndpoint endpoint.Endpoint
	{
		listfriendEndpoint = grpctransport.NewClient(
			conn,
			"kitty.Share",
			"ListFriend",
			EncodeGRPCListFriendRequest,
			DecodeGRPCListFriendResponse,
			pb.ShareListFriendReply{},
			clientOptions...,
		).Endpoint()
	}

	var claimrewardEndpoint endpoint.Endpoint
	{
		claimrewardEndpoint = grpctransport.NewClient(
			conn,
			"kitty.Share",
			"ClaimReward",
			EncodeGRPCClaimRewardRequest,
			DecodeGRPCClaimRewardResponse,
			pb.ShareGenericReply{},
			clientOptions...,
		).Endpoint()
	}

	return svc.Endpoints{
		InviteByUrlEndpoint:       invitebyurlEndpoint,
		InviteByTokenEndpoint:     invitebytokenEndpoint,
		AddInvitationCodeEndpoint: addinvitationcodeEndpoint,
		ListFriendEndpoint:        listfriendEndpoint,
		ClaimRewardEndpoint:       claimrewardEndpoint,
	}, nil
}

// GRPC Client Decode

// DecodeGRPCInviteByUrlResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC invitebyurl reply to a user-domain invitebyurl response. Primarily useful in a client.
func DecodeGRPCInviteByUrlResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.ShareDataReply)
	return reply, nil
}

// DecodeGRPCInviteByTokenResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC invitebytoken reply to a user-domain invitebytoken response. Primarily useful in a client.
func DecodeGRPCInviteByTokenResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.ShareDataReply)
	return reply, nil
}

// DecodeGRPCAddInvitationCodeResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC addinvitationcode reply to a user-domain addinvitationcode response. Primarily useful in a client.
func DecodeGRPCAddInvitationCodeResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.ShareGenericReply)
	return reply, nil
}

// DecodeGRPCListFriendResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC listfriend reply to a user-domain listfriend response. Primarily useful in a client.
func DecodeGRPCListFriendResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.ShareListFriendReply)
	return reply, nil
}

// DecodeGRPCClaimRewardResponse is a transport/grpc.DecodeResponseFunc that converts a
// gRPC claimreward reply to a user-domain claimreward response. Primarily useful in a client.
func DecodeGRPCClaimRewardResponse(_ context.Context, grpcReply interface{}) (interface{}, error) {
	reply := grpcReply.(*pb.ShareGenericReply)
	return reply, nil
}

// GRPC Client Encode

// EncodeGRPCInviteByUrlRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain invitebyurl request to a gRPC invitebyurl request. Primarily useful in a client.
func EncodeGRPCInviteByUrlRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.ShareEmptyRequest)
	return req, nil
}

// EncodeGRPCInviteByTokenRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain invitebytoken request to a gRPC invitebytoken request. Primarily useful in a client.
func EncodeGRPCInviteByTokenRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.ShareEmptyRequest)
	return req, nil
}

// EncodeGRPCAddInvitationCodeRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain addinvitationcode request to a gRPC addinvitationcode request. Primarily useful in a client.
func EncodeGRPCAddInvitationCodeRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.ShareAddInvitationRequest)
	return req, nil
}

// EncodeGRPCListFriendRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain listfriend request to a gRPC listfriend request. Primarily useful in a client.
func EncodeGRPCListFriendRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.ShareListFriendRequest)
	return req, nil
}

// EncodeGRPCClaimRewardRequest is a transport/grpc.EncodeRequestFunc that converts a
// user-domain claimreward request to a gRPC claimreward request. Primarily useful in a client.
func EncodeGRPCClaimRewardRequest(_ context.Context, request interface{}) (interface{}, error) {
	req := request.(*pb.ShareClaimRewardRequest)
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
