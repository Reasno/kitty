//+build wireinject

package module

import (
	"github.com/Reasno/kitty/app/handlers"
	"github.com/Reasno/kitty/app/repository"
	"github.com/Reasno/kitty/pkg/config"
	"github.com/Reasno/kitty/pkg/contract"
	kittyhttp "github.com/Reasno/kitty/pkg/khttp"
	"github.com/Reasno/kitty/pkg/otredis"
	"github.com/Reasno/kitty/pkg/ots3"
	"github.com/Reasno/kitty/pkg/sms"
	"github.com/Reasno/kitty/pkg/wechat"
	pb "github.com/Reasno/kitty/proto"
	"github.com/go-kit/kit/log"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
)

var DbSet = wire.NewSet(
	provideDialector,
	provideGormConfig,
	provideGormDB,
)

var OpenTracingSet = wire.NewSet(
	provideJaegerLogAdapter,
	provideOpentracing,
)

var AppServerSet = wire.NewSet(
	provideSmsConfig,
	DbSet,
	OpenTracingSet,
	provideKeyManager,
	provideHttpClient,
	provideUploadManager,
	provideRedis,
	provideWechatConfig,
	wechat.NewTransport,
	sms.NewTransport,
	repository.NewUserRepo,
	repository.NewCodeRepo,
	repository.NewFileRepo,
	repository.NewExtraRepo,
	config.ProvideAppName,
	config.ProvideEnv,
	handlers.NewAppService,
	wire.Bind(new(redis.Cmdable), new(redis.UniversalClient)),
	wire.Bind(new(contract.SmsSender), new(*sms.Transport)),
	wire.Bind(new(contract.Keyer), new(otredis.KeyManager)),
	wire.Bind(new(contract.Uploader), new(*ots3.Manager)),
	wire.Bind(new(contract.HttpDoer), new(*kittyhttp.Client)),
	wire.Bind(new(contract.Env), new(config.Env)),
	wire.Bind(new(contract.AppName), new(config.AppName)),
	wire.Bind(new(handlers.UserRepository), new(*repository.UserRepo)),
	wire.Bind(new(handlers.CodeRepository), new(*repository.CodeRepo)),
	wire.Bind(new(handlers.FileRepository), new(*repository.FileRepo)),
	wire.Bind(new(handlers.ExtraRepository), new(*repository.ExtraRepo)),
)

func injectModule(reader contract.ConfigReader, logger log.Logger) (*Module, func(), error) {
	panic(wire.Build(
		AppServerSet,
		provideKafkaProducerFactory,
		provideUserBus,
		provideEventBus,
		provideSecurityConfig,
		provideHistogramMetrics,
		provideEndpointsMiddleware,
		provideModule,
		handlers.NewMonitoredAppService,
		wire.Bind(new(handlers.EventBus), new(*eventBus)),
		wire.Bind(new(handlers.UserBus), new(*userBus)),
		wire.Bind(new(pb.AppServer), new(*handlers.MonitoredAppService)),
	))
}
