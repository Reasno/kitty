// Code generated by truss. DO NOT EDIT.
// Rerunning truss will overwrite this file.
// Version: 831b290599
// Version Date: 2020-11-16T05:27:36Z

package svc

// This file contains methods to make individual endpoints from services,
// request and response types to serve those endpoints, as well as encoders and
// decoders for those types, for all of our supported transport serialization
// formats.

import (
	"context"
	"fmt"

	"github.com/go-kit/kit/endpoint"

	pb "glab.tagtic.cn/ad_gains/kitty/proto"
)

// Endpoints collects all of the endpoints that compose an add service. It's
// meant to be used as a helper struct, to collect all of the endpoints into a
// single parameter.
//
// In a server, it's useful for functions that need to operate on a per-endpoint
// basis. For example, you might pass an Endpoints to a function that produces
// an http.Handler, with each method (endpoint) wired up to a specific path. (It
// is probably a mistake in design to invoke the Service methods on the
// Endpoints struct in a server.)
//
// In a client, it's useful to collect individually constructed endpoints into a
// single type that implements the Service interface. For example, you might
// construct individual endpoints using transport/http.NewClient, combine them into an Endpoints, and return it to the caller as a Service.
type Endpoints struct {
	InviteByUrlEndpoint       endpoint.Endpoint
	InviteByTokenEndpoint     endpoint.Endpoint
	AddInvitationCodeEndpoint endpoint.Endpoint
	ListFriendEndpoint        endpoint.Endpoint
	ClaimRewardEndpoint       endpoint.Endpoint
}

func NewEndpoints(service pb.ShareServer) Endpoints {

	// Endpoint domain.
	var (
		invitebyurlEndpoint       = MakeInviteByUrlEndpoint(service)
		invitebytokenEndpoint     = MakeInviteByTokenEndpoint(service)
		addinvitationcodeEndpoint = MakeAddInvitationCodeEndpoint(service)
		listfriendEndpoint        = MakeListFriendEndpoint(service)
		claimrewardEndpoint       = MakeClaimRewardEndpoint(service)
	)

	endpoints := Endpoints{
		InviteByUrlEndpoint:       invitebyurlEndpoint,
		InviteByTokenEndpoint:     invitebytokenEndpoint,
		AddInvitationCodeEndpoint: addinvitationcodeEndpoint,
		ListFriendEndpoint:        listfriendEndpoint,
		ClaimRewardEndpoint:       claimrewardEndpoint,
	}

	return endpoints
}

// Endpoints

func (e Endpoints) InviteByUrl(ctx context.Context, in *pb.ShareEmptyRequest) (*pb.ShareDataReply, error) {
	response, err := e.InviteByUrlEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	return response.(*pb.ShareDataReply), nil
}

func (e Endpoints) InviteByToken(ctx context.Context, in *pb.ShareEmptyRequest) (*pb.ShareDataReply, error) {
	response, err := e.InviteByTokenEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	return response.(*pb.ShareDataReply), nil
}

func (e Endpoints) AddInvitationCode(ctx context.Context, in *pb.ShareAddInvitationRequest) (*pb.ShareGenericReply, error) {
	response, err := e.AddInvitationCodeEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	return response.(*pb.ShareGenericReply), nil
}

func (e Endpoints) ListFriend(ctx context.Context, in *pb.ShareListFriendRequest) (*pb.ShareListFriendReply, error) {
	response, err := e.ListFriendEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	return response.(*pb.ShareListFriendReply), nil
}

func (e Endpoints) ClaimReward(ctx context.Context, in *pb.ShareClaimRewardRequest) (*pb.ShareGenericReply, error) {
	response, err := e.ClaimRewardEndpoint(ctx, in)
	if err != nil {
		return nil, err
	}
	return response.(*pb.ShareGenericReply), nil
}

// Make Endpoints

func MakeInviteByUrlEndpoint(s pb.ShareServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.ShareEmptyRequest)
		v, err := s.InviteByUrl(ctx, req)
		if err != nil {
			return nil, err
		}
		return v, nil
	}
}

func MakeInviteByTokenEndpoint(s pb.ShareServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.ShareEmptyRequest)
		v, err := s.InviteByToken(ctx, req)
		if err != nil {
			return nil, err
		}
		return v, nil
	}
}

