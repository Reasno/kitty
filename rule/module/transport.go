package module

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/pkg/errors"
	"glab.tagtic.cn/ad_gains/kitty/pkg/kerr"
	"glab.tagtic.cn/ad_gains/kitty/rule/dto"
)

var decoder = schema.NewDecoder()

func MakeHTTPHandler(endpoints Endpoints, options ...httptransport.ServerOption) http.Handler {
	serverOptions := []httptransport.ServerOption{
		httptransport.ServerBefore(headersToContext),
		httptransport.ServerErrorEncoder(kerr.ErrorEncoder),
	}
	serverOptions = append(serverOptions, options...)

	m := mux.NewRouter()

	m.Methods("GET", "POST").Path("/v1/calculate/{rule}").Handler(httptransport.NewServer(
		endpoints.calculateRulesEndpoints,
		DecodeCalculateRuleRequest,
		httptransport.EncodeJSONResponse,
		serverOptions...,
	))

	m.Methods("GET").Path("/v1/rule/{rule}").Handler(httptransport.NewServer(
		endpoints.getRulesEndpoint,
		DecodeGetRuleRequest,
		httptransport.EncodeJSONResponse,
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
func decodePayload(payload *dto.Payload, r *http.Request) (err error) {
	if r.Method == "POST" {
		buf, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return errors.Wrapf(err, "cannot read body of http request")
		}
		err = json.Unmarshal(buf, payload)
		if err != nil {
			return errors.Wrap(err, "cannot json unmarshal")
		}
		err = json.Unmarshal(buf, &payload.B)
		if err != nil {
			return errors.Wrap(err, "cannot json unmarshal")
		}
		return nil
	}
	err = decoder.Decode(payload, r.URL.Query())
	// store extra queries here.
	payload.Q = r.URL.Query()
	if err != nil {
		return errors.Wrap(err, "cannot decode payload in query")
	}
	return nil
}

func DecodeCalculateRuleRequest(_ context.Context, r *http.Request) (interface{}, error) {
	defer r.Body.Close()
	var payload dto.Payload
	if err := decodePayload(&payload, r); err != nil {
		return nil, err
	}
	payload.Ip = realIP(r)
	params := mux.Vars(r)
	var req = calculateRulesRequest{
		ruleName: params["rule"],
		payload:  &payload,
	}
	return &req, nil
}

func DecodeGetRuleRequest(_ context.Context, r *http.Request) (interface{}, error) {
	params := mux.Vars(r)
	var req = getRulesRequest{
		ruleName: params["rule"],
	}
	return &req, nil
}

func DecodeUpdateRuleRequest(_ context.Context, r *http.Request) (interface{}, error) {
	defer r.Body.Close()
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot read body of http request")
	}
	params := mux.Vars(r)
	_, dryRun := r.URL.Query()["verify"]
	var req = updateRulesRequest{
		ruleName: params["rule"],
		dryRun:   dryRun,
		data:     buf,
	}
	return &req, nil
}

func DecodePreflightRequest(_ context.Context, r *http.Request) (interface{}, error) {
	defer r.Body.Close()
	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot read body of http request")
	}
	params := mux.Vars(r)
	var req = preflightRequest{
		ruleName: params["rule"],
		hash:     string(buf),
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
