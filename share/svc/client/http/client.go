// Code generated by truss. DO NOT EDIT.
// Rerunning truss will overwrite this file.
// Version: 1a83e32aba
// Version Date: 2020-12-03T10:58:42Z

// Package http provides an HTTP client for the Share service.
package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/golang/protobuf/jsonpb"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/pkg/errors"

	// This Service
	pb "glab.tagtic.cn/ad_gains/kitty/proto"
	"glab.tagtic.cn/ad_gains/kitty/share/svc"
)

var (
	_ = endpoint.Chain
	_ = httptransport.NewClient
	_ = fmt.Sprint
	_ = bytes.Compare
	_ = ioutil.NopCloser
)

// New returns a service backed by an HTTP server living at the remote
// instance. We expect instance to come from a service discovery system, so
// likely of the form "host:port".
func New(instance string, options ...httptransport.ClientOption) (pb.ShareServer, error) {

	if !strings.HasPrefix(instance, "http") {
		instance = "http://" + instance
	}
	u, err := url.Parse(instance)
	if err != nil {
		return nil, err
	}
	_ = u

	var InviteByUrlZeroEndpoint endpoint.Endpoint
	{
		InviteByUrlZeroEndpoint = httptransport.NewClient(
			"GET",
			copyURL(u, "/url"),
			EncodeHTTPInviteByUrlZeroRequest,
			DecodeHTTPInviteByUrlResponse,
			options...,
		).Endpoint()
	}
	var InviteByTokenZeroEndpoint endpoint.Endpoint
	{
		InviteByTokenZeroEndpoint = httptransport.NewClient(
			"GET",
			copyURL(u, "/code"),
			EncodeHTTPInviteByTokenZeroRequest,
			DecodeHTTPInviteByTokenResponse,
			options...,
		).Endpoint()
	}
	var AddInvitationCodeZeroEndpoint endpoint.Endpoint
	{
		AddInvitationCodeZeroEndpoint = httptransport.NewClient(
			"PUT",
			copyURL(u, "/code"),
			EncodeHTTPAddInvitationCodeZeroRequest,
			DecodeHTTPAddInvitationCodeResponse,
			options...,
		).Endpoint()
	}
	var ListFriendZeroEndpoint endpoint.Endpoint
	{
		ListFriendZeroEndpoint = httptransport.NewClient(
			"GET",
			copyURL(u, "/list"),
			EncodeHTTPListFriendZeroRequest,
			DecodeHTTPListFriendResponse,
			options...,
		).Endpoint()
	}
	var ClaimRewardZeroEndpoint endpoint.Endpoint
	{
		ClaimRewardZeroEndpoint = httptransport.NewClient(
			"POST",
			copyURL(u, "/claim"),
			EncodeHTTPClaimRewardZeroRequest,
			DecodeHTTPClaimRewardResponse,
			options...,
		).Endpoint()
	}
	var PushSignEventZeroEndpoint endpoint.Endpoint
	{
		PushSignEventZeroEndpoint = httptransport.NewClient(
			"POST",
			copyURL(u, "/event/sign"),
			EncodeHTTPPushSignEventZeroRequest,
			DecodeHTTPPushSignEventResponse,
			options...,
		).Endpoint()
	}
	var PushTaskEventZeroEndpoint endpoint.Endpoint
	{
		PushTaskEventZeroEndpoint = httptransport.NewClient(
			"POST",
			copyURL(u, "/event/task"),
			EncodeHTTPPushTaskEventZeroRequest,
			DecodeHTTPPushTaskEventResponse,
			options...,
		).Endpoint()
	}

	return svc.Endpoints{
		InviteByUrlEndpoint:       InviteByUrlZeroEndpoint,
		InviteByTokenEndpoint:     InviteByTokenZeroEndpoint,
		AddInvitationCodeEndpoint: AddInvitationCodeZeroEndpoint,
		ListFriendEndpoint:        ListFriendZeroEndpoint,
		ClaimRewardEndpoint:       ClaimRewardZeroEndpoint,
		PushSignEventEndpoint:     PushSignEventZeroEndpoint,
		PushTaskEventEndpoint:     PushTaskEventZeroEndpoint,
	}, nil
}

func copyURL(base *url.URL, path string) *url.URL {
	next := *base
	next.Path = path
	return &next
}

