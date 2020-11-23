package module

import (
	"github.com/go-kit/kit/log"
	"glab.tagtic.cn/ad_gains/kitty/pkg/contract"
)

type Module struct {
	Conf   contract.ConfigReader
	Logger log.Logger
}

func New(cfgFile string) *Module {
	conf, err := ProvideConfig(cfgFile)
	if err != nil {
		panic(err)
	}
	logger := ProvideLogger(conf)
	return &Module{
		Conf:   conf,
		Logger: logger,
	}
}
