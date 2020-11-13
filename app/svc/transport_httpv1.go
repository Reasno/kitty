package svc

// This file provides server-side bindings for the HTTP transport.
// It utilizes the transport/http.Server.

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gogo/protobuf/proto"
	"github.com/golang/protobuf/jsonpb"

	"context"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"

	// This service
	pb "glab.tagtic.cn/ad_gains/kitty/proto"
)

var (
	_ = fmt.Sprint
	_ = bytes.Compare
	_ = strconv.Atoi
	_ = httptransport.NewServer
	_ = ioutil.NopCloser
	_ = pb.NewAppClient
	_ = io.Copy
	_ = errors.Wrap
)

// MakeHTTPHandler returns a handler that makes a set of endpoints available
// on predefined paths.
func MakeHTTPHandlerV1(endpoints Endpoints, options ...httptransport.ServerOption) http.Handler {
	serverOptions := []httptransport.ServerOption{
		httptransport.ServerBefore(headersToContext),
		httptransport.ServerErrorEncoder(errorEncoder),
		httptransport.ServerAfter(httptransport.SetContentType(contentType)),
	}
	serverOptions = append(serverOptions, options...)
	m := mux.NewRouter()

	m.Methods("POST").Path("/login").Handler(httptransport.NewServer(
		endpoints.LoginEndpoint,
		DecodeHTTPLoginZeroRequestV1,
		EncodeHTTPGenericResponseV1,
		serverOptions...,
	))

	m.Methods("GET").Path("/code").Handler(httptransport.NewServer(
		endpoints.GetCodeEndpoint,
		DecodeHTTPGetCodeZeroRequestV1,
		EncodeHTTPGenericResponseV1,
		serverOptions...,
	))

	m.Methods("GET").Path("/info/{id}").Handler(httptransport.NewServer(
		endpoints.GetInfoEndpoint,
		DecodeHTTPGetInfoZeroRequestV1,
		EncodeHTTPGenericResponseV1,
		serverOptions...,
	))

	m.Methods("POST").Path("/info").Handler(httptransport.NewServer(
		endpoints.UpdateInfoEndpoint,
		DecodeHTTPUpdateInfoZeroRequestV1,
		EncodeHTTPGenericResponseV1,
		serverOptions...,
	))

	m.Methods("POST").Path("/bind").Handler(httptransport.NewServer(
		endpoints.BindEndpoint,
		DecodeHTTPBindZeroRequestV1,
		EncodeHTTPGenericResponseV1,
		serverOptions...,
	))

	m.Methods("POST").Path("/unbind").Handler(httptransport.NewServer(
		endpoints.UnbindEndpoint,
		DecodeHTTPUnbindZeroRequestV1,
		EncodeHTTPGenericResponseV1,
		serverOptions...,
	))

	m.Methods("POST").Path("/refresh").Handler(httptransport.NewServer(
		endpoints.RefreshEndpoint,
		DecodeHTTPRefreshZeroRequestV1,
		EncodeHTTPGenericResponseV1,
		serverOptions...,
	))
	return m
}

// Server Decode

// DecodeHTTPLoginZeroRequest is a transport/http.DecodeRequestFunc that
// decodes a JSON-encoded login request from the HTTP request
// body. Primarily useful in a server.
func DecodeHTTPLoginZeroRequestV1(_ context.Context, r *http.Request) (interface{}, error) {
	defer r.Body.Close()
	var req pb.UserLoginRequest
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

// DecodeHTTPGetCodeZeroRequest is a transport/http.DecodeRequestFunc that
// decodes a JSON-encoded getcode request from the HTTP request
// body. Primarily useful in a server.
func DecodeHTTPGetCodeZeroRequestV1(_ context.Context, r *http.Request) (interface{}, error) {
	defer r.Body.Close()
	var req pb.GetCodeRequest
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

	if MobileGetCodeStrArr, ok := queryParams["mobile"]; ok {
		MobileGetCodeStr := MobileGetCodeStrArr[0]
		MobileGetCode := MobileGetCodeStr
		req.Mobile = MobileGetCode
	}

	if PackageNameGetCodeStrArr, ok := queryParams["packageName"]; ok {
		PackageNameGetCodeStr := PackageNameGetCodeStrArr[0]
		PackageNameGetCode := PackageNameGetCodeStr
		req.PackageName = PackageNameGetCode
	}

	return &req, err
}

// DecodeHTTPGetInfoZeroRequest is a transport/http.DecodeRequestFunc that
// decodes a JSON-encoded getinfo request from the HTTP request
// body. Primarily useful in a server.
func DecodeHTTPGetInfoZeroRequestV1(_ context.Context, r *http.Request) (interface{}, error) {
	defer r.Body.Close()
	var req pb.UserInfoRequest
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

	IdGetInfoStr := pathParams["id"]
	IdGetInfo, err := strconv.ParseUint(IdGetInfoStr, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("Error while extracting IdGetInfo from path, pathParams: %v", pathParams))
	}
	req.Id = IdGetInfo

	if WechatGetInfoStrArr, ok := queryParams["wechat"]; ok {
		WechatGetInfoStr := WechatGetInfoStrArr[0]
		WechatGetInfo, err := strconv.ParseBool(WechatGetInfoStr)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("Error while extracting WechatGetInfo from query, queryParams: %v", queryParams))
		}
		req.Wechat = WechatGetInfo
	}

	if TaobaoGetInfoStrArr, ok := queryParams["taobao"]; ok {
		TaobaoGetInfoStr := TaobaoGetInfoStrArr[0]
		TaobaoGetInfo, err := strconv.ParseBool(TaobaoGetInfoStr)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("Error while extracting TaobaoGetInfo from query, queryParams: %v", queryParams))
		}
		req.Taobao = TaobaoGetInfo
	}

	return &req, err
}

// DecodeHTTPUpdateInfoZeroRequest is a transport/http.DecodeRequestFunc that
// decodes a JSON-encoded updateinfo request from the HTTP request
// body. Primarily useful in a server.
func DecodeHTTPUpdateInfoZeroRequestV1(_ context.Context, r *http.Request) (interface{}, error) {
	defer r.Body.Close()
	var req pb.UserInfoUpdateRequest
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

// DecodeHTTPBindZeroRequest is a transport/http.DecodeRequestFunc that
// decodes a JSON-encoded bind request from the HTTP request
// body. Primarily useful in a server.
func DecodeHTTPBindZeroRequestV1(_ context.Context, r *http.Request) (interface{}, error) {
	defer r.Body.Close()
	var req pb.UserBindRequest
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

// DecodeHTTPUnbindZeroRequest is a transport/http.DecodeRequestFunc that
// decodes a JSON-encoded unbind request from the HTTP request
// body. Primarily useful in a server.
func DecodeHTTPUnbindZeroRequestV1(_ context.Context, r *http.Request) (interface{}, error) {
	defer r.Body.Close()
	var req pb.UserUnbindRequest
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

// DecodeHTTPRefreshZeroRequest is a transport/http.DecodeRequestFunc that
// decodes a JSON-encoded refresh request from the HTTP request
// body. Primarily useful in a server.
func DecodeHTTPRefreshZeroRequestV1(_ context.Context, r *http.Request) (interface{}, error) {
	defer r.Body.Close()
	var req pb.UserRefreshRequest
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
func EncodeHTTPGenericResponseV1(_ context.Context, w http.ResponseWriter, response interface{}) error {
	marshaller := jsonpb.Marshaler{
		EmitDefaults: false,
		OrigName:     true,
	}

	return marshaller.Marshal(w, response.(proto.Message))
}
