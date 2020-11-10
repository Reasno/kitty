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
	wechat.NewWechaterFactory,
	wechat.NewWechaterFacade,
	sms.NewTransportFactory,
	sms.NewSenderFacade,
	repository.NewUserRepo,
	repository.NewCodeRepo,
	repository.NewFileRepo,
	repository.NewExtraRepo,
	config.ProvideAppName,
	config.ProvideEnv,
	handlers.NewAppService,
	handlers.ProvideAppServer,
	wire.Bind(new(redis.Cmdable), new(redis.UniversalClient)),
	wire.Bind(new(contract.Keyer), new(otredis.KeyManager)),
	wire.Bind(new(contract.Uploader), new(*ots3.Manager)),
	wire.Bind(new(contract.HttpDoer), new(*kittyhttp.Client)),
	wire.Bind(new(wechat.Wechater), new(*wechat.WechaterFacade)),
	wire.Bind(new(contract.SmsSender), new(*sms.SenderFacade)),
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
		wire.Bind(new(handlers.EventBus), new(*eventBus)),
		wire.Bind(new(handlers.UserBus), new(*userBus)),
	))
}
