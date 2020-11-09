package cmd

import (
	"github.com/Reasno/kitty/app/module"
	"github.com/Reasno/kitty/pkg/config"
	"github.com/Reasno/kitty/pkg/container"
	kittyhttp "github.com/Reasno/kitty/pkg/khttp"
	kitty_log "github.com/Reasno/kitty/pkg/klog"
	"github.com/Reasno/kitty/rule"
	"github.com/go-kit/kit/log"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/spf13/cobra"
)

var moduleContainer container.ModuleContainer

func initModules() {
	moduleContainer = container.NewModuleContainer()
	appModuleConfig := conf.Cut("app")
	moduleContainer.Register(module.New(appModuleConfig, logger))
	ruleModuleConfig := conf.Cut("rule")
	moduleContainer.Register(rule.New(ruleModuleConfig, logger))
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

func initConfig(_ *cobra.Command, _ []string) error {
	k := koanf.New(".")
	if cfgFile == "" {
		cfgFile = "./config/kitty.yaml"
	}

	err := k.Load(file.Provider("./config/kitty.yaml"), yaml.Parser())
	if err != nil {
		panic(err)
	}

	conf = config.NewKoanfAdapter(k)
	return nil
}

func initLogger(cmd *cobra.Command, _ []string) error {
	logger = kitty_log.NewLogger(config.Env(conf.String("global.env")))
	logger = log.With(logger, "subcommand", cmd.Use)
	return nil
}
