package module

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/metrics"
	"github.com/opentracing/opentracing-go"
	"glab.tagtic.cn/ad_gains/kitty/pkg/contract"
	"glab.tagtic.cn/ad_gains/kitty/pkg/kmiddleware"
	"glab.tagtic.cn/ad_gains/kitty/share/svc"
)

func provideEndpointsMiddleware(l log.Logger, securityConfig *kmiddleware.SecurityConfig, hist metrics.Histogram, tracer opentracing.Tracer, env contract.Env, appName contract.AppName) overallMiddleware {
	return func(in svc.Endpoints) svc.Endpoints {
		in.WrapAllExcept(kmiddleware.NewValidationMiddleware())
		in.WrapAllExcept(kmiddleware.NewLoggingMiddleware(l, env.IsLocal()))
		in.WrapAllLabeledExcept(kmiddleware.NewLabeledMetricsMiddleware(hist, appName.String()))
		in.WrapAllLabeledExcept(kmiddleware.NewTraceMiddleware(tracer, env.String()))
		in.WrapAllExcept(kmiddleware.NewConfigMiddleware())
		in.WrapAllExcept(kmiddleware.NewAuthenticationMiddleware(securityConfig))
		in.WrapAllExcept(kmiddleware.NewErrorMarshallerMiddleware())
		return in
	}
}
