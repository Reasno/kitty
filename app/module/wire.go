//+build wireinject

package module

import (
	"github.com/go-kit/kit/log"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"glab.tagtic.cn/ad_gains/kitty/app/entity"
	"glab.tagtic.cn/ad_gains/kitty/app/handlers"
	"glab.tagtic.cn/ad_gains/kitty/app/listener"
	"glab.tagtic.cn/ad_gains/kitty/app/repository"
	"glab.tagtic.cn/ad_gains/kitty/pkg/config"
	"glab.tagtic.cn/ad_gains/kitty/pkg/contract"
	"glab.tagtic.cn/ad_gains/kitty/pkg/event"
	kittyhttp "glab.tagtic.cn/ad_gains/kitty/pkg/khttp"
	kclient "glab.tagtic.cn/ad_gains/kitty/pkg/kkafka/client"
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

var NameAndEnvSet = wire.NewSet(
	config.ProvideAppName,
	config.ProvideEnv,
	wire.Bind(new(contract.Env), new(config.Env)),
	wire.Bind(new(contract.AppName), new(config.AppName)),
)

var OpenTracingSet = wire.NewSet(
	ProvideJaegerLogAdapter,
	ProvideOpentracing,
)

var AppServerSet = wire.NewSet(
	provideSmsConfig,
	DbSet,
	OpenTracingSet,
	NameAndEnvSet,
	provideKeyManager,
	ProvideHttpClient,
	ProvideUploadManager,
	ProvideDispatcher,
	ProvideRedis,
	provideWechatConfig,
	provideUserBus,
	providePublisherOptions,
	ProvideKafkaFactory,
	provideEventBus,
	wechat.NewWechaterFactory,
	wechat.NewWechaterFacade,
	sms.NewTransportFactory,
	sms.NewSenderFacade,
	repository.NewUserRepo,
	repository.NewCodeRepo,
	repository.NewFileRepo,
	repository.NewExtraRepo,
	repository.NewUniqueID,
	handlers.NewAppService,
	handlers.ProvideAppServer,
	wire.Bind(new(redis.Cmdable), new(redis.UniversalClient)),
	wire.Bind(new(contract.Keyer), new(otredis.KeyManager)),
	wire.Bind(new(contract.Uploader), new(*ots3.Manager)),
	wire.Bind(new(contract.HttpDoer), new(*kittyhttp.Client)),
	wire.Bind(new(listener.UserBus), new(*kclient.DataStore)),
	wire.Bind(new(listener.EventBus), new(*kclient.EventStore)),
	wire.Bind(new(contract.Dispatcher), new(*event.Dispatcher)),
	wire.Bind(new(wechat.Wechater), new(*wechat.WechaterFacade)),
	wire.Bind(new(contract.SmsSender), new(*sms.SenderFacade)),
	wire.Bind(new(handlers.UserRepository), new(*repository.UserRepo)),
	wire.Bind(new(handlers.CodeRepository), new(*repository.CodeRepo)),
	wire.Bind(new(handlers.FileRepository), new(*repository.FileRepo)),
	wire.Bind(new(entity.IDAssigner), new(*repository.UniqueID)),
)

func injectModule(reader contract.ConfigReader, logger log.Logger, dynConf config.DynamicConfigReader) (*Module, func(), error) {
	panic(wire.Build(
		AppServerSet,
		ProvideSecurityConfig,
		ProvideHistogramMetrics,
		provideEndpointsMiddleware,
		provideProducerMiddleware,
		provideModule,
	))
}
