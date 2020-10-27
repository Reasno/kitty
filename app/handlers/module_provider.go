package handlers

import (
	"github.com/Reasno/kitty/app/entity"
	"github.com/Reasno/kitty/app/svc"
	"github.com/Reasno/kitty/pkg/contract"
	pb "github.com/Reasno/kitty/proto"
	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/tracing/opentracing"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	stdopentracing "github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"gorm.io/gorm"
	"net/http"
)

type AppModule struct {
	logger log.Logger
	db *gorm.DB
	tracer stdopentracing.Tracer
	cleanup   func()
	endpoints svc.Endpoints
}


func New(appModuleConfig contract.ConfigReader) *AppModule {
	appModule, cleanup, err := injectModule(appModuleConfig)
	if err != nil {
		panic(err)
	}
	appModule.cleanup = cleanup
	return appModule
}

func (a *AppModule) ProvideMigration() error {
	err := a.db.AutoMigrate(&entity.User{})
	if err != nil {
		return err
	}
	err = a.db.AutoMigrate(&entity.Device{})
	if err != nil {
		return err
	}
	return nil
}

func (a *AppModule) ProvideCloser() {
	a.cleanup()
}

func (a *AppModule) ProvideGrpc(server *grpc.Server) {
	pb.RegisterAppServer(server, svc.MakeGRPCServer(a.endpoints,
		grpctransport.ServerBefore(opentracing.GRPCToContext(
			a.tracer, "app", a.logger),
		),
		grpctransport.ServerBefore(jwt.GRPCToContext()),
	))
}

func (a *AppModule) ProvideHttp(router *mux.Router) {
	router.PathPrefix("/app/").Handler(http.StripPrefix("/app", svc.MakeHTTPHandler(a.endpoints,
		httptransport.ServerBefore(opentracing.HTTPToContext(
			a.tracer, "app", a.logger)),
		httptransport.ServerBefore(jwt.HTTPToContext()),
	)))
}