// CtxValuesToSend configures the http client to pull the specified keys out of
// the context and add them to the http request as headers.  Note that keys
// will have net/http.CanonicalHeaderKey called on them before being send over
// the wire and that is the form they will be available in the server context.
func CtxValuesToSend(keys ...string) httptransport.ClientOption {
	return httptransport.ClientBefore(func(ctx context.Context, r *http.Request) context.Context {
		for _, k := range keys {
			if v, ok := ctx.Value(k).(string); ok {
				r.Header.Set(k, v)
			}
		}
		return ctx
	})
}

// HTTP Client Decode

// DecodeHTTPInviteByUrlResponse is a transport/http.DecodeResponseFunc that decodes
// a JSON-encoded ShareDataUrlReply response from the HTTP response body.
// If the response has a non-200 status code, we will interpret that as an
// error and attempt to decode the specific error message from the response
// body. Primarily useful in a client.
func DecodeHTTPInviteByUrlResponse(_ context.Context, r *http.Response) (interface{}, error) {
	defer r.Body.Close()
	buf, err := ioutil.ReadAll(r.Body)
	if err == io.EOF {
		return nil, errors.New("response http body empty")
	}
	if err != nil {
		return nil, errors.Wrap(err, "cannot read http body")
	}

	if r.StatusCode != http.StatusOK {
		return nil, errors.Wrapf(errorDecoder(buf), "status code: '%d'", r.StatusCode)
	}

	var resp pb.ShareDataUrlReply
	if err = jsonpb.UnmarshalString(string(buf), &resp); err != nil {
		return nil, errorDecoder(buf)
	}

	return &resp, nil
}

// DecodeHTTPInviteByTokenResponse is a transport/http.DecodeResponseFunc that decodes
// a JSON-encoded ShareDataTokenReply response from the HTTP response body.
// If the response has a non-200 status code, we will interpret that as an
// error and attempt to decode the specific error message from the response
// body. Primarily useful in a client.
func DecodeHTTPInviteByTokenResponse(_ context.Context, r *http.Response) (interface{}, error) {
	defer r.Body.Close()
	buf, err := ioutil.ReadAll(r.Body)
	if err == io.EOF {
		return nil, errors.New("response http body empty")
	}
	if err != nil {
		return nil, errors.Wrap(err, "cannot read http body")
	}

	if r.StatusCode != http.StatusOK {
		return nil, errors.Wrapf(errorDecoder(buf), "status code: '%d'", r.StatusCode)
	}

	var resp pb.ShareDataTokenReply
	if err = jsonpb.UnmarshalString(string(buf), &resp); err != nil {
		return nil, errorDecoder(buf)
	}

	return &resp, nil
}

// DecodeHTTPAddInvitationCodeResponse is a transport/http.DecodeResponseFunc that decodes
// a JSON-encoded ShareGenericReply response from the HTTP response body.
// If the response has a non-200 status code, we will interpret that as an
// error and attempt to decode the specific error message from the response
// body. Primarily useful in a client.
func DecodeHTTPAddInvitationCodeResponse(_ context.Context, r *http.Response) (interface{}, error) {
	defer r.Body.Close()
	buf, err := ioutil.ReadAll(r.Body)
	if err == io.EOF {
		return nil, errors.New("response http body empty")
	}
	if err != nil {
		return nil, errors.Wrap(err, "cannot read http body")
	}

	if r.StatusCode != http.StatusOK {
		return nil, errors.Wrapf(errorDecoder(buf), "status code: '%d'", r.StatusCode)
	}

	var resp pb.ShareGenericReply
	if err = jsonpb.UnmarshalString(string(buf), &resp); err != nil {
		return nil, errorDecoder(buf)
	}

	return &resp, nil
}

// DecodeHTTPListFriendResponse is a transport/http.DecodeResponseFunc that decodes
// a JSON-encoded ShareListFriendReply response from the HTTP response body.
// If the response has a non-200 status code, we will interpret that as an
// error and attempt to decode the specific error message from the response
// body. Primarily useful in a client.
func DecodeHTTPListFriendResponse(_ context.Context, r *http.Response) (interface{}, error) {
	defer r.Body.Close()
	buf, err := ioutil.ReadAll(r.Body)
	if err == io.EOF {
		return nil, errors.New("response http body empty")
	}
	if err != nil {
		return nil, errors.Wrap(err, "cannot read http body")
	}

	if r.StatusCode != http.StatusOK {
		return nil, errors.Wrapf(errorDecoder(buf), "status code: '%d'", r.StatusCode)
	}

	var resp pb.ShareListFriendReply
	if err = jsonpb.UnmarshalString(string(buf), &resp); err != nil {
		return nil, errorDecoder(buf)
	}

	return &resp, nil
}

