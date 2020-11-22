package cmd

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"glab.tagtic.cn/ad_gains/kitty/pkg/config"
	"glab.tagtic.cn/ad_gains/kitty/pkg/contract"
	"glab.tagtic.cn/ad_gains/kitty/pkg/klog"

	app "glab.tagtic.cn/ad_gains/kitty/app/module"
	"glab.tagtic.cn/ad_gains/kitty/pkg/container"
	kittyhttp "glab.tagtic.cn/ad_gains/kitty/pkg/khttp"
	"glab.tagtic.cn/ad_gains/kitty/pkg/ots3"
	"glab.tagtic.cn/ad_gains/kitty/rule/client"
	rule "glab.tagtic.cn/ad_gains/kitty/rule/module"
	share "glab.tagtic.cn/ad_gains/kitty/share/module"
)

var moduleContainer container.ModuleContainer

func initModules() {
	moduleContainer = container.NewModuleContainer()
	ruleModule := rule.New(name("rule"))
	engine, _ := client.NewRuleEngine(client.WithRepository(ruleModule.GetRepository()))

	moduleContainer.Register(coreModule)
	moduleContainer.Register(ruleModule)
	moduleContainer.Register(app.New(nameD("app", engine)))
	moduleContainer.Register(share.New(nameD("app", engine)))
	moduleContainer.Register(ots3.New(name("s3")))
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

func warn(msg string) {
	_ = level.Warn(coreModule.Logger).Log("msg", msg)
}

func er(err error) {
	_ = level.Error(coreModule.Logger).Log("err", err)
}

func debug(msg string) {
	_ = level.Debug(coreModule.Logger).Log("msg", msg)
}

func info(msg string) {
	_ = level.Info(coreModule.Logger).Log("msg", msg)
}

func conf() contract.ConfigReader {
	return coreModule.StaticConf
}

// name unpacks the core module to several dependencies for other modules
func name(name string) (contract.ConfigReader, log.Logger) {
	m := coreModule
	conf := m.StaticConf.Cut(name)
	logger := log.With(m.Logger, "module", conf.String("name"))
	logger = level.NewFilter(logger, klog.LevelFilter(conf.String("level")))
	return conf, logger
}

// nameD like name, but also provide config.DynamicConfigReader
func nameD(name string, client *client.RuleEngine) (contract.ConfigReader, log.Logger, config.DynamicConfigReader) {
	m := coreModule
	conf := m.StaticConf.Cut(name)
	logger := log.With(m.Logger, "module", conf.String("name"))
	logger = level.NewFilter(logger, klog.LevelFilter(conf.String("level")))
	dyn := client.Of(
		fmt.Sprintf("%s-%s",
			conf.String("name"),
			conf.String("env"),
		),
	)
	return conf, logger, dyn
}
