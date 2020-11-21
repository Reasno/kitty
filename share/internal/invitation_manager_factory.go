package internal

import (
	"github.com/go-kit/kit/log"
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
		OrientationEvents: conf.Strings("orientation_events"),
		Url:               conf.String("url"),
		Reward: struct {
			Level1 int `yaml:"level1"`
			Level2 int `yaml:"level2"`
		}{
			conf.Int("reward.level1"),
			conf.Int("reward.level2"),
		},
		TaskId: conf.String("task_id"),
	}
	return &InvitationManager{conf: &sc, rr: i.Rr, tokenizer: i.T, xtaskClient: i.C, logger: i.Logger}
}
