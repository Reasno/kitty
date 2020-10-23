package register

import (
	"github.com/Reasno/kitty/app/entity"
	"github.com/Reasno/kitty/app/handlers"
	"github.com/Reasno/kitty/app/svc"
	"github.com/Reasno/kitty/app/svc/server"
	pb "github.com/Reasno/kitty/proto"
	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/tracing/opentracing"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
)

func RegisterApp(httpProviders *[]func(router *mux.Router), grpcProvider *[]func(*grpc.Server)) {
	appServer := handlers.NewService()
	endpoints := server.NewEndpoints(appServer)
	*httpProviders = append(*httpProviders, func(r *mux.Router) {
		r.Handle("/", svc.MakeHTTPHandler(endpoints,
			httptransport.ServerBefore(opentracing.HTTPToContext(
				handlers.InjectOpentracingTracer(), "app", handlers.ProvideLogger())),
			httptransport.ServerBefore(jwt.HTTPToContext()),
		))
	})
	*grpcProvider = append(*grpcProvider, func(s *grpc.Server) {
		pb.RegisterAppServer(s, svc.MakeGRPCServer(endpoints,
			grpctransport.ServerBefore(opentracing.GRPCToContext(
				handlers.InjectOpentracingTracer(), "app", handlers.ProvideLogger()),
			),
			grpctransport.ServerBefore(jwt.GRPCToContext()),
		))
	})
}

func RegisterAppMigrations(migrationProvider *[]func() error) {
	*migrationProvider = append(*migrationProvider, func() error {
		db, err := handlers.InjectDb()
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
	})
}
