package module

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	"github.com/oklog/run"
	"glab.tagtic.cn/ad_gains/kitty/pkg/config"
	"glab.tagtic.cn/ad_gains/kitty/pkg/contract"
	pb "glab.tagtic.cn/ad_gains/kitty/proto"
	"glab.tagtic.cn/ad_gains/kitty/share/consumer"
	"google.golang.org/grpc"
)

type Module struct {
	appName       contract.AppName
	cleanup       func()
	handler       http.Handler
	grpcServer    GrpcShareServer
	eventReciever consumer.EventReceiver
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

func (a *Module) ProvideRunGroup(group *run.Group) {
	{
		ctx, cancel := context.WithCancel(context.Background())
		group.Add(func() error {
			return a.eventReciever.SubscribeCheckin(ctx)
		}, func(err error) {
			cancel()
		})
	}

	{
		ctx, cancel := context.WithCancel(context.Background())
		group.Add(func() error {
			return a.eventReciever.SubscribeTask(ctx)
		}, func(err error) {
			cancel()
		})
	}
}
