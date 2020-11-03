package rule

import (
	"context"
	"encoding/json"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"strings"
)

func MakeHTTPHandler(endpoints Endpoints, options ...httptransport.ServerOption) http.Handler {
	serverOptions := []httptransport.ServerOption{
		httptransport.ServerBefore(headersToContext),
	}
	serverOptions = append(serverOptions, options...)

	m := mux.NewRouter()

	m.Methods("POST").Path("/v1/calculate/{rule}").Handler(httptransport.NewServer(
		endpoints.calculateRulesEndpoints,
		DecodeCalculateRuleRequest,
		httptransport.EncodeJSONResponse,
		serverOptions...,
	))

	m.Methods("GET").Path("/v1/rule/{rule}").Handler(httptransport.NewServer(
		endpoints.getRulesEndpoint,
		DecodeGetRuleRequest,
		EncodeYamlResponse,
		serverOptions...,
	))

	m.Methods("POST").Path("/v1/rule/{rule}").Handler(httptransport.NewServer(
		endpoints.updateRulesEndpoint,
		DecodeUpdateRuleRequest,
		httptransport.EncodeJSONResponse,
		serverOptions...,
	))

	m.Methods("POST").Path("/v1/preflight/{rule}").Handler(httptransport.NewServer(
		endpoints.preflightEndpoint,
		DecodePreflightRequest,
		httptransport.EncodeJSONResponse,
		serverOptions...,
	))
	
	return m
}

func EncodeYamlResponse(_ context.Context, writer http.ResponseWriter, i interface{}) error {
	writer.Header().Set("Content-Type", "application/x-yaml; charset=utf-8")
	if headerer, ok := i.(httptransport.Headerer); ok {
		for k, values := range headerer.Headers() {
			for _, v := range values {
				writer.Header().Add(k, v)
			}
		}
	}
	code := http.StatusOK
	if sc, ok := i.(httptransport.StatusCoder); ok {
		code = sc.StatusCode()
	}
	writer.WriteHeader(code)
	if code == http.StatusNoContent {
		return nil
	}
	writer.Write(i.(ByteResponse))
	return nil
}

func DecodeCalculateRuleRequest (_ context.Context, r *http.Request) (interface{}, error)  {
	defer r.Body.Close()
	var payload Payload
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot read body of http request")
	}
	err = json.Unmarshal(buf, &payload)
	if err != nil {
		return nil, errors.Wrap(err, "cannot json unmarshal")
	}
	payload.Ip = realIP(r)
	params := mux.Vars(r)
	var req = calculateRulesRequest{
		ruleName: params["rule"],
		payload: &payload,
	}
	return &req, nil
}

func DecodeGetRuleRequest (_ context.Context, r *http.Request) (interface{}, error)  {
	params := mux.Vars(r)
	var req = getRulesRequest{
		ruleName: params["rule"],
	}
	return &req, nil
}


func DecodeUpdateRuleRequest (_ context.Context, r *http.Request) (interface{}, error)  {
	defer r.Body.Close()
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot read body of http request")
	}
	params := mux.Vars(r)
	_, dryRun := r.URL.Query()["verify"]
	var req = updateRulesRequest{
		ruleName: params["rule"],
		dryRun: dryRun,
		data: buf,
	}
	return &req, nil
}

func DecodePreflightRequest (_ context.Context, r *http.Request) (interface{}, error)  {
	defer r.Body.Close()
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot read body of http request")
	}
	params := mux.Vars(r)
	var req = preflightRequest{
		ruleName: params["rule"],
		hash: string(buf),
	}
	return &req, nil
}

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
