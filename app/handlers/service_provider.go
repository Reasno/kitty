package handlers

import (
	"github.com/Reasno/kitty/app/entity"
	"github.com/Reasno/kitty/app/svc"
	pb "github.com/Reasno/kitty/proto"
	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/tracing/opentracing"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
)

type AppProvider struct {
	cleanup   func()
	endpoints svc.Endpoints
}


func New() *AppProvider {
	appServer, cleanup, err := NewService()
	if err != nil {
		panic(err)
	}
	return &AppProvider{
		cleanup,
		NewEndpoints(appServer),
	}
}

func (a *AppProvider) ProvideMigration() error {
	db, err := injectDb()
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&entity.User{})
	if err != nil {
		return err
	}
	err = db.AutoMigrate(&entity.Device{})
	if err != nil {
		return err
	}
	return nil
}

func (a *AppProvider) ProvideCloser() {
	a.cleanup()
}

func (a *AppProvider) ProvideGrpc(server *grpc.Server) {
	tracer, _, _ := injectOpentracingTracer()
	logger, _ := injectLogger()
	pb.RegisterAppServer(server, svc.MakeGRPCServer(a.endpoints,
		grpctransport.ServerBefore(opentracing.GRPCToContext(
			tracer, "app", logger),
		),
		grpctransport.ServerBefore(jwt.GRPCToContext()),
	))
}

func (a *AppProvider) ProvideHttp(router *mux.Router) {
	tracer, _, _ := injectOpentracingTracer()
	logger, _ := injectLogger()

	router.PathPrefix("/v1/").Handler(svc.MakeHTTPHandler(a.endpoints,
		httptransport.ServerBefore(opentracing.HTTPToContext(
			tracer, "app", logger)),
		httptransport.ServerBefore(jwt.HTTPToContext()),
	))
}
