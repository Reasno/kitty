//+build wireinject

package module

import (
	"github.com/go-kit/kit/log"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"glab.tagtic.cn/ad_gains/kitty/app/handlers"
	"glab.tagtic.cn/ad_gains/kitty/app/repository"
	"glab.tagtic.cn/ad_gains/kitty/pkg/config"
	"glab.tagtic.cn/ad_gains/kitty/pkg/contract"
	kittyhttp "glab.tagtic.cn/ad_gains/kitty/pkg/khttp"
	"glab.tagtic.cn/ad_gains/kitty/pkg/otredis"
	"glab.tagtic.cn/ad_gains/kitty/pkg/ots3"
	"glab.tagtic.cn/ad_gains/kitty/pkg/sms"
	"glab.tagtic.cn/ad_gains/kitty/pkg/wechat"
)

var DbSet = wire.NewSet(
	ProvideDialector,
	ProvideGormConfig,
	ProvideGormDB,
)

var OpenTracingSet = wire.NewSet(
	ProvideJaegerLogAdapter,
	ProvideOpentracing,
)

var NameAndEnvSet = wire.NewSet(
	config.ProvideAppName,
	config.ProvideEnv,
	wire.Bind(new(contract.Env), new(config.Env)),
	wire.Bind(new(contract.AppName), new(config.AppName)),
)

var AppServerSet = wire.NewSet(
	provideSmsConfig,
	DbSet,
	OpenTracingSet,
	NameAndEnvSet,
	provideKeyManager,
	ProvideHttpClient,
	ProvideUploadManager,
	ProvideRedis,
	provideWechatConfig,
	wechat.NewWechaterFactory,
	wechat.NewWechaterFacade,
	sms.NewTransportFactory,
	sms.NewSenderFacade,
	repository.NewUserRepo,
	repository.NewCodeRepo,
	repository.NewFileRepo,
	repository.NewExtraRepo,
	handlers.NewAppService,
	handlers.ProvideAppServer,
	wire.Bind(new(redis.Cmdable), new(redis.UniversalClient)),
	wire.Bind(new(contract.Keyer), new(otredis.KeyManager)),
	wire.Bind(new(contract.Uploader), new(*ots3.Manager)),
	wire.Bind(new(contract.HttpDoer), new(*kittyhttp.Client)),
	wire.Bind(new(wechat.Wechater), new(*wechat.WechaterFacade)),
	wire.Bind(new(contract.SmsSender), new(*sms.SenderFacade)),
	wire.Bind(new(handlers.UserRepository), new(*repository.UserRepo)),
	wire.Bind(new(handlers.CodeRepository), new(*repository.CodeRepo)),
)

func injectModule(reader contract.ConfigReader, logger log.Logger, dynConf config.DynamicConfigReader) (*Module, func(), error) {
	panic(wire.Build(
		AppServerSet,
		provideKafkaProducerFactory,
		provideUserBus,
		provideEventBus,
		ProvideSecurityConfig,
		ProvideHistogramMetrics,
		provideEndpointsMiddleware,
		provideModule,
		wire.Bind(new(handlers.EventBus), new(*eventBus)),
		wire.Bind(new(handlers.UserBus), new(*userBus)),
	))
}
