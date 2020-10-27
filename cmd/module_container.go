package cmd

import (
	"github.com/Reasno/kitty/app/handlers"
	"github.com/Reasno/kitty/pkg/config"
	"github.com/Reasno/kitty/pkg/container"
	kittyhttp "github.com/Reasno/kitty/pkg/http"
)

var moduleContainer container.ModuleContainer

func initModules() {
	moduleContainer = container.NewModuleContainer()
	appModuleConfig, err := config.ProvideChildConfig("app", "global")
	if err != nil {
		panic(err)
	}
	moduleContainer.Register(handlers.New(appModuleConfig))
	moduleContainer.Register(container.HttpFunc(kittyhttp.Doc))
	moduleContainer.Register(container.HttpFunc(kittyhttp.HealthCheck))
	moduleContainer.Register(container.HttpFunc(kittyhttp.Metrics))
	moduleContainer.Register(container.HttpFunc(kittyhttp.Debug))
}

func shutdownModules() {
	for _, f := range moduleContainer.CloserProviders {
		f()
	}
}
