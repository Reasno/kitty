package register

import (
	"github.com/Reasno/kitty/app/handlers"
	"github.com/Reasno/kitty/app/svc"
	"github.com/Reasno/kitty/app/svc/server"
	pb "github.com/Reasno/kitty/proto"
	"google.golang.org/grpc"
	"net/http"
)

func RegisterApp(httpProviders *[]func() http.Handler, grpcProvider *[]func(*grpc.Server))  {
	appServer := handlers.NewService()
	endpoints := server.NewEndpoints(appServer)
	*httpProviders = append(*httpProviders, func() http.Handler {
		return svc.MakeHTTPHandler(endpoints)
	})
	*grpcProvider = append(*grpcProvider, func(s *grpc.Server) {
		pb.RegisterAppServer(s, svc.MakeGRPCServer(endpoints))
	})
}
