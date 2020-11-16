//+build wireinject

package module

import (
	"github.com/go-kit/kit/log"
	"github.com/google/wire"
	"glab.tagtic.cn/ad_gains/kitty/pkg/config"
	"glab.tagtic.cn/ad_gains/kitty/pkg/contract"
	"glab.tagtic.cn/ad_gains/kitty/rule/service"
)

var serviceSet = wire.NewSet(
	provideEtcdClient,
	provideRepository,
	service.ProvideService,
)

func injectModule(reader contract.ConfigReader, logger log.Logger) (*Module, func(), error) {
	panic(wire.Build(
		serviceSet,
		newEndpoints,
		provideModule,
		provideHistogramMetrics,
		config.ProvideAppName,
		config.ProvideEnv,
		wire.Bind(new(contract.Env), new(config.Env)),
		wire.Bind(new(contract.AppName), new(config.AppName)),
	))
}
