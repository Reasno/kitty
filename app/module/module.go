package module

import (
	"context"
	"net/http"

	"github.com/oklog/run"
	"glab.tagtic.cn/ad_gains/kitty/pkg/kerr"
	"glab.tagtic.cn/ad_gains/kitty/pkg/kgrpc"
	"glab.tagtic.cn/ad_gains/kitty/pkg/khttp"
	"glab.tagtic.cn/ad_gains/kitty/pkg/kkafka"

	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/tracing/opentracing"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	stdopentracing "github.com/opentracing/opentracing-go"
	"glab.tagtic.cn/ad_gains/kitty/app/repository"
	"glab.tagtic.cn/ad_gains/kitty/app/svc"
	"glab.tagtic.cn/ad_gains/kitty/pkg/config"
	"glab.tagtic.cn/ad_gains/kitty/pkg/contract"
	pb "glab.tagtic.cn/ad_gains/kitty/proto"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type Module struct {
	appName     contract.AppName
	conf        contract.ConfigReader
	factory     *kkafka.KafkaFactory
	logger      log.Logger
	db          *gorm.DB
	tracer      stdopentracing.Tracer
	cleanup     func()
	endpoints   svc.Endpoints
	kafkaServer kkafka.Server
}

func (a *Module) ProvideRunGroup(group *run.Group) {
	serverOptions := []kkafka.SubscriberOption{
		kkafka.SubscriberBefore(kkafka.KafkaToContext(a.tracer, "app", a.logger)),
		kkafka.SubscriberBefore(kkafka.Trust()),
		kkafka.SubscriberErrorHandler(kkafka.ErrHandler(a.logger)),
	}
	kafkaServer := svc.MakeKafkaServer(a.endpoints, a.factory, a.conf, serverOptions...)
	ctx, cancel := context.WithCancel(context.Background())
	group.Add(func() error {
		return kafkaServer.Serve(ctx)
	}, func(err error) {
		cancel()
	})
}

func New(appModuleConfig contract.ConfigReader, logger log.Logger, dynConf config.DynamicConfigReader) *Module {
	appModule, cleanup, err := injectModule(appModuleConfig, logger, dynConf)
	if err != nil {
		panic(err)
	}
	appModule.cleanup = cleanup
	return appModule
}

func (a *Module) ProvideMigration() error {
	m := repository.ProvideMigrator(a.db, a.appName)
	return m.Migrate()
}

func (a *Module) ProvideSeed() error {
	s := repository.ProvideSeeder(a.db)
	return s.Seed()
}

func (a *Module) ProvideRollback(id string) error {
	m := repository.ProvideMigrator(a.db, a.appName)
	if id == "-1" {
		return m.RollbackLast()
	}
	return m.RollbackTo(id)
}

func (a *Module) ProvideCloser() {
	a.cleanup()
}

func (a *Module) ProvideGrpc(server *grpc.Server) {
	pb.RegisterAppServer(server, svc.MakeGRPCServer(a.endpoints,
		grpctransport.ServerBefore(
			opentracing.GRPCToContext(a.tracer, "app", a.logger),
			jwt.GRPCToContext(),
			kgrpc.IpToContext(),
		),
		grpctransport.ServerBefore(jwt.GRPCToContext()),
	))
}

func (a *Module) ProvideHttp(router *mux.Router) {
	router.PathPrefix("/app/v1/").Handler(http.StripPrefix("/app/v1", svc.MakeHTTPHandlerV1(a.endpoints,
		httptransport.ServerBefore(
			opentracing.HTTPToContext(a.tracer, "app", a.logger),
			jwt.HTTPToContext(),
			khttp.IpToContext(),
		),
		httptransport.ServerErrorEncoder(kerr.ErrorEncoder),
	)))
	router.PathPrefix("/app/v2/").Handler(http.StripPrefix("/app/v2", svc.MakeHTTPHandler(a.endpoints,
		httptransport.ServerBefore(
			opentracing.HTTPToContext(a.tracer, "app", a.logger),
			jwt.HTTPToContext(),
			khttp.IpToContext(),
		),
		httptransport.ServerErrorEncoder(kerr.ErrorEncoder),
	)))
}
