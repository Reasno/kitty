package internal

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"glab.tagtic.cn/ad_gains/kitty/pkg/contract"
)

type InvitationManagerFactory struct {
	Rr     RelationRepository
	T      EncodeDecoder
	C      XTaskRequester
	Logger log.Logger
}

func (i *InvitationManagerFactory) NewManager(conf contract.ConfigReader) *InvitationManager {
	sc := ShareConfig{
		Url: conf.String("url"),
		Reward: struct {
			Level1 int `yaml:"level1"`
			Level2 int `yaml:"level2"`
		}{
			conf.Int("reward.level1"),
			conf.Int("reward.level2"),
		},
		TaskId: conf.String("taskId"),
	}
	err := conf.Unmarshal("orientationEvents", &sc.OrientationEvents)
	if err != nil {
		level.Warn(i.Logger).Log("err", err.Error())
	}
	return &InvitationManager{conf: &sc, rr: i.Rr, tokenizer: i.T, xtaskClient: i.C, logger: i.Logger}
}
