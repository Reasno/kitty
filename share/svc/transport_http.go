// Code generated by truss. DO NOT EDIT.
// Rerunning truss will overwrite this file.
// Version: 831b290599
// Version Date: 2020-11-16T05:27:36Z

package svc

// This file provides server-side bindings for the HTTP transport.
// It utilizes the transport/http.Server.

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/gogo/protobuf/proto"
	"github.com/golang/protobuf/jsonpb"

	"context"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	// This service
	pb "glab.tagtic.cn/ad_gains/kitty/proto"
)

const contentType = "application/json; charset=utf-8"

var (
	_ = fmt.Sprint
	_ = bytes.Compare
	_ = strconv.Atoi
	_ = httptransport.NewServer
	_ = ioutil.NopCloser
	_ = pb.NewShareClient
	_ = io.Copy
	_ = errors.Wrap
)

// MakeHTTPHandler returns a handler that makes a set of endpoints available
// on predefined paths.
func MakeHTTPHandler(endpoints Endpoints, options ...httptransport.ServerOption) http.Handler {
	serverOptions := []httptransport.ServerOption{
		httptransport.ServerBefore(headersToContext),
		httptransport.ServerErrorEncoder(errorEncoder),
		httptransport.ServerAfter(httptransport.SetContentType(contentType)),
	}
	serverOptions = append(serverOptions, options...)
	m := mux.NewRouter()

	m.Methods("GET").Path("/url").Handler(httptransport.NewServer(
		endpoints.InviteByUrlEndpoint,
		DecodeHTTPInviteByUrlZeroRequest,
		EncodeHTTPGenericResponse,
		serverOptions...,
	))

	m.Methods("GET").Path("/code").Handler(httptransport.NewServer(
		endpoints.InviteByTokenEndpoint,
		DecodeHTTPInviteByTokenZeroRequest,
		EncodeHTTPGenericResponse,
		serverOptions...,
	))

	m.Methods("PUT").Path("/code").Handler(httptransport.NewServer(
		endpoints.AddInvitationCodeEndpoint,
		DecodeHTTPAddInvitationCodeZeroRequest,
		EncodeHTTPGenericResponse,
		serverOptions...,
	))

	m.Methods("GET").Path("/list").Handler(httptransport.NewServer(
		endpoints.ListFriendEndpoint,
		DecodeHTTPListFriendZeroRequest,
		EncodeHTTPGenericResponse,
		serverOptions...,
	))

	m.Methods("POST").Path("/claim").Handler(httptransport.NewServer(
		endpoints.ClaimRewardEndpoint,
		DecodeHTTPClaimRewardZeroRequest,
		EncodeHTTPGenericResponse,
		serverOptions...,
	))

	m.Methods("POST").Path("/event/sign").Handler(httptransport.NewServer(
		endpoints.PushSignEventEndpoint,
		DecodeHTTPPushSignEventZeroRequest,
		EncodeHTTPGenericResponse,
		serverOptions...,
	))

	m.Methods("POST").Path("/event/task").Handler(httptransport.NewServer(
		endpoints.PushTaskEventEndpoint,
		DecodeHTTPPushTaskEventZeroRequest,
		EncodeHTTPGenericResponse,
		serverOptions...,
	))
	return m
}

// ErrorEncoder writes the error to the ResponseWriter, by default a content
// type of application/json, a body of json with key "error" and the value
// error.Error(), and a status code of 500. If the error implements Headerer,
// the provided headers will be applied to the response. If the error
// implements json.Marshaler, and the marshaling succeeds, the JSON encoded
// form of the error will be used. If the error implements StatusCoder, the
// provided StatusCode will be used instead of 500.
func errorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	body, _ := json.Marshal(errorWrapper{Error: err.Error()})
	if marshaler, ok := err.(json.Marshaler); ok {
		if jsonBody, marshalErr := marshaler.MarshalJSON(); marshalErr == nil {
			body = jsonBody
		}
	}
	w.Header().Set("Content-Type", contentType)
	if headerer, ok := err.(httptransport.Headerer); ok {
		for k := range headerer.Headers() {
			w.Header().Set(k, headerer.Headers().Get(k))
		}
	}
	code := http.StatusInternalServerError
	if sc, ok := err.(httptransport.StatusCoder); ok {
		code = sc.StatusCode()
	}
	w.WriteHeader(code)
	w.Write(body)
}

type errorWrapper struct {
	Error string `json:"error"`
}

// httpError satisfies the Headerer and StatusCoder interfaces in
// package github.com/go-kit/kit/transport/http.
type httpError struct {
	error
	statusCode int
	headers    map[string][]string
}

func (h httpError) StatusCode() int {
	return h.statusCode
}

func (h httpError) Headers() http.Header {
	return h.headers
}

// Server Decode

