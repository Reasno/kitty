package module

import (
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log/level"
	"github.com/go-redis/redis/v8"
	"github.com/opentracing/opentracing-go/ext"
	"glab.tagtic.cn/ad_gains/kitty/pkg/event"
	kclient "glab.tagtic.cn/ad_gains/kitty/pkg/kkafka/client"
	"glab.tagtic.cn/ad_gains/kitty/pkg/kmiddleware"
	"glab.tagtic.cn/ad_gains/kitty/share/listener"
	"net/http"
	"time"

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

func provideDispatcher(icBus listener.InvitationCodeBus) *event.Dispatcher {
	dispatcher := event.Dispatcher{}
	dispatcher.Subscribe(listener.InvitationCodeAdded{
		Bus: icBus,
	})
	return &dispatcher
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

func providePublisherOptions(tracer stdopentracing.Tracer, logger log.Logger) []kkafka.PublisherOption {
	return []kkafka.PublisherOption{
		kkafka.PublisherBefore(kkafka.ContextToKafka(tracer, logger)),
	}
}

type producerMiddleware func(operationName string) endpoint.Middleware

func provideProducerMiddleware(tracer stdopentracing.Tracer, logger log.Logger) producerMiddleware {
	return func(operationName string) endpoint.Middleware {
		return endpoint.Chain(
			kmiddleware.NewAsyncMiddleware(logger),
			kmiddleware.TraceProducer(tracer, operationName, ext.SpanKindProducerEnum),
			kmiddleware.NewTimeoutMiddleware(time.Second),
		)
	}
}

func provideInvitationCodeBus(factory *kkafka.KafkaFactory, conf contract.ConfigReader, option []kkafka.PublisherOption, mw producerMiddleware) *kclient.DataStore {
	return kclient.NewDataStore(conf.String("kafka.shareInvitationCodeBus"), factory, option, mw("kafka.Share"))
}

func ProvideRedis(logging log.Logger, conf contract.ConfigReader) (redis.UniversalClient, func()) {
	client := redis.NewUniversalClient(
		&redis.UniversalOptions{
			Addrs:    conf.Strings("redis.addrs"),
			DB:       conf.Int("redis.database"),
			Password: conf.String("redis.password"),
		})
	return client, func() {
		if err := client.Close(); err != nil {
			level.Error(logging).Log("err", err.Error())
		}
	}
}
