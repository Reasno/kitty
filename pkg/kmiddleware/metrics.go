package kmiddleware

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/metrics"
	"time"
)

func NewLabeledMetricsMiddleware(his metrics.Histogram, module string) LabeledMiddleware {
	return func(name string, e endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			defer func(begin time.Time) {
				his.With("module", module, "method", name).Observe(time.Since(begin).Seconds())
			}(time.Now())
			return e(ctx, request)
		}
	}
}

func NewMetricsMiddleware(his metrics.Histogram, module, method string) endpoint.Middleware {
	return func(e endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			defer func(begin time.Time) {
				his.With("module", module, "method", method).Observe(time.Since(begin).Seconds())
			}(time.Now())
			return e(ctx, request)
		}
	}
}
