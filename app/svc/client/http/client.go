// Code generated by truss. DO NOT EDIT.
// Rerunning truss will overwrite this file.
// Version: 831b290599
// Version Date: 2020-11-16T05:27:36Z

// Package http provides an HTTP client for the App service.
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
	"glab.tagtic.cn/ad_gains/kitty/app/svc"
	pb "glab.tagtic.cn/ad_gains/kitty/proto"
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
func New(instance string, options ...httptransport.ClientOption) (pb.AppServer, error) {

	if !strings.HasPrefix(instance, "http") {
		instance = "http://" + instance
	}
	u, err := url.Parse(instance)
	if err != nil {
		return nil, err
	}
	_ = u

	var LoginZeroEndpoint endpoint.Endpoint
	{
		LoginZeroEndpoint = httptransport.NewClient(
			"POST",
			copyURL(u, "/login"),
			EncodeHTTPLoginZeroRequest,
			DecodeHTTPLoginResponse,
			options...,
		).Endpoint()
	}
	var GetCodeZeroEndpoint endpoint.Endpoint
	{
		GetCodeZeroEndpoint = httptransport.NewClient(
			"GET",
			copyURL(u, "/code"),
			EncodeHTTPGetCodeZeroRequest,
			DecodeHTTPGetCodeResponse,
			options...,
		).Endpoint()
	}
	var GetInfoZeroEndpoint endpoint.Endpoint
	{
		GetInfoZeroEndpoint = httptransport.NewClient(
			"GET",
			copyURL(u, "/info/"),
			EncodeHTTPGetInfoZeroRequest,
			DecodeHTTPGetInfoResponse,
			options...,
		).Endpoint()
	}
	var UpdateInfoZeroEndpoint endpoint.Endpoint
	{
		UpdateInfoZeroEndpoint = httptransport.NewClient(
			"POST",
			copyURL(u, "/info"),
			EncodeHTTPUpdateInfoZeroRequest,
			DecodeHTTPUpdateInfoResponse,
			options...,
		).Endpoint()
	}
	var BindZeroEndpoint endpoint.Endpoint
	{
		BindZeroEndpoint = httptransport.NewClient(
			"POST",
			copyURL(u, "/bind"),
			EncodeHTTPBindZeroRequest,
			DecodeHTTPBindResponse,
			options...,
		).Endpoint()
	}
	var UnbindZeroEndpoint endpoint.Endpoint
	{
		UnbindZeroEndpoint = httptransport.NewClient(
			"POST",
			copyURL(u, "/unbind"),
			EncodeHTTPUnbindZeroRequest,
			DecodeHTTPUnbindResponse,
			options...,
		).Endpoint()
	}
	var RefreshZeroEndpoint endpoint.Endpoint
	{
		RefreshZeroEndpoint = httptransport.NewClient(
			"POST",
			copyURL(u, "/refresh"),
			EncodeHTTPRefreshZeroRequest,
			DecodeHTTPRefreshResponse,
			options...,
		).Endpoint()
	}

	return svc.Endpoints{
		LoginEndpoint:      LoginZeroEndpoint,
		GetCodeEndpoint:    GetCodeZeroEndpoint,
		GetInfoEndpoint:    GetInfoZeroEndpoint,
		UpdateInfoEndpoint: UpdateInfoZeroEndpoint,
		BindEndpoint:       BindZeroEndpoint,
		UnbindEndpoint:     UnbindZeroEndpoint,
		RefreshEndpoint:    RefreshZeroEndpoint,
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

// DecodeHTTPLoginResponse is a transport/http.DecodeResponseFunc that decodes
// a JSON-encoded UserInfoReply response from the HTTP response body.
// If the response has a non-200 status code, we will interpret that as an
// error and attempt to decode the specific error message from the response
// body. Primarily useful in a client.
func DecodeHTTPLoginResponse(_ context.Context, r *http.Response) (interface{}, error) {
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

	var resp pb.UserInfoReply
	if err = jsonpb.UnmarshalString(string(buf), &resp); err != nil {
		return nil, errorDecoder(buf)
	}

	return &resp, nil
}

// DecodeHTTPGetCodeResponse is a transport/http.DecodeResponseFunc that decodes
// a JSON-encoded GenericReply response from the HTTP response body.
// If the response has a non-200 status code, we will interpret that as an
// error and attempt to decode the specific error message from the response
// body. Primarily useful in a client.
func DecodeHTTPGetCodeResponse(_ context.Context, r *http.Response) (interface{}, error) {
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

	var resp pb.GenericReply
	if err = jsonpb.UnmarshalString(string(buf), &resp); err != nil {
		return nil, errorDecoder(buf)
	}

	return &resp, nil
}

// DecodeHTTPGetInfoResponse is a transport/http.DecodeResponseFunc that decodes
// a JSON-encoded UserInfoReply response from the HTTP response body.
// If the response has a non-200 status code, we will interpret that as an
// error and attempt to decode the specific error message from the response
// body. Primarily useful in a client.
func DecodeHTTPGetInfoResponse(_ context.Context, r *http.Response) (interface{}, error) {
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

	var resp pb.UserInfoReply
	if err = jsonpb.UnmarshalString(string(buf), &resp); err != nil {
		return nil, errorDecoder(buf)
	}

	return &resp, nil
}

// DecodeHTTPUpdateInfoResponse is a transport/http.DecodeResponseFunc that decodes
// a JSON-encoded UserInfoReply response from the HTTP response body.
// If the response has a non-200 status code, we will interpret that as an
// error and attempt to decode the specific error message from the response
// body. Primarily useful in a client.
func DecodeHTTPUpdateInfoResponse(_ context.Context, r *http.Response) (interface{}, error) {
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

	var resp pb.UserInfoReply
	if err = jsonpb.UnmarshalString(string(buf), &resp); err != nil {
		return nil, errorDecoder(buf)
	}

	return &resp, nil
}

// DecodeHTTPBindResponse is a transport/http.DecodeResponseFunc that decodes
// a JSON-encoded UserInfoReply response from the HTTP response body.
// If the response has a non-200 status code, we will interpret that as an
// error and attempt to decode the specific error message from the response
// body. Primarily useful in a client.
func DecodeHTTPBindResponse(_ context.Context, r *http.Response) (interface{}, error) {
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

	var resp pb.UserInfoReply
	if err = jsonpb.UnmarshalString(string(buf), &resp); err != nil {
		return nil, errorDecoder(buf)
	}

	return &resp, nil
}

// DecodeHTTPUnbindResponse is a transport/http.DecodeResponseFunc that decodes
// a JSON-encoded UserInfoReply response from the HTTP response body.
// If the response has a non-200 status code, we will interpret that as an
// error and attempt to decode the specific error message from the response
// body. Primarily useful in a client.
func DecodeHTTPUnbindResponse(_ context.Context, r *http.Response) (interface{}, error) {
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

	var resp pb.UserInfoReply
	if err = jsonpb.UnmarshalString(string(buf), &resp); err != nil {
		return nil, errorDecoder(buf)
	}

	return &resp, nil
}

// DecodeHTTPRefreshResponse is a transport/http.DecodeResponseFunc that decodes
// a JSON-encoded UserInfoReply response from the HTTP response body.
// If the response has a non-200 status code, we will interpret that as an
// error and attempt to decode the specific error message from the response
// body. Primarily useful in a client.
func DecodeHTTPRefreshResponse(_ context.Context, r *http.Response) (interface{}, error) {
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

	var resp pb.UserInfoReply
	if err = jsonpb.UnmarshalString(string(buf), &resp); err != nil {
		return nil, errorDecoder(buf)
	}

	return &resp, nil
}

// HTTP Client Encode

// EncodeHTTPLoginZeroRequest is a transport/http.EncodeRequestFunc
// that encodes a login request into the various portions of
// the http request (path, query, and body).
func EncodeHTTPLoginZeroRequest(_ context.Context, r *http.Request, request interface{}) error {
	strval := ""
	_ = strval
	req := request.(*pb.UserLoginRequest)
	_ = req

	r.Header.Set("transport", "HTTPJSON")
	r.Header.Set("request-url", r.URL.Path)

	// Set the path parameters
	path := strings.Join([]string{
		"",
		"login",
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
	toRet := request.(*pb.UserLoginRequest)

	toRet.Mobile = req.Mobile

	toRet.Code = req.Code

	toRet.Wechat = req.Wechat

	toRet.Device = req.Device

	toRet.Channel = req.Channel

	toRet.VersionCode = req.VersionCode

	toRet.PackageName = req.PackageName

	toRet.ThirdPartyId = req.ThirdPartyId

	toRet.InviteCode = req.InviteCode

	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(toRet); err != nil {
		return errors.Wrapf(err, "couldn't encode body as json %v", toRet)
	}
	r.Body = ioutil.NopCloser(&buf)
	return nil
}

// EncodeHTTPGetCodeZeroRequest is a transport/http.EncodeRequestFunc
// that encodes a getcode request into the various portions of
// the http request (path, query, and body).
func EncodeHTTPGetCodeZeroRequest(_ context.Context, r *http.Request, request interface{}) error {
	strval := ""
	_ = strval
	req := request.(*pb.GetCodeRequest)
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

	values.Add("mobile", fmt.Sprint(req.Mobile))

	values.Add("packageName", fmt.Sprint(req.PackageName))

	r.URL.RawQuery = values.Encode()
	return nil
}

// EncodeHTTPGetInfoZeroRequest is a transport/http.EncodeRequestFunc
// that encodes a getinfo request into the various portions of
// the http request (path, query, and body).
func EncodeHTTPGetInfoZeroRequest(_ context.Context, r *http.Request, request interface{}) error {
	strval := ""
	_ = strval
	req := request.(*pb.UserInfoRequest)
	_ = req

	r.Header.Set("transport", "HTTPJSON")
	r.Header.Set("request-url", r.URL.Path)

	// Set the path parameters
	path := strings.Join([]string{
		"",
		"info",
		fmt.Sprint(req.Id),
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

	values.Add("wechat", fmt.Sprint(req.Wechat))

	values.Add("taobao", fmt.Sprint(req.Taobao))

	r.URL.RawQuery = values.Encode()
	return nil
}

// EncodeHTTPUpdateInfoZeroRequest is a transport/http.EncodeRequestFunc
// that encodes a updateinfo request into the various portions of
// the http request (path, query, and body).
func EncodeHTTPUpdateInfoZeroRequest(_ context.Context, r *http.Request, request interface{}) error {
	strval := ""
	_ = strval
	req := request.(*pb.UserInfoUpdateRequest)
	_ = req

	r.Header.Set("transport", "HTTPJSON")
	r.Header.Set("request-url", r.URL.Path)

	// Set the path parameters
	path := strings.Join([]string{
		"",
		"info",
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
	toRet := request.(*pb.UserInfoUpdateRequest)

	toRet.UserName = req.UserName

	toRet.HeadImg = req.HeadImg

	toRet.Gender = req.Gender

	toRet.Birthday = req.Birthday

	toRet.ThirdPartyId = req.ThirdPartyId

	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(toRet); err != nil {
		return errors.Wrapf(err, "couldn't encode body as json %v", toRet)
	}
	r.Body = ioutil.NopCloser(&buf)
	return nil
}

// EncodeHTTPBindZeroRequest is a transport/http.EncodeRequestFunc
// that encodes a bind request into the various portions of
// the http request (path, query, and body).
func EncodeHTTPBindZeroRequest(_ context.Context, r *http.Request, request interface{}) error {
	strval := ""
	_ = strval
	req := request.(*pb.UserBindRequest)
	_ = req

	r.Header.Set("transport", "HTTPJSON")
	r.Header.Set("request-url", r.URL.Path)

	// Set the path parameters
	path := strings.Join([]string{
		"",
		"bind",
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
	toRet := request.(*pb.UserBindRequest)

	toRet.Mobile = req.Mobile

	toRet.Code = req.Code

	toRet.Wechat = req.Wechat

	toRet.OpenId = req.OpenId

	toRet.TaobaoExtra = req.TaobaoExtra

	toRet.WechatExtra = req.WechatExtra

	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(toRet); err != nil {
		return errors.Wrapf(err, "couldn't encode body as json %v", toRet)
	}
	r.Body = ioutil.NopCloser(&buf)
	return nil
}

// EncodeHTTPUnbindZeroRequest is a transport/http.EncodeRequestFunc
// that encodes a unbind request into the various portions of
// the http request (path, query, and body).
func EncodeHTTPUnbindZeroRequest(_ context.Context, r *http.Request, request interface{}) error {
	strval := ""
	_ = strval
	req := request.(*pb.UserUnbindRequest)
	_ = req

	r.Header.Set("transport", "HTTPJSON")
	r.Header.Set("request-url", r.URL.Path)

	// Set the path parameters
	path := strings.Join([]string{
		"",
		"unbind",
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
	toRet := request.(*pb.UserUnbindRequest)

	toRet.Mobile = req.Mobile

	toRet.Wechat = req.Wechat

	toRet.Taobao = req.Taobao

	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(toRet); err != nil {
		return errors.Wrapf(err, "couldn't encode body as json %v", toRet)
	}
	r.Body = ioutil.NopCloser(&buf)
	return nil
}

// EncodeHTTPRefreshZeroRequest is a transport/http.EncodeRequestFunc
// that encodes a refresh request into the various portions of
// the http request (path, query, and body).
func EncodeHTTPRefreshZeroRequest(_ context.Context, r *http.Request, request interface{}) error {
	strval := ""
	_ = strval
	req := request.(*pb.UserRefreshRequest)
	_ = req

	r.Header.Set("transport", "HTTPJSON")
	r.Header.Set("request-url", r.URL.Path)

	// Set the path parameters
	path := strings.Join([]string{
		"",
		"refresh",
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
	toRet := request.(*pb.UserRefreshRequest)

	toRet.Device = req.Device

	toRet.Channel = req.Channel

	toRet.VersionCode = req.VersionCode

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
