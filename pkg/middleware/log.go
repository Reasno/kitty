package middleware

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"
)

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func NewLoggingMiddleware(logger log.Logger) LabeledMiddleware {
	return func(s string, endpoint endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			defer level.Info(logger).Log("method", s, "req", request, "response", response)
			response, err = endpoint(ctx, request)
			if err != nil {
				level.Warn(logger).Log("err", err)
			}
			return response, err
		}
	}
}
