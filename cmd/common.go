package cmd

import (
	"fmt"

	"github.com/go-kit/kit/log"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/spf13/cobra"
	app "glab.tagtic.cn/ad_gains/kitty/app/module"
	"glab.tagtic.cn/ad_gains/kitty/pkg/config"
	"glab.tagtic.cn/ad_gains/kitty/pkg/container"
	kittyhttp "glab.tagtic.cn/ad_gains/kitty/pkg/khttp"
	kitty_log "glab.tagtic.cn/ad_gains/kitty/pkg/klog"
	"glab.tagtic.cn/ad_gains/kitty/pkg/ots3"
	"glab.tagtic.cn/ad_gains/kitty/rule/client"
	module2 "glab.tagtic.cn/ad_gains/kitty/rule/module"
	share "glab.tagtic.cn/ad_gains/kitty/share/module"
)

var moduleContainer container.ModuleContainer

func initModules() {
	moduleContainer = container.NewModuleContainer()
	ruleModuleConfig := conf.Cut("rule")
	ruleModule := module2.New(ruleModuleConfig, logger)
	moduleContainer.Register(ruleModule)

	dynConf, err := client.NewRuleEngine(client.WithRepository(ruleModule.GetRepository()))
	if err != nil {
		panic(err)
	}

	appModuleConfig := conf.Cut("app")
	appModuleDynConfig := dynConf.Of(
		fmt.Sprintf("%s-%s",
			appModuleConfig.String("name"),
			appModuleConfig.String("env")),
	)
	moduleContainer.Register(app.New(appModuleConfig, logger, appModuleDynConfig))
	moduleContainer.Register(share.New(appModuleConfig, logger, appModuleDynConfig))
	moduleContainer.Register(ots3.New(conf.Cut("global"), logger))
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
