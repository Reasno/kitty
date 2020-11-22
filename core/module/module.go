package module

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"glab.tagtic.cn/ad_gains/kitty/pkg/config"
	"glab.tagtic.cn/ad_gains/kitty/pkg/contract"
	"glab.tagtic.cn/ad_gains/kitty/pkg/klog"
	"glab.tagtic.cn/ad_gains/kitty/rule/client"
)

type Module struct {
	StaticConf contract.ConfigReader
	DynConf    config.DynamicConfigReader
	Logger     log.Logger
}

func New(cfgFile string) *Module {
	conf, err := ProvideConfig(cfgFile)
	if err != nil {
		panic(err)
	}
	logger := ProvideLogger(conf)
	return &Module{
		StaticConf: conf,
		Logger:     logger,
	}
}

func (m Module) Make(name string) (contract.ConfigReader, log.Logger) {
	conf := m.StaticConf.Cut(name)
	logger := log.With(m.Logger, "module", conf.String("name"))
	logger = level.NewFilter(logger, klog.LevelFilter(conf.String("level")))
	return conf, logger
}

func (m Module) MakeWithEngine(name string, client *client.RuleEngine) (contract.ConfigReader, log.Logger, config.DynamicConfigReader) {
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
