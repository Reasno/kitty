package module

import (
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	"glab.tagtic.cn/ad_gains/kitty/pkg/config"
	"glab.tagtic.cn/ad_gains/kitty/pkg/contract"
	pb "glab.tagtic.cn/ad_gains/kitty/proto"
	"google.golang.org/grpc"
)

type Module struct {
	appName    contract.AppName
	cleanup    func()
	handler    http.Handler
	grpcServer GrpcShareServer
}

func New(appModuleConfig contract.ConfigReader, logger log.Logger, dynConf config.DynamicConfigReader) *Module {
	appModule, cleanup, err := injectModule(appModuleConfig, logger, dynConf)
	if err != nil {
		panic(err)
	}
	appModule.cleanup = cleanup
	return appModule
}

func (a *Module) ProvideCloser() {
	a.cleanup()
}

func (a *Module) ProvideGrpc(server *grpc.Server) {
	pb.RegisterShareServer(server, a.grpcServer)
}

func (a *Module) ProvideHttp(router *mux.Router) {
	router.PathPrefix("/share/v1/").Handler(http.StripPrefix("/share/v1", a.handler))
}
