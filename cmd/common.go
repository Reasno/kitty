package cmd

import (
	"github.com/go-kit/kit/log"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/spf13/cobra"
	"glab.tagtic.cn/ad_gains/kitty/app/module"
	"glab.tagtic.cn/ad_gains/kitty/pkg/config"
	"glab.tagtic.cn/ad_gains/kitty/pkg/container"
	kittyhttp "glab.tagtic.cn/ad_gains/kitty/pkg/khttp"
	kitty_log "glab.tagtic.cn/ad_gains/kitty/pkg/klog"
	"glab.tagtic.cn/ad_gains/kitty/pkg/ots3"
	"glab.tagtic.cn/ad_gains/kitty/rule"
)

var moduleContainer container.ModuleContainer

func initModules() {
	moduleContainer = container.NewModuleContainer()
	appModuleConfig := conf.Cut("app")
	moduleContainer.Register(module.New(appModuleConfig, logger))
	ruleModuleConfig := conf.Cut("rule")
	moduleContainer.Register(rule.New(ruleModuleConfig, logger))
	moduleContainer.Register(ots3.New(appModuleConfig, logger))
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