// DecodeHTTPClaimRewardResponse is a transport/http.DecodeResponseFunc that decodes
// a JSON-encoded ShareGenericReply response from the HTTP response body.
// If the response has a non-200 status code, we will interpret that as an
// error and attempt to decode the specific error message from the response
// body. Primarily useful in a client.
func DecodeHTTPClaimRewardResponse(_ context.Context, r *http.Response) (interface{}, error) {
	defer r.Body.Close()
	buf, err := ioutil.ReadAll(r.Body)
	if err == io.EOF {
		return nil, errors.New("response http body empty")
	}
	if err != nil {
		return nil, errors.Wrap(err, "cannot read http body")
	}

	if r.StatusCode != http.StatusOK {
		return nil, errors.Wrapf(errorDecoder(buf), "status code: '%d'", r.StatusCode)
	}

	var resp pb.ShareGenericReply
	if err = jsonpb.UnmarshalString(string(buf), &resp); err != nil {
		return nil, errorDecoder(buf)
	}

	return &resp, nil
}

// DecodeHTTPPushSignEventResponse is a transport/http.DecodeResponseFunc that decodes
// a JSON-encoded ShareGenericReply response from the HTTP response body.
// If the response has a non-200 status code, we will interpret that as an
// error and attempt to decode the specific error message from the response
// body. Primarily useful in a client.
func DecodeHTTPPushSignEventResponse(_ context.Context, r *http.Response) (interface{}, error) {
	defer r.Body.Close()
	buf, err := ioutil.ReadAll(r.Body)
	if err == io.EOF {
		return nil, errors.New("response http body empty")
	}
	if err != nil {
		return nil, errors.Wrap(err, "cannot read http body")
	}

	if r.StatusCode != http.StatusOK {
		return nil, errors.Wrapf(errorDecoder(buf), "status code: '%d'", r.StatusCode)
	}

	var resp pb.ShareGenericReply
	if err = jsonpb.UnmarshalString(string(buf), &resp); err != nil {
		return nil, errorDecoder(buf)
	}

	return &resp, nil
}

// DecodeHTTPPushTaskEventResponse is a transport/http.DecodeResponseFunc that decodes
// a JSON-encoded ShareGenericReply response from the HTTP response body.
// If the response has a non-200 status code, we will interpret that as an
// error and attempt to decode the specific error message from the response
// body. Primarily useful in a client.
func DecodeHTTPPushTaskEventResponse(_ context.Context, r *http.Response) (interface{}, error) {
	defer r.Body.Close()
	buf, err := ioutil.ReadAll(r.Body)
	if err == io.EOF {
		return nil, errors.New("response http body empty")
	}
	if err != nil {
		return nil, errors.Wrap(err, "cannot read http body")
	}

	if r.StatusCode != http.StatusOK {
		return nil, errors.Wrapf(errorDecoder(buf), "status code: '%d'", r.StatusCode)
	}

	var resp pb.ShareGenericReply
	if err = jsonpb.UnmarshalString(string(buf), &resp); err != nil {
		return nil, errorDecoder(buf)
	}

	return &resp, nil
}

// HTTP Client Encode

// EncodeHTTPInviteByUrlZeroRequest is a transport/http.EncodeRequestFunc
// that encodes a invitebyurl request into the various portions of
// the http request (path, query, and body).
func EncodeHTTPInviteByUrlZeroRequest(_ context.Context, r *http.Request, request interface{}) error {
	strval := ""
	_ = strval
	req := request.(*pb.ShareEmptyRequest)
	_ = req

	r.Header.Set("transport", "HTTPJSON")
	r.Header.Set("request-url", r.URL.Path)

	// Set the path parameters
	path := strings.Join([]string{
		"",
		"url",
	}, "/")
	u, err := url.Parse(path)
	if err != nil {
		return errors.Wrapf(err, "couldn't unmarshal path %q", path)
	}
	r.URL.RawPath = u.RawPath
	r.URL.Path = u.Path

	// Set the query parameters
	values := r.URL.Query()
	var tmp []byte
	_ = tmp

	r.URL.RawQuery = values.Encode()
	return nil
}

