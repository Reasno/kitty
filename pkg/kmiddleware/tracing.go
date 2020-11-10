package kmiddleware

import (
	"context"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/tracing/opentracing"
	stdtracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"glab.tagtic.cn/ad_gains/kitty/pkg/kjwt"
	pb "glab.tagtic.cn/ad_gains/kitty/proto"
)

type LabeledMiddleware func(string, endpoint.Endpoint) endpoint.Endpoint

func NewTraceMiddleware(tracer stdtracing.Tracer, env string) LabeledMiddleware {
	return func(s string, endpoint endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			name := fmt.Sprintf("%s(%s)", s, env)
			e := opentracing.TraceServer(tracer, name)(endpoint)
			span := stdtracing.SpanFromContext(ctx)
			claim := kjwt.GetClaim(ctx)
			if r, ok := request.(pb.UserLoginRequest); ok {
				claim.Suuid = r.Device.Suuid

			}
			span.SetTag("env", env)
			span.SetTag("package.name", claim.PackageName)
			span.SetTag("suuid", claim.Suuid)
			span.SetTag("user.id", claim.UserId)
			resp, err := e(ctx, request)
			if err != nil {
				ext.Error.Set(span, true)
				span.LogKV("error", err.Error())
			}
			return resp, err
		}
	}
}
