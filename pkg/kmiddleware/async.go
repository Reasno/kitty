package kmiddleware

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/opentracing/opentracing-go"
)

func NewAsyncMiddleware(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			span := opentracing.SpanFromContext(ctx)
			go func() {
				ctx := opentracing.ContextWithSpan(context.Background(), span)
				_, err = next(ctx, request)
				level.Warn(logger).Log("err", err.Error())
			}()
			return nil, nil
		}
	}
}
