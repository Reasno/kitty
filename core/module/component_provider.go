package module

import (
	"github.com/go-kit/kit/log"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/pkg/errors"
	"glab.tagtic.cn/ad_gains/kitty/pkg/config"
	"glab.tagtic.cn/ad_gains/kitty/pkg/contract"
	kitty_log "glab.tagtic.cn/ad_gains/kitty/pkg/klog"
)

func ProvideConfig(cfgFile string) (contract.ConfigReader, error) {
	k := koanf.New(".")
	if cfgFile == "" {
		cfgFile = "./config/kitty.yaml"
	}

	err := k.Load(file.Provider("./config/kitty.yaml"), yaml.Parser())
	if err != nil {
		return nil, errors.Wrap(err, "unable to load config")
	}

	return config.NewKoanfAdapter(k), nil
}

func ProvideLogger(conf contract.ConfigReader) log.Logger {
	logger := kitty_log.NewLogger(config.Env(conf.String("global.env")))
	return logger
}