// EncodeHTTPInviteByTokenZeroRequest is a transport/http.EncodeRequestFunc
// that encodes a invitebytoken request into the various portions of
// the http request (path, query, and body).
func EncodeHTTPInviteByTokenZeroRequest(_ context.Context, r *http.Request, request interface{}) error {
	strval := ""
	_ = strval
	req := request.(*pb.ShareEmptyRequest)
	_ = req

	r.Header.Set("transport", "HTTPJSON")
	r.Header.Set("request-url", r.URL.Path)

	// Set the path parameters
	path := strings.Join([]string{
		"",
		"code",
	}, "/")
	u, err := url.Parse(path)
	if err != nil {
		return errors.Wrapf(err, "couldn't unmarshal path %q", path)
	}
	r.URL.RawPath = u.RawPath
	r.URL.Path = u.Path

	// Set the query parameters
	values := r.URL.Query()
	var tmp []byte
	_ = tmp

	r.URL.RawQuery = values.Encode()
	return nil
}

// EncodeHTTPAddInvitationCodeZeroRequest is a transport/http.EncodeRequestFunc
// that encodes a addinvitationcode request into the various portions of
// the http request (path, query, and body).
func EncodeHTTPAddInvitationCodeZeroRequest(_ context.Context, r *http.Request, request interface{}) error {
	strval := ""
	_ = strval
	req := request.(*pb.ShareAddInvitationRequest)
	_ = req

	r.Header.Set("transport", "HTTPJSON")
	r.Header.Set("request-url", r.URL.Path)

	// Set the path parameters
	path := strings.Join([]string{
		"",
		"code",
	}, "/")
	u, err := url.Parse(path)
	if err != nil {
		return errors.Wrapf(err, "couldn't unmarshal path %q", path)
	}
	r.URL.RawPath = u.RawPath
	r.URL.Path = u.Path

	// Set the query parameters
	values := r.URL.Query()
	var tmp []byte
	_ = tmp

	r.URL.RawQuery = values.Encode()
	// Set the body parameters
	var buf bytes.Buffer
	toRet := request.(*pb.ShareAddInvitationRequest)

	toRet.InviteCode = req.InviteCode

	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(toRet); err != nil {
		return errors.Wrapf(err, "couldn't encode body as json %v", toRet)
	}
	r.Body = ioutil.NopCloser(&buf)
	return nil
}

// EncodeHTTPListFriendZeroRequest is a transport/http.EncodeRequestFunc
// that encodes a listfriend request into the various portions of
// the http request (path, query, and body).
func EncodeHTTPListFriendZeroRequest(_ context.Context, r *http.Request, request interface{}) error {
	strval := ""
	_ = strval
	req := request.(*pb.ShareListFriendRequest)
	_ = req

	r.Header.Set("transport", "HTTPJSON")
	r.Header.Set("request-url", r.URL.Path)

	// Set the path parameters
	path := strings.Join([]string{
		"",
		"list",
	}, "/")
	u, err := url.Parse(path)
	if err != nil {
		return errors.Wrapf(err, "couldn't unmarshal path %q", path)
	}
	r.URL.RawPath = u.RawPath
	r.URL.Path = u.Path

	// Set the query parameters
	values := r.URL.Query()
	var tmp []byte
	_ = tmp

	values.Add("depth", fmt.Sprint(req.Depth))

	r.URL.RawQuery = values.Encode()
	return nil
}

// EncodeHTTPClaimRewardZeroRequest is a transport/http.EncodeRequestFunc
// that encodes a claimreward request into the various portions of
// the http request (path, query, and body).
func EncodeHTTPClaimRewardZeroRequest(_ context.Context, r *http.Request, request interface{}) error {
	strval := ""
	_ = strval
	req := request.(*pb.ShareClaimRewardRequest)
	_ = req

	r.Header.Set("transport", "HTTPJSON")
	r.Header.Set("request-url", r.URL.Path)

	// Set the path parameters
	path := strings.Join([]string{
		"",
		"claim",
	}, "/")
	u, err := url.Parse(path)
	if err != nil {
		return errors.Wrapf(err, "couldn't unmarshal path %q", path)
	}
	r.URL.RawPath = u.RawPath
	r.URL.Path = u.Path

	// Set the query parameters
	values := r.URL.Query()
	var tmp []byte
	_ = tmp

	r.URL.RawQuery = values.Encode()
	// Set the body parameters
	var buf bytes.Buffer
	toRet := request.(*pb.ShareClaimRewardRequest)

	toRet.ApprenticeId = req.ApprenticeId

	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(toRet); err != nil {
		return errors.Wrapf(err, "couldn't encode body as json %v", toRet)
	}
	r.Body = ioutil.NopCloser(&buf)
	return nil
}

