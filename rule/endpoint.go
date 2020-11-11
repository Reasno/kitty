package rule

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"glab.tagtic.cn/ad_gains/kitty/pkg/contract"
	"glab.tagtic.cn/ad_gains/kitty/pkg/kmiddleware"
)

type GenericResponse struct {
	Code    uint32 `json:"code"`
	Message string `json:"message,omitempty"`
	Data    Data   `json:"data,omitempty"`
}

func (g GenericResponse) String() string {
	ss, _ := json.Marshal(g)
	return string(ss)
}

type StringResponse struct {
	Code    uint32 `json:"code"`
	Message string `json:"message,omitempty"`
	Data    string `json:"data,omitempty"`
}

func (s StringResponse) String() string {
	ss, _ := json.Marshal(s)
	return string(ss)
}

type calculateRulesRequest struct {
	ruleName string
	payload  *Payload
}

func (c calculateRulesRequest) String() string {
	return fmt.Sprintf("%s: %s", c.ruleName, c.payload)
}

type getRulesRequest struct {
	ruleName string
}

func (g getRulesRequest) String() string {
	return fmt.Sprintf("%s", g.ruleName)
}

type updateRulesRequest struct {
	ruleName string
	data     []byte
	dryRun   bool
}

func (u updateRulesRequest) String() string {
	return fmt.Sprintf("%s: %s(%t)", u.ruleName, u.data, u.dryRun)
}

type preflightRequest struct {
	ruleName string
	hash     string
}

func (p preflightRequest) String() string {
	return fmt.Sprintf("%s: %s", p.ruleName, p.hash)
}

type Endpoints struct {
	calculateRulesEndpoints endpoint.Endpoint
	getRulesEndpoint        endpoint.Endpoint
	updateRulesEndpoint     endpoint.Endpoint
	preflightEndpoint       endpoint.Endpoint
}

func newEndpoints(s Service, hist metrics.Histogram, logger log.Logger, appName contract.AppName, env contract.Env) Endpoints {
	l := kmiddleware.NewLoggingMiddleware(logger, env.IsLocal())
	e := kmiddleware.NewErrorMarshallerMiddleware()
	mw := func(name string) endpoint.Middleware {
		return endpoint.Chain(e, kmiddleware.NewMetricsMiddleware(hist, appName.String(), name), l)
	}
	return Endpoints{
		calculateRulesEndpoints: mw("CalculateRules")(MakeCalculateRulesEndpoint(s)),
		getRulesEndpoint:        mw("GetRules")(MakeGetRulesEndpoint(s)),
		updateRulesEndpoint:     mw("UpdateRules")(MakeUpdateRulesEndpoint(s)),
		preflightEndpoint:       mw("Preflight")(MakePreflightEndpoint(s)),
	}
}

func MakeCalculateRulesEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*calculateRulesRequest)
		v, err := s.CalculateRules(ctx, req.ruleName, req.payload)
		if err != nil {
			return nil, err
		}
		return GenericResponse{Data: v}, nil
	}
}

func MakeGetRulesEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*getRulesRequest)
		v, err := s.GetRules(ctx, req.ruleName)
		if err != nil {
			return nil, err
		}
		return StringResponse{Data: string(v)}, nil
	}
}

func MakeUpdateRulesEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*updateRulesRequest)
		err = s.UpdateRules(ctx, req.ruleName, req.data, req.dryRun)
		var invalid *ErrInvalidRules
		if errors.As(err, &invalid) {
			return GenericResponse{Message: invalid.Error(), Code: 3}, nil
		}
		if err != nil {
			return nil, err
		}
		return GenericResponse{}, nil
	}
}

func MakePreflightEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*preflightRequest)
		err = s.Preflight(ctx, req.ruleName, req.hash)
		if err != nil {
			return nil, err
		}
		return GenericResponse{}, nil
	}
}
