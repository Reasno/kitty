// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package rule

import (
	"github.com/Reasno/kitty/pkg/config"
	"github.com/Reasno/kitty/pkg/contract"
	"github.com/go-kit/kit/log"
	"github.com/google/wire"
)

// Injectors from wire.go:

func injectModule(reader contract.ConfigReader, logger log.Logger) (*Module, func(), error) {
	client, cleanup, err := provideEtcdClient(reader)
	if err != nil {
		return nil, nil, err
	}
	ruleRepository, err := provideRepository(client, logger, reader)
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	ruleService := &service{
		logger: logger,
		repo:   ruleRepository,
	}
	appName := config.ProvideAppName(reader)
	env := config.ProvideEnv(reader)
	histogram := provideHistogramMetrics(appName, env)
	endpoints := newEndpoints(ruleService, histogram, logger, appName, env)
	module := provideModule(ruleRepository, endpoints)
	return module, func() {
		cleanup()
	}, nil
}

// wire.go:

var serviceSet = wire.NewSet(
	provideEtcdClient,
	provideRepository, wire.Bind(new(Repository), new(*repository)), wire.Bind(new(Service), new(*service)), wire.Struct(new(service), "*"),
)
