package register

import (
	"github.com/Reasno/kitty/app/handlers"
	"github.com/Reasno/kitty/app/svc"
	"github.com/Reasno/kitty/app/svc/server"
	pb "github.com/Reasno/kitty/proto"
	"github.com/go-kit/kit/tracing/opentracing"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	httptransport "github.com/go-kit/kit/transport/http"
	"google.golang.org/grpc"
	"net/http"
)

func RegisterApp(httpProviders *[]func() http.Handler, grpcProvider *[]func(*grpc.Server))  {
	appServer := handlers.NewService()
	endpoints := server.NewEndpoints(appServer)
	*httpProviders = append(*httpProviders, func() http.Handler {
		return svc.MakeHTTPHandler(endpoints,
			httptransport.ServerBefore(opentracing.HTTPToContext(
				handlers.InjectOpentracingTracer(), "app", handlers.ProvideLogger())),
		)
	})
	*grpcProvider = append(*grpcProvider, func(s *grpc.Server) {
		pb.RegisterAppServer(s, svc.MakeGRPCServer(endpoints,
			grpctransport.ServerBefore(opentracing.GRPCToContext(
				handlers.InjectOpentracingTracer(), "app", handlers.ProvideLogger()))))
	})
}
