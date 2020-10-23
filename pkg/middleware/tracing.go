package middleware

import (
	"github.com/go-kit/kit/endpoint"

	"github.com/go-kit/kit/tracing/opentracing"
	stdtracing "github.com/opentracing/opentracing-go"
)

type LabeledMiddleware func(string, endpoint.Endpoint) endpoint.Endpoint

func NewTraceMiddleware(tracer stdtracing.Tracer, service string) LabeledMiddleware {
	return func(s string, endpoint endpoint.Endpoint) endpoint.Endpoint {
		return opentracing.TraceServer(tracer, service+"."+s)(endpoint)
	}
}
