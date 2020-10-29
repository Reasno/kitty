// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package handlers

import (
	"github.com/Reasno/kitty/app/repository"
	"github.com/Reasno/kitty/pkg/contract"
	"github.com/Reasno/kitty/pkg/http"
	"github.com/Reasno/kitty/pkg/ots3"
	"github.com/Reasno/kitty/pkg/sms"
	"github.com/Reasno/kitty/pkg/wechat"
	"github.com/Reasno/kitty/proto"
	"github.com/go-redis/redis/v8"
	"github.com/google/wire"
)

// Injectors from wire.go:

func injectModule(reader contract.ConfigReader) (*AppModule, func(), error) {
	dialector, err := provideDialector(reader)
	if err != nil {
		return nil, nil, err
	}
	logger := provideLogger(reader)
	config := provideGormConfig(logger, reader)
	db, cleanup, err := provideGormDB(dialector, config)
	if err != nil {
		return nil, nil, err
	}
	jaegerLogger := provideJaegerLogAdapter(logger)
	tracer, cleanup2, err := provideOpentracing(jaegerLogger, reader)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	securityConfig := provideSecurityConfig(reader)
	histogram := provideHistogramMetrics(reader)
	handlersOverallMiddleware := provideEndpointsMiddleware(logger, securityConfig, histogram, tracer)
	userRepo := repository.NewUserRepo(db)
	universalClient, cleanup3 := provideRedis(logger, reader)
	codeRepo := repository.NewCodeRepo(universalClient)
	client := provideHttpClient(tracer)
	transportConfig := provideSmsConfig(client, reader)
	transport := sms.NewTransport(transportConfig)
	wechatConfig := provideWechatConfig(reader, client)
	wechatTransport := wechat.NewTransport(wechatConfig)
	manager := provideUploadManager(tracer, reader, client)
	fileRepo := repository.NewFileRepo(manager, client)
	handlersAppService := appService{
		conf:     reader,
		log:      logger,
		ur:       userRepo,
		cr:       codeRepo,
		sender:   transport,
		wechat:   wechatTransport,
		uploader: manager,
		fr:       fileRepo,
	}
	appModule := provideModule(db, tracer, logger, handlersOverallMiddleware, handlersAppService)
	return appModule, func() {
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
	provideLogger,
	provideSmsConfig,
	DbSet,
	OpenTracingSet,
	provideHttpClient,
	provideUploadManager,
	provideRedis,
	provideWechatConfig, wechat.NewTransport, sms.NewTransport, repository.NewUserRepo, repository.NewCodeRepo, repository.NewFileRepo, wire.Struct(new(appService), "*"), wire.Bind(new(redis.Cmdable), new(redis.UniversalClient)), wire.Bind(new(contract.SmsSender), new(*sms.Transport)), wire.Bind(new(contract.Uploader), new(*ots3.Manager)), wire.Bind(new(contract.HttpDoer), new(*http.Client)), wire.Bind(new(kitty.AppServer), new(appService)), wire.Bind(new(UserRepository), new(*repository.UserRepo)), wire.Bind(new(CodeRepository), new(*repository.CodeRepo)), wire.Bind(new(FileRepository), new(*repository.FileRepo)),
)
