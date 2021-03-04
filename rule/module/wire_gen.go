// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package module

import (
	"github.com/go-kit/kit/log"
	"github.com/google/wire"
	"glab.tagtic.cn/ad_gains/kitty/app/module"
	"glab.tagtic.cn/ad_gains/kitty/pkg/config"
	"glab.tagtic.cn/ad_gains/kitty/pkg/contract"
	"glab.tagtic.cn/ad_gains/kitty/rule/service"
)

// Injectors from wire.go:

func injectModule(reader contract.ConfigReader, logger log.Logger) (*Module, func(), error) {
	client, cleanup, err := provideEtcdClient(reader)
	if err != nil {
		return nil, nil, err
	}
	repository, err := provideRepository(client, logger)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	jaegerLogger := module.ProvideJaegerLogAdapter(logger)
	tracer, cleanup2, err := module.ProvideOpentracing(jaegerLogger, reader)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	env := config.ProvideEnv(reader)
	dmpServer, err := provideDmpServer(reader, tracer, logger, env)
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	serviceService := service.ProvideService(logger, repository, dmpServer)
	appName := config.ProvideAppName(reader)
	histogram := provideHistogramMetrics(appName, env)
	endpoints := newEndpoints(serviceService, histogram, logger, appName, env, tracer)
	moduleModule := provideModule(repository, endpoints)
	return moduleModule, func() {
		cleanup2()
		cleanup()
	}, nil
}

// wire.go:

var serviceSet = wire.NewSet(
	provideEtcdClient,
	provideRepository, service.ProvideService,
)
