package module

import (
	"net/http"

	"glab.tagtic.cn/ad_gains/kitty/pkg/kerr"
	"glab.tagtic.cn/ad_gains/kitty/pkg/kgrpc"
	"glab.tagtic.cn/ad_gains/kitty/pkg/khttp"

	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/tracing/opentracing"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	stdopentracing "github.com/opentracing/opentracing-go"
	"glab.tagtic.cn/ad_gains/kitty/app/repository"
	"glab.tagtic.cn/ad_gains/kitty/app/svc"
	"glab.tagtic.cn/ad_gains/kitty/pkg/config"
	"glab.tagtic.cn/ad_gains/kitty/pkg/contract"
	"glab.tagtic.cn/ad_gains/kitty/pkg/klog"
	pb "glab.tagtic.cn/ad_gains/kitty/proto"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type Module struct {
	appName   contract.AppName
	logger    log.Logger
	db        *gorm.DB
	tracer    stdopentracing.Tracer
	cleanup   func()
	endpoints svc.Endpoints
}

func New(appModuleConfig contract.ConfigReader, logger log.Logger, dynConf config.DynamicConfigReader) *Module {
	appModule, cleanup, err := injectModule(setUp(appModuleConfig, logger, dynConf))
	if err != nil {
		panic(err)
	}
	appModule.cleanup = cleanup
	return appModule
}

func setUp(appModuleConfig contract.ConfigReader, logger log.Logger, dynConf config.DynamicConfigReader) (contract.ConfigReader, log.Logger, config.DynamicConfigReader) {
	appLogger := log.With(logger, "module", config.ProvideAppName(appModuleConfig).String())
	appLogger = level.NewFilter(logger, klog.LevelFilter(appModuleConfig.String("level")))
	return appModuleConfig, appLogger, dynConf
}

func (a *Module) ProvideMigration() error {
	m := repository.ProvideMigrator(a.db, a.appName)
	return m.Migrate()
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
	router.PathPrefix("/app/").Handler(http.StripPrefix("/app", svc.MakeHTTPHandler(a.endpoints,
		httptransport.ServerBefore(
			opentracing.HTTPToContext(a.tracer, "app", a.logger),
			jwt.HTTPToContext(),
			khttp.IpToContext(),
		),
		httptransport.ServerErrorEncoder(kerr.ErrorEncoder),
	)))
}
