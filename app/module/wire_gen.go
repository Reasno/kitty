// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package module

import (
	"github.com/go-kit/kit/log"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
	"glab.tagtic.cn/ad_gains/kitty/app/handlers"
	"glab.tagtic.cn/ad_gains/kitty/app/repository"
	"glab.tagtic.cn/ad_gains/kitty/pkg/config"
	"glab.tagtic.cn/ad_gains/kitty/pkg/contract"
	"glab.tagtic.cn/ad_gains/kitty/pkg/khttp"
	"glab.tagtic.cn/ad_gains/kitty/pkg/otredis"
	"glab.tagtic.cn/ad_gains/kitty/pkg/ots3"
	"glab.tagtic.cn/ad_gains/kitty/pkg/sms"
	"glab.tagtic.cn/ad_gains/kitty/pkg/wechat"
)

// Injectors from wire.go:

func injectModule(reader contract.ConfigReader, logger log.Logger, dynConf config.DynamicConfigReader) (*Module, func(), error) {
	dialector, err := provideDialector(reader)
	if err != nil {
		return nil, nil, err
	}
	gormConfig := provideGormConfig(logger, reader)
	jaegerLogger := provideJaegerLogAdapter(logger)
	tracer, cleanup, err := provideOpentracing(jaegerLogger, reader)
	if err != nil {
		return nil, nil, err
	}
	db, cleanup2, err := provideGormDB(dialector, gormConfig, tracer)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	securityConfig := provideSecurityConfig(reader)
	appName := config.ProvideAppName(reader)
	env := config.ProvideEnv(reader)
	histogram := provideHistogramMetrics(appName, env)
	moduleOverallMiddleware := provideEndpointsMiddleware(logger, securityConfig, histogram, tracer, env, appName)
	kafkaProducerFactory, cleanup3 := provideKafkaProducerFactory(reader, logger, tracer)
	moduleUserBus := provideUserBus(kafkaProducerFactory, reader)
	moduleEventBus := provideEventBus(kafkaProducerFactory, reader)
	client := provideHttpClient(tracer)
	manager := provideUploadManager(tracer, reader, client)
	fileRepo := repository.NewFileRepo(manager, client)
	userRepo := repository.NewUserRepo(db, fileRepo)
	universalClient, cleanup4 := provideRedis(logger, reader, tracer)
	keyManager := provideKeyManager(appName, env)
	codeRepo := repository.NewCodeRepo(universalClient, keyManager, env)
	extraRepo := repository.NewExtraRepo(universalClient, keyManager)
	senderFactory := sms.NewTransportFactory(reader, client)
	senderFacade := sms.NewSenderFacade(senderFactory, dynConf)
	wechaterFactory := wechat.NewWechaterFactory(reader, client)
	wechaterFacade := wechat.NewWechaterFacade(wechaterFactory, dynConf)
	appService := handlers.NewAppService(reader, logger, userRepo, codeRepo, extraRepo, senderFacade, wechaterFacade)
	appServer := handlers.ProvideAppServer(moduleUserBus, moduleEventBus, appService)
	module := provideModule(db, tracer, logger, moduleOverallMiddleware, appServer, appName)
	return module, func() {
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
	}, nil
}

// wire.go:

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
	provideWechatConfig, wechat.NewWechaterFactory, wechat.NewWechaterFacade, sms.NewTransportFactory, sms.NewSenderFacade, repository.NewUserRepo, repository.NewCodeRepo, repository.NewFileRepo, repository.NewExtraRepo, config.ProvideAppName, config.ProvideEnv, handlers.NewAppService, handlers.ProvideAppServer, wire.Bind(new(redis.Cmdable), new(redis.UniversalClient)), wire.Bind(new(contract.Keyer), new(otredis.KeyManager)), wire.Bind(new(contract.Uploader), new(*ots3.Manager)), wire.Bind(new(contract.HttpDoer), new(*khttp.Client)), wire.Bind(new(wechat.Wechater), new(*wechat.WechaterFacade)), wire.Bind(new(contract.SmsSender), new(*sms.SenderFacade)), wire.Bind(new(contract.Env), new(config.Env)), wire.Bind(new(contract.AppName), new(config.AppName)), wire.Bind(new(handlers.UserRepository), new(*repository.UserRepo)), wire.Bind(new(handlers.CodeRepository), new(*repository.CodeRepo)), wire.Bind(new(handlers.ExtraRepository), new(*repository.ExtraRepo)),
)