// DecodeHTTPInviteByUrlZeroRequest is a transport/http.DecodeRequestFunc that
// decodes a JSON-encoded invitebyurl request from the HTTP request
// body. Primarily useful in a server.
func DecodeHTTPInviteByUrlZeroRequest(_ context.Context, r *http.Request) (interface{}, error) {
	defer r.Body.Close()
	var req pb.ShareEmptyRequest
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot read body of http request")
	}
	if len(buf) > 0 {
		// AllowUnknownFields stops the unmarshaler from failing if the JSON contains unknown fields.
		unmarshaller := jsonpb.Unmarshaler{
			AllowUnknownFields: true,
		}
		if err = unmarshaller.Unmarshal(bytes.NewBuffer(buf), &req); err != nil {
			const size = 8196
			if len(buf) > size {
				buf = buf[:size]
			}
			return nil, httpError{errors.Wrapf(err, "request body '%s': cannot parse non-json request body", buf),
				http.StatusBadRequest,
				nil,
			}
		}
	}

	pathParams := mux.Vars(r)
	_ = pathParams

	queryParams := r.URL.Query()
	_ = queryParams

	return &req, err
}

// DecodeHTTPInviteByTokenZeroRequest is a transport/http.DecodeRequestFunc that
// decodes a JSON-encoded invitebytoken request from the HTTP request
// body. Primarily useful in a server.
func DecodeHTTPInviteByTokenZeroRequest(_ context.Context, r *http.Request) (interface{}, error) {
	defer r.Body.Close()
	var req pb.ShareEmptyRequest
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot read body of http request")
	}
	if len(buf) > 0 {
		// AllowUnknownFields stops the unmarshaler from failing if the JSON contains unknown fields.
		unmarshaller := jsonpb.Unmarshaler{
			AllowUnknownFields: true,
		}
		if err = unmarshaller.Unmarshal(bytes.NewBuffer(buf), &req); err != nil {
			const size = 8196
			if len(buf) > size {
				buf = buf[:size]
			}
			return nil, httpError{errors.Wrapf(err, "request body '%s': cannot parse non-json request body", buf),
				http.StatusBadRequest,
				nil,
			}
		}
	}

	pathParams := mux.Vars(r)
	_ = pathParams

	queryParams := r.URL.Query()
	_ = queryParams

	return &req, err
}

// DecodeHTTPAddInvitationCodeZeroRequest is a transport/http.DecodeRequestFunc that
// decodes a JSON-encoded addinvitationcode request from the HTTP request
// body. Primarily useful in a server.
func DecodeHTTPAddInvitationCodeZeroRequest(_ context.Context, r *http.Request) (interface{}, error) {
	defer r.Body.Close()
	var req pb.ShareAddInvitationRequest
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot read body of http request")
	}
	if len(buf) > 0 {
		// AllowUnknownFields stops the unmarshaler from failing if the JSON contains unknown fields.
		unmarshaller := jsonpb.Unmarshaler{
			AllowUnknownFields: true,
		}
		if err = unmarshaller.Unmarshal(bytes.NewBuffer(buf), &req); err != nil {
			const size = 8196
			if len(buf) > size {
				buf = buf[:size]
			}
			return nil, httpError{errors.Wrapf(err, "request body '%s': cannot parse non-json request body", buf),
				http.StatusBadRequest,
				nil,
			}
		}
	}

	pathParams := mux.Vars(r)
	_ = pathParams

	queryParams := r.URL.Query()
	_ = queryParams

	return &req, err
}

// DecodeHTTPListFriendZeroRequest is a transport/http.DecodeRequestFunc that
// decodes a JSON-encoded listfriend request from the HTTP request
// body. Primarily useful in a server.
func DecodeHTTPListFriendZeroRequest(_ context.Context, r *http.Request) (interface{}, error) {
	defer r.Body.Close()
	var req pb.ShareListFriendRequest
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot read body of http request")
	}
	if len(buf) > 0 {
		// AllowUnknownFields stops the unmarshaler from failing if the JSON contains unknown fields.
		unmarshaller := jsonpb.Unmarshaler{
			AllowUnknownFields: true,
		}
		if err = unmarshaller.Unmarshal(bytes.NewBuffer(buf), &req); err != nil {
			const size = 8196
			if len(buf) > size {
				buf = buf[:size]
			}
			return nil, httpError{errors.Wrapf(err, "request body '%s': cannot parse non-json request body", buf),
				http.StatusBadRequest,
				nil,
			}
		}
	}

	pathParams := mux.Vars(r)
	_ = pathParams

	queryParams := r.URL.Query()
	_ = queryParams

	if DepthListFriendStrArr, ok := queryParams["depth"]; ok {
		DepthListFriendStr := DepthListFriendStrArr[0]
		DepthListFriend, err := strconv.ParseInt(DepthListFriendStr, 10, 32)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("Error while extracting DepthListFriend from query, queryParams: %v", queryParams))
		}
		req.Depth = int32(DepthListFriend)
	}

	return &req, err
}

