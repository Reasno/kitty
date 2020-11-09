//+build wireinject

package rule

import (
	"github.com/go-kit/kit/log"
	"github.com/google/wire"
	"glab.tagtic.cn/ad_gains/kitty/pkg/config"
	"glab.tagtic.cn/ad_gains/kitty/pkg/contract"
)

var serviceSet = wire.NewSet(
	provideEtcdClient,
	provideRepository,
	wire.Bind(new(Repository), new(*repository)),
	wire.Bind(new(Service), new(*service)),
	wire.Struct(new(service), "*"),
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
