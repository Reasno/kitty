package module

import (
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/gorilla/mux"
	"glab.tagtic.cn/ad_gains/kitty/pkg/config"
	"glab.tagtic.cn/ad_gains/kitty/pkg/contract"
	"glab.tagtic.cn/ad_gains/kitty/pkg/klog"
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

func (a *Module) ProvideCloser() {
	a.cleanup()
}

func (a *Module) ProvideGrpc(server *grpc.Server) {
	pb.RegisterShareServer(server, a.grpcServer)
}

func (a *Module) ProvideHttp(router *mux.Router) {
	router.PathPrefix("/share/v1/").Handler(http.StripPrefix("/share/v1", a.handler))
}