// DecodeHTTPClaimRewardZeroRequest is a transport/http.DecodeRequestFunc that
// decodes a JSON-encoded claimreward request from the HTTP request
// body. Primarily useful in a server.
func DecodeHTTPClaimRewardZeroRequest(_ context.Context, r *http.Request) (interface{}, error) {
	defer r.Body.Close()
	var req pb.ShareClaimRewardRequest
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot read body of http request")
	}
	if len(buf) > 0 {
		// AllowUnknownFields stops the unmarshaler from failing if the JSON contains unknown fields.
		unmarshaller := jsonpb.Unmarshaler{
			AllowUnknownFields: true,
		}
		if err = unmarshaller.Unmarshal(bytes.NewBuffer(buf), &req); err != nil {
			const size = 8196
			if len(buf) > size {
				buf = buf[:size]
			}
			return nil, httpError{errors.Wrapf(err, "request body '%s': cannot parse non-json request body", buf),
				http.StatusBadRequest,
				nil,
			}
		}
	}

	pathParams := mux.Vars(r)
	_ = pathParams

	queryParams := r.URL.Query()
	_ = queryParams

	return &req, err
}

// DecodeHTTPPushSignEventZeroRequest is a transport/http.DecodeRequestFunc that
// decodes a JSON-encoded pushsignevent request from the HTTP request
// body. Primarily useful in a server.
func DecodeHTTPPushSignEventZeroRequest(_ context.Context, r *http.Request) (interface{}, error) {
	defer r.Body.Close()
	var req pb.SignEvent
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot read body of http request")
	}
	if len(buf) > 0 {
		// AllowUnknownFields stops the unmarshaler from failing if the JSON contains unknown fields.
		unmarshaller := jsonpb.Unmarshaler{
			AllowUnknownFields: true,
		}
		if err = unmarshaller.Unmarshal(bytes.NewBuffer(buf), &req); err != nil {
			const size = 8196
			if len(buf) > size {
				buf = buf[:size]
			}
			return nil, httpError{errors.Wrapf(err, "request body '%s': cannot parse non-json request body", buf),
				http.StatusBadRequest,
				nil,
			}
		}
	}

	pathParams := mux.Vars(r)
	_ = pathParams

	queryParams := r.URL.Query()
	_ = queryParams

	return &req, err
}

// DecodeHTTPPushTaskEventZeroRequest is a transport/http.DecodeRequestFunc that
// decodes a JSON-encoded pushtaskevent request from the HTTP request
// body. Primarily useful in a server.
func DecodeHTTPPushTaskEventZeroRequest(_ context.Context, r *http.Request) (interface{}, error) {
	defer r.Body.Close()
	var req pb.TaskEvent
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot read body of http request")
	}
	if len(buf) > 0 {
		// AllowUnknownFields stops the unmarshaler from failing if the JSON contains unknown fields.
		unmarshaller := jsonpb.Unmarshaler{
			AllowUnknownFields: true,
		}
		if err = unmarshaller.Unmarshal(bytes.NewBuffer(buf), &req); err != nil {
			const size = 8196
			if len(buf) > size {
				buf = buf[:size]
			}
			return nil, httpError{errors.Wrapf(err, "request body '%s': cannot parse non-json request body", buf),
				http.StatusBadRequest,
				nil,
			}
		}
	}

	pathParams := mux.Vars(r)
	_ = pathParams

	queryParams := r.URL.Query()
	_ = queryParams

	return &req, err
}

// EncodeHTTPGenericResponse is a transport/http.EncodeResponseFunc that encodes
// the response as JSON to the response writer. Primarily useful in a server.
func EncodeHTTPGenericResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	marshaller := jsonpb.Marshaler{
		EmitDefaults: false,
		OrigName:     true,
	}

	return marshaller.Marshal(w, response.(proto.Message))
}

// Helper functions

func headersToContext(ctx context.Context, r *http.Request) context.Context {
	for k := range r.Header {
		// The key is added both in http format (k) which has had
		// http.CanonicalHeaderKey called on it in transport as well as the
		// strings.ToLower which is the grpc metadata format of the key so
		// that it can be accessed in either format
		ctx = context.WithValue(ctx, k, r.Header.Get(k))
		ctx = context.WithValue(ctx, strings.ToLower(k), r.Header.Get(k))
	}

	// Tune specific change.
	// also add the request url
	ctx = context.WithValue(ctx, "request-url", r.URL.Path)
	ctx = context.WithValue(ctx, "transport", "HTTPJSON")

	return ctx
}
