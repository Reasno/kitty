package handlers

import (
	"github.com/Reasno/kitty/pkg/kerr"
	"net/http"

	"github.com/Reasno/kitty/app/repository"
	"github.com/Reasno/kitty/app/svc"
	"github.com/Reasno/kitty/pkg/config"
	"github.com/Reasno/kitty/pkg/contract"
	"github.com/Reasno/kitty/pkg/klog"
	pb "github.com/Reasno/kitty/proto"
	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-kit/kit/tracing/opentracing"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	stdopentracing "github.com/opentracing/opentracing-go"
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

func New(appModuleConfig contract.ConfigReader, logger log.Logger) *Module {
	appModule, cleanup, err := injectModule(setUp(appModuleConfig, logger))
	if err != nil {
		panic(err)
	}
	appModule.cleanup = cleanup
	return appModule
}

func setUp(appModuleConfig contract.ConfigReader, logger log.Logger) (contract.ConfigReader, log.Logger) {
	appLogger := log.With(logger, "module", config.ProvideAppName(appModuleConfig).String())
	appLogger = level.NewFilter(logger, klog.LevelFilter(appModuleConfig.String("level")))
	return appModuleConfig, appLogger
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
		grpctransport.ServerBefore(opentracing.GRPCToContext(
			a.tracer, "app", a.logger),
		),
		grpctransport.ServerBefore(jwt.GRPCToContext()),
	))
}

func (a *Module) ProvideHttp(router *mux.Router) {
	router.PathPrefix("/app/").Handler(http.StripPrefix("/app", svc.MakeHTTPHandler(a.endpoints,
		httptransport.ServerBefore(opentracing.HTTPToContext(
			a.tracer, "app", a.logger)),
		httptransport.ServerBefore(jwt.HTTPToContext()),
		httptransport.ServerErrorEncoder(kerr.ErrorEncoder),
	)))
}
