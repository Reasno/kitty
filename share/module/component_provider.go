package module

import (
	"fmt"
	"net/http"

	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/tracing/opentracing"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	httptransport "github.com/go-kit/kit/transport/http"
	stdopentracing "github.com/opentracing/opentracing-go"
	"glab.tagtic.cn/ad_gains/kitty/pkg/contract"
	code "glab.tagtic.cn/ad_gains/kitty/pkg/invitecode"
	"glab.tagtic.cn/ad_gains/kitty/pkg/kerr"
	"glab.tagtic.cn/ad_gains/kitty/pkg/kgrpc"
	"glab.tagtic.cn/ad_gains/kitty/pkg/khttp"
	"glab.tagtic.cn/ad_gains/kitty/pkg/kkafka"
	kitty "glab.tagtic.cn/ad_gains/kitty/proto"
	"glab.tagtic.cn/ad_gains/kitty/share/svc"
)

func provideTokenizer(conf contract.ConfigReader) *code.Tokenizer {
	return code.NewTokenizer(conf.String("salt"))
}

func provideEndpoints(middleware overallMiddleware, server kitty.ShareServer) svc.Endpoints {
	return middleware(svc.NewEndpoints(server))
}

type overallMiddleware func(endpoints svc.Endpoints) svc.Endpoints

func provideModule(server GrpcShareServer, handler http.Handler, kafkaServer kkafka.Server, appName contract.AppName) *Module {
	return &Module{
		appName:     appName,
		grpcServer:  server,
		handler:     handler,
		kafkaServer: kafkaServer,
	}
}

func provideHttp(endpoints svc.Endpoints, tracer stdopentracing.Tracer, logger log.Logger, appName contract.AppName) http.Handler {
	return svc.MakeHTTPHandler(endpoints,
		httptransport.ServerBefore(
			opentracing.HTTPToContext(tracer, appName.String(), logger),
			jwt.HTTPToContext(),
			khttp.IpToContext(),
		),
		httptransport.ServerErrorEncoder(kerr.ErrorEncoder),
	)
}

type GrpcShareServer kitty.ShareServer

func provideGrpc(endpoints svc.Endpoints, tracer stdopentracing.Tracer, logger log.Logger, appName contract.AppName) GrpcShareServer {
	return svc.MakeGRPCServer(endpoints,
		grpctransport.ServerBefore(
			opentracing.GRPCToContext(tracer, appName.String(), logger),
			jwt.GRPCToContext(),
			kgrpc.IpToContext(),
		),
		grpctransport.ServerBefore(jwt.GRPCToContext()),
	)
}

func provideKafkaServer(endpoints svc.Endpoints, factory *kkafka.KafkaFactory, conf contract.ConfigReader, tracer stdopentracing.Tracer, env contract.Env, logger log.Logger) kkafka.Server {
	serverOptions := []kkafka.SubscriberOption{
		kkafka.SubscriberBefore(kkafka.KafkaToContext(tracer, fmt.Sprintf("kafka(%s)", env.String()), logger)),
		kkafka.SubscriberBefore(kkafka.Trust()),
		kkafka.SubscriberErrorHandler(kkafka.ErrHandler(logger)),
	}
	return svc.MakeKafkaServer(endpoints, factory, conf, serverOptions...)
}
