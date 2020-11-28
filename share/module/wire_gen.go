// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package module

import (
	"github.com/go-kit/kit/log"
	"github.com/google/wire"
	"glab.tagtic.cn/ad_gains/kitty/app/module"
	"glab.tagtic.cn/ad_gains/kitty/app/repository"
	"glab.tagtic.cn/ad_gains/kitty/pkg/config"
	"glab.tagtic.cn/ad_gains/kitty/pkg/contract"
	"glab.tagtic.cn/ad_gains/kitty/pkg/khttp"
	"glab.tagtic.cn/ad_gains/kitty/pkg/ots3"
	"glab.tagtic.cn/ad_gains/kitty/share/consumer"
	"glab.tagtic.cn/ad_gains/kitty/share/handlers"
	"glab.tagtic.cn/ad_gains/kitty/share/internal"
)

// Injectors from wire.go:

func injectModule(reader contract.ConfigReader, logger log.Logger, dynConf config.DynamicConfigReader) (*Module, func(), error) {
	securityConfig := module.ProvideSecurityConfig(reader)
	appName := config.ProvideAppName(reader)
	env := config.ProvideEnv(reader)
	histogram := module.ProvideHistogramMetrics(appName, env)
	jaegerLogger := module.ProvideJaegerLogAdapter(logger)
	tracer, cleanup, err := module.ProvideOpentracing(jaegerLogger, reader)
	if err != nil {
		return nil, nil, err
	}
	moduleOverallMiddleware := provideEndpointsMiddleware(logger, securityConfig, histogram, tracer, env, appName)
	dialector, err := module.ProvideDialector(reader)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	gormConfig := module.ProvideGormConfig(logger, reader)
	db, cleanup2, err := module.ProvideGormDB(dialector, gormConfig, tracer)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	relationRepo := repository.NewRelationRepo(db)
	tokenizer := provideTokenizer(reader)
	client := module.ProvideHttpClient(tracer)
	xTaskRequester, err := internal.NewXTaskRequester(reader, client)
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	invitationManagerFactory := internal.InvitationManagerFactory{
		Rr:     relationRepo,
		T:      tokenizer,
		C:      xTaskRequester,
		Logger: logger,
	}
	invitationManagerFacade := &internal.InvitationManagerFacade{
		Name:    appName,
		Factory: invitationManagerFactory,
		DynConf: dynConf,
	}
	manager := module.ProvideUploadManager(tracer, reader, client)
	fileRepo := repository.NewFileRepo(manager, client)
	userRepo := repository.NewUserRepo(db, fileRepo)
	shareService := handlers.NewShareService(invitationManagerFacade, userRepo)
	shareServer := handlers.ProvideShareServer(shareService)
	endpoints := provideEndpoints(moduleOverallMiddleware, shareServer)
	grpcShareServer := provideGrpc(endpoints, tracer, logger, appName)
	handler := provideHttp(endpoints, tracer, logger, appName)
	kafkaFactory, cleanup3 := module.ProvideKafkaFactory(reader, logger, tracer)
	middleware := provideKafkaMiddleware(tracer)
	eventReceiver := consumer.EventReceiver{
		AppName: appName,
		Conf:    reader,
		Manager: invitationManagerFacade,
		Factory: kafkaFactory,
		MW:      middleware,
	}
	moduleModule := provideModule(grpcShareServer, handler, eventReceiver, appName)
	return moduleModule, func() {
		cleanup3()
		cleanup2()
		cleanup()
	}, nil
}

// wire.go:

var ShareServiceSet = wire.NewSet(module.DbSet, module.OpenTracingSet, module.NameAndEnvSet, module.ProvideSecurityConfig, module.ProvideKafkaFactory, module.ProvideHistogramMetrics, module.ProvideHttpClient, module.ProvideUploadManager, repository.NewUserRepo, repository.NewRelationRepo, repository.NewFileRepo, provideTokenizer, internal.NewXTaskRequester, handlers.NewShareService, handlers.ProvideShareServer, provideKafkaMiddleware, wire.Struct(new(internal.InvitationManagerFactory), "*"), wire.Struct(new(internal.InvitationManagerFacade), "*"), wire.Struct(new(consumer.EventReceiver), "*"), wire.Bind(new(handlers.UserRepository), new(*repository.UserRepo)), wire.Bind(new(internal.RelationRepository), new(*repository.RelationRepo)), wire.Bind(new(handlers.InvitationManager), new(*internal.InvitationManagerFacade)), wire.Bind(new(consumer.InvitationManager), new(*internal.InvitationManagerFacade)), wire.Bind(new(contract.Uploader), new(*ots3.Manager)), wire.Bind(new(contract.HttpDoer), new(*khttp.Client)), wire.Bind(new(internal.EncodeDecoder), new(*internal.Tokenizer)))