func MakeAddInvitationCodeEndpoint(s pb.ShareServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.ShareAddInvitationRequest)
		v, err := s.AddInvitationCode(ctx, req)
		if err != nil {
			return nil, err
		}
		return v, nil
	}
}

func MakeListFriendEndpoint(s pb.ShareServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.ShareListFriendRequest)
		v, err := s.ListFriend(ctx, req)
		if err != nil {
			return nil, err
		}
		return v, nil
	}
}

func MakeClaimRewardEndpoint(s pb.ShareServer) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*pb.ShareClaimRewardRequest)
		v, err := s.ClaimReward(ctx, req)
		if err != nil {
			return nil, err
		}
		return v, nil
	}
}

// WrapAllExcept wraps each Endpoint field of struct Endpoints with a
// go-kit/kit/endpoint.Middleware.
// Use this for applying a set of middlewares to every endpoint in the service.
// Optionally, endpoints can be passed in by name to be excluded from being wrapped.
// WrapAllExcept(middleware, "Status", "Ping")
func (e *Endpoints) WrapAllExcept(middleware endpoint.Middleware, excluded ...string) {
	included := map[string]struct{}{
		"InviteByUrl":       {},
		"InviteByToken":     {},
		"AddInvitationCode": {},
		"ListFriend":        {},
		"ClaimReward":       {},
	}

	for _, ex := range excluded {
		if _, ok := included[ex]; !ok {
			panic(fmt.Sprintf("Excluded endpoint '%s' does not exist; see middlewares/endpoints.go", ex))
		}
		delete(included, ex)
	}

	for inc := range included {
		if inc == "InviteByUrl" {
			e.InviteByUrlEndpoint = middleware(e.InviteByUrlEndpoint)
		}
		if inc == "InviteByToken" {
			e.InviteByTokenEndpoint = middleware(e.InviteByTokenEndpoint)
		}
		if inc == "AddInvitationCode" {
			e.AddInvitationCodeEndpoint = middleware(e.AddInvitationCodeEndpoint)
		}
		if inc == "ListFriend" {
			e.ListFriendEndpoint = middleware(e.ListFriendEndpoint)
		}
		if inc == "ClaimReward" {
			e.ClaimRewardEndpoint = middleware(e.ClaimRewardEndpoint)
		}
	}
}

// LabeledMiddleware will get passed the endpoint name when passed to
// WrapAllLabeledExcept, this can be used to write a generic metrics
// middleware which can send the endpoint name to the metrics collector.
type LabeledMiddleware func(string, endpoint.Endpoint) endpoint.Endpoint

// WrapAllLabeledExcept wraps each Endpoint field of struct Endpoints with a
// LabeledMiddleware, which will receive the name of the endpoint. See
// LabeldMiddleware. See method WrapAllExept for details on excluded
// functionality.
func (e *Endpoints) WrapAllLabeledExcept(middleware func(string, endpoint.Endpoint) endpoint.Endpoint, excluded ...string) {
	included := map[string]struct{}{
		"InviteByUrl":       {},
		"InviteByToken":     {},
		"AddInvitationCode": {},
		"ListFriend":        {},
		"ClaimReward":       {},
	}

	for _, ex := range excluded {
		if _, ok := included[ex]; !ok {
			panic(fmt.Sprintf("Excluded endpoint '%s' does not exist; see middlewares/endpoints.go", ex))
		}
		delete(included, ex)
	}

	for inc := range included {
		if inc == "InviteByUrl" {
			e.InviteByUrlEndpoint = middleware("InviteByUrl", e.InviteByUrlEndpoint)
		}
		if inc == "InviteByToken" {
			e.InviteByTokenEndpoint = middleware("InviteByToken", e.InviteByTokenEndpoint)
		}
		if inc == "AddInvitationCode" {
			e.AddInvitationCodeEndpoint = middleware("AddInvitationCode", e.AddInvitationCodeEndpoint)
		}
		if inc == "ListFriend" {
			e.ListFriendEndpoint = middleware("ListFriend", e.ListFriendEndpoint)
		}
		if inc == "ClaimReward" {
			e.ClaimRewardEndpoint = middleware("ClaimReward", e.ClaimRewardEndpoint)
		}
	}
}
