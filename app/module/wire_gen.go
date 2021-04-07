// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

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
	"glab.tagtic.cn/ad_gains/kitty/pkg/khttp"
	"glab.tagtic.cn/ad_gains/kitty/pkg/kkafka/client"
	"glab.tagtic.cn/ad_gains/kitty/pkg/otredis"
	"glab.tagtic.cn/ad_gains/kitty/pkg/ots3"
	"glab.tagtic.cn/ad_gains/kitty/pkg/sms"
	"glab.tagtic.cn/ad_gains/kitty/pkg/wechat"
)

// Injectors from wire.go:

func injectModule(reader contract.ConfigReader, logger log.Logger, dynConf config.DynamicConfigReader) (*Module, func(), error) {
	dialector, err := ProvideDialector(reader)
	if err != nil {
		return nil, nil, err
	}
	gormConfig := ProvideGormConfig(logger, reader)
	jaegerLogger := ProvideJaegerLogAdapter(logger)
	tracer, cleanup, err := ProvideOpentracing(jaegerLogger, reader)
	if err != nil {
		return nil, nil, err
	}
	db, cleanup2, err := ProvideGormDB(dialector, gormConfig, tracer)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	securityConfig := ProvideSecurityConfig(reader)
	appName := config.ProvideAppName(reader)
	env := config.ProvideEnv(reader)
	histogram := ProvideHistogramMetrics(appName, env)
	moduleOverallMiddleware := provideEndpointsMiddleware(logger, securityConfig, histogram, tracer, env, appName)
	client := ProvideHttpClient(tracer)
	manager := ProvideUploadManager(tracer, reader, client)
	fileRepo := repository.NewFileRepo(manager, client)
	universalClient, cleanup3 := ProvideRedis(logger, reader, tracer)
	uniqueID := repository.NewUniqueID(universalClient, reader)
	userRepo := repository.NewUserRepo(db, fileRepo, uniqueID)
	keyManager := provideKeyManager(appName, env)
	codeRepo := repository.NewCodeRepo(universalClient, keyManager, env)
	senderFactory := sms.NewTransportFactory(reader, client)
	senderFacade := sms.NewSenderFacade(senderFactory, dynConf)
	wechaterFactory := wechat.NewWechaterFactory(reader, client)
	wechaterFacade := wechat.NewWechaterFacade(wechaterFactory, dynConf)
	kafkaFactory, cleanup4 := ProvideKafkaFactory(reader, logger)
	v := providePublisherOptions(tracer, logger)
	moduleProducerMiddleware := provideProducerMiddleware(tracer, logger)
	dataStore := provideUserBus(kafkaFactory, reader, v, moduleProducerMiddleware)
	eventStore := provideEventBus(kafkaFactory, reader, v, moduleProducerMiddleware)
	dispatcher := ProvideDispatcher(dataStore, eventStore)
	appService := handlers.NewAppService(reader, logger, userRepo, codeRepo, fileRepo, senderFacade, wechaterFacade, dispatcher)
	appServer := handlers.ProvideAppServer(appService)
	module := provideModule(db, tracer, logger, moduleOverallMiddleware, appServer, appName, reader, kafkaFactory)
	return module, func() {
		cleanup4()
		cleanup3()
		cleanup2()
		cleanup()
	}, nil
}

// wire.go:

var DbSet = wire.NewSet(
	ProvideDialector,
	ProvideGormConfig,
	ProvideGormDB,
)

var NameAndEnvSet = wire.NewSet(config.ProvideAppName, config.ProvideEnv, wire.Bind(new(contract.Env), new(config.Env)), wire.Bind(new(contract.AppName), new(config.AppName)))

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
	provideEventBus, wechat.NewWechaterFactory, wechat.NewWechaterFacade, sms.NewTransportFactory, sms.NewSenderFacade, repository.NewUserRepo, repository.NewCodeRepo, repository.NewFileRepo, repository.NewExtraRepo, repository.NewUniqueID, handlers.NewAppService, handlers.ProvideAppServer, wire.Bind(new(redis.Cmdable), new(redis.UniversalClient)), wire.Bind(new(contract.Keyer), new(otredis.KeyManager)), wire.Bind(new(contract.Uploader), new(*ots3.Manager)), wire.Bind(new(contract.HttpDoer), new(*khttp.Client)), wire.Bind(new(listener.UserBus), new(*client.DataStore)), wire.Bind(new(listener.EventBus), new(*client.EventStore)), wire.Bind(new(contract.Dispatcher), new(*event.Dispatcher)), wire.Bind(new(wechat.Wechater), new(*wechat.WechaterFacade)), wire.Bind(new(contract.SmsSender), new(*sms.SenderFacade)), wire.Bind(new(handlers.UserRepository), new(*repository.UserRepo)), wire.Bind(new(handlers.CodeRepository), new(*repository.CodeRepo)), wire.Bind(new(handlers.FileRepository), new(*repository.FileRepo)), wire.Bind(new(entity.IDAssigner), new(*repository.UniqueID)),
)