// EncodeHTTPPushSignEventZeroRequest is a transport/http.EncodeRequestFunc
// that encodes a pushsignevent request into the various portions of
// the http request (path, query, and body).
func EncodeHTTPPushSignEventZeroRequest(_ context.Context, r *http.Request, request interface{}) error {
	strval := ""
	_ = strval
	req := request.(*pb.SignEvent)
	_ = req

	r.Header.Set("transport", "HTTPJSON")
	r.Header.Set("request-url", r.URL.Path)

	// Set the path parameters
	path := strings.Join([]string{
		"",
		"event",
		"sign",
	}, "/")
	u, err := url.Parse(path)
	if err != nil {
		return errors.Wrapf(err, "couldn't unmarshal path %q", path)
	}
	r.URL.RawPath = u.RawPath
	r.URL.Path = u.Path

	// Set the query parameters
	values := r.URL.Query()
	var tmp []byte
	_ = tmp

	r.URL.RawQuery = values.Encode()
	// Set the body parameters
	var buf bytes.Buffer
	toRet := request.(*pb.SignEvent)

	toRet.Id = req.Id

	toRet.UserId = req.UserId

	toRet.PackageName = req.PackageName

	toRet.Channel = req.Channel

	toRet.Name = req.Name

	toRet.EventName = req.EventName

	toRet.Score = req.Score

	toRet.DateTime = req.DateTime

	toRet.ThirdPartyId = req.ThirdPartyId

	toRet.IsDouble = req.IsDouble

	toRet.Ext = req.Ext

	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(toRet); err != nil {
		return errors.Wrapf(err, "couldn't encode body as json %v", toRet)
	}
	r.Body = ioutil.NopCloser(&buf)
	return nil
}

// EncodeHTTPPushTaskEventZeroRequest is a transport/http.EncodeRequestFunc
// that encodes a pushtaskevent request into the various portions of
// the http request (path, query, and body).
func EncodeHTTPPushTaskEventZeroRequest(_ context.Context, r *http.Request, request interface{}) error {
	strval := ""
	_ = strval
	req := request.(*pb.TaskEvent)
	_ = req

	r.Header.Set("transport", "HTTPJSON")
	r.Header.Set("request-url", r.URL.Path)

	// Set the path parameters
	path := strings.Join([]string{
		"",
		"event",
		"task",
	}, "/")
	u, err := url.Parse(path)
	if err != nil {
		return errors.Wrapf(err, "couldn't unmarshal path %q", path)
	}
	r.URL.RawPath = u.RawPath
	r.URL.Path = u.Path

	// Set the query parameters
	values := r.URL.Query()
	var tmp []byte
	_ = tmp

	r.URL.RawQuery = values.Encode()
	// Set the body parameters
	var buf bytes.Buffer
	toRet := request.(*pb.TaskEvent)

	toRet.Id = req.Id

	toRet.UserId = req.UserId

	toRet.PackageName = req.PackageName

	toRet.Channel = req.Channel

	toRet.Name = req.Name

	toRet.EventName = req.EventName

	toRet.Score = req.Score

	toRet.DateTime = req.DateTime

	toRet.ThirdPartyId = req.ThirdPartyId

	toRet.DoneNum = req.DoneNum

	toRet.TotalNum = req.TotalNum

	toRet.IsDone = req.IsDone

	toRet.Ext = req.Ext

	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(toRet); err != nil {
		return errors.Wrapf(err, "couldn't encode body as json %v", toRet)
	}
	r.Body = ioutil.NopCloser(&buf)
	return nil
}

func errorDecoder(buf []byte) error {
	var w errorWrapper
	if err := json.Unmarshal(buf, &w); err != nil {
		const size = 8196
		if len(buf) > size {
			buf = buf[:size]
		}
		return fmt.Errorf("response body '%s': cannot parse non-json request body", buf)
	}

	return errors.New(w.Error)
}

type errorWrapper struct {
	Error string `json:"error"`
}
